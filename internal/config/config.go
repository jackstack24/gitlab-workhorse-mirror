package config

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

type TomlURL struct {
	url.URL
}

func (u *TomlURL) UnmarshalText(text []byte) error {
	temp, err := url.Parse(string(text))
	u.URL = *temp
	return err
}

type TomlDuration struct {
	time.Duration
}

func (d *TomlDuration) UnmarshalTest(text []byte) error {
	temp, err := time.ParseDuration(string(text))
	d.Duration = temp
	return err
}

type ObjectStorageConfig struct {
	Enabled  bool
	Provider string

	S3Config S3Config `toml:"s3"`
}

type S3Config struct {
	AwsAccessKeyID     string `toml:"aws_access_key_id"`
	AwsSecretAccessKey string `toml:"aws_secret_access_key"`
	Region             string
	Bucket             string
	PathStyle          bool `toml:"path_style"`
	Endpoint           string
}

type RedisConfig struct {
	URL             TomlURL
	Sentinel        []TomlURL
	SentinelMaster  string
	Password        string
	DB              *int
	ReadTimeout     *TomlDuration
	WriteTimeout    *TomlDuration
	KeepAlivePeriod *TomlDuration
	MaxIdle         *int
	MaxActive       *int
}

type Config struct {
	Redis                    *RedisConfig                    `toml:"redis"`
	Backend                  *url.URL                        `toml:"-"`
	CableBackend             *url.URL                        `toml:"-"`
	Version                  string                          `toml:"-"`
	DocumentRoot             string                          `toml:"-"`
	DevelopmentMode          bool                            `toml:"-"`
	Socket                   string                          `toml:"-"`
	CableSocket              string                          `toml:"-"`
	ProxyHeadersTimeout      time.Duration                   `toml:"-"`
	APILimit                 uint                            `toml:"-"`
	APIQueueLimit            uint                            `toml:"-"`
	APIQueueTimeout          time.Duration                   `toml:"-"`
	APICILongPollingDuration time.Duration                   `toml:"-"`
	ObjectStorages           map[string]*ObjectStorageConfig `toml:"object_storage"`
}

// LoadConfig from a file
func LoadConfig(filename string) (*Config, error) {
	cfg := &Config{}
	if _, err := toml.DecodeFile(filename, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) FindObjectStorageConfig(name string) (*ObjectStorageConfig, error) {
	cfg, ok := c.ObjectStorages[name]

	if ok {
		return cfg, nil
	}

	return nil, fmt.Errorf("unable to find object storage %s", name)
}

func (c *ObjectStorageConfig) IsAWS() bool {
	return strings.EqualFold(c.Provider, "AWS") || strings.EqualFold(c.Provider, "S3")
}

func (c *ObjectStorageConfig) IsValid() bool {
	if c.Enabled && c.IsAWS() {
		return c.S3Config.Bucket != "" && c.S3Config.Region != ""
	}

	return false
}
