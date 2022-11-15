package aws

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/galaxy-future/BridgX/internal/logs"
	"github.com/galaxy-future/BridgX/pkg/cloud"
)

func (p *AWSCloud) ListObjects(endpoint, bucketName, prefix string) ([]cloud.ObjectProperties, error) {
	svc := s3.New(p.sess, &aws.Config{Endpoint: aws.String(endpoint)})
	continuationToken := ""
	res := make([]cloud.ObjectProperties, 0)
	for {
		input := &s3.ListObjectsV2Input{
			Bucket: aws.String(bucketName),
			Prefix: aws.String(prefix),
		}
		if continuationToken != "" {
			input.ContinuationToken = aws.String(continuationToken)
		}
		output, err := svc.ListObjectsV2(input)
		if err != nil {
			fmt.Println(err.Error())
			logs.Logger.Errorf("[ListObjects] error: %s", err.Error())
			return nil, err
		}
		for _, object := range output.Contents {
			res = append(res, cloud.ObjectProperties{
				Name: aws.StringValue(object.Key),
			})
		}
		if aws.BoolValue(output.IsTruncated) {
			continuationToken = aws.StringValue(output.NextContinuationToken)
			continue
		}
		break
	}
	return res, nil

}

func (p *AWSCloud) ListBucket(endpoint string) ([]cloud.BucketProperties, error) {
	svc := s3.New(p.sess, &aws.Config{Endpoint: aws.String(endpoint)})
	output, err := svc.ListBuckets(nil)
	if err != nil {
		logs.Logger.Errorf("[ListBucket] error: %s", err.Error())
		return nil, err
	}
	buckets := make([]cloud.BucketProperties, 0, len(output.Buckets))
	for _, bucket := range output.Buckets {
		buckets = append(buckets, cloud.BucketProperties{
			Name: aws.StringValue(bucket.Name),
		})
	}
	return buckets, nil
}

func (p *AWSCloud) GetOssDownloadUrl(endpoint, bucketName, region string) string {
	str := strings.Split(endpoint, "//")
	if len(str) != 2 {
		if strings.Contains(endpoint, "com.cn") {
			return fmt.Sprintf("https://%s.s3.%s.amazonaws.com.cn", bucketName, region)
		}
		return fmt.Sprintf("https://%s.s3.%s.amazonaws.com", bucketName, region)
	}
	return fmt.Sprintf("https://%s.%s", bucketName, str[1])
}

func (p *AWSCloud) GetObjectDownloadUrl(bucketName, objectKey string) (string, error) {
	svc := s3.New(p.sess)
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
