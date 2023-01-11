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
	var Images []cloud.Image

	if req.ImageType == cloud.ImageShared {
		// 查询共享镜像列表
		request := &model.ListShareImageRequest{}
		response, err := p.imsClient.ListShareImage(request)
		if err != nil {
			return cloud.DescribeImagesResponse{}, errors.New(response.ErrorMessage)
		}
		if response.State == model.ListShareImageResponseStateEnumOk {
			for _, v := range *response.Body.Content {
				var image cloud.Image
				image.ImageId = v.ImageId
				image.OsName = v.OsName
				image.Size = int(*v.Size)
				image.OsType = string(v.OsType)
				Images = append(Images, image)
			}
			return cloud.DescribeImagesResponse{Images: Images}, nil
		}
	} else if req.ImageType == cloud.ImageGlobal {
		// 查询自定义镜像列表
		request := &model.ListImageRespRequest{}
		response, err := p.imsClient.ListImageResp(request)
		if err != nil {
			return cloud.DescribeImagesResponse{}, errors.New(response.ErrorMessage)
		}
		if response.State == model.ListImageRespResponseStateEnumOk {
			for _, v := range *response.Body.Content {
				var image cloud.Image
				image.ImageId = v.ImageId
				image.OsName = v.OsName
				image.Size = int(*v.Size)
				image.OsType = string(v.OsType)
				image.Platform = string(v.PublicImageType)
				Images = append(Images, image)
			}
			return cloud.DescribeImagesResponse{Images: Images}, nil
		}

	} else if req.ImageType == cloud.ImagePrivate {
		// 查询镜像信息
		request := &model.IMSgetImageRespRequest{}
		request.IMSgetImageRespPath.ImageId = req.ImageType
		response, err := p.imsClient.IMSgetImageResp(request)
		if err != nil {
			return cloud.DescribeImagesResponse{}, errors.New(response.ErrorMessage)
		}
		if response.State == model.IMSgetImageRespResponseStateEnumOk {
			var image cloud.Image
			image.ImageId = response.Body.ImageId
			image.OsName = response.Body.OsName
			image.Size = int(*response.Body.Size)
			image.OsType = string(response.Body.OsType)
			image.Platform = string(response.Body.PublicImageType)
			Images = append(Images, image)
			return cloud.DescribeImagesResponse{Images: Images}, nil
		}
	}

	return cloud.DescribeImagesResponse{Images: Images}, nil
}
