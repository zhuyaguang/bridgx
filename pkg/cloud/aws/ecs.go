package aws

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/galaxy-future/BridgX/internal/logs"
	"github.com/galaxy-future/BridgX/pkg/cloud"
	"github.com/galaxy-future/BridgX/pkg/utils"
)

func (p *AwsCloud) BatchCreate(m cloud.Params, num int) ([]string, error) {
	var tags = make([]*ec2.Tag, 0, len(m.Tags))
	for _, tag := range m.Tags {
		tags = append(tags, &ec2.Tag{
			Key:   aws.String(tag.Key),
			Value: aws.String(tag.Value),
		})
	}
	var blockDeviceMappings = make([]*ec2.BlockDeviceMapping, 0, len(m.Disks.DataDisk))
	blockDeviceMappings = append(blockDeviceMappings, &ec2.BlockDeviceMapping{
		//DeviceName: "",
		Ebs: &ec2.EbsBlockDevice{
			DeleteOnTermination: aws.Bool(true),
			VolumeType:          aws.String(m.Disks.SystemDisk.Category),
			VolumeSize:          aws.Int64(int64(m.Disks.SystemDisk.Size)),
		},
	})
	for _, disk := range m.Disks.DataDisk {
		blockDeviceMappings = append(blockDeviceMappings, &ec2.BlockDeviceMapping{
			Ebs: &ec2.EbsBlockDevice{
				DeleteOnTermination: aws.Bool(true),
				VolumeType:          aws.String(disk.Category),
				VolumeSize:          aws.Int64(int64(disk.Size)),
			},
		})
	}
	input := &ec2.RunInstancesInput{
		BlockDeviceMappings: blockDeviceMappings,
		ImageId:             aws.String(m.ImageId),
		InstanceType:        aws.String(m.InstanceType),
		MaxCount:            aws.Int64(int64(num)),
		MinCount:            aws.Int64(int64(num)),
		SecurityGroupIds: []*string{
			aws.String(m.Network.SecurityGroup),
		},
		SubnetId: aws.String(m.Network.SubnetId),
		TagSpecifications: []*ec2.TagSpecification{
			{
				ResourceType: aws.String("instance"),
				Tags:         tags,
			},
		},
	}
	if m.DryRun {
		input.DryRun = aws.Bool(m.DryRun)
	}

	result, err := p.ec2Client.RunInstances(input)
	if err != nil {
		logs.Logger.Errorf("BatchCreate AwsCloud failed.err: [%v] req[%v]", err, m)
		return []string{}, err
	}
	var instanceIds = make([]string, 0, len(result.Instances))
	for _, instance := range result.Instances {
		instanceIds = append(instanceIds, *instance.InstanceId)
	}
	return instanceIds, nil
}

func (p *AwsCloud) GetInstances(ids []string) (instances []cloud.Instance, err error) {
	idNum := len(ids)
	if idNum < 1 {
		return []cloud.Instance{}, nil
	}
	var awsInstances = make([]*ec2.Instance, 0, len(ids))
	batchIds := utils.StringSliceSplit(ids, _pageSize)
	for _, onceIds := range batchIds {
		input := &ec2.DescribeInstancesInput{
			InstanceIds: aws.StringSlice(onceIds),
		}
		result, err := p.ec2Client.DescribeInstances(input)
		if err != nil {
			logs.Logger.Errorf("GetInstances AwsCloud failed.err: [%v] req[%v]", err, ids)
			return []cloud.Instance{}, nil
		}
		for _, reservation := range result.Reservations {
			awsInstances = append(awsInstances, reservation.Instances...)
		}
	}
	for _, instance := range awsInstances {
		var securityGroupIds = make([]string, 0, len(instance.SecurityGroups))
		for _, securityGroup := range instance.SecurityGroups {
			securityGroupIds = append(securityGroupIds, *securityGroup.GroupId)
		}
		instances = append(instances, cloud.Instance{
			Id: *instance.InstanceId,
			//CostWay: ,
			Provider: cloud.AwsCloud,
			IpInner:  *instance.PrivateIpAddress,
			IpOuter:  *instance.PublicIpAddress,
			Network: &cloud.Network{
				VpcId:         *instance.VpcId,
				SubnetId:      *instance.SubnetId,
				SecurityGroup: strings.Join(securityGroupIds, ","),
				//InternetChargeType: ,
				//InternetMaxBandwidthOut: ,
				//InternetIpType:
			},
			ImageId: *instance.ImageId,
			Status:  _ecsStatus[*instance.State.Name],
			//ExpireAt: in,
		})
	}
	return instances, nil
}

func (p *AwsCloud) GetInstancesByTags(regionId string, tags []cloud.Tag) (instances []cloud.Instance, err error) {
	pageSize := _pageSize * 10
	var awsInstances = make([]*ec2.Instance, 0, pageSize)
	var filters = make([]*ec2.Filter, 0, len(tags))
	for _, tag := range tags {
		filters = append(filters, &ec2.Filter{
			Name:   aws.String("tag:" + tag.Key),
			Values: []*string{aws.String(tag.Value)},
		})
	}
	input := &ec2.DescribeInstancesInput{
		Filters:    filters,
		MaxResults: aws.Int64(int64(pageSize)),
	}
	//auto page
	p.ec2Client.DescribeInstancesPages(input, func(output *ec2.DescribeInstancesOutput, b bool) bool {
		for _, reservation := range output.Reservations {
			awsInstances = append(awsInstances, reservation.Instances...)
		}
		return output.NextToken != nil
	})
	for _, instance := range awsInstances {
		var securityGroupIds = make([]string, 0, len(instance.SecurityGroups))
		for _, securityGroup := range instance.SecurityGroups {
			securityGroupIds = append(securityGroupIds, *securityGroup.GroupId)
		}
		instances = append(instances, cloud.Instance{
			Id: *instance.InstanceId,
			//CostWay: ,
			Provider: cloud.AwsCloud,
			IpInner:  *instance.PrivateIpAddress,
			IpOuter:  *instance.PublicIpAddress,
			Network: &cloud.Network{
				VpcId:         *instance.VpcId,
				SubnetId:      *instance.SubnetId,
				SecurityGroup: strings.Join(securityGroupIds, ","),
				//InternetChargeType: ,
				//InternetMaxBandwidthOut: ,
				//InternetIpType:
			},
			ImageId: *instance.ImageId,
			Status:  _ecsStatus[*instance.State.Name],
			//ExpireAt: in,
		})
	}
	return instances, nil
}

func (p *AwsCloud) GetInstancesByCluster(regionId, clusterName string) (instances []cloud.Instance, err error) {
	return p.GetInstancesByTags(regionId, []cloud.Tag{{
		Key:   cloud.ClusterName,
		Value: clusterName,
	}})
}

func (p *AwsCloud) BatchDelete(ids []string, regionId string) error {
	idNum := len(ids)
	if idNum < 1 {
		return _errInstanceIdsEmpty
	}
	pageSize := _pageSize * 10
	batchIds := utils.StringSliceSplit(ids, int64(pageSize))
	for _, onceIds := range batchIds {
		input := &ec2.TerminateInstancesInput{
			InstanceIds: aws.StringSlice(onceIds),
		}
		//TODO
		_, err := p.ec2Client.TerminateInstances(input)
		if err != nil {
			logs.Logger.Errorf("BatchDelete AwsCloud failed.err: [%v] req[%v]", err, ids)
			return err
		}
	}
	return nil
}

func (p *AwsCloud) StartInstances(ids []string) error {
	idNum := len(ids)
	if idNum < 1 {
		return _errInstanceIdsEmpty
	}
	pageSize := _pageSize * 10
	batchIds := utils.StringSliceSplit(ids, int64(pageSize))
	for _, onceIds := range batchIds {
		input := &ec2.StartInstancesInput{
			InstanceIds: aws.StringSlice(onceIds),
		}
		//TODO
		_, err := p.ec2Client.StartInstances(input)
		if err != nil {
			logs.Logger.Errorf("StartInstances AwsCloud failed.err: [%v] req[%v]", err, ids)
			return err
		}
	}
	return nil
}

func (p *AwsCloud) StopInstances(ids []string) error {
	idNum := len(ids)
	if idNum < 1 {
		return _errInstanceIdsEmpty
	}
	pageSize := _pageSize * 10
	batchIds := utils.StringSliceSplit(ids, int64(pageSize))
	for _, onceIds := range batchIds {
		input := &ec2.StopInstancesInput{
			InstanceIds: aws.StringSlice(onceIds),
		}
		//TODO
		_, err := p.ec2Client.StopInstances(input)
		if err != nil {
			logs.Logger.Errorf("StopInstances AwsCloud failed.err: [%v] req[%v]", err, ids)
			return err
		}
	}
	return nil
}

func (p *AwsCloud) GetZones(req cloud.GetZonesRequest) (cloud.GetZonesResponse, error) {
	var filters = make([]*ec2.Filter, 0, 1)
	input := &ec2.DescribeAvailabilityZonesInput{
		Filters: append(filters, &ec2.Filter{
			Name:   aws.String("region-name"),
			Values: []*string{&req.RegionId},
		}),
	}
	result, err := p.ec2Client.DescribeAvailabilityZones(input)
	if err != nil {
		logs.Logger.Errorf("GetZones AwsCloud failed.err: [%v] req[%v]", err, req)
		return cloud.GetZonesResponse{}, nil
	}
	var zones = make([]cloud.Zone, 0, len(result.AvailabilityZones))
	for _, zone := range result.AvailabilityZones {
		zones = append(zones, cloud.Zone{
			ZoneId:    *zone.ZoneId,
			LocalName: *zone.ZoneName,
		})
	}
	return cloud.GetZonesResponse{Zones: zones}, nil
}

//DescribeAvailableResource miss status
func (p *AwsCloud) DescribeAvailableResource(req cloud.DescribeAvailableResourceRequest) (cloud.DescribeAvailableResourceResponse, error) {
	pageSize := _pageSize * 10
	var awsInstanceTypes = make([]*ec2.InstanceTypeOffering, 0, pageSize)
	var filters = make([]*ec2.Filter, 0, 1)
	var instanceTypes = make([]string, 0, pageSize)
	zoneInsTypeMap := make(map[string][]cloud.InstanceType, 64)
	input := &ec2.DescribeInstanceTypeOfferingsInput{
		LocationType: aws.String(_locationTypeNameZoneId),
		MaxResults:   aws.Int64(int64(pageSize)),
	}
	if req.ZoneId != "" {
		input.Filters = append(filters, &ec2.Filter{
			Name:   aws.String(_filterNameLocation),
			Values: []*string{&req.ZoneId},
		})
		zoneInsTypeMap[req.ZoneId] = []cloud.InstanceType{}
	} else {
		zones, err := p.GetZones(cloud.GetZonesRequest{
			RegionId: req.RegionId,
		})
		if err != nil {
			return cloud.DescribeAvailableResourceResponse{}, nil
		}
		var zoneIds = make([]*string, 0, len(zones.Zones))
		for _, zone := range zones.Zones {
			zoneIds = append(zoneIds, aws.String(zone.ZoneId))
			zoneInsTypeMap[zone.ZoneId] = []cloud.InstanceType{}
		}
		input.Filters = append(filters, &ec2.Filter{
			Name:   aws.String(_filterNameLocation),
			Values: zoneIds,
		})
	}

	//auto page
	err := p.ec2Client.DescribeInstanceTypeOfferingsPages(input, func(output *ec2.DescribeInstanceTypeOfferingsOutput, b bool) bool {
		awsInstanceTypes = append(awsInstanceTypes, output.InstanceTypeOfferings...)
		return output.NextToken != nil
	})
	if err != nil {
		logs.Logger.Errorf("DescribeAvailableResource AwsCloud failed.err: [%v] req[%v]", err, req)
		return cloud.DescribeAvailableResourceResponse{}, nil
	}
	for _, instanceType := range awsInstanceTypes {
		instanceTypes = append(instanceTypes, aws.StringValue(instanceType.InstanceType))
		if types, ok := zoneInsTypeMap[*instanceType.Location]; ok {
			types = append(types, cloud.InstanceType{
				InstanceInfo: cloud.InstanceInfo{
					InsTypeName: aws.StringValue(instanceType.InstanceType),
				},
				Status: cloud.InsTypeAvailable,
			})
			zoneInsTypeMap[aws.StringValue(instanceType.Location)] = types
		}
	}
	instanceTypeInfos, err := p.DescribeInstanceTypes(cloud.DescribeInstanceTypesRequest{TypeName: instanceTypes})
	if err != nil {
		return cloud.DescribeAvailableResourceResponse{}, nil
	}
	var infoMap = make(map[string]cloud.InstanceInfo, len(instanceTypeInfos.Infos))
	for _, info := range instanceTypeInfos.Infos {
		infoMap[info.InsTypeName] = info
	}
	for zoneId, infos := range zoneInsTypeMap {
		for i, info := range infos {
			instanceInfo := infoMap[info.InsTypeName]
			infos[i].Core = instanceInfo.Core
			infos[i].Memory = instanceInfo.Memory
			infos[i].Family = instanceInfo.Family
		}
		zoneInsTypeMap[zoneId] = infos
	}
	return cloud.DescribeAvailableResourceResponse{InstanceTypes: zoneInsTypeMap}, nil
}

//DescribeInstanceTypes miss Family
func (p *AwsCloud) DescribeInstanceTypes(req cloud.DescribeInstanceTypesRequest) (cloud.DescribeInstanceTypesResponse, error) {
	var instanceTypeInfos = make([]cloud.InstanceInfo, 0, _pageSize)
	var awsInstanceTypeInfos = make([]*ec2.InstanceTypeInfo, 0, _pageSize)
	batchIds := utils.StringSliceSplit(req.TypeName, _pageSize)
	for _, onceIds := range batchIds {
		input := &ec2.DescribeInstanceTypesInput{
			InstanceTypes: aws.StringSlice(onceIds),
		}
		result, err := p.ec2Client.DescribeInstanceTypes(input)
		if err != nil {
			logs.Logger.Errorf("DescribeInstanceTypes AwsCloud failed.err: [%v] req[%v]", err, req)
			return cloud.DescribeInstanceTypesResponse{}, nil
		}
		awsInstanceTypeInfos = append(awsInstanceTypeInfos, result.InstanceTypes...)
	}
	for _, instanceTypeInfo := range awsInstanceTypeInfos {
		instanceTypeInfos = append(instanceTypeInfos, cloud.InstanceInfo{
			Core:   int(aws.Int64Value(instanceTypeInfo.VCpuInfo.DefaultCores)),
			Memory: int(aws.Int64Value(instanceTypeInfo.MemoryInfo.SizeInMiB) / 1024),
			//Family:      "",
			InsTypeName: *instanceTypeInfo.InstanceType,
		})
	}
	return cloud.DescribeInstanceTypesResponse{Infos: instanceTypeInfos}, nil
}
