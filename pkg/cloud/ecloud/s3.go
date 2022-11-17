package ecloud

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"time"

	"github.com/galaxy-future/BridgX/pkg/cloud"
)

func getOssEndpoint(region string) string {
	return "eos-ningbo-1.cmecloud.cn"
	// return fmt.Sprintf("https://eos-%s.cmecloud.cn", region)
}

func (p *ECloud) ListObjects(endpoint, bucketName, prefix string) ([]cloud.ObjectProperties, error) {
	var ObjectPropertiesList []cloud.ObjectProperties
	svc := s3.New(p.eosSession)

	params := &s3.ListObjectsInput{
		Bucket: aws.String(bucketName),
		Prefix: &prefix,
	}
	// 列举文件时，最多一次性列举 1000 个文件
	resp, err := svc.ListObjects(params)
	if err != nil {
		fmt.Printf("Unable to list items in bucket %q, %v\n", bucketName, err)
		return ObjectPropertiesList, err
	}
	for _, item := range resp.Contents {
		objectProperties := cloud.ObjectProperties{Name: *item.Key}
		ObjectPropertiesList = append(ObjectPropertiesList, objectProperties)
	}
	return ObjectPropertiesList, nil
}

func (p *ECloud) ListBucket(endpoint string) ([]cloud.BucketProperties, error) {
	var BucketPropertiesList []cloud.BucketProperties
	svc := s3.New(p.eosSession)
	result, err := svc.ListBuckets(nil)
	if err != nil {
		fmt.Printf("Unable to list buckets, %v\n", err)
		return BucketPropertiesList, err
	}

	for _, b := range result.Buckets {
		bucketProperties := cloud.BucketProperties{Name: *b.Name}
		BucketPropertiesList = append(BucketPropertiesList, bucketProperties)

	}
	return BucketPropertiesList, nil
}

func (p *ECloud) GetOssDownloadUrl(s string, s2 string, s3 string) string {
	// TODO implement me
	panic("implement me")
}

func (p *ECloud) GetObjectDownloadUrl(bucketName, objectKey string) (string, error) {
	svc := s3.New(p.eosSession)
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
