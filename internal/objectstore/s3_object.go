package objectstore

import (
	"context"
	"io"
	"time"

	"gitlab.com/gitlab-org/gitlab-workhorse/internal/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"gitlab.com/gitlab-org/labkit/log"
)

type S3Object struct {
	session    *session.Session
	s3Config   config.S3Config
	objectName string
	uploader
}

func NewS3Object(ctx context.Context, objectName string, s3Config config.S3Config, deadline time.Time) (*S3Object, error) {
	pr, pw := io.Pipe()
	objectStorageUploadsOpen.Inc()
	uploadCtx, cancelFn := context.WithDeadline(ctx, deadline)

	o := &S3Object{
		uploader: newUploader(uploadCtx, pw),
		s3Config: s3Config,
	}

	go o.trackUploadTime()
	go o.cleanup(ctx)

	go func() {
		defer cancelFn()
		defer objectStorageUploadsOpen.Dec()
		defer func() {
			// This will be returned as error to the next write operation on the pipe
			pr.CloseWithError(o.uploadError)
		}()

		config := &aws.Config{
			Region:           aws.String(s3Config.Region),
			S3ForcePathStyle: aws.Bool(s3Config.PathStyle),
		}

		// In case IAM profiles aren't being used, use the static credentials
		if s3Config.AwsAccessKeyID != "" || s3Config.AwsSecretAccessKey != "" {
			config.Credentials = credentials.NewStaticCredentials(s3Config.AwsAccessKeyID, s3Config.AwsSecretAccessKey, "")
		}

		if s3Config.Endpoint != "" {
			config.Endpoint = aws.String(s3Config.Endpoint)
		}

		sess, err := session.NewSession(config)
		if err != nil {
			o.uploadError = err
			log.WithError(err).Errorf("error creating S3 session: %v", err)
			return
		}

		o.session = sess
		o.objectName = objectName
		uploader := s3manager.NewUploader(o.session)

		_, err = uploader.UploadWithContext(uploadCtx, &s3manager.UploadInput{
			Bucket: aws.String(s3Config.Bucket),
			Key:    aws.String(objectName),
			Body:   pr,
		})
		if err != nil {
			o.uploadError = err
			objectStorageUploadRequestsRequestFailed.Inc()
			log.WithError(err).Errorf("error uploading S3 file: %v", err)
			return
		}
	}()

	return o, nil
}

func (o *S3Object) trackUploadTime() {
	started := time.Now()
	<-o.ctx.Done()
	objectStorageUploadTime.Observe(time.Since(started).Seconds())
}

func (o *S3Object) cleanup(ctx context.Context) {
	// wait for the upload to finish
	<-o.ctx.Done()

	if o.uploadError != nil {
		objectStorageUploadRequestsRequestFailed.Inc()
		o.delete()
		return
	}

	// We have now successfully uploaded the file to object storage. Another
	// goroutine will hand off the object to gitlab-rails.
	<-ctx.Done()

	// gitlab-rails is now done with the object so it's time to delete it.
	o.delete()
}

func (o *S3Object) delete() {
	if o.session == nil || o.objectName == "" {
		return
	}

	svc := s3.New(o.session)
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(o.s3Config.Bucket),
		Key:    aws.String(o.objectName),
	}

	// Note we can't use the request context because in a successful
	// case, the original request has already completed.
	deleteCtx, cancel := context.WithTimeout(context.Background(), 60*time.Second) // lint:allow context.Background
	defer cancel()

	_, err := svc.DeleteObjectWithContext(deleteCtx, input)

	if err != nil {
		log.WithError(err).Errorf("error deleting S3 object: %v", err)
	}
}
