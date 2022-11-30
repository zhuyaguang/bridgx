package tencent

import (
	"log"
	"testing"

	"github.com/galaxy-future/BridgX/pkg/cloud"
	"github.com/stretchr/testify/assert"
)

const (
	_AK        = ""
	_SK        = ""
	_PublicKey = ""
)

var client *TencentCloud
var clientErr error

func TestMain(m *testing.M) {
	client, clientErr = New(_AK, _SK, "ap-guangzhou")
	if clientErr != nil {
		log.Println(clientErr)
	}
	m.Run()
}

func TestTencentCloud_CreateKeyPair(t *testing.T) {

	response, err := client.CreateKeyPair(cloud.CreateKeyPairRequest{
		KeyPairName: "test002",
	})
	if err != nil {
		t.Log(err)
	}
	t.Log(response)
}

func TestTencentCloud_ImportKeyPair(t *testing.T) {
	response, err := client.ImportKeyPair(cloud.ImportKeyPairRequest{
		KeyPairName: "test005",
		PublicKey:   _PublicKey,
	})
	if err != nil {
		t.Log(err)
	}
	t.Log(response.KeyPairId)
}

func TestTencentCloud_DescribeKeyPairs(t *testing.T) {
	response, err := client.DescribeKeyPairs(cloud.DescribeKeyPairsRequest{
		PageNumber: 1,
		PageSize:   _pageSize,
	})
	if err != nil {
		t.Log(err)
	}
	t.Log(response.TotalCount)
	assert.Equal(t, response.TotalCount, 5)
	t.Log(response.KeyPairs)
}
