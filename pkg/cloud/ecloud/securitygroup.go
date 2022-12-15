package ecloud

import (
	"errors"
	"github.com/galaxy-future/BridgX/pkg/cloud"
)

func (p *ECloud) CreateSecurityGroup(req cloud.CreateSecurityGroupRequest) (cloud.CreateSecurityGroupResponse, error) {
	// TODO implement me
	return cloud.CreateSecurityGroupResponse{}, errors.New("implement me")
}

func (p *ECloud) AddIngressSecurityGroupRule(req cloud.AddSecurityGroupRuleRequest) error {
	// TODO implement me
	return errors.New("implement me")
}

func (p *ECloud) AddEgressSecurityGroupRule(req cloud.AddSecurityGroupRuleRequest) error {
	// TODO implement me
	return errors.New("implement me")
}

func (p *ECloud) DescribeSecurityGroups(req cloud.DescribeSecurityGroupsRequest) (cloud.DescribeSecurityGroupsResponse, error) {
	// TODO implement me
	return cloud.DescribeSecurityGroupsResponse{}, errors.New("implement me")
}
