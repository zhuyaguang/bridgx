package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/galaxy-future/BridgX/internal/logs"
	"github.com/galaxy-future/BridgX/pkg/cloud"
	"github.com/pkg/errors"
)

func (p *AWSCloud) AllocateEip(req cloud.AllocateEipRequest) (ids []string, err error) {
	input := &ec2.AllocateAddressInput{}
	if req.Name != "" {
		input.TagSpecifications = []*ec2.TagSpecification{
			{
				ResourceType: aws.String(_resourceTypeEip),
				Tags: append([]*ec2.Tag{}, &ec2.Tag{
					Key:   aws.String(_tagKeyEipName),
					Value: aws.String(req.Name),
				}),
			},
		}
	}
	idChan := make(chan string, req.Num)
	errChan := make(chan error, req.Num)
	for i := 0; i < req.Num; i++ {
		go func() {
			result, err := p.ec2Client.AllocateAddress(input)
			if err != nil {
				logs.Logger.Errorf("AllocateEip AwsCloud failed.err:[%v] req:[%v]", err, req)
				errChan <- err
				return
			}
			idChan <- *result.AllocationId
		}()
	}
	for i := 0; i < req.Num; i++ {
		select {
		case err = <-errChan:
		case id := <-idChan:
			ids = append(ids, id)
		}
	}
	return ids, err
}

func (p *AWSCloud) GetEips(ids []string, regionId string) (map[string]cloud.Eip, error) {
	input := &ec2.DescribeAddressesInput{
		AllocationIds: aws.StringSlice(ids),
	}
	result, err := p.ec2Client.DescribeAddresses(input)
	if err != nil {
		logs.Logger.Errorf("GetEips AwsCloud failed.err:[%v] req:[%v]", err, ids)
		return nil, err
	}
	eipMap := make(map[string]cloud.Eip, len(result.Addresses))
	for _, v := range result.Addresses {
		eipMap[aws.StringValue(v.AllocationId)] = eip2Cloud(v)
	}
	return eipMap, nil
}

func eip2Cloud(eip *ec2.Address) cloud.Eip {
	var name string
	if len(eip.Tags) > 0 {
		name = aws.StringValue(eip.Tags[0].Value)
	}
	return cloud.Eip{
		Id:         aws.StringValue(eip.AllocationId),
		Name:       name,
		Ip:         aws.StringValue(eip.PublicIp),
		InstanceId: aws.StringValue(eip.InstanceId),
	}
}

func (p *AWSCloud) ReleaseEip(ids []string) (err error) {
	num := len(ids)
	if num < 1 {
		return nil
	}

	idChan := make(chan string, num)
	errChan := make(chan error, num)
	for _, id := range ids {
		go func(id string) {
			input := &ec2.ReleaseAddressInput{
				AllocationId: aws.String(id),
			}
			_, err := p.ec2Client.ReleaseAddress(input)
			if err != nil {
				logs.Logger.Errorf("ReleaseEip AwsCloud failed.err:[%v] req:[%v]", err, id)
				errChan <- err
				return
			}
			idChan <- id
		}(id)
	}
	for i := 0; i < num; i++ {
		select {
		case err = <-errChan:
		case <-idChan:
		}
	}
	return err
}

func (p *AWSCloud) AssociateEip(id, instanceId, vpcId string) error {
	gateway, err := p.describeInternetGatewayByVpcId(vpcId)
	if err != nil {
		return err
	}
	if gateway == nil {
		gateway, err = p.createInternetGateway()
		if err != nil {
			return err
		}
		err = p.attachInternetGateway(aws.StringValue(gateway.InternetGatewayId), vpcId)
		if err != nil {
			return err
		}
	}
	input := &ec2.AssociateAddressInput{
		AllocationId: aws.String(id),
		InstanceId:   aws.String(instanceId),
	}
	_, err = p.ec2Client.AssociateAddress(input)
	if err != nil {
		logs.Logger.Errorf("AssociateEip AwsCloud failed.err:[%v] id:[%v] instanceId:[%v]", err, id, instanceId)
		return err
	}
	return nil
}

func (p *AWSCloud) DisassociateEip(ip string) error {
	input := &ec2.DisassociateAddressInput{
		PublicIp: aws.String(ip),
	}
	_, err := p.ec2Client.DisassociateAddress(input)
	if err != nil {
		logs.Logger.Errorf("AssociateEip AwsCloud failed.err:[%v] ip:[%v]", err, ip)
		return err
	}
	return nil
}

func (p *AWSCloud) createInternetGateway() (*ec2.InternetGateway, error) {
	input := &ec2.CreateInternetGatewayInput{}
	gateway, err := p.ec2Client.CreateInternetGateway(input)
	if err != nil {
		logs.Logger.Errorf("createInternetGateway AwsCloud failed.err:[%v] req:[%v]", err, input)
		return nil, err
	}
	return gateway.InternetGateway, nil
}

func (p *AWSCloud) describeInternetGatewayByVpcId(vpcId string) (*ec2.InternetGateway, error) {
	pageSize := _pageSize * 10
	input := &ec2.DescribeInternetGatewaysInput{
		Filters: append([]*ec2.Filter{}, &ec2.Filter{
			Name:   aws.String(_filterNameAttachmentVpcId),
			Values: []*string{aws.String(vpcId)},
		}),
		MaxResults: aws.Int64(int64(pageSize)),
	}
	var internetGateways = make([]*ec2.InternetGateway, 0)
	err := p.ec2Client.DescribeInternetGatewaysPages(input, func(output *ec2.DescribeInternetGatewaysOutput, b bool) bool {
		internetGateways = append(internetGateways, output.InternetGateways...)
		return output.NextToken != nil
	})
	if err != nil {
		logs.Logger.Errorf("describeInternetGateways AwsCloud failed.err:[%v] req:[%v]", err, input)
		return nil, err
	}
	if len(internetGateways) == 0 {
		return nil, nil
	}
	return internetGateways[0], nil
}

func (p *AWSCloud) describeRouteTaleByVpcId(vpcId string) (*ec2.RouteTable, error) {
	input := &ec2.DescribeRouteTablesInput{
		Filters: append([]*ec2.Filter{}, &ec2.Filter{
			Name:   aws.String(_filterNameVpcId),
			Values: []*string{aws.String(vpcId)},
		}),
		MaxResults: aws.Int64(int64(_pageSize)),
	}
	var routeTables = make([]*ec2.RouteTable, 0)
	err := p.ec2Client.DescribeRouteTablesPages(input, func(output *ec2.DescribeRouteTablesOutput, b bool) bool {
		routeTables = append(routeTables, output.RouteTables...)
		return output.NextToken != nil
	})
	if err != nil {
		logs.Logger.Errorf("describeRouteTaleByVpcId AwsCloud failed.err:[%v] req:[%v]", err, input)
		return nil, err
	}
	if len(routeTables) == 0 {
		return nil, nil
	}
	return routeTables[0], nil
}

func (p *AWSCloud) createRouteTable(vpcId string) (*ec2.RouteTable, error) {
	input := &ec2.CreateRouteTableInput{
		VpcId: aws.String(vpcId),
	}
	output, err := p.ec2Client.CreateRouteTable(input)
	if err != nil {
		logs.Logger.Errorf("createRouteTable AwsCloud failed.err:[%v] req:[%v]", err, input)
		return nil, err
	}
	return output.RouteTable, nil
}

func (p *AWSCloud) createRoute(routeTableId, internetGatewayId string) error {
	input := &ec2.CreateRouteInput{
		RouteTableId:         aws.String(routeTableId),
		GatewayId:            aws.String(internetGatewayId),
		DestinationCidrBlock: aws.String("0.0.0.0/0"),
	}
	_, err := p.ec2Client.CreateRoute(input)
	if err != nil {
		logs.Logger.Errorf("createRoute AwsCloud failed.err:[%v] req:[%v]", err, input)
		return err
	}
	return nil
}

func (p *AWSCloud) attachInternetGateway(internetGatewayId, vpcId string) error {
	input := &ec2.AttachInternetGatewayInput{
		InternetGatewayId: aws.String(internetGatewayId),
		VpcId:             aws.String(vpcId),
	}
	_, err := p.ec2Client.AttachInternetGateway(input)
	if err != nil {
		logs.Logger.Errorf("attachInternetGateway AwsCloud failed.err:[%v] req:[%v]", err, input)
		return err
	}
	return nil
}

func (p *AWSCloud) DescribeEip(req cloud.DescribeEipRequest) (cloud.DescribeEipResponse, error) {
	input := &ec2.DescribeAddressesInput{}
	if req.InstanceId != "" {
		input.Filters = append([]*ec2.Filter{}, &ec2.Filter{
			Name:   aws.String(_filterNameInstanceId),
			Values: []*string{aws.String(req.InstanceId)},
		})
	}
	result, err := p.ec2Client.DescribeAddresses(input)
	if err != nil {
		logs.Logger.Errorf("GetEips AwsCloud failed.err:[%v] req:[%v]", err, req)
		return cloud.DescribeEipResponse{}, err
	}
	eipList := make([]cloud.Eip, 0, len(result.Addresses))
	for _, address := range result.Addresses {
		eipList = append(eipList, eip2Cloud(address))
	}
	return cloud.DescribeEipResponse{TotalCount: len(eipList), List: eipList}, nil
}

func (p *AWSCloud) ConvertPublicIpToEip(req cloud.ConvertPublicIpToEipRequest) error {

	return errors.New("not Implemented") // do not support in aws
}
