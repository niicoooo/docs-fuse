package docsConn

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
)

/*
func Dump(resp *http.Response) {
	dump1, err := httputil.DumpRequest(resp.Request, true)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("dump request: %q\n", dump1)
	dump2, err := httputil.DumpResponse(resp, true)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("dump response: %q\n", dump2)
	return
}
*/

type DocsConn struct {
	username         string
	password         string
	cookie           *http.Cookie
	authTokenVersion int
	apiUrl           string
	mux              sync.Mutex
}

func (this *DocsConn) connect(version int) (*http.Cookie, int, *DocsError) {
	this.mux.Lock()
	defer this.mux.Unlock()

	if version < this.authTokenVersion {
		return this.cookie, this.authTokenVersion, nil
	}

	data := url.Values{}
	data.Add("username", this.username)
	data.Add("password", this.password)
	u, _ := url.ParseRequestURI(this.apiUrl)
	u.Path = "/api/user/login"
	urlStr := u.String()

	client := &http.Client{}

	r, _ := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	resp, err := client.Do(r)
	if err != nil {
		return nil, version, newDocsError(0, "%s", error.Error(err))
	}
	if resp.StatusCode != 200 {
		resp, err := client.Do(r)
		if err != nil {
			return nil, version, newDocsError(0, "%s", error.Error(err))
		}
		if resp.StatusCode != 200 {
			return nil, version, newDocsError(resp.StatusCode, "%s", resp.Status)
		}
	}
	if len(resp.Cookies()) != 1 {
		return nil, version, newDocsError(0, "unexpected cookie structure (1)")
	}
	if strings.Compare(resp.Cookies()[0].Name, "auth_token") != 0 {
		return nil, version, newDocsError(0, "unexpected cookie structure (2)")
	}

	this.cookie = resp.Cookies()[0]
	this.authTokenVersion++

	return this.cookie, this.authTokenVersion, nil
}

func NewDocsConn(apiUrl string, username string, password string) (*DocsConn, error) {
	conn := DocsConn{
		username:         username,
		password:         password,
		authTokenVersion: 0,
		apiUrl:           apiUrl,
	}
	_, _, err := conn.connect(0)
	if err != nil {
		return nil, fmt.Errorf("Unable to connect: %s", err.s)
	}
	return &conn, nil
}

func (this *DocsConn) get(ressource string, query string) ([]byte, *DocsError) {
	u, _ := url.ParseRequestURI(this.apiUrl)
	u.Path = ressource
	u.RawQuery = query
	urlStr := u.String()
	client := &http.Client{}
	r, _ := http.NewRequest("GET", urlStr, nil)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	cookie, version, err := this.connect(0)
	if err != nil {
		return nil, newDocsError(0, "Unable to reconnect: %s (%s)", err.s, urlStr)
	}
	r.AddCookie(cookie)
	resp, _ := client.Do(r)

	if resp.StatusCode != 200 {
		cookie, _, err = this.connect(version)
		if err != nil {
			return nil, newDocsError(0, "Unable to reconnect(2): %s (%s)", err.s, urlStr)
		}
		r.AddCookie(cookie)
		resp, _ = client.Do(r)
		if resp.StatusCode != 200 {
			return nil, newDocsError(resp.StatusCode, "%s (%s)", resp.Status, urlStr)
		}
	}
	body, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		log.Fatal("ioutil.ReadAll error: %s", err)
	}
	return body, nil
}

type DocumentList struct {
	Documents []struct {
		Id      string
		File_id string
		Title   string
	}
}

func (this *DocsConn) GetDocumentList() (*DocumentList, error) {
	body, err := this.get("api/document/list", "")
	if err != nil {
		return nil, fmt.Errorf("GetDocumentList get error: %s", err.s)
	}
	var result DocumentList
	err2 := json.Unmarshal(body, &result)
	if err2 != nil {
		return nil, fmt.Errorf("GetDocumentList unmarshal error: %s", error.Error(err2))
	}
	return &result, nil
}

func (this *DocsConn) GetDocument(id string) ([]byte, error) {
	body, err := this.get("api/document/"+id, "")
	if err != nil {
		return nil, fmt.Errorf("GetDocument get error: %s", err.s)
	}
	return body, nil
}

type FileList struct {
	Files []struct {
		Id          string
		Name        string
		Processing  bool
		Version     int
		Mimetype    string
		Document_id string
		Create_date int
		Size        int
	}
}

func (this *DocsConn) GetFileList(docId string) (*FileList, []byte, error) {
	body, err := this.get("api/file/list", "id="+docId)
	if err != nil {
		return nil, nil, fmt.Errorf("GetFileList get error: %s", err.s)
	}
	var result FileList
	err2 := json.Unmarshal(body, &result)
	if err2 != nil {
		return nil, nil, fmt.Errorf("GetFileList unmarshal error: %s", error.Error(err2))
	}
	return &result, body, nil
}

func (this *DocsConn) GetFileData(id string) ([]byte, error) {
	body, err := this.get("/api/file/"+id+"/data", "")
	if err != nil {
		return nil, fmt.Errorf("GetFileData get error: %s", err.s)
	}
	return body, nil
}

type TagList struct {
	Tags []struct {
		Id     string
		Name   string
		Color  string
		Parent string
	} `json:"tags"`
}

func (this *DocsConn) GetTags() (*TagList, error) {
	body, err := this.get("/api/tag/list", "")
	if err != nil {
		return nil, fmt.Errorf("GetTags get error: %s", err.s)
	}
	var result TagList
	err2 := json.Unmarshal(body, &result)
	if err2 != nil {
		return nil, fmt.Errorf("GetTags unmarshal error: %s", error.Error(err2))
	}
	return &result, nil
}

func (this *DocsConn) GetDocumentListByTag(tagName string) (*DocumentList, error) {
	body, err := this.get("api/document/list", "search=tag:"+tagName)
	if err != nil {
		return nil, fmt.Errorf("GetDocumentList get error: %s", err.s)
	}
	var result DocumentList
	err2 := json.Unmarshal(body, &result)
	if err2 != nil {
		return nil, fmt.Errorf("GetDocumentList unmarshal error: %s", error.Error(err2))
	}
	return &result, nil
}
