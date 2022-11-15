package tencent

import "testing"

const (
	_ak = ""
	_sk = ""
)

func TestListBucket(t *testing.T) {
	client, err := New(_ak, _sk, "ap-beijing")
	if err != nil {
		t.Log(err.Error())
		return
	}
	res, err := client.ListBucket(getCosEndpoint("dt-demo-1308988865", "ap-beijing"))
	for _, b := range res {
		t.Log(b.Name)
	}

}

func TestListListObjects(t *testing.T) {
	client, err := New(_ak, _sk, "ap-beijing")
	if err != nil {
		t.Log(err.Error())
		return
	}
	res, err := client.ListObjects(getCosEndpoint("dt-demo-1308988865", "ap-beijing"), "", "luhuajun/")
	if err != nil {
		t.Log(err.Error())
		return
	}
	for _, c := range res {
		t.Logf("%s\n", c.Name)
	}
}

func TestGetObjectDownloadUrl(t *testing.T) {
	client, err := New(_ak, _sk, "ap-beijing")
	if err != nil {
		t.Log(err.Error())
		return
	}
	res, _ := client.GetObjectDownloadUrl(getCosEndpoint("dt-demo-1308988865", "ap-beijing"), "luhuajun/123.txt")
	t.Log(res)
}
