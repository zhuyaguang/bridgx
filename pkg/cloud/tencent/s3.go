package tencent

import (
	"context"
	"fmt"
	"github.com/galaxy-future/BridgX/pkg/cloud"
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/url"
)

func getCosEndpoint(bucketName, region string) string {
	//cosURL := "https://" + bucketName + ".cos." + region + ".myqcloud.com"
	return fmt.Sprintf("https://%s.cos.%s.myqcloud.com", bucketName, region)
}

// ListObjects
// endpoint Cos 的静态地址
// bucketName 被忽略
func (p *TencentCloud) ListObjects(endpoint, bucketName, prefix string) (res []cloud.ObjectProperties, err error) {
	u, _ := url.Parse(endpoint)
	b := &cos.BaseURL{BucketURL: u}
	p.cosClient.BaseURL = b

	opt := &cos.BucketGetOptions{
		Prefix: prefix,
	}
	v, _, err := p.cosClient.Bucket.Get(context.Background(), opt)
	if err != nil {
		return []cloud.ObjectProperties{}, err
	}
	var objectList []cloud.ObjectProperties
	for _, c := range v.Contents {
		objectElemt := cloud.ObjectProperties{
			Name: c.Key,
		}
		objectList = append(objectList, objectElemt)
		fmt.Printf("%s, %d\n", c.Key, c.Size)
	}
	return objectList, nil
}

func (p *TencentCloud) ListBucket(endpoint string) (res []cloud.BucketProperties, err error) {
	s, _, err := p.cosClient.Service.Get(context.Background())
	if err != nil {
		return []cloud.BucketProperties{}, err
	}
	var bucketList []cloud.BucketProperties
	for _, b := range s.Buckets {
		bucket := cloud.BucketProperties{
			Name: b.Name,
		}
		bucketList = append(bucketList, bucket)
		fmt.Printf("%#v\n", b)
	}
	return bucketList, nil
}

func (p *TencentCloud) GetOssDownloadUrl(endpoint, bucketName, region string) string {
	return ""
}

func (p *TencentCloud) GetObjectDownloadUrl(endpoint, objectKey string) (string, error) {
	u, _ := url.Parse(endpoint)
	b := &cos.BaseURL{BucketURL: u}
	p.cosClient.BaseURL = b
	oURL := p.cosClient.Object.GetObjectURL(objectKey)
	fmt.Println(oURL)
	return oURL.String(), nil
}
