package ecloud

import "github.com/galaxy-future/BridgX/pkg/cloud"

func (p *ECloud) BatchCreate(m cloud.Params, num int) (instanceIds []string, err error) {
	// TODO implement me
	panic("implement me")
}

func (p *ECloud) GetInstances(ids []string) (instances []cloud.Instance, err error) {
	// TODO implement me
	panic("implement me")
}

func (p *ECloud) GetInstancesByTags(region string, tags []cloud.Tag) (instances []cloud.Instance, err error) {
	// TODO implement me
	panic("implement me")
}

func (p *ECloud) GetInstancesByCluster(regionId, clusterName string) (instances []cloud.Instance, err error) {
	// TODO implement me
	panic("implement me")
}

func (p *ECloud) BatchDelete(ids []string, regionId string) error {
	// TODO implement me
	panic("implement me")
}

func (p *ECloud) StartInstances(ids []string) error {
	// TODO implement me
	panic("implement me")
}

func (p *ECloud) StopInstances(ids []string) error {
	// TODO implement me
	panic("implement me")
}

func (p *ECloud) CreateVPC(req cloud.CreateVpcRequest) (cloud.CreateVpcResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (p *ECloud) GetVPC(req cloud.GetVpcRequest) (cloud.GetVpcResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (p *ECloud) CreateSwitch(req cloud.CreateSwitchRequest) (cloud.CreateSwitchResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (p *ECloud) GetSwitch(req cloud.GetSwitchRequest) (cloud.GetSwitchResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (p *ECloud) CreateSecurityGroup(req cloud.CreateSecurityGroupRequest) (cloud.CreateSecurityGroupResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (p *ECloud) AddIngressSecurityGroupRule(req cloud.AddSecurityGroupRuleRequest) error {
	// TODO implement me
	panic("implement me")
}

func (p *ECloud) AddEgressSecurityGroupRule(req cloud.AddSecurityGroupRuleRequest) error {
	// TODO implement me
	panic("implement me")
}

func (p *ECloud) DescribeSecurityGroups(req cloud.DescribeSecurityGroupsRequest) (cloud.DescribeSecurityGroupsResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (p *ECloud) GetRegions() (cloud.GetRegionsResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (p *ECloud) GetZones(req cloud.GetZonesRequest) (cloud.GetZonesResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (p *ECloud) DescribeAvailableResource(req cloud.DescribeAvailableResourceRequest) (cloud.DescribeAvailableResourceResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (p *ECloud) DescribeInstanceTypes(req cloud.DescribeInstanceTypesRequest) (cloud.DescribeInstanceTypesResponse, error) {
	// TODO implement me
	panic("implement me")
}
