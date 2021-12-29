package huawei

import (
	"strings"

	"github.com/galaxy-future/BridgX/pkg/cloud"
	"github.com/spf13/cast"
)

// CreateVPC 返回缺少RequestId
func (p *AwsCloud) CreateVPC(req cloud.CreateVpcRequest) (cloud.CreateVpcResponse, error) {
	return cloud.CreateVpcResponse{}, nil
}

func (p *AwsCloud) GetVPC(req cloud.GetVpcRequest) (cloud.GetVpcResponse, error) {
	return cloud.GetVpcResponse{}, nil
}

func (p *AwsCloud) DescribeVpcs(req cloud.DescribeVpcsRequest) (cloud.DescribeVpcsResponse, error) {

	return cloud.DescribeVpcsResponse{}, nil
}

// CreateSwitch add GatewayIp,miss RequestId
func (p *AwsCloud) CreateSwitch(req cloud.CreateSwitchRequest) (cloud.CreateSwitchResponse, error) {

	return cloud.CreateSwitchResponse{}, nil
}

func (p *AwsCloud) GetSwitch(req cloud.GetSwitchRequest) (cloud.GetSwitchResponse, error) {

	return cloud.GetSwitchResponse{}, nil
}

func (p *AwsCloud) DescribeSwitches(req cloud.DescribeSwitchesRequest) (cloud.DescribeSwitchesResponse, error) {

	return cloud.DescribeSwitchesResponse{}, nil
}

//func (p *AwsCloud) getUsedIpNum(switchIds []string) (map[string]int, error) {
//	return nil, nil
//}

func getSubnetTotalIpNum(cidr string) int {
	index := strings.Index(cidr, "/")
	if index < 0 {
		return 0
	}
	num := cast.ToInt(cidr[index+1:])
	if num < 1 || num > 31 {
		return 0
	}

	return 1 << (32 - num)
}
