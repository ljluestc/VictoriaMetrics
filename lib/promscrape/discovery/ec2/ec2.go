package ec2

import (
	"flag"
	"fmt"
	"time"

	"github.com/VictoriaMetrics/VictoriaMetrics/lib/awsapi"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/promauth"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/promutil"
)

// SDCheckInterval defines interval for targets refresh.
var SDCheckInterval = flag.Duration("promscrape.ec2SDCheckInterval", time.Minute, "Interval for checking for changes in ec2. "+
	"This works only if ec2_sd_configs is configured in '-promscrape.config' file. "+
	"See https://docs.victoriametrics.com/victoriametrics/sd_configs/#ec2_sd_configs for details")

// SDConfig represents service discovery config for ec2.
//
// See https://prometheus.io/docs/prometheus/latest/configuration/configuration/#ec2_sd_config
type SDConfig struct {
	Region      string           `yaml:"region,omitempty"`
	Endpoint    string           `yaml:"endpoint,omitempty"`
	STSEndpoint string           `yaml:"sts_endpoint,omitempty"`
	AccessKey   string           `yaml:"access_key,omitempty"`
	SecretKey   *promauth.Secret `yaml:"secret_key,omitempty"`
	Profile     string           `yaml:"profile,omitempty"` // New field for named AWS profile
	RoleARN     string           `yaml:"role_arn,omitempty"`
	// RefreshInterval time.Duration `yaml:"refresh_interval"`
	// refresh_interval is obtained from `-promscrape.ec2SDCheck