package ecloud

import (
	"github.com/galaxy-future/BridgX/internal/logs"
	"github.com/galaxy-future/BridgX/pkg/cloud"
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	logs.Init()
}

func TestEcloud_DescribeKeyPairs(t *testing.T) {
	p, err := New(_AK, _SK, _regionId)
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}
	describeKeyPairs, err := p.DescribeKeyPairs(cloud.DescribeKeyPairsRequest{
		PageSize:   10,
		PageNumber: 1})

	if err != nil {
		t.Errorf("err:%v", err)
		return
	}
	t.Log(describeKeyPairs.KeyPairs)
	assert.Equal(t, describeKeyPairs.TotalCount, len(describeKeyPairs.KeyPairs))
}

func TestEcloud_CreateKeyPair(t *testing.T) {
	p, err := New(_AK, _SK, _regionId)
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}
	createKeyPair, err := p.CreateKeyPair(cloud.CreateKeyPairRequest{
		KeyPairName: "test002",
	})

	if err != nil {
		t.Errorf("err:%v", err)
		return
	}
	t.Log(createKeyPair)
	assert.Equal(t, createKeyPair.KeyPairName, "test002")
}
