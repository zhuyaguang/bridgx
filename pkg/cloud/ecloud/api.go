package ecloud

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/galaxy-future/BridgX/pkg/cloud"
	"gitlab.ecloud.com/ecloud/ecloudsdkcore/config"
	"gitlab.ecloud.com/ecloud/ecloudsdkecs"
	"gitlab.ecloud.com/ecloud/ecloudsdkeip"
	"gitlab.ecloud.com/ecloud/ecloudsdkims"
	"gitlab.ecloud.com/ecloud/ecloudsdkims/model"
)

type ECloud struct {
	eipClient  *ecloudsdkeip.Client
	ecsClient  *ecloudsdkecs.Client
	imsClient  *ecloudsdkims.Client
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

	imsClient := ecloudsdkims.NewClient(&config.Config{
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
		imsClient:  imsClient,
		eosSession: eosSession,
	}, nil
}

func (p *ECloud) ProviderType() string {
	return cloud.ECloud
}

func (p *ECloud) DescribeImages(req cloud.DescribeImagesRequest) (cloud.DescribeImagesResponse, error) {
	request := &model.IMSgetImageRespRequest{}
	iMSgetImageRespPath := &model.IMSgetImageRespPath{}
	iMSgetImageRespPath.ImageId = req.ImageID
	request.IMSgetImageRespPath = iMSgetImageRespPath
	response, err := p.imsClient.IMSgetImageResp(request)
	if err == nil {
		fmt.Printf("%+v\n", response)
		if response.State == model.IMSgetImageRespResponseStateEnumOk {
			var images cloud.Image
			images.ImageName = response.Body.Name
			images.ImageId = response.Body.ImageId
			images.OsName = response.Body.OsName
			images.Size = int(*response.Body.Size)
			images.OsType = string(response.Body.OsType)
			images.Platform = string(response.Body.PublicImageType)
			return cloud.DescribeImagesResponse{Images: []cloud.Image{images}}, nil
		} else {
			err := errors.New(response.ErrorMessage)
			return cloud.DescribeImagesResponse{}, err
		}
	}
	return cloud.DescribeImagesResponse{}, err
}
