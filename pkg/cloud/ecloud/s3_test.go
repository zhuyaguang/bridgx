package ecloud

import "testing"

func TestListObjects(t *testing.T) {
	client, err := New(_AK, _SK, _regionId)
	if err != nil {
		t.Log(err.Error())
		return
	}
	res, err := client.ListObjects("", "", "")
	for _, b := range res {
		t.Log(b.Name)
	}
}

func TestListBucket(t *testing.T) {
	client, err := New(_AK, _SK, _regionId)
	if err != nil {
		t.Log(err.Error())
		return
	}
	res, err := client.ListBucket("")
	for _, b := range res {
		t.Log(b.Name)
	}
}

func TestGetObjectDownloadUrl(t *testing.T) {
	client, err := New(_AK, _SK, _regionId)
	if err != nil {
		t.Log(err.Error())
		return
	}
	res, err := client.GetObjectDownloadUrl("", "")
	t.Log(res)
}
