package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"path/filepath"
	"time"

	// . "github.com/getperf/gcagent/common"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type WebService struct {
	ServiceUrl string
	ArchiveDir string
	Cfg        *Config
}

func NewWebServiceBase(url string, archiveDir string, cfg *Config) *WebService {
	handle := WebService{
		ServiceUrl: url,
		ArchiveDir: archiveDir,
		Cfg:        cfg,
	}
	return &handle
}

func NewWebService(cfg *Config) (*WebService, error) {
	serviceUrl := cfg.Schedule.WebServiceUrl
	// tlsConfig := cfg.GetTSLConfig()
	webService := NewWebServiceBase(serviceUrl, cfg.BackupDir, cfg)
	_, err := url.Parse(serviceUrl)
	return webService, err
}

func (h *WebService) NewDatastore(node, job string) *Datastore {
	return NewDatastoreBase(h.ArchiveDir, node, job, time.Now())
}

func (h *WebService) ServHttp(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	u, err := url.Parse(h.ServiceUrl)
	if err != nil {
		return errors.Wrap(err, "exporter service, pasing url")
	}

	router := mux.NewRouter()
	router.HandleFunc("/store/{node}/{job}", h.DownloadDatastoreZip)
	router.HandleFunc("/store/{node}/{job}/{since}", h.DownloadDatastoreZip)
	router.HandleFunc("/store", h.GetLatestDatastoreSets)
	// router.HandleFunc("/upload/{file}", h.Upload)

	// Use the PathPrefix file share method on your router
	router.PathPrefix("/zip/").Handler(http.StripPrefix("/zip/",
		http.FileServer(http.Dir(h.ArchiveDir))))

	s := &http.Server{
		Addr:              u.Host,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}
	// if err := s.ListenAndServe(); err != nil {
	if err := Listen(s, h.Cfg); err != nil {
		return errors.Wrap(err, "exporter listen and serve")
	}
	return nil
}

func (config *Config) RunWebService(ctx context.Context) {
	ws, err := NewWebService(config)
	if err != nil {
		log.Error(errors.Wrap(err, "initialize http service"))
	}
	if err := ws.ServHttp(ctx); err != nil {
		log.Error(errors.Wrap(err, "start http service"))
	}
}

func (h *WebService) FindDatastoreZip(vars map[string]string) (string, error) {
	log.Info("download service : ", vars)
	node, ok := vars["node"]
	if !ok {
		return "", errors.New("required node not found")
	}
	job, ok := vars["job"]
	if !ok {
		return "", errors.New("required job not found")
	}
	ds := NewDatastoreBase(h.ArchiveDir, node, job, time.Now())
	log.Info("Datastore:", ds)
	sinceLabel := vars["since"]
	if sinceLabel == "" {
		return ds.LatestZip() // "arc_hoge__Windows.zip"
	} else {
		return ds.OldestZip(sinceLabel)
	}
}

func (h *WebService) DownloadDatastoreZipTemp(w http.ResponseWriter, r *http.Request) {
	// URL ディレクトリパスを検索条件にして zip ファイルを検索
	zipname, err := h.FindDatastoreZip(mux.Vars(r))
	if handleError(w, err, "download zip", 500) {
		return
	}
	zippath := filepath.Join(h.ArchiveDir, zipname)
	file, err := os.Open(zippath)
	if handleError(w, err, "download zip not found", 400) {
		return
	}
	defer file.Close()

	fieldName := "file"
	body := new(bytes.Buffer)
	multipartWriter := multipart.NewWriter(body)

	contentType := func() string {
		defer func() {
			file.Seek(0, 0)
		}()

		fileData, err := ioutil.ReadAll(file)
		if err != nil {
			return "application/octet-stream"
		}
		return http.DetectContentType(fileData)
	}()

	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`, fieldName, zipname))
	header.Set("Content-Type", contentType)
	part, _ := multipartWriter.CreatePart(header)

	_, _ = io.Copy(part, file)

	log.Info("download : ", zipname)
}

func (h *WebService) DownloadDatastoreZip(w http.ResponseWriter, r *http.Request) {
	// URL ディレクトリパスを検索条件にして zip ファイルを検索
	zipname, err := h.FindDatastoreZip(mux.Vars(r))
	if handleError(w, err, "download zip", 500) {
		return
	}
	zippath := filepath.Join(h.ArchiveDir, zipname)
	file, err := os.Open(zippath)
	if handleError(w, err, "download zip not found", 400) {
		return
	}
	defer file.Close()

	// リクエストボディのデータを受け取るio.Writerを生成する。
	body := &bytes.Buffer{}

	// データのmultipartエンコーディングを管理するmultipart.Writerを生成する。
	// ランダムなbase-16バウンダリが生成される。
	mw := multipart.NewWriter(body) // (body)

	// ファイルに使うパートを生成する。
	// ヘッダ以外はデータは書き込まれない。
	// fieldnameとzipnameの値がヘッダに含められる。
	// ファイルデータを書き込むio.Writerが返却される。
	fieldname := "file"
	fw, err := mw.CreateFormFile(fieldname, zipname)
	contentType := mw.FormDataContentType()
	log.Info("CONTENTTYPE:", contentType)
	w.Header().Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`, fieldname, zipname))
	w.Header().Set("Content-Type", contentType)

	// fwで作ったパートにファイルのデータを書き込む
	_, err = io.Copy(fw, file)
	if handleError(w, err, "write MIME multipart", 500) {
		return
	}

	// リクエストのContent-Typeヘッダに使う値を取得する（バウンダリを含む）
	err = mw.Close()
	if handleError(w, err, "close MIME multipart", 500) {
		return
	}
	// w.WriteHeader(200)
	w.Write(body.Bytes())
	log.Info("download : ", zipname)
}

func (h *WebService) GetLatestDatastoreSets(w http.ResponseWriter, r *http.Request) {
	// URL ディレクトリパスを検索条件にして zip ファイルを検索
	ds := NewDatastoreBase(h.ArchiveDir, "", "", time.Now())
	datastoreSets, err := ds.GetDatastoreSets()
	if handleError(w, err, "search datastore zip key", 500) {
		return
	}
	jsonBytes, err := json.Marshal(datastoreSets)
	if handleError(w, err, "search datastore zip key", 500) {
		return
	}
	w.Write(jsonBytes)
}

func (h *WebService) Upload(w http.ResponseWriter, r *http.Request) {
	log.Info("upload service : ", h.ArchiveDir)
	vars := mux.Vars(r)
	zipFile := vars["file"]

	fmt.Printf("request %v\n", r)
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// リクエストの情報を出力する
	// requestDump, err := httputil.DumpRequest(r, true)
	// if handleError(w, err, "upload handle", 500) {
	// 	return
	// }
	// log.Println(string(requestDump))

	// "file"というフィールド名に一致する最初のファイルが返却される
	// マルチパートフォームのデータはパースされていない場合ここでパースされる
	formFile, _, err := r.FormFile("file")
	if handleError(w, err, "upload handle", 500) {
		return
	}
	defer formFile.Close()

	// データを保存するファイルを開く
	zipPath := filepath.Join(h.ArchiveDir, zipFile)
	saveFile, err := os.Create(zipPath)
	if handleError(w, err, "upload handle", 500) {
		return
	}
	defer saveFile.Close()

	// ファイルにデータを書き込む
	_, err = io.Copy(saveFile, formFile)
	if handleError(w, err, "upload handle", 500) {
		return
	}

	w.WriteHeader(http.StatusCreated)
	log.Info("upload : ", zipFile)
}

func handleError(w http.ResponseWriter, err error, msg string, code int) bool {
	if err != nil {
		err = errors.Wrap(err, msg)
		log.Error(err)
		http.Error(w, err.Error(), code)
		return true
	} else {
		return false
	}
}
