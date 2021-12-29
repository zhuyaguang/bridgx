package huawei

import (
	"github.com/galaxy-future/BridgX/pkg/cloud"
)

// CreateSecurityGroup 将VpcId写入Description，方便查找
func (p *AwsCloud) CreateSecurityGroup(req cloud.CreateSecurityGroupRequest) (cloud.CreateSecurityGroupResponse, error) {

	return cloud.CreateSecurityGroupResponse{}, nil
}

// AddIngressSecurityGroupRule 入参各云得统一
func (p *AwsCloud) AddIngressSecurityGroupRule(req cloud.AddSecurityGroupRuleRequest) error {
	return nil
}

func (p *AwsCloud) AddEgressSecurityGroupRule(req cloud.AddSecurityGroupRuleRequest) error {
	return nil
}

func (p *AwsCloud) DescribeSecurityGroups(req cloud.DescribeSecurityGroupsRequest) (cloud.DescribeSecurityGroupsResponse, error) {

	return cloud.DescribeSecurityGroupsResponse{}, nil
}

func (p *AwsCloud) DescribeGroupRules(req cloud.DescribeGroupRulesRequest) (cloud.DescribeGroupRulesResponse, error) {
	return cloud.DescribeGroupRulesResponse{}, nil
}

func (p *AwsCloud) addSecGrpRule(req cloud.AddSecurityGroupRuleRequest, direction string) error {
	return nil
}
