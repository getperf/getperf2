package agent

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

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
	if err == nil {
		return !fi.Mode().IsDir(), nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
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

// // CheckDiskFree は指定したディレクトリのディスク使用量[%]を取得します。
// func CheckDiskFree(dir string, capacity *int) error {
// 	return nil
// }

// CheckDiskUtil はディスク容量のチェックをします。
// func (config Config)CheckDiskUtil() (bool, error) {
// 	return true, nil
// }

// CheckPathInHome はパスがホーム下を指定しているか、".."が含まれないかをチェックします。
func (config Config) CheckPathInHome(path string) (bool, error) {

	// 	if ( strstr( path, config->home ) != path)
	// 		return gpfError( "path error (home) %s", path );

	// 	/* ".."が2回以上含まれている場合はNG */

	// #if defined _WINDOWS
	// 	if ( strstr( path, "../.." ) != NULL )
	// 		return gpfError( "path error (../..) %s", path );
	// #else
	// 	if ( strstr( path, "..\\.." ) != NULL )
	// 		return gpfError( "path error (..\\..) %s", path );
	// #endif

	return true, nil
}

// Usage はヘルプメッセージを出力します。
func Usage(msgs *[]string) {
}

// BackupConfig は構成ファイルのバックアップをします。
func BackupConfig(srcDir string, targetDir string, filename string) error {
	return nil
}
