package agent

import (
	// "io/ioutil"
	// "os"
	// "path/filepath"
	// "testing"

	// . "github.com/getperf/gcagent/common"
	// "bytes"
	// "io"
	// "mime/multipart"
	// "net/http"
	// "net/http/httptest"
	// "net/url"
	// "os"
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/gorilla/mux"
)

var cfg = Config{}
var ws = NewWebServiceBase("http://0.0.0.0:59000",
	"../testdata/webservice/backup", &cfg)

func TestExporterUrl(t *testing.T) {
	// now := time.Date(2018, 1, 2, 3, 4, 5, 0, time.UTC)
	ds := ws.NewDatastore("windowsconf", "hoge")
	t.Log(ds.LatestZip())
	t.Log(ws.ServiceUrl)
	t.Log(ws.ArchiveDir)

	u, err := url.Parse(ws.ServiceUrl)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("Scheme:", u.Scheme)       // http
		t.Log("Host:", u.Host)           // Host: bing.com:8080
		t.Log("Hostname:", u.Hostname()) // bing.com
		t.Log("Path:", u.Path)           // search
	}
}

func TestFindInventoryZip(t *testing.T) {
	vars := map[string]string{
		"node": "hogehoge",
		"job":  "Windows",
	}
	zip, err := ws.FindDatastoreZip(vars)
	if err != nil {
		t.Error("find inventory zip ", err)
	}
	if zip != "arc_hogehoge__Windows_20200520_1000.zip" {
		t.Error("find inventory zip ")
	}
}

func TestFindPerfMetricZip(t *testing.T) {
	var tests = []struct {
		node   string
		job    string
		since  string
		result string
	}{
		{"hogehoge", "Windows", "20200519_1000", "arc_hogehoge__Windows_20200519_1120.zip"},
		{"fuga", "Windows", "20200519_1000", ""},
		{"hogehoge", "HogeWindows", "20200519_1000", ""},
	}
	for _, test := range tests {
		vars := map[string]string{
			"node":  test.node,
			"job":   test.job,
			"since": test.since,
		}
		zip, _ := ws.FindDatastoreZip(vars)
		if zip != test.result {
			t.Errorf("FindDatastoreZip : %v : %s", test, zip)
		}
	}
}

func TestDownloadHandler(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/store/{node}/{job}", ws.DownloadDatastoreZip)

	req := httptest.NewRequest("POST", "/store/hogehoge/Windows", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	// コンテンツヘッダーの出力
	// Content-Disposition: form-data; name="file"; filename="arc_hoge__Windows.zip"
	if status := rec.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestGetLatestSets(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/store", ws.GetLatestDatastoreSets)

	req := httptest.NewRequest("GET", "/store", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	t.Log(rec.Body.String())
	// if rec.Body.String() != `[{"host":"hogehoge","stat_name":"Windows"}]` {
	// 	t.Error("unmatch")
	// }
}

func TestDownloadNgHandler(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/store/{node}/{job}", ws.DownloadDatastoreZip)

	req := httptest.NewRequest("GET", "/store/hoge/Windows", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if status := rec.Code; status == http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	t.Log(rec.Body.String())
}

func TestUploadHandler(t *testing.T) {
	fieldname := "file"
	filename := "../testdata/arc_hogehoge__Windows.zip"
	file, err := os.Open(filename)
	if err != nil {
		t.Error("test upload file open : ", err)
	}
	router := mux.NewRouter()
	router.HandleFunc("/upload/{file}", ws.Upload)

	body := &bytes.Buffer{}
	mw := multipart.NewWriter(body)
	fw, err := mw.CreateFormFile(fieldname, filename)
	_, err = io.Copy(fw, file)
	if err != nil {
		t.Error("test upload copy file to body : ", err)
	}
	contentType := mw.FormDataContentType()
	err = mw.Close()
	if err != nil {
		t.Error("test upload close body : ", err)
	}

	req := httptest.NewRequest("POST", "/upload/hogehoge.zip", body)
	req.Header.Set("Content-Type", contentType)

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if status := rec.Code; status == http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
