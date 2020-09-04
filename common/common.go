package common

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type DateFormat int

const (
	DEFAULT DateFormat = iota
	YYYYMMDD
	HHMISS
	YYYYMMDD_HHMISS
	DIR
)

// type JobStatus int

// const (
// 	JOB_INIT JobStatus = iota
// 	JOB_SUCCESS
// 	JOB_WARN
// 	JOB_ERROR
// )

type DiskStatus struct {
	All  uint64 `json:"all"`
	Used uint64 `json:"used"`
	Free uint64 `json:"free"`
}

// GetHostname はホスト名のエイリアスを取得します。
func GetHostname() (string, error) {
	hostName, err := os.Hostname()
	if err != nil {
		return hostName, fmt.Errorf("get hostname %s", err)
	}
	return GetHostnameAlias(hostName)
}

// GetHostnameAlias はホスト名の「.」以降の文字列はカットし、大文字は小文字に変換します。
func GetHostnameAlias(hostName string) (string, error) {
	if i := strings.Index(hostName, "."); i > 0 {
		hostName = hostName[:i]
	}
	hostName = strings.ToLower(hostName)
	return hostName, nil
}

// GetBaseDirは実行パスのベースディレクトリを返します。
func GetBaseDir() string {
	baseDir := "."
	// /tmp/go-build...ではないコンパイル済みバイナリからの実行かチェック
	exe, err := os.Executable()
	if err == nil && strings.Index(exe, "go-build") == -1 {
		baseDir = filepath.Dir(exe)
	} else {
		log.Warn("failed to get program name")
		exe = ""
	}
	return baseDir
}

// GetParentPathAbs は実行パスから上位のディレクトリを絶対パスに変換して返します。
func GetParentPath(inPath string, parentLevel int) string {
	for parentLevel > 0 {
		parentLevel--
		inPath = filepath.Dir(inPath)
	}
	return inPath
}

// GetParentPathAbs は実行パスから上位のディレクトリを絶対パスに変換して返します。
func GetParentAbsPath(inPath string, parentLevel int) (string, error) {
	inPath, err := filepath.Abs(inPath)
	if err != nil {
		return inPath, fmt.Errorf("get parent absolute path %s : %s", inPath, err)
	}
	return GetParentPath(inPath, parentLevel), nil
}

// CheckDirectory はディレクトリの存在確認をします。
func CheckDirectory(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err == nil {
		return fi.Mode().IsDir(), nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// CheckFile はファイルの存在確認をします。
func CheckFile(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, fmt.Errorf("not found %s", path)
		} else {
			return false, errors.Wrap(err, "check file stat")
		}
	} else {
		if fi.Mode().IsDir() {
			return false, fmt.Errorf("check file, %s is directory", path)
		}
	}
	return true, nil
}

func CheckExists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

// CreateAndOpenFile ファイルパスから書き込み用ファイルをオープンします。
// ディレクトリが存在しない場合は事前にディレクトリを作成します。
func CreateAndOpenFile(filePath string) (*os.File, error) {
	fileDir := filepath.Dir(filePath)
	if err := os.MkdirAll(fileDir, 0777); err != nil {
		return nil, errors.Wrap(err, "create node directory")
	}
	return os.OpenFile(filePath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
}

// HandleError はエラーがある場合に、メッセージを付加したエラーをログ出力し、
// エラーとして返します
func HandleError(w io.Writer, inErr error, message string) error {
	if inErr != nil {
		msg := fmt.Sprintf("%s : %s\n", message, inErr)
		// _, err := fmt.Fprintf(w, "%s : %s\n", message, inErr)
		// fmt.Fprint(os.Stderr, msg)
		_, err := fmt.Fprint(w, msg)
		if err != nil {
			log.Errorf("write log error : %s", err)
		}
		return errors.Wrap(inErr, message)
	} else {
		return nil
	}
}

// HandleError はエラーがある場合に、メッセージを付加したエラーをログ出力し、
// エラーとして返します
func HandleErrorWithAlert(w io.Writer, inErr error, message string) error {
	err := HandleError(w, inErr, message)
	if err != nil {
		log.Error(err)
	}
	return err
}

// CheckDirectoryIsNull は中身が空のディレクトリかチェックします
func CheckDirectoryIsNull(filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return errors.Wrap(err, "check directory is null")
	}
	files, err := filepath.Glob(filePath + "/*")
	if err != nil {
		return errors.Wrap(err, "check directory is null")
	}
	if len(files) != 0 {
		return fmt.Errorf("not null directory : %s", filePath)
	}
	return nil
}

// RemoveAndCreateDir はディレクトリを再作成します。
func RemoveAndCreateDir(filePath string) error {
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		if err := os.RemoveAll(filePath); err != nil {
			return errors.Wrap(err, "initialize directory")
		}
	}
	if err := os.MkdirAll(filePath, 0777); err != nil {
		return errors.Wrap(err, "initialize directory")
	}
	return nil
}

// CopyFile はファイルのコピーをします。
func CopyFile(srcPath, targetPath string) error {
	src, err := os.Open(srcPath)
	defer src.Close()
	if err != nil {
		return fmt.Errorf("copy from source %s : %s", srcPath, err)
	}

	dst, err := os.Create(targetPath)
	defer dst.Close()
	if err != nil {
		return fmt.Errorf("copy to target %s : %s", targetPath, err)
	}

	_, err = io.Copy(dst, src)
	if err != nil {
		return fmt.Errorf("copy file : %s", err)
	}
	return nil
}

// GetCurrentTime は指定したフォーマット形式で何秒前の現在時刻を取得します。
func GetCurrentTime(sec int, dateFormat DateFormat) string {
	now := time.Now().Add(-1 * time.Second * time.Duration(sec))
	return GetTimeString(dateFormat, now)
}

// gpfDGetTimeString は指定したフォーマット形式で時刻を変換します。
//    GPF_DATE_FORMAT_DEFAULT         0
//    GPF_DATE_FORMAT_YYYYMMDD        1
//    GPF_DATE_FORMAT_HHMISS          2
//    GPF_DATE_FORMAT_YYYYMMDD_HHMISS 3
//    GPF_DATE_FORMAT_DIR             4
func GetTimeString(dateFormat DateFormat, t time.Time) string {
	var format string
	switch dateFormat {
	case DEFAULT:
		format = t.Format("2006/01/02 15:04:05")
	case YYYYMMDD:
		format = t.Format("20060102")
	case HHMISS:
		format = t.Format("150405")
	case YYYYMMDD_HHMISS:
		format = t.Format("20060102_150405")
	case DIR:
		format = filepath.Join(t.Format("20060102"), t.Format("150405"))
	}
	return format
}

// TrimPathSeparator はパス名(/tmp/log/data/)から前後のセパレータを
// 取り除きます(tmp/log/data)
func TrimPathSeparator(path string) string {
	return strings.Trim(path, string(os.PathSeparator))
}

// ; Log level. None 0, FATAL 1, CRIT 2, ERR 3, WARN 4, NOTICE 5, INFO 6, DBG 7
// LOG_LEVEL = 5
func SetLogLevel(level int) error {
	switch level {
	case 0, 1:
		log.SetLevel(log.FatalLevel)
	case 2, 3:
		log.SetLevel(log.ErrorLevel)
	case 4, 5:
		log.SetLevel(log.WarnLevel)
	case 6:
		log.SetLevel(log.InfoLevel)
	case 7:
		log.SetLevel(log.DebugLevel)
	default:
		return fmt.Errorf("unkown log level %d", level)
	}
	return nil
}

// // CheckDiskFree は指定したディレクトリのディスク使用量[%]を取得します。
// func CheckDiskFree(dir string, capacity *int) error {
// 	return nil
// }

// CheckDiskUtil はディスク容量のチェックをします。
// func (config Config)CheckDiskUtil() (bool, error) {
// 	return true, nil
// }
