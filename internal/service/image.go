package service

import (
	"context"

	"github.com/galaxy-future/BridgX/internal/types"
	"github.com/galaxy-future/BridgX/pkg/cloud"
)

type GetImagesRequest struct {
	Account   *types.OrgKeys
	Provider  string
	RegionId  string
	InsType   string
	ImageType string
}

func GetImages(ctx context.Context, req GetImagesRequest) ([]cloud.Image, error) {
	ak := getFirstAk(req.Account, req.Provider)
	p, err := getProvider(req.Provider, ak, req.RegionId)
	if err != nil {
		return []cloud.Image{}, err
	}
	imagesRes, err := p.DescribeImages(cloud.DescribeImagesRequest{
		RegionId:  req.RegionId,
		InsType:   req.InsType,
		ImageType: req.ImageType,
	})
	if err != nil {
		return []cloud.Image{}, err
	}
	return imagesRes.Images, nil
}
