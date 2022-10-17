package tencent

import "github.com/galaxy-future/BridgX/pkg/cloud"

func (p *TencentCloud) ListObjects(endpoint, bucketName, prefix string) (res []cloud.ObjectProperties, err error) {

	return
}

func (p *TencentCloud) ListBucket(endpoint string) (res []cloud.BucketProperties, err error) {

	return
}

func (p *TencentCloud) GetOssDownloadUrl(endpoint, bucketName, region string) string {

	return ""
}

func (p *TencentCloud) GetObjectDownloadUrl(bucketName, objectKey string) (string, error) {

	return "", nil
}
