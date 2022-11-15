package huawei

import "github.com/galaxy-future/BridgX/pkg/cloud"

func (p *HuaweiCloud) ListObjects(endpoint, bucketName, prefix string) (res []cloud.ObjectProperties, err error) {

	return
}

func (p *HuaweiCloud) ListBucket(endpoint string) (res []cloud.BucketProperties, err error) {

	return
}

func (p *HuaweiCloud) GetOssDownloadUrl(endpoint, bucketName, region string) string {

	return ""
}

func (p *HuaweiCloud) GetObjectDownloadUrl(BucketName, ObjectKey string) (string, error) {

	return "", nil
}
