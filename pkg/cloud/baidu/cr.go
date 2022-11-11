package baidu

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/galaxy-future/BridgX/pkg/cloud"
	"github.com/pkg/errors"
)

func (p *BaiduCloud) ContainerInstanceList(region string, pageNumber, pageSize int) ([]cloud.RegistryInstance, int, error) {
	ri := []cloud.RegistryInstance{}
	//use utc time

	timeStamp := time.Now().UTC().Format("2006-01-02T15:04:05Z")
	authStringPrefix := fmt.Sprintf("bce-auth-v1/%s/%s/10000", p.ak, timeStamp)
	h := hmac.New(sha256.New, []byte(p.sk))
	h.Write([]byte(authStringPrefix))
	signKey := hex.EncodeToString(h.Sum(nil))

	uri := fmt.Sprintf("/v1/instances\npageNo=%d&pageSize=%d", pageNumber, pageSize)

	canonicalRequest := "GET\n" + uri + "\nhost:ccr.bd.baidubce.com"

	hs := hmac.New(sha256.New, []byte(signKey))
	hs.Write([]byte(canonicalRequest))
	signature := hex.EncodeToString(hs.Sum(nil))
	authorization := authStringPrefix + "/host/" + signature
	ser := &http.Client{}
	requrl := fmt.Sprintf("http://ccr.bd.baidubce.com/v1/instances?pageNo=%d&pageSize=%d", pageNumber, pageSize)

	req, _ := http.NewRequest("GET", requrl, nil)

	req.Header.Set("Host", "ccr.bd.baidubce.com")
	req.Header.Set("Authorization", authorization)
	resp, err := ser.Do(req)
	if err != nil {
		return ri, 0, err
	}

	m := make(map[string]interface{})
	data, err := ioutil.ReadAll(resp.Body)

	_ = json.Unmarshal(data, &m)
	mm := m["instances"].([]interface{})
	for _, v := range mm {
		n := v.(map[string]interface{})
		repo := cloud.RegistryInstance{
			InstanceId:   n["id"].(string),
			InstanceName: n["name"].(string),
		}
		ri = append(ri, repo)
	}
	return ri, 0, nil
}

func (p *BaiduCloud) EnterpriseNamespaceList(region, instanceId string, pageNumber, pageSize int) ([]cloud.Namespace, int, error) {
	namespaces := []cloud.Namespace{}
	//use utc time

	timeStamp := time.Now().UTC().Format("2006-01-02T15:04:05Z")
	authStringPrefix := fmt.Sprintf("bce-auth-v1/%s/%s/10000", p.ak, timeStamp)
	h := hmac.New(sha256.New, []byte(p.sk))
	h.Write([]byte(authStringPrefix))
	signKey := hex.EncodeToString(h.Sum(nil))

	uri := fmt.Sprintf("/v1/instances/%s/projects\npageNo=%d&pageSize=%d", instanceId, pageNumber, pageSize)

	canonicalRequest := "GET\n" + uri + "\nhost:ccr.bd.baidubce.com"
	hs := hmac.New(sha256.New, []byte(signKey))
	hs.Write([]byte(canonicalRequest))
	signature := hex.EncodeToString(hs.Sum(nil))
	authorization := authStringPrefix + "/host/" + signature
	ser := &http.Client{}

	req, _ := http.NewRequest("GET", fmt.Sprintf("http://ccr.bd.baidubce.com/v1/instances/%s/projects?pageNo=%d&pageSize=%d", instanceId, pageNumber, pageSize), nil)

	req.Header.Set("Host", "ccr.bd.baidubce.com")
	req.Header.Set("Authorization", authorization)
	resp, err := ser.Do(req)
	if err != nil {
		return namespaces, 0, err
	}

	m := make(map[string]interface{})
	data, err := ioutil.ReadAll(resp.Body)
	_ = json.Unmarshal(data, &m)
	mm := m["projects"].([]interface{})
	for _, v := range mm {
		n := v.(map[string]interface{})
		namespace := cloud.Namespace{
			Name: n["projectName"].(string),
		}
		namespaces = append(namespaces, namespace)
	}

	return namespaces, 0, nil
}

func (p *BaiduCloud) PersonalNamespaceList(region string) ([]cloud.Namespace, error) {
	namespacs := []cloud.Namespace{}
	//use utc time
	timeStamp := time.Now().UTC().Format("2006-01-02T15:04:05Z")
	authStringPrefix := fmt.Sprintf("bce-auth-v1/%s/%s/10000", p.ak, timeStamp)
	h := hmac.New(sha256.New, []byte(p.sk))
	h.Write([]byte(authStringPrefix))
	signKey := hex.EncodeToString(h.Sum(nil))

	canonicalRequest := "GET\n/v1/ccr/projects\n\nhost:ccr.baidubce.com\nx-bce-date:" + url.QueryEscape(timeStamp)
	hs := hmac.New(sha256.New, []byte(signKey))
	hs.Write([]byte(canonicalRequest))
	signature := hex.EncodeToString(hs.Sum(nil))

	authorization := authStringPrefix + "/host;x-bce-date/" + signature
	ser := &http.Client{}

	req, _ := http.NewRequest("GET", "http://ccr.baidubce.com/v1/ccr/projects", nil)
	req.Header.Set("Host", "ccr.baidubce.com")
	req.Header.Set("Authorization", authorization)
	req.Header.Set("x-bce-date", timeStamp)
	resp, err := ser.Do(req)
	if err != nil {
		return nil, err
	}
	m := make(map[string]interface{})
	data, err := ioutil.ReadAll(resp.Body)
	_ = json.Unmarshal(data, &m)

	mm := m["result"].([]interface{})
	for _, v := range mm {
		n := v.(map[string]interface{})
		//fmt.Println(n["projectId"].(float64))
		namespace := cloud.Namespace{
			Name: n["projectName"].(string),
		}
		namespacs = append(namespacs, namespace)
	}

	return namespacs, nil
}

func (p *BaiduCloud) EnterpriseRepositoryList(region, instanceId, namespace string, pageNumber, pageSize int) ([]cloud.Repository, int, error) {
	repos := []cloud.Repository{}
	//use utc time

	timeStamp := time.Now().UTC().Format("2006-01-02T15:04:05Z")
	authStringPrefix := fmt.Sprintf("bce-auth-v1/%s/%s/10000", p.ak, timeStamp)
	h := hmac.New(sha256.New, []byte(p.sk))
	h.Write([]byte(authStringPrefix))
	signKey := hex.EncodeToString(h.Sum(nil))

	uri := fmt.Sprintf("/v1/instances/%s/projects/%s/repositories\npageNo=%d&pageSize=%d", instanceId, namespace, pageNumber, pageSize)

	canonicalRequest := "GET\n" + uri + "\nhost:ccr.bd.baidubce.com"

	hs := hmac.New(sha256.New, []byte(signKey))
	hs.Write([]byte(canonicalRequest))
	signature := hex.EncodeToString(hs.Sum(nil))
	authorization := authStringPrefix + "/host/" + signature
	ser := &http.Client{}
	requrl := fmt.Sprintf("http://ccr.bd.baidubce.com/v1/instances/%s/projects/%s/repositories?pageNo=%d&pageSize=%d", instanceId, namespace, pageNumber, pageSize)

	req, _ := http.NewRequest("GET", requrl, nil)

	req.Header.Set("Host", "ccr.bd.baidubce.com")
	req.Header.Set("Authorization", authorization)
	resp, err := ser.Do(req)
	if err != nil {
		return repos, 0, err
	}

	m := make(map[string]interface{})
	data, err := ioutil.ReadAll(resp.Body)

	_ = json.Unmarshal(data, &m)
	mm := m["items"].([]interface{})
	for _, v := range mm {
		n := v.(map[string]interface{})
		repo := cloud.Repository{
			Name: n["repositoryName"].(string),
			ID:   n["repositoryPath"].(string),
		}
		repos = append(repos, repo)
	}

	return repos, 0, nil
}

func (p *BaiduCloud) PersonalRepositoryList(region, namespace string, pageNumber, pageSize int) ([]cloud.Repository, int, error) {
	repos := []cloud.Repository{}
	//use utc time
	timeStamp := time.Now().UTC().Format("2006-01-02T15:04:05Z")
	authStringPrefix := fmt.Sprintf("bce-auth-v1/%s/%s/10000", p.ak, timeStamp)
	h := hmac.New(sha256.New, []byte(p.sk))
	h.Write([]byte(authStringPrefix))
	signKey := hex.EncodeToString(h.Sum(nil))

	uri := fmt.Sprintf("/v1/ccr/repositories/user\npageNo=%d&pageSize=%d", pageNumber, pageSize)

	canonicalRequest := "GET\n" + uri + "\nhost:ccr.baidubce.com\nx-bce-date:" + url.QueryEscape(timeStamp)
	hs := hmac.New(sha256.New, []byte(signKey))
	hs.Write([]byte(canonicalRequest))
	signature := hex.EncodeToString(hs.Sum(nil))
	authorization := authStringPrefix + "/host;x-bce-date/" + signature
	ser := &http.Client{}
	requrl := fmt.Sprintf("http://ccr.baidubce.com/v1/ccr/repositories/user?pageNo=%d&pageSize=%d", pageNumber, pageSize)

	req, _ := http.NewRequest("GET", requrl, nil)

	req.Header.Set("Host", "ccr.baidubce.com")
	req.Header.Set("Authorization", authorization)
	req.Header.Set("x-bce-date", timeStamp)
	resp, err := ser.Do(req)
	if err != nil {
		return repos, 0, err
	}

	if resp.StatusCode == 404 {
		return nil, 0, errors.New("404 no found")
	}
	m := make(map[string]interface{})
	data, err := ioutil.ReadAll(resp.Body)

	_ = json.Unmarshal(data, &m)
	mm := m["result"].([]interface{})
	for _, v := range mm {
		n := v.(map[string]interface{})
		id := strconv.FormatFloat(n["repositoryId"].(float64), 'f', 0, 64)
		repo := cloud.Repository{
			Name: n["repositoryName"].(string),
			ID:   id,
		}
		repos = append(repos, repo)
	}

	return repos, 0, nil
}

func (b BaiduCloud) PersonalImageList(region, projectID, repoName string, pageNum, pageSize int) ([]cloud.DockerArtifact, int, error) {
	imagesVersion := []cloud.DockerArtifact{}
	//use utc time
	timeStamp := time.Now().UTC().Format("2006-01-02T15:04:05Z")
	authStringPrefix := fmt.Sprintf("bce-auth-v1/%s/%s/10000", b.ak, timeStamp)
	h := hmac.New(sha256.New, []byte(b.sk))
	h.Write([]byte(authStringPrefix))
	signKey := hex.EncodeToString(h.Sum(nil))
	uri := fmt.Sprintf("/v1/ccr/repositories/tags\nprojectId=%s&repoName=%s", projectID, repoName)

	canonicalRequest := "GET\n" + uri + "\nhost:ccr.baidubce.com\nx-bce-date:" + url.QueryEscape(timeStamp)
	hs := hmac.New(sha256.New, []byte(signKey))
	hs.Write([]byte(canonicalRequest))
	signature := hex.EncodeToString(hs.Sum(nil))
	authorization := authStringPrefix + "/host;x-bce-date/" + signature
	ser := &http.Client{}
	requrl := fmt.Sprintf("http://ccr.baidubce.com/v1/ccr/repositories/tags?projectId=%s&repoName=%s", projectID, repoName)

	req, _ := http.NewRequest("GET", requrl, nil)
	req.Header.Set("Host", "ccr.baidubce.com")
	req.Header.Set("Authorization", authorization)
	req.Header.Set("x-bce-date", timeStamp)
	resp, err := ser.Do(req)
	if err != nil {
		return imagesVersion, 0, err
	}
	if resp.StatusCode == 404 {
		return nil, 0, errors.New("404 no found")
	}
	m := make(map[string]interface{})
	data, err := ioutil.ReadAll(resp.Body)

	_ = json.Unmarshal(data, &m)
	mm := m["result"].([]interface{})
	for _, v := range mm {
		n := v.(map[string]interface{})
		version := cloud.DockerArtifact{
			Name: n["name"].(string),
		}
		imagesVersion = append(imagesVersion, version)
	}

	return imagesVersion, 0, nil
}

// list images version
func (p *BaiduCloud) EnterpriseImageList(region, instanceId, repoId, namespace, repoName string, pageNumber, pageSize int) ([]cloud.DockerArtifact, int, error) {
	imagesVersion := []cloud.DockerArtifact{}
	//use utc time

	timeStamp := time.Now().UTC().Format("2006-01-02T15:04:05Z")
	authStringPrefix := fmt.Sprintf("bce-auth-v1/%s/%s/10000", p.ak, timeStamp)
	h := hmac.New(sha256.New, []byte(p.sk))
	h.Write([]byte(authStringPrefix))
	signKey := hex.EncodeToString(h.Sum(nil))
	uri := fmt.Sprintf("/v1/instances/%s/projects/%s/repositories/%s/tags", instanceId, namespace, repoName)

	canonicalRequest := "GET\n" + uri + "\n\nhost:ccr.bd.baidubce.com"
	hs := hmac.New(sha256.New, []byte(signKey))
	hs.Write([]byte(canonicalRequest))
	signature := hex.EncodeToString(hs.Sum(nil))
	authorization := authStringPrefix + "/host/" + signature
	ser := &http.Client{}

	req, _ := http.NewRequest("GET", fmt.Sprintf("http://ccr.bd.baidubce.com/v1/instances/%s/projects/%s/repositories/%s/tags", instanceId, namespace, repoName), nil)

	req.Header.Set("Host", "ccr.bd.baidubce.com")
	req.Header.Set("Authorization", authorization)
	resp, err := ser.Do(req)
	if err != nil {
		return imagesVersion, 0, err
	}
	//fmt.Println(resp.StatusCode, resp.Header)
	m := make(map[string]interface{})
	data, err := ioutil.ReadAll(resp.Body)
	_ = json.Unmarshal(data, &m)
	//fmt.Println(string(data))
	mm := m["items"].([]interface{})
	for _, v := range mm {
		n := v.(map[string]interface{})
		tag := cloud.DockerArtifact{
			Name: n["tagName"].(string),
		}
		imagesVersion = append(imagesVersion, tag)
	}

	return imagesVersion, 0, nil
}
