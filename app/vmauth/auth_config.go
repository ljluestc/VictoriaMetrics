package main

type UserInfo struct {
	Username               string
	Password               string
	BearerToken            string
	URLPrefix              *URLPrefix
	URLMaps                []URLMap
	MetricLabels           map[string]string
	MaxConcurrentRequests  int
	TLSInsecureSkipVerify  *bool
	TLSServerName          string
	TLSCAFile              string
	TLSCertFile            string
	TLSKeyFile             string
	DefaultURL             *URLPrefix
	RetryStatusCodes       []int
	LoadBalancingPolicy    string
	DropSrcPathPrefixParts int
	DiscoverBackendIPs     *bool
	HeadersConf            HeadersConf
}
