package tencent

import "testing"

func TestContainerInstanceList(t *testing.T) {
	client, err := New(_ak, _sk, "ap-beijing")
	if err != nil {
		t.Log(err.Error())
		return
	}
	res, count, err := client.ContainerInstanceList("", 0, 0)
	t.Log(count)
	for _, b := range res {
		t.Log(b.InstanceId, b.InstanceName)
	}
}

func TestEnterpriseNamespaceList(t *testing.T) {
	client, err := New(_ak, _sk, "ap-beijing")
	if err != nil {
		t.Log(err.Error())
		return
	}
	res, count, err := client.EnterpriseNamespaceList("", "test", 0, 0)
	t.Log(count)
	for _, b := range res {
		t.Log(b.Name)
	}
}

func TestPersonalNamespaceList(t *testing.T) {
	client, err := New(_ak, _sk, "ap-beijing")
	if err != nil {
		t.Log(err.Error())
		return
	}
	res, err := client.PersonalNamespaceList("")
	for _, b := range res {
		t.Log(b.Name)
	}
}

func TestEnterpriseRepositoryList(t *testing.T) {
	client, err := New(_ak, _sk, "ap-beijing")
	if err != nil {
		t.Log(err.Error())
		return
	}
	res, count, err := client.EnterpriseRepositoryList("", "test", "", 0, 0)
	t.Log(count)
	for _, b := range res {
		t.Log(b.Name, b.ID)
	}
}

func TestPersonalRepositoryList(t *testing.T) {
	client, err := New(_ak, _sk, "ap-beijing")
	if err != nil {
		t.Log(err.Error())
		return
	}
	res, count, err := client.PersonalRepositoryList("", "", 0, 0)
	t.Log(count)
	for _, b := range res {
		t.Log(b.Name, b.ID)
	}
}

func TestEnterpriseImageList(t *testing.T) {
	client, err := New(_ak, _sk, "ap-beijing")
	if err != nil {
		t.Log(err.Error())
		return
	}
	res, count, err := client.EnterpriseImageList("", "", "test", "test", "test", 0, 0)
	t.Log(count)
	for _, b := range res {
		t.Log(b.Name)
	}
}

func TestPersonalImageList(t *testing.T) {
	client, err := New(_ak, _sk, "ap-beijing")
	if err != nil {
		t.Log(err.Error())
		return
	}
	res, count, err := client.PersonalImageList("", "", "test", 0, 0)
	t.Log(count)
	for _, b := range res {
		t.Log(b.Name)
	}
}
