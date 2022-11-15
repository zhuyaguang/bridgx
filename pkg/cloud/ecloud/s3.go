package ecloud

import (
	// "github.com/aws/aws-sdk-go/service/s3"
	"github.com/galaxy-future/BridgX/pkg/cloud"
)

func getOssEndpoint(region string) string {
	return "eos-ningbo-1.cmecloud.cn"
	// return fmt.Sprintf("https://eos-%s.cmecloud.cn", region)
}

func (p *ECloud) ListObjects(endpoint, bucketName, prefix string) ([]cloud.ObjectProperties, error) {
	// TODO implement me
	panic("implement me")
}

func (p *ECloud) ListBucket(endpoint string) ([]cloud.BucketProperties, error) {
	// TODO implement me
	panic("implement me")
}

func (p *ECloud) GetOssDownloadUrl(s string, s2 string, s3 string) string {
	// TODO implement me
	panic("implement me")
}

func (p *ECloud) GetObjectDownloadUrl(bucketName, objectKey string) (string, error) {
	// TODO implement me
	panic("implement me")
}
