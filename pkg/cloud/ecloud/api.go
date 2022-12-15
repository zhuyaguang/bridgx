package ecloud

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/galaxy-future/BridgX/pkg/cloud"
	"gitlab.ecloud.com/ecloud/ecloudsdkcore/config"
	"gitlab.ecloud.com/ecloud/ecloudsdkecs"
	"gitlab.ecloud.com/ecloud/ecloudsdkeip"
)

type ECloud struct {
	eipClient  *ecloudsdkeip.Client
	ecsClient  *ecloudsdkecs.Client
	eosSession *session.Session
}

func New(ak, sk, regionId string) (*ECloud, error) {
	eipClient := ecloudsdkeip.NewClient(&config.Config{
		AccessKey: ak,
		SecretKey: sk,
		PoolId:    regionId,
	})
	ecsClient := ecloudsdkecs.NewClient(&config.Config{
		AccessKey: ak,
		SecretKey: sk,
		PoolId:    regionId,
	})

	endPoint := getOssEndpoint(regionId)
	disableSSL := false
	sessionConfig := &aws.Config{
		Region:           aws.String(regionId),
		Endpoint:         &endPoint,
		S3ForcePathStyle: aws.Bool(true),
		Credentials:      credentials.NewStaticCredentials(ak, sk, ""),
		DisableSSL:       &disableSSL,
	}
	eosSession, err := session.NewSession(sessionConfig)
	if err != nil {
		return nil, err
	}

	return &ECloud{
		eipClient:  eipClient,
		ecsClient:  ecsClient,
		eosSession: eosSession,
	}, nil
}

func (p *ECloud) ProviderType() string {
	return cloud.ECloud
}

func (p *ECloud) DescribeImages(req cloud.DescribeImagesRequest) (cloud.DescribeImagesResponse, error) {
	// TODO implement me
	return cloud.DescribeImagesResponse{}, errors.New("implement me")
}
