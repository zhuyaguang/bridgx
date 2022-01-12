package tests

import (
	"testing"

	"github.com/galaxy-future/BridgX/pkg/cloud/aws"

	"github.com/galaxy-future/BridgX/pkg/cloud"
	jsoniter "github.com/json-iterator/go"
)

func getAwsClient() (*aws.AwsCloud, error) {
	client, err := aws.New("xxx", "xxx", "cn-north-1")
	if err != nil {
		return nil, err
	}
	return client, nil
}

func TestAwsGetRegions(t *testing.T) {
	client, err := getAwsClient()
	if err != nil {
		t.Log(err)
		return
	}
	result, err := client.GetRegions()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	t.Log(result)
}

func TestAwsGetZones(t *testing.T) {
	client, err := getAwsClient()
	if err != nil {
		t.Log(err)
		return
	}
	result, err := client.GetZones(cloud.GetZonesRequest{
		RegionId: "cn-north-1",
	})
	if err != nil {
		t.Log(err.Error())
		return
	}
	t.Log(result)
}

func TestAwsDescribeAvailableResource(t *testing.T) {
	client, err := getAwsClient()
	if err != nil {
		t.Log(err)
		return
	}
	result, err := client.DescribeAvailableResource(cloud.DescribeAvailableResourceRequest{
		RegionId: "cn-north-1",
		//ZoneId:   "cnn1-az4",
	})
	if err != nil {
		t.Log(err.Error())
		return
	}
	resStr, _ := jsoniter.MarshalToString(result)
	t.Log(resStr)
}

func TestAwsDescribeImages(t *testing.T) {
	client, err := getAwsClient()
	if err != nil {
		t.Log(err)
		return
	}
	result, err := client.DescribeImages(cloud.DescribeImagesRequest{
		RegionId: "cn-north-1",
		//InsType: "",
		ImageType: "global",
	})
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	t.Log(result)
}

func TestAwsDescribeInstanceTypesRequest(t *testing.T) {
	client, err := getAwsClient()
	if err != nil {
		t.Log(err)
		return
	}
	result, err := client.DescribeInstanceTypes(cloud.DescribeInstanceTypesRequest{
		TypeName: []string{"g3.4xlarge", "m5.8xlarge", "t3.large", "m5.2xlarge"},
	})
	if err != nil {
		t.Log(err.Error())
		return
	}
	t.Log(result)
}

func TestAwsCreateSecurityGroup(t *testing.T) {
	client, err := getAwsClient()
	if err != nil {
		t.Log(err)
		return
	}

	req := cloud.CreateSecurityGroupRequest{
		RegionId:          "cn-north-1",
		SecurityGroupName: "cxd-security-group",
		VpcId:             "vpc-05ac2dba6c6e3a683",
		SecurityGroupType: "",
	}
	res, err := client.CreateSecurityGroup(req)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	resStr, _ := jsoniter.MarshalToString(res)
	t.Log(resStr)
}

func TestAwsDescribeSecurityGroups(t *testing.T) {
	client, err := getAwsClient()
	if err != nil {
		t.Log(err)
		return
	}
	//sg-08ec8406adec6883a
	req := cloud.DescribeSecurityGroupsRequest{
		RegionId: "cn-north-1",
		//VpcId:    "vpc-05ac2dba6c6e3a683",
	}
	res, err := client.DescribeSecurityGroups(req)
	if err != nil {
		t.Log(err.Error())
		return
	}
	resStr, _ := jsoniter.MarshalToString(res)
	t.Log(resStr)
}

func TestAwsDescribeGroupRules(t *testing.T) {
	client, err := getAwsClient()
	if err != nil {
		t.Log(err)
		return
	}
	req := cloud.DescribeGroupRulesRequest{
		RegionId:        "cn-north-1",
		SecurityGroupId: "sg-01942a6ccaffab8fa",
	}
	res, err := client.DescribeGroupRules(req)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	resStr, _ := jsoniter.MarshalToString(res)
	t.Log(resStr)
}

func TestAwsAddSecGrpRule(t *testing.T) {
	client, err := getAwsClient()
	if err != nil {
		t.Log(err)
		return
	}

	req := cloud.AddSecurityGroupRuleRequest{
		SecurityGroupId: "sg-010c2d1733ab35cb4",
		GroupId:         "",
		VpcId:           "vpc-05137e7d412f0382b",
		IpProtocol:      "udp",
		PortFrom:        121,
		PortTo:          122,
		CidrIp:          "192.168.1.0/24",
	}
	err = client.AddIngressSecurityGroupRule(req)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	req = cloud.AddSecurityGroupRuleRequest{
		SecurityGroupId: "sg-010c2d1733ab35cb4",
		GroupId:         "",
		VpcId:           "vpc-05137e7d412f0382b",
		IpProtocol:      "tcp",
		PortFrom:        121,
		PortTo:          122,
		CidrIp:          "192.168.1.1/24",
	}
	err = client.AddEgressSecurityGroupRule(req)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
}

func TestAwsCreateVpc(t *testing.T) {
	client, err := getAwsClient()
	if err != nil {
		t.Log(err)
		return
	}

	var resStr string
	vpc, err := client.CreateVPC(cloud.CreateVpcRequest{
		VpcName:   "chenxudong-vpc",
		CidrBlock: "10.0.0.0/16",
	})
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	resStr, _ = jsoniter.MarshalToString(vpc)
	t.Log(resStr)
}

func TestAwsGetVpc(t *testing.T) {
	client, err := getAwsClient()
	if err != nil {
		t.Log(err)
		return
	}

	var resStr string
	vpc, err := client.GetVPC(cloud.GetVpcRequest{
		RegionId: "cn-north-1",
		VpcId:    "vpc-05ac2dba6c6e3a683",
		//VpcName: "chenxudong-vpc",
	})
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	resStr, _ = jsoniter.MarshalToString(vpc)
	t.Log(resStr)
}

func TestAwsDescribeVpcs(t *testing.T) {
	client, err := getAwsClient()
	if err != nil {
		t.Log(err)
		return
	}

	var resStr string
	vpc, err := client.DescribeVpcs(cloud.DescribeVpcsRequest{
		RegionId: "cn-north-1",
	})
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	resStr, _ = jsoniter.MarshalToString(vpc)
	t.Log(resStr)
}

func TestAwsCreateSwitch(t *testing.T) {
	client, err := getAwsClient()
	if err != nil {
		t.Log(err)
		return
	}

	var res interface{}
	res, err = client.CreateSwitch(cloud.CreateSwitchRequest{
		RegionId:    "cn-north-1",
		ZoneId:      "cnn1-az1",
		CidrBlock:   "172.31.0.0/16",
		VSwitchName: "cxd-subnet",
		VpcId:       "vpc-05ac2dba6c6e3a683",
		GatewayIp:   "10.8.63.254",
	})
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	resStr, _ := jsoniter.MarshalToString(res)
	t.Log(resStr)
}

func TestAwsGetSwitch(t *testing.T) {
	client, err := getAwsClient()
	if err != nil {
		t.Log(err)
		return
	}

	var res interface{}
	res, err = client.GetSwitch(cloud.GetSwitchRequest{
		SwitchId: "subnet-0afa54a295bae7448",
	})
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	resStr, _ := jsoniter.MarshalToString(res)
	t.Log(resStr)
}

func TestAwsDescribeSwitches(t *testing.T) {
	client, err := getAwsClient()
	if err != nil {
		t.Log(err)
		return
	}

	var res interface{}
	res, err = client.DescribeSwitches(cloud.DescribeSwitchesRequest{
		VpcId: "vpc-05ac2dba6c6e3a683",
	})
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	resStr, _ := jsoniter.MarshalToString(res)
	t.Log(resStr)
}

func TestAwsShowVpc(t *testing.T) {
	client, err := getAwsClient()
	if err != nil {
		t.Log(err)
		return
	}

	var res interface{}
	var resStr string
	vpcId := ""
	swId := ""
	res, err = client.GetVPC(cloud.GetVpcRequest{
		VpcId: vpcId,
	})
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	resStr, _ = jsoniter.MarshalToString(res)
	t.Log(resStr)

	res, err = client.DescribeVpcs(cloud.DescribeVpcsRequest{})
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	resStr, _ = jsoniter.MarshalToString(res)
	t.Log(resStr)

	res, err = client.GetSwitch(cloud.GetSwitchRequest{
		SwitchId: swId,
	})
	if err != nil {
		t.Log(err.Error())
		return
	}
	resStr, _ = jsoniter.MarshalToString(res)
	t.Log(resStr)

	res, err = client.DescribeSwitches(cloud.DescribeSwitchesRequest{
		VpcId: vpcId,
	})
	if err != nil {
		t.Log(err.Error())
		return
	}
	resStr, _ = jsoniter.MarshalToString(res)
	t.Log(resStr)
}

func TestAwsBatchCreate(t *testing.T) {
	client, err := getAwsClient()
	if err != nil {
		t.Log(err)
		return
	}

	var res interface{}
	param := cloud.Params{
		InstanceType: "t2.micro",
		ImageId:      "ami-05248307900d52e3a",
		Network: &cloud.Network{
			VpcId:                   "vpc-05ac2dba6c6e3a683",
			SubnetId:                "subnet-07d13eb9dc425a7ef",
			SecurityGroup:           "sg-0cfbf25f5a5edd343",
			InternetChargeType:      cloud.BandwidthPayByTraffic,
			InternetMaxBandwidthOut: 0,
			InternetIpType:          "5_bgp",
		},
		Disks: &cloud.Disks{
			SystemDisk: cloud.DiskConf{Size: 10, Category: "standard"},
			DataDisk:   []cloud.DiskConf{{Size: 10, Category: "standard"}},
		},
		//Charge: &cloud.Charge{
		//	ChargeType: cloud.InstanceChargeTypePostPaid,
		//	Period:     1,
		//	PeriodUnit: "Month",
		//},
		//Password: "",
		Tags: []cloud.Tag{
			{
				Key:   cloud.TaskId,
				Value: "12345",
			},
			{
				Key:   cloud.ClusterName,
				Value: "cluster-cxd",
			},
		},
		//DryRun: true,
	}
	res, err = client.BatchCreate(param, 1)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	resStr, _ := jsoniter.MarshalToString(res)
	t.Log(resStr)
}

func TestAwsGetInstances(t *testing.T) {
	client, err := getAwsClient()
	if err != nil {
		t.Log(err)
		return
	}

	var res interface{}
	res, err = client.GetInstances([]string{"i-0dc582750a8bec681"})
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	resStr, _ := jsoniter.MarshalToString(res)
	t.Log(resStr)
}

func TestAwsGetInstancesByTags(t *testing.T) {
	client, err := getAwsClient()
	if err != nil {
		t.Log(err)
		return
	}

	var res interface{}
	res, err = client.GetInstancesByTags("", []cloud.Tag{{Key: cloud.TaskId, Value: "12345"}})
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	resStr, _ := jsoniter.MarshalToString(res)
	t.Log(resStr)
}

func TestAwsBatchDelete(t *testing.T) {
	client, err := getAwsClient()
	if err != nil {
		t.Log(err)
		return
	}

	err = client.BatchDelete([]string{"i-04f6969c4e64c2f4c", "i-012e3dad520396b7b", "i-0dc582750a8bec681"}, "")
	if err != nil {
		t.Errorf(err.Error())
		return
	}
}

func TestAwsStartInstances(t *testing.T) {
	client, err := getAwsClient()
	if err != nil {
		t.Log(err)
		return
	}

	err = client.StartInstances([]string{"i-0ec1b07d7330c90a1", "i-0bb1517d5a6d77fdb", "i-0ca37c8b6424de153", "i-08a04b43d895683ac", "i-0d91c65c32dd0cdcf", "i-0fdb3872a43b072b1", "i-06a39390d74a6678c", "i-026746cf9ac888eb8"})
	if err != nil {
		t.Errorf(err.Error())
		return
	}
}

func TestAwsStopInstances(t *testing.T) {
	client, err := getAwsClient()
	if err != nil {
		t.Log(err)
		return
	}

	err = client.StopInstances([]string{"i-04f6969c4e64c2f4c", "i-012e3dad520396b7b", "i-0dc582750a8bec681"})
	if err != nil {
		t.Errorf(err.Error())
		return
	}
}
