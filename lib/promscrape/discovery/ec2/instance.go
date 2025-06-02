package ec2

import (
	"context"
	"encoding/xml"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/VictoriaMetrics/VictoriaMetrics/lib/awsapi"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/logger"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/promscrape/discoveryutil"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/promutil"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
)
// getInstancesLabels returns labels for ec2 instances obtained from the given cfg
func getInstancesLabels(cfg *apiConfig) ([]*promutil.Labels, error) {
	rs, err := getReservations(cfg)
	if err != nil {
		return nil, err
	}
	azMap := getAZMap(cfg)
	region := cfg.awsConfig.GetRegion()
	var ms []*promutil.Labels
	for _, r := range rs {
		for _, inst := range r.InstanceSet.Items {
			ms = inst.appendTargetLabels(ms, r.OwnerID, region, cfg.port, azMap)
		}
	}
	return ms, nil
}

func getReservations(cfg *apiConfig) ([]Reservation, error) {
	// See https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_DescribeInstances.html
	var rs []Reservation
	pageToken := ""
	instanceFilters := awsapi.GetFiltersQueryString(cfg.instanceFilters)
	for {
		data, err := cfg.awsConfig.GetEC2APIResponse("DescribeInstances", instanceFilters, pageToken)
		if err != nil {
			return nil, fmt.Errorf("cannot obtain instances: %w", err)
		}
		ir, err := parseInstancesResponse(data)
		if err != nil {
			return nil, fmt.Errorf("cannot parse instance list: %w", err)
		}
		rs = append(rs, ir.ReservationSet.Items...)
		if len(ir.NextPageToken) == 0 {
			return rs, nil
		}
		pageToken = ir.NextPageToken
	}
}

// InstancesResponse represents response to https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_DescribeInstances.html
type InstancesResponse struct {
	ReservationSet ReservationSet `xml:"reservationSet"`
	NextPageToken  string         `xml:"nextToken"`
}

// ReservationSet represetns ReservationSet from https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_DescribeInstances.html
type ReservationSet struct {
	Items []Reservation `xml:"item"`
}

// Reservation represents Reservation from https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_Reservation.html
type Reservation struct {
	OwnerID     string      `xml:"ownerId"`
	InstanceSet InstanceSet `xml:"instancesSet"`
}

// InstanceSet represents InstanceSet from https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_Reservation.html
type InstanceSet struct {
	Items []Instance `xml:"item"`
}

// Instance represents Instance from https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_Instance.html
type Instance struct {
	PrivateIPAddress    string              `xml:"privateIpAddress"`
	Architecture        string              `xml:"architecture"`
	Placement           Placement           `xml:"placement"`
	ImageID             string              `xml:"imageId"`
	ID                  string              `xml:"instanceId"`
	Lifecycle           string              `xml:"instanceLifecycle"`
	State               InstanceState       `xml:"instanceState"`
	Type                string              `xml:"instanceType"`
	Platform            string              `xml:"platform"`
	SubnetID            string              `xml:"subnetId"`
	PrivateDNSName      string              `xml:"privateDnsName"`
	PublicDNSName       string              `xml:"dnsName"`
	PublicIPAddress     string              `xml:"ipAddress"`
	VPCID               string              `xml:"vpcId"`
	NetworkInterfaceSet NetworkInterfaceSet `xml:"networkInterfaceSet"`
	TagSet              TagSet              `xml:"tagSet"`
}

// Placement represents Placement from https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_Placement.html
type Placement struct {
	AvailabilityZone string `xml:"availabilityZone"`
}

// InstanceState represents InstanceState from https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_InstanceState.html
type InstanceState struct {
	Name string `xml:"name"`
}

// NetworkInterfaceSet represents NetworkInterfaceSet from https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_Instance.html
type NetworkInterfaceSet struct {
	Items []NetworkInterface `xml:"item"`
}

// NetworkInterface represents NetworkInterface from https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_InstanceNetworkInterface.html
type NetworkInterface struct {
	SubnetID         string           `xml:"subnetId"`
	IPv6AddressesSet Ipv6AddressesSet `xml:"ipv6AddressesSet"`
}

// Ipv6AddressesSet represents ipv6AddressesSet from https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_InstanceNetworkInterface.html
type Ipv6AddressesSet struct {
	Items []string `xml:"item"`
}

// TagSet represents TagSet from https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_Instance.html
type TagSet struct {
	Items []Tag `xml:"item"`
}

// Tag represents Tag from https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_Tag.html
type Tag struct {
	Key   string `xml:"key"`
	Value string `xml:"value"`
}

func parseInstancesResponse(data []byte) (*InstancesResponse, error) {
	var v InstancesResponse
	if err := xml.Unmarshal(data, &v); err != nil {
		return nil, fmt.Errorf("cannot unmarshal InstancesResponse from %q: %w", data, err)
	}
	return &v, nil
}

func (inst *Instance) appendTargetLabels(ms []*promutil.Labels, ownerID, region string, port int, azMap map[string]string) []*promutil.Labels {
	if len(inst.PrivateIPAddress) == 0 {
		// Cannot scrape instance without private IP address
		return ms
	}
	addr := discoveryutil.JoinHostPort(inst.PrivateIPAddress, port)
	m := promutil.NewLabels(24)
	m.Add("__address__", addr)
	m.Add("__meta_ec2_architecture", inst.Architecture)
	m.Add("__meta_ec2_ami", inst.ImageID)
	m.Add("__meta_ec2_availability_zone", inst.Placement.AvailabilityZone)
	m.Add("__meta_ec2_availability_zone_id", azMap[inst.Placement.AvailabilityZone])
	m.Add("__meta_ec2_instance_id", inst.ID)
	m.Add("__meta_ec2_instance_lifecycle", inst.Lifecycle)
	m.Add("__meta_ec2_instance_state", inst.State.Name)
	m.Add("__meta_ec2_instance_type", inst.Type)
	m.Add("__meta_ec2_owner_id", ownerID)
	m.Add("__meta_ec2_platform", inst.Platform)
	m.Add("__meta_ec2_primary_subnet_id", inst.SubnetID)
	m.Add("__meta_ec2_private_dns_name", inst.PrivateDNSName)
	m.Add("__meta_ec2_private_ip", inst.PrivateIPAddress)
	m.Add("__meta_ec2_public_dns_name", inst.PublicDNSName)
	m.Add("__meta_ec2_public_ip", inst.PublicIPAddress)
	m.Add("__meta_ec2_region", region)
	m.Add("__meta_ec2_vpc_id", inst.VPCID)
	if len(inst.VPCID) > 0 {
		subnets := make([]string, 0, len(inst.NetworkInterfaceSet.Items))
		seenSubnets := make(map[string]bool, len(inst.NetworkInterfaceSet.Items))
		var ipv6Addrs []string
		for _, ni := range inst.NetworkInterfaceSet.Items {
			if len(ni.SubnetID) == 0 {
				continue
			}
			// Deduplicate VPC Subnet IDs maintaining the order of the network interfaces returned by EC2.
			if !seenSubnets[ni.SubnetID] {
				seenSubnets[ni.SubnetID] = true
				subnets = append(subnets, ni.SubnetID)
			}
			// Collect ipv6 addresses
			ipv6Addrs = append(ipv6Addrs, ni.IPv6AddressesSet.Items...)
		}
		// We surround the separated list with the separator as well. This way regular expressions
		// in relabeling rules don't have to consider tag positions.
		m.Add("__meta_ec2_subnet_id", ","+strings.Join(subnets, ",")+",")
		if len(ipv6Addrs) > 0 {
			m.Add("__meta_ec2_ipv6_addresses", ","+strings.Join(ipv6Addrs, ",")+",")
		}
	}
	for _, t := range inst.TagSet.Items {
		if len(t.Key) == 0 || len(t.Value) == 0 {
			continue
		}
		m.Add(discoveryutil.SanitizeLabelName("__meta_ec2_tag_"+t.Key), t.Value)
	}
	ms = append(ms, m)
	return ms
}

// Discovery discovers EC2 instances.
type Discovery struct {
	// ...existing code...
	profile string // New field for profile
	// ...existing code...
}

// NewDiscovery creates a new EC2 discovery.
func NewDiscovery(cfg *EC2SDConfig) (*Discovery, error) {
	var awsCfg aws.Config
	var err error

	if cfg.Profile != "" {
		awsCfg, err = config.LoadDefaultConfig(context.Background(),
			config.WithSharedConfigProfile(cfg.Profile),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to load AWS profile %q: %w", cfg.Profile, err)
		}
	} else {
		// ...existing code...
	}

	if cfg.RoleARN != "" {
		stsClient := sts.NewFromConfig(awsCfg)
		creds := stscreds.NewAssumeRoleProvider(stsClient, cfg.RoleARN)
		awsCfg.Credentials = aws.NewCredentialsCache(creds)
	}

	client := ec2.NewFromConfig(awsCfg, func(o *ec2.Options) {
		// ...existing code...
	})

	// ...existing code...

	return &Discovery{
		client: client,
		// ...existing code...
		profile: cfg.Profile,
	}, nil
}
