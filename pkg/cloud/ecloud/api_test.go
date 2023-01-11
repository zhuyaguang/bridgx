package ecloud

import (
	"github.com/galaxy-future/BridgX/pkg/cloud"
	"testing"
)

func TestECloud_DescribeImages(t *testing.T) {
	client, err := New(_AK, _SK, _regionId)
	if err != nil {
		t.Log(err.Error())
		return
	}

	res, err := client.DescribeImages(cloud.DescribeImagesRequest{RegionId: "shanghai", ImageType: "shared", InsType: "InsType"})
	t.Log(res.Images)
}
