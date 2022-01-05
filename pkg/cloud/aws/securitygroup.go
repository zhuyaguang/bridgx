package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/galaxy-future/BridgX/internal/logs"
	"github.com/galaxy-future/BridgX/pkg/cloud"
	"github.com/spf13/cast"
)

// CreateSecurityGroup output missing field: RequestId
func (p *AwsCloud) CreateSecurityGroup(req cloud.CreateSecurityGroupRequest) (cloud.CreateSecurityGroupResponse, error) {
	input := &ec2.CreateSecurityGroupInput{
		Description: aws.String(req.SecurityGroupName),
		GroupName:   aws.String(req.SecurityGroupName),
		VpcId:       aws.String(req.VpcId),
	}

	output, err := p.ec2Client.CreateSecurityGroup(input)
	if err != nil {
		logs.Logger.Errorf("CreateSecurityGroup AwsCloud failed.err: [%v] req[%v]", err, req)
		return cloud.CreateSecurityGroupResponse{}, err
	}
	return cloud.CreateSecurityGroupResponse{SecurityGroupId: *output.GroupId}, nil
}

// AddIngressSecurityGroupRule 入参各云得统一
func (p *AwsCloud) AddIngressSecurityGroupRule(req cloud.AddSecurityGroupRuleRequest) error {
	input := &ec2.AuthorizeSecurityGroupIngressInput{
		GroupId:    aws.String(req.SecurityGroupId),
		IpProtocol: aws.String(req.IpProtocol),
		IpPermissions: []*ec2.IpPermission{
			{
				FromPort: aws.Int64(int64(req.PortFrom)),
				//IpProtocol: aws.String("tcp"),
				IpRanges: []*ec2.IpRange{
					{
						CidrIp: aws.String(req.CidrIp),
						//Description: aws.String("SSH access from the LA office"),
					},
				},
				ToPort: aws.Int64(int64(req.PortTo)),
				UserIdGroupPairs: []*ec2.UserIdGroupPair{
					{
						GroupId: aws.String(req.GroupId),
						VpcId:   aws.String(req.VpcId),
						//Description: aws.String("HTTP access from other instances"),
					},
				},
			},
		},
	}
	_, err := p.ec2Client.AuthorizeSecurityGroupIngress(input)
	if err != nil {
		logs.Logger.Errorf("AddIngressSecurityGroupRule AwsCloud failed.err: [%v] req[%v]", err, req)
		return err
	}
	//if !*result.Return {
	//	logs.Logger.Errorf("AddIngressSecurityGroupRule AwsCloud failed. req[%v]", req)
	//	return err
	//}
	return nil
}

func (p *AwsCloud) AddEgressSecurityGroupRule(req cloud.AddSecurityGroupRuleRequest) error {
	input := &ec2.AuthorizeSecurityGroupEgressInput{
		GroupId:    aws.String(req.SecurityGroupId),
		IpProtocol: aws.String(req.IpProtocol),
		IpPermissions: []*ec2.IpPermission{
			{
				FromPort: aws.Int64(int64(req.PortFrom)),
				//IpProtocol: aws.String("tcp"),
				IpRanges: []*ec2.IpRange{
					{
						CidrIp: aws.String(req.CidrIp),
						//Description: aws.String("SSH access from the LA office"),
					},
				},
				ToPort: aws.Int64(int64(req.PortTo)),
				UserIdGroupPairs: []*ec2.UserIdGroupPair{
					{
						GroupId: aws.String(req.GroupId),
						VpcId:   aws.String(req.VpcId),
						//Description: aws.String("HTTP access from other instances"),
					},
				},
			},
		},
	}
	_, err := p.ec2Client.AuthorizeSecurityGroupEgress(input)
	if err != nil {
		logs.Logger.Errorf("AddEgressSecurityGroupRule AwsCloud failed.err: [%v] req[%v]", err, req)
		return err
	}
	//if !*result.Return {
	//	logs.Logger.Errorf("AddEgressSecurityGroupRule AwsCloud failed. req[%v]", req)
	//	return errors.New("")
	//}
	return nil
}

// DescribeSecurityGroups output missing field: CreateAt
func (p *AwsCloud) DescribeSecurityGroups(req cloud.DescribeSecurityGroupsRequest) (cloud.DescribeSecurityGroupsResponse, error) {
	pageSize := _pageSize * 10
	var awsSecurityGroups = make([]*ec2.SecurityGroup, 0, pageSize)
	input := &ec2.DescribeSecurityGroupsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []*string{&req.VpcId},
			},
		},
		MaxResults: aws.Int64(int64(pageSize)),
	}
	err := p.ec2Client.DescribeSecurityGroupsPages(input, func(output *ec2.DescribeSecurityGroupsOutput, b bool) bool {
		awsSecurityGroups = append(awsSecurityGroups, output.SecurityGroups...)
		return output.NextToken != nil
	})
	if err != nil {
		logs.Logger.Errorf("DescribeSecurityGroups AwsCloud failed.err: [%v] req[%v]", err, req)
		return cloud.DescribeSecurityGroupsResponse{}, err
	}
	if len(awsSecurityGroups) == 0 {
		logs.Logger.Warnf("DescribeSecurityGroups AwsCloud failed. req[%v] len(awsSubnets) is zero", req)
		return cloud.DescribeSecurityGroupsResponse{}, nil
	}
	var securityGroups = make([]cloud.SecurityGroup, 0, len(awsSecurityGroups))
	for _, group := range awsSecurityGroups {
		securityGroups = append(securityGroups, cloud.SecurityGroup{
			SecurityGroupId:   aws.StringValue(group.GroupId),
			SecurityGroupName: aws.StringValue(group.GroupName),
			SecurityGroupType: "normal",
			VpcId:             req.VpcId,
			RegionId:          req.RegionId,
			//CreateAt: "",
		})
	}
	return cloud.DescribeSecurityGroupsResponse{Groups: securityGroups}, nil
}

// DescribeGroupRules output missing field: CreateAt
func (p *AwsCloud) DescribeGroupRules(req cloud.DescribeGroupRulesRequest) (cloud.DescribeGroupRulesResponse, error) {
	pageSize := _pageSize * 10
	var awsSecurityGroupRules = make([]*ec2.SecurityGroupRule, 0, pageSize)
	input := &ec2.DescribeSecurityGroupRulesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("group-id"),
				Values: []*string{&req.SecurityGroupId},
			},
		},
		MaxResults: aws.Int64(int64(pageSize)),
	}
	err := p.ec2Client.DescribeSecurityGroupRulesPages(input, func(output *ec2.DescribeSecurityGroupRulesOutput, b bool) bool {
		awsSecurityGroupRules = append(awsSecurityGroupRules, output.SecurityGroupRules...)
		return output.NextToken != nil
	})
	if err != nil {
		logs.Logger.Errorf("DescribeGroupRules AwsCloud failed.err: [%v] req[%v]", err, req)
		return cloud.DescribeGroupRulesResponse{}, err
	}
	if len(awsSecurityGroupRules) == 0 {
		logs.Logger.Errorf("DescribeGroupRules AwsCloud failed. req[%v] len(awsSecurityGroupRules) is zero", req)
		return cloud.DescribeGroupRulesResponse{}, nil
	}
	var rules = make([]cloud.SecurityGroupRule, 0, len(awsSecurityGroupRules))
	for _, rule := range awsSecurityGroupRules {
		var vpcId string
		if rule.ReferencedGroupInfo != nil {
			vpcId = aws.StringValue(rule.ReferencedGroupInfo.VpcId)
		}
		rules = append(rules, cloud.SecurityGroupRule{
			VpcId:           vpcId,
			SecurityGroupId: aws.StringValue(rule.GroupId),
			PortRange:       formatPortRange(aws.Int64Value(rule.FromPort), aws.Int64Value(rule.ToPort)),
			Protocol:        aws.StringValue(rule.IpProtocol),
			Direction:       _secGrpRuleDirection[aws.BoolValue(rule.IsEgress)],
			GroupId:         aws.StringValue(rule.GroupId),
			CidrIp:          aws.StringValue(rule.CidrIpv4),
			PrefixListId:    aws.StringValue(rule.PrefixListId),
			//CreateAt:
		})
	}
	return cloud.DescribeGroupRulesResponse{Rules: rules}, nil
}

func formatPortRange(fromPort, toPort int64) string {
	if fromPort == toPort {
		return cast.ToString(fromPort)
	}
	return fmt.Sprintf("%d-%d", fromPort, toPort)
}
