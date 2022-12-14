package tencent

import (
	"fmt"

	"github.com/galaxy-future/BridgX/pkg/cloud"
	"github.com/galaxy-future/BridgX/pkg/utils"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

//默认projectId 0
const _projectId int64 = 0

func (p *TencentCloud) CreateKeyPair(req cloud.CreateKeyPairRequest) (cloud.CreateKeyPairResponse, error) {
	request := cvm.NewCreateKeyPairRequest()
	request.KeyName = common.StringPtr(req.KeyPairName)
	request.ProjectId = common.Int64Ptr(_projectId)
	response, err := p.cvmClient.CreateKeyPair(request)
	if err != nil {
		return cloud.CreateKeyPairResponse{}, err
	}
	if response == nil || response.Response == nil || response.Response.KeyPair == nil {
		return cloud.CreateKeyPairResponse{}, fmt.Errorf("CreateKeyPair TencentCloud resp nil, req:[%v]", req)
	}
	return keypair2CreateResponse(response.Response.KeyPair), err
}

func (p *TencentCloud) ImportKeyPair(req cloud.ImportKeyPairRequest) (cloud.ImportKeyPairResponse, error) {
	request := cvm.NewImportKeyPairRequest()
	request.PublicKey = common.StringPtr(req.PublicKey)
	request.KeyName = common.StringPtr(req.KeyPairName)
	request.ProjectId = common.Int64Ptr(_projectId)
	request.PublicKey = common.StringPtr(req.PublicKey)
	response, err := p.cvmClient.ImportKeyPair(request)
	if err != nil {
		return cloud.ImportKeyPairResponse{}, err
	}
	if response == nil || response.Response == nil || response.Response.KeyId == nil {
		return cloud.ImportKeyPairResponse{}, fmt.Errorf("ImportKeyPair TencentCloud resp nil, req:[%v]", req)
	}
	return cloud.ImportKeyPairResponse{
		KeyPairId: *response.Response.KeyId,
	}, err

}

func (p *TencentCloud) DescribeKeyPairs(req cloud.DescribeKeyPairsRequest) (cloud.DescribeKeyPairsResponse, error) {
	request := cvm.NewDescribeKeyPairsRequest()
	if req.PageNumber < 1 {
		req.PageNumber = 1
	}
	if req.PageSize < 1 || req.PageSize > _pageSize {
		req.PageSize = _pageSize
	}
	request.Offset = common.Int64Ptr(int64((req.PageNumber - 1) * req.PageSize))
	request.Limit = common.Int64Ptr(int64(req.PageSize))
	response, err := p.cvmClient.DescribeKeyPairs(request)
	if err != nil {
		return cloud.DescribeKeyPairsResponse{}, err
	}
	if response != nil && response.Response != nil && response.Response.KeyPairSet != nil {
		ret := cloud.DescribeKeyPairsResponse{
			KeyPairs:   []cloud.KeyPair{},
			TotalCount: int(utils.Int64Value(response.Response.TotalCount)),
		}
		for _, v := range response.Response.KeyPairSet {
			ret.KeyPairs = append(ret.KeyPairs, keypair2Cloud(v))
		}
		return ret, nil
	}
	return cloud.DescribeKeyPairsResponse{}, err
}

func keypair2Cloud(keyPair *cvm.KeyPair) cloud.KeyPair {
	return cloud.KeyPair{
		KeyPairName: *keyPair.KeyName,
		KeyPairId:   *keyPair.KeyId,
	}
}

func keypair2CreateResponse(keyPair *cvm.KeyPair) cloud.CreateKeyPairResponse {
	return cloud.CreateKeyPairResponse{
		KeyPairName: *keyPair.KeyName,
		KeyPairId:   *keyPair.KeyId,
		PublicKey:   *keyPair.PublicKey,
		PrivateKey:  *keyPair.PrivateKey,
	}
}
