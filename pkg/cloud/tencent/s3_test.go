package tencent

import "testing"

func TestListBucket(t *testing.T) {
	client, err := New("ak", "sk", "ap-shanghai")
	if err != nil {
		t.Log(err.Error())
		return
	}
	res, err := client.ListBucket(getCosEndpoint("zhuyaguang-1308110266", "ap-shanghai"))
	for _, b := range res {
		t.Log(b.Name)
	}

}

func TestListListObjects(t *testing.T) {
	client, err := New("ak", "sk", "ap-shanghai")
	if err != nil {
		t.Log(err.Error())
		return
	}
	res, err := client.ListObjects(getCosEndpoint("zhuyaguang-1308110266", "ap-shanghai"), "", "img/")
	if err != nil {
		t.Log(err.Error())
		return
	}
	for _, c := range res {
		t.Logf("%s\n", c.Name)
	}
}

func TestGetObjectDownloadUrl(t *testing.T) {
	client, err := New("ak", "sk", "ap-shanghai")
	if err != nil {
		t.Log(err.Error())
		return
	}
	res, _ := client.GetObjectDownloadUrl(getCosEndpoint("zhuyaguang-1308110266", "ap-shanghai"), "img/20220916184018.png")
	t.Log(res)
}
