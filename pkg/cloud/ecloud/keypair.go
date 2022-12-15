package ecloud

import (
	"errors"

	"github.com/galaxy-future/BridgX/internal/logs"
	"github.com/galaxy-future/BridgX/pkg/cloud"
	"github.com/galaxy-future/BridgX/pkg/utils"
	"gitlab.ecloud.com/ecloud/ecloudsdkecs/model"
)

func (p *ECloud) DescribeGroupRules(req cloud.DescribeGroupRulesRequest) (cloud.DescribeGroupRulesResponse, error) {

	return cloud.DescribeGroupRulesResponse{}, errors.New("implement me")
}

func (p *ECloud) CreateKeyPair(req cloud.CreateKeyPairRequest) (cloud.CreateKeyPairResponse, error) {
	vmCreateKeypairRequest := &model.VmCreateKeypairRequest{
		VmCreateKeypairBody: &model.VmCreateKeypairBody{
			Name: req.KeyPairName,
		},
	}
	vmCreateKeypair, err := p.ecsClient.VmCreateKeypair(vmCreateKeypairRequest)
	if err != nil {
		logs.Logger.Errorf("CreateKeyPair Ecloud failed.err:[%v] req:[%v]", err, req)
		return cloud.CreateKeyPairResponse{}, err
	}
	if vmCreateKeypair.State != _State_OK {
		errMsg := vmCreateKeypair.ErrorMessage
		logs.Logger.Errorf("CreateKeyPair resp state not ok :[%v] req:[%v]", errMsg, req)
		return cloud.CreateKeyPairResponse{}, errors.New(errMsg)
	}
	if vmCreateKeypair.Body == nil {
		errMsg := "response.Body is null"
		logs.Logger.Errorf("CreateKeyPair Ecloud failed.err:[%v] req:[%v]", errMsg, req)
		return cloud.CreateKeyPairResponse{}, errors.New(errMsg)
	}
	vmGetKeyPairDetailRequest := &model.VmGetKeyPairDetailRequest{
		VmGetKeyPairDetailPath: &model.VmGetKeyPairDetailPath{
			KeypairName: req.KeyPairName,
		},
		VmGetKeyPairDetailQuery: &model.VmGetKeyPairDetailQuery{
			Region: req.RegionId,
		},
	}
	vmGetKeyPairDetail, err1 := p.ecsClient.VmGetKeyPairDetail(vmGetKeyPairDetailRequest)
	if err1 != nil {
		logs.Logger.Errorf("GetKeyPairDetail Ecloud failed.err:[%v] req:[%v]", err1, req)
		return cloud.CreateKeyPairResponse{}, err
	}
	if vmGetKeyPairDetail.Body == nil {
		errMsg := "response.Body is null"
		logs.Logger.Errorf("GetKeyPairDetail Ecloud failed.err:[%v] req:[%v]", errMsg, req)
		return cloud.CreateKeyPairResponse{}, errors.New(errMsg)
	}
	return cloud.CreateKeyPairResponse{
		KeyPairId:   vmGetKeyPairDetail.Body.Id,
		KeyPairName: vmGetKeyPairDetail.Body.Name,
		PrivateKey:  vmGetKeyPairDetail.Body.PrivateKey,
		PublicKey:   vmGetKeyPairDetail.Body.PublicKey,
	}, nil
}

func (p *ECloud) ImportKeyPair(req cloud.ImportKeyPairRequest) (cloud.ImportKeyPairResponse, error) {
	// TODO implement me
	return cloud.ImportKeyPairResponse{}, errors.New("implement me")
}

func (p *ECloud) DescribeKeyPairs(req cloud.DescribeKeyPairsRequest) (cloud.DescribeKeyPairsResponse, error) {
	vmListKeyPairRequest := &model.VmListKeyPairRequest{
		VmListKeyPairQuery: &model.VmListKeyPairQuery{
			Page: utils.Int32(int32(req.PageNumber)),
			Size: utils.Int32(int32(req.PageSize)),
		},
	}
	response, err := p.ecsClient.VmListKeyPair(vmListKeyPairRequest)
	if err != nil {
		logs.Logger.Errorf("DescribeKeyPairs Ecloud failed.err:[%v] req:[%v]", err, req)
		return cloud.DescribeKeyPairsResponse{}, err
	}
	if response.Body == nil {
		errMsg := "response.Body is null"
		logs.Logger.Errorf("DescribeKeyPairs Ecloud failed.err:[%v] req:[%v]", errMsg, req)
		return cloud.DescribeKeyPairsResponse{}, errors.New(errMsg)
	}
	rsp := cloud.DescribeKeyPairsResponse{
		TotalCount: int(*response.Body.Total),
	}
	if response.Body.Content != nil && len(*response.Body.Content) > 0 {
		for _, pair := range *response.Body.Content {
			rsp.KeyPairs = append(rsp.KeyPairs, cloud.KeyPair{
				KeyPairName: pair.Name,
				KeyPairId:   pair.Id,
			})
		}
	}

	return rsp, nil
}
