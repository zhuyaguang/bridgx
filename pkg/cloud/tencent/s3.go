package tencent

import (
	"context"
	"fmt"
	"github.com/galaxy-future/BridgX/pkg/cloud"
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/url"
)

//GetOssDownloadUrl

func (p *TencentCloud) ListObjects(request cloud.ListObjectsRequest) (cloud.ListObjectsResponse,error) {
	cosURL := "https://" + request.BucketName + ".cos." + request.CosRegion + ".myqcloud.com"
	u, _ := url.Parse(cosURL)
	b := &cos.BaseURL{BucketURL: u}
	p.cosClient.BaseURL = b

	opt := &cos.BucketGetOptions{
		Prefix:  request.Prefix,
		MaxKeys: request.MaxKeys,
	}
	v, _, err := p.cosClient.Bucket.Get(context.Background(), opt)
	if err != nil {
		return cloud.ListObjectsResponse{},err
	}

	for _, c := range v.Contents {
		fmt.Printf("%s, %d\n", c.Key, c.Size)
	}
	return cloud.ListObjectsResponse{CosObjects: v.Contents},nil

}

func (p *TencentCloud) ListBucket() (cloud.ListBucketResponse, error) {
	s, _, err := p.cosClient.Service.Get(context.Background())
	if err != nil {
		return cloud.ListBucketResponse{}, err
	}

	for _, b := range s.Buckets {
		fmt.Printf("%#v\n", b)
	}
	return cloud.ListBucketResponse{ CosBucket: s.Buckets}, nil
}

//func GetOssDownloadUrl() string {
//
//}

func (p *TencentCloud) GetObjectDownloadUrl(req cloud.GetObjectDownloadUrlRequest) cloud.GetObjectDownloadUrlResponse {
	cosURL := "https://" + req.BucketName + ".cos." + req.CosRegion + ".myqcloud.com"
	u, _ := url.Parse(cosURL)
	b := &cos.BaseURL{BucketURL: u}
	p.cosClient.BaseURL = b
	oURL := p.cosClient.Object.GetObjectURL(req.Key)
	fmt.Println(oURL)
	return cloud.GetObjectDownloadUrlResponse{URL: oURL}
}
