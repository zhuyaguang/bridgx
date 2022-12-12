package ecloud

import (
	"fmt"
	"github.com/galaxy-future/BridgX/internal/logs"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/galaxy-future/BridgX/pkg/cloud"
)

func getOssEndpoint(region string) string {
	// 华东-上海1	shanghai1	eos-shanghai-1.cmecloud.cn
	regionName := region[0 : len(region)-1]
	regionID := region[len(region)-1:]
	return fmt.Sprintf("https://eos-%s-%s.cmecloud.cn", regionName, regionID)
}

func (p *ECloud) ListObjects(endpoint, bucketName, prefix string) ([]cloud.ObjectProperties, error) {
	var objectPropertiesList []cloud.ObjectProperties
	p.eosSession.Config.Endpoint = &endpoint
	svc := s3.New(p.eosSession)

	params := &s3.ListObjectsInput{
		Bucket: aws.String(bucketName),
		Prefix: &prefix,
	}
	// 列举文件时，最多一次性列举 1000 个文件
	resp, err := svc.ListObjects(params)
	if err != nil {
		logs.Logger.Errorf("Unable to list items in bucket %q, %v\n", bucketName, err)
		return objectPropertiesList, err
	}
	for _, item := range resp.Contents {
		objectProperties := cloud.ObjectProperties{Name: *item.Key}
		objectPropertiesList = append(objectPropertiesList, objectProperties)
	}
	return objectPropertiesList, nil
}

func (p *ECloud) ListBucket(endpoint string) ([]cloud.BucketProperties, error) {
	var bucketPropertiesList []cloud.BucketProperties
	p.eosSession.Config.Endpoint = &endpoint
	svc := s3.New(p.eosSession)

	result, err := svc.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		logs.Logger.Errorf("Unable to list buckets, %v\n", err)
		return bucketPropertiesList, err
	}

	for _, b := range result.Buckets {
		bucketProperties := cloud.BucketProperties{Name: *b.Name}
		bucketPropertiesList = append(bucketPropertiesList, bucketProperties)
	}
	return bucketPropertiesList, nil
}

func (p *ECloud) GetOssDownloadUrl(endpoint, bucketName, region string) string {
	str := strings.Split(endpoint, "//")
	if len(str) != 2 {
		regionName := region[0 : len(region)-1]
		regionID := region[len(region)-1:]
		return fmt.Sprintf("https://eos-%s-%s.cmecloud.cn", regionName, regionID)
	}
	return fmt.Sprintf("https://%s", str[1])
}

func (p *ECloud) GetObjectDownloadUrl(bucketName, objectKey string) (string, error) {
	svc := s3.New(p.eosSession)

	params := &s3.GetBucketLocationInput{
		Bucket: aws.String(bucketName),
	}
	bucketLocation, err := svc.GetBucketLocation(params)
	if err != nil {
		logs.Logger.Errorf("Unable to get bucketLocation %q, %v", bucketName, err)
	}
	*p.eosSession.Config.Endpoint = bucketLocation.String()

	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	url, err := req.Presign(time.Hour)
	if err != nil {
		return "", err
	}
	return url, nil
}
