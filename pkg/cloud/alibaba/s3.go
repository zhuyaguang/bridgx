package alibaba

import (
	"fmt"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/galaxy-future/BridgX/internal/logs"
	"github.com/galaxy-future/BridgX/pkg/cloud"
)

func (p *AlibabaCloud) ListObjects(endpoint, bucketName, prefix string) (ls []cloud.ObjectProperties, err error) {
	p.ossClient.Config.Endpoint = endpoint
	bucket, err := p.ossClient.Bucket(bucketName)
	if err != nil {
		return
	}
	continueToken := ""
	for {
		lsRes, lErr := bucket.ListObjectsV2(oss.ContinuationToken(continueToken), oss.Prefix(prefix))
		if lErr != nil {
			logs.Logger.Errorf("")
		}
		for _, l := range lsRes.Objects {
			ls = append(ls, cloud.ObjectProperties{
				Name: l.Key,
			})
		}
		if lsRes.IsTruncated {
			continueToken = lsRes.NextContinuationToken
		} else {
			break
		}
	}
	return
}

func getOssEndpoint(region string) string {
	return fmt.Sprintf("https://oss-%s.aliyuncs.com", region)
}

func (p *AlibabaCloud) GetOssDownloadUrl(endpoint, bucketName, region string) string {
	str := strings.Split(endpoint, "//")
	if len(str) != 2 {
		return fmt.Sprintf("https://%s.oss-%s.aliyuncs.com", bucketName, region)
	}
	return fmt.Sprintf("https://%s.%s", bucketName, str[1])
}

func (p *AlibabaCloud) GetObjectDownloadUrl(BucketName, ObjectKey string) (string, error) {
	bucket, err := p.ossClient.Bucket(BucketName)
	if err != nil {
		return "", err
	}
	url, err := bucket.SignURL(ObjectKey, oss.HTTPGet, 3600)
	if err != nil {
		return "", err
	}
	return url, nil
}

func (p *AlibabaCloud) ListBucket(endpoint string) (buckets []cloud.BucketProperties, err error) {
	marker := ""
	for {
		p.ossClient.Config.Endpoint = endpoint
		lsRes, lErr := p.ossClient.ListBuckets(oss.Marker(marker))
		if lErr != nil {
			logs.Logger.Errorf("[ListBucket] error: %s", lErr.Error())
		}
		for _, l := range lsRes.Buckets {
			buckets = append(buckets, cloud.BucketProperties{
				Name: l.Name,
			})
		}

		if lsRes.IsTruncated {
			marker = lsRes.NextMarker
		} else {
			break
		}
	}
	return
}
