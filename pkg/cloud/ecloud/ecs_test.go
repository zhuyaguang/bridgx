package ecloud

import (
	"github.com/galaxy-future/BridgX/pkg/cloud"
	"testing"
)

func TestECloud_BatchCreate(t *testing.T) {
	client, err := New("_ak", "_sk", "ap-beijing")
	if err != nil {
		t.Log(err.Error())
		return
	}

	res, err := client.BatchCreate(cloud.Params{
		Provider:     "",
		InstanceType: "common",
		ImageId:      "",
		Network:      nil,
		Zone:         "",
		Region:       "ap-beijing",
		Disks:        nil,
		Charge: &cloud.Charge{
			Period:     0,
			PeriodUnit: cloud.Year,
			ChargeType: "",
		},
		Password:    "",
		Tags:        nil,
		DryRun:      false,
		KeyPairId:   "",
		KeyPairName: "",
	}, 3)

	t.Log(res)
}

func TestGetInstances(t *testing.T) {
	client, err := New(_AK, _SK, "ap-beijing")
	if err != nil {
		t.Log(err.Error())
		return
	}
	res, err := client.GetInstances([]string{"id1", "id2"})
	if err != nil {
		t.Log(err.Error())
		return
	}
	for _, r := range res {
		t.Log(r)
	}
}

func TestGetInstanceStatus(t *testing.T) {
	client, err := New(_AK, _SK, "ap-beijing")
	if err != nil {
		t.Log(err.Error())
		return
	}
	res, err := client.GetInstanceStatus("testID")
	if err != nil {
		t.Log(err.Error())
		return
	}
	t.Log(res)
}

func TestBatchDelete(t *testing.T) {
	client, err := New(_AK, _SK, "ap-beijing")
	if err != nil {
		t.Log(err.Error())
		return
	}
	err = client.BatchDelete([]string{"id1", "id2"}, "ap-beijing")
	if err != nil {
		t.Log(err.Error())
		return
	}
}

func TestStartInstances(t *testing.T) {
	client, err := New(_AK, _SK, "ap-beijing")
	if err != nil {
		t.Log(err.Error())
		return
	}
	err = client.StartInstances([]string{"id1", "id2"})
	if err != nil {
		t.Log(err.Error())
		return
	}
}

func TestStopInstances(t *testing.T) {
	client, err := New(_AK, _SK, "ap-beijing")
	if err != nil {
		t.Log(err.Error())
		return
	}
	err = client.StopInstances([]string{"id1", "id2"})
	if err != nil {
		t.Log(err.Error())
		return
	}
}
