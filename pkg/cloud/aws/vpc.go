package aws

import (
	"github.com/galaxy-future/BridgX/internal/logs"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"

	"github.com/galaxy-future/BridgX/pkg/cloud"
)

//CreateVPC miss requestId
func (p *AwsCloud) CreateVPC(req cloud.CreateVpcRequest) (cloud.CreateVpcResponse, error) {
	input := &ec2.CreateVpcInput{
		CidrBlock: aws.String(req.CidrBlock),
	}
	result, err := p.ec2Client.CreateVpc(input)
	if err != nil {
		logs.Logger.Errorf("CreateVPC AwsCloud failed.err: [%v] req[%v]", err, req)
		return cloud.CreateVpcResponse{}, err
	}
	return cloud.CreateVpcResponse{VpcId: aws.StringValue(result.Vpc.VpcId)}, nil
}

//GetVPC aws request miss VpcName
func (p *AwsCloud) GetVPC(req cloud.GetVpcRequest) (cloud.GetVpcResponse, error) {
	input := &ec2.DescribeVpcsInput{
		VpcIds: []*string{aws.String(req.VpcId)},
	}
	result, err := p.ec2Client.DescribeVpcs(input)
	if err != nil {
		logs.Logger.Errorf("GetVPC AwsCloud failed.err: [%v] req[%v]", err, req)
		return cloud.GetVpcResponse{}, err
	}
	if result == nil || len(result.Vpcs) == 0 {
		logs.Logger.Errorf("GetVPC AwsCloud failed. req[%v] result[%v]", req, result)
		return cloud.GetVpcResponse{}, nil
	}
	awsVpc := result.Vpcs[0]
	vpc := cloud.VPC{
		VpcId: aws.StringValue(awsVpc.VpcId),
		//VpcName: ,
		CidrBlock: aws.StringValue(awsVpc.CidrBlock),
		//SwitchIds: ,
		RegionId: req.RegionId,
		Status:   _vpcStatus[aws.StringValue(awsVpc.State)],
		//CreateAt:
	}
	return cloud.GetVpcResponse{Vpc: vpc}, nil
}

func (p *AwsCloud) DescribeVpcs(req cloud.DescribeVpcsRequest) (cloud.DescribeVpcsResponse, error) {
	pageSize := _pageSize * 10
	var awsVpcs = make([]*ec2.Vpc, 0, pageSize)
	input := &ec2.DescribeVpcsInput{
		MaxResults: aws.Int64(int64(pageSize)),
	}
	err := p.ec2Client.DescribeVpcsPages(input, func(output *ec2.DescribeVpcsOutput, b bool) bool {
		awsVpcs = append(awsVpcs, output.Vpcs...)
		return output.NextToken != nil
	})
	if err != nil {
		logs.Logger.Errorf("DescribeVpcs AwsCloud failed.err: [%v] req[%v]", err, req)
		return cloud.DescribeVpcsResponse{}, err
	}
	var vpcs = make([]cloud.VPC, 0, len(awsVpcs))
	for _, vpc := range awsVpcs {
		vpcs = append(vpcs, cloud.VPC{
			VpcId: aws.StringValue(vpc.VpcId),
			//VpcName:
			CidrBlock: aws.StringValue(vpc.CidrBlock),
			//SwitchIds:
			RegionId: req.RegionId,
			Status:   _vpcStatus[aws.StringValue(vpc.State)],
			//CreateAt:
		})
	}
	return cloud.DescribeVpcsResponse{Vpcs: vpcs}, nil
}

// CreateSwitch
func (p *AwsCloud) CreateSwitch(req cloud.CreateSwitchRequest) (cloud.CreateSwitchResponse, error) {
	input := &ec2.CreateSubnetInput{
		AvailabilityZoneId: aws.String(req.ZoneId),
		CidrBlock:          aws.String(req.CidrBlock),
		VpcId:              aws.String(req.VpcId),
	}
	output, err := p.ec2Client.CreateSubnet(input)
	if err != nil {
		logs.Logger.Errorf("CreateSwitch AwsCloud failed.err: [%v] req[%v]", err, req)
		return cloud.CreateSwitchResponse{}, err
	}
	if output == nil || output.Subnet == nil {
		logs.Logger.Warnf("CreateSwitch AwsCloud failed. req[%v] output[%v]", err, req)
		return cloud.CreateSwitchResponse{}, err
	}
	return cloud.CreateSwitchResponse{SwitchId: aws.StringValue(output.Subnet.SubnetId)}, nil
}

func (p *AwsCloud) GetSwitch(req cloud.GetSwitchRequest) (cloud.GetSwitchResponse, error) {
	input := &ec2.DescribeSubnetsInput{
		SubnetIds: append([]*string{}, &req.SwitchId),
	}
	output, err := p.ec2Client.DescribeSubnets(input)
	if err != nil {
		logs.Logger.Errorf("GetSwitch AwsCloud failed.err: [%v] req[%v]", err, req)
		return cloud.GetSwitchResponse{}, err
	}
	if output == nil || len(output.Subnets) == 0 {
		logs.Logger.Errorf("GetSwitch AwsCloud failed. req[%v] output[%v]", req, output)
		return cloud.GetSwitchResponse{}, nil
	}
	awsSubnet := output.Subnets[0]
	subnet := cloud.Switch{
		VpcId:                   aws.StringValue(awsSubnet.VpcId),
		SwitchId:                aws.StringValue(awsSubnet.SubnetId),
		Name:                    aws.StringValue(awsSubnet.SubnetArn),
		IsDefault:               _subnetIsDefault[aws.BoolValue(awsSubnet.DefaultForAz)],
		AvailableIpAddressCount: int(aws.Int64Value(awsSubnet.AvailableIpAddressCount)),
		VStatus:                 _subnetStatus[aws.StringValue(awsSubnet.State)],
		//CreateAt:
		ZoneId:    aws.StringValue(awsSubnet.AvailabilityZoneId),
		CidrBlock: aws.StringValue(awsSubnet.CidrBlock),
		//GatewayIp:
	}
	return cloud.GetSwitchResponse{Switch: subnet}, nil
}

func (p *AwsCloud) DescribeSwitches(req cloud.DescribeSwitchesRequest) (cloud.DescribeSwitchesResponse, error) {
	pageSize := _pageSize * 10
	var awsSubnets = make([]*ec2.Subnet, 0, pageSize)
	input := &ec2.DescribeSubnetsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []*string{&req.VpcId},
			},
		},
		MaxResults: aws.Int64(int64(pageSize)),
	}
	err := p.ec2Client.DescribeSubnetsPages(input, func(output *ec2.DescribeSubnetsOutput, b bool) bool {
		awsSubnets = append(awsSubnets, output.Subnets...)
		return output.NextToken != nil
	})
	if err != nil {
		logs.Logger.Errorf("DescribeSwitches AwsCloud failed.err: [%v] req[%v]", err, req)
		return cloud.DescribeSwitchesResponse{}, err
	}
	if len(awsSubnets) == 0 {
		logs.Logger.Errorf("DescribeSwitches AwsCloud failed. req[%v] len(awsSubnets) is zero", req)
		return cloud.DescribeSwitchesResponse{}, nil
	}
	var subnets = make([]cloud.Switch, 0, len(awsSubnets))
	for _, awsSubnet := range awsSubnets {
		subnets = append(subnets, cloud.Switch{
			VpcId:                   aws.StringValue(awsSubnet.VpcId),
			SwitchId:                aws.StringValue(awsSubnet.SubnetId),
			Name:                    aws.StringValue(awsSubnet.SubnetArn),
			IsDefault:               _subnetIsDefault[aws.BoolValue(awsSubnet.DefaultForAz)],
			AvailableIpAddressCount: int(aws.Int64Value(awsSubnet.AvailableIpAddressCount)),
			VStatus:                 _subnetStatus[aws.StringValue(awsSubnet.State)],
			//CreateAt:
			ZoneId:    aws.StringValue(awsSubnet.AvailabilityZoneId),
			CidrBlock: aws.StringValue(awsSubnet.CidrBlock),
			//GatewayIp:
		})
	}
	return cloud.DescribeSwitchesResponse{Switches: subnets}, nil
}
