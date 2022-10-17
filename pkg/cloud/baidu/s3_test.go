package baidu

import (
	"fmt"
	"testing"
)

var s3Account = account{
	ak:       "",
	sk:       "",
	regionID: "",
}

func TestListBucket(t *testing.T) {
	b, err := New(s3Account.ak, s3Account.sk, s3Account.regionID)
	if err != nil {
		fmt.Println(err)
	} else {
		res, err := b.ListBucket("bd")
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(res)
		}

	}
}

func TestListObjects(t *testing.T) {
	b, err := New(s3Account.ak, s3Account.sk, s3Account.regionID)
	if err != nil {
		fmt.Println(err)
	} else {
		res, err := b.ListObjects("bj", "test1111xiewei", "")
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(res)
		}

	}
}

func TestGetObjectDownloadUrl(t *testing.T) {
	b, err := New(s3Account.ak, s3Account.sk, s3Account.regionID)
	if err != nil {
		fmt.Println(err)
	} else {
		res := b.GetOssDownloadUrl("test1111xiewei", "", "")
		fmt.Println(res)

	}
}
