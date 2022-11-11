package baidu

import "github.com/galaxy-future/BridgX/pkg/cloud"

func (b BaiduCloud) ListObjects(endpoint, bucketName, prefix string) (objects []cloud.ObjectProperties, err error) {
	if res, err := b.bosClient.ListBuckets(); err != nil {
		return nil, err
	} else {
		for _, buck := range res.Buckets {
			if buck.Location == endpoint && buck.Name == bucketName {
				listObjectResult, err := b.bosClient.ListObjects(bucketName, nil)
				if err != nil {
					return nil, err
				}
				for _, obj := range listObjectResult.Contents {
					object := cloud.ObjectProperties{
						Name: obj.Key,
					}
					objects = append(objects, object)
				}
				break
			}
		}
	}
	return
}
func (b BaiduCloud) ListBucket(endpoint string) ([]cloud.BucketProperties, error) {
	buckets := []cloud.BucketProperties{}
	if res, err := b.bosClient.ListBuckets(); err != nil {
		return nil, err
	} else {
		for _, b := range res.Buckets {
			if b.Location == endpoint {
				bucket := cloud.BucketProperties{
					Name: b.Name,
				}
				buckets = append(buckets, bucket)
			}
		}
	}
	return buckets, nil
}

func (b BaiduCloud) GetOssDownloadUrl(endpoint, bucketName, region string) string {
	// todo
	return ""
}

func (b BaiduCloud) GetObjectDownloadUrl(bucketName, objectName string) (string, error) {
	url := b.bosClient.BasicGeneratePresignedUrl(bucketName, objectName, 300)
	return url, nil
}
