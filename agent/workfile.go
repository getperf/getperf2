package agent

import (
    "bufio"
    "fmt"
    "os"
    "path/filepath"
    "strconv"
    // log "github.com/sirupsen/logrus"
)

const (
    initialTokenSize = 4096
    maxTokenSize     = 500 * 1024 // scanner.Err()でエラーになる場合、ここの数字を増やす
)

// GetWorkfilePath はファイル名が'_'で始まる場合は共有ディレクトリ下のパスを、
// そうでない場合はローカルディレクトリ下のパスを返します。
func (config Config) GetWorkfilePath(filename string) string {
    workDir := config.WorkDir
    if filename[0] == '_' {
        workDir = config.WorkCommonDir
    }
    workPath := filepath.Join(workDir, filename)
    return workPath
}

func (config Config) GetArchivefilePath(filename string) string {
    return filepath.Join(config.ArchiveDir, filename)
}

// ReadFile はファイルを指定した行数分読み込みます。
func ReadFileHead(filename string, maxRow int) ([]string, error) {
    var lines []string
    file, err := os.Open(filename)
    if err != nil {
        return lines, fmt.Errorf("read work file %s : %s", filename, err)
    }
    defer file.Close()
    scanner := bufio.NewScanner(file)
    buf := make([]byte, initialTokenSize)
    //bufのサイズが足りなかった場合、maxTokenSizeまでバッファを拡大してスキャンする
    scanner.Buffer(buf, maxTokenSize)

    row := 0
    for scanner.Scan() {
        row++
        if maxRow != 0 && row > maxRow { // maxRowが0の場合は全行読み込み
            break
        }
        lines = append(lines, scanner.Text())
    }
    // bufferが足りなかった場合などは、scanner.Err()でエラーを検出
    err = scanner.Err()
    if err != nil {
        return nil, fmt.Errorf("read work file %s : %s", filename, err)
    }
    return lines, nil
}

// ReadWorkFile はワークファイルを指定した行数分読み込みます。
func (config Config) ReadWorkFileHead(filename string, maxRow int) ([]string, error) {
    return ReadFileHead(config.GetWorkfilePath(filename), maxRow)
}

// ReadWorkFile はワークファイルを全行読み込みます。
func (config Config) ReadWorkFile(filename string) ([]string, error) {
    return config.ReadWorkFileHead(filename, 0)
}

// WriteWorkFile はークファイルの書き込みをします。
func (config Config) WriteWorkFile(filename string, lines []string) error {
    file, err := os.Create(config.GetWorkfilePath(filename))
    if err != nil {
        return fmt.Errorf("write work file %s: %s", filename, err)
    }
    defer file.Close()
    for _, line := range lines {
        _, err := file.WriteString(line + "\n")
        if err != nil {
            return fmt.Errorf("write work file %s: %s", filename, err)
        }
    }
    return nil
}

func (config Config) WriteLineWorkFile(filename, line string) error {
    return config.WriteWorkFile(filename, []string{line})
}

// ReadWorkFileNumber はワークファイルから数値を読み込みます。
func (config Config) ReadWorkFileNumber(filename string) (int, error) {
    lines, err := config.ReadWorkFile(filename)
    if err != nil {
        return 0, fmt.Errorf("read work file number %s", err)
    }
    i, err2 := strconv.Atoi(lines[0])
    if err2 != nil {
        return 0, fmt.Errorf("read work file number %s", err2)
    }
    return i, nil
}

// WriteWorkFileNumber はワークファイルへの数値の書き込みます。
func (config Config) WriteWorkFileNumber(filename string, num int) error {
    return config.WriteWorkFile(filename, []string{strconv.Itoa(num)})
}

// CheckWorkFile はワークファイルの確認の有無を確認します。
func (config Config) CheckWorkFile(filename string) (bool, error) {
    _, err := os.Stat(config.GetWorkfilePath(filename))
    if err == nil {
        return true, nil
    }
    if os.IsNotExist(err) {
        return false, nil
    }
    return false, err
}

// RemoveWorkFile はワークファイルの削除をします。
func (config Config) RemoveWorkFile(filename string) error {
    filepath := config.GetWorkfilePath(filename)
    _, err := os.Stat(filepath)
    if err == nil {
        if err2 := os.Remove(filepath); err2 != nil {
            return err2
        }
    }
    return nil
}

func (config Config) RemoveArchiveFile(filename string) error {
    filepath := config.GetArchivefilePath(filename)
    _, err := os.Stat(filepath)
    if err == nil {
        if err2 := os.Remove(filepath); err2 != nil {
            return err2
        }
    }
    return nil
}

func (config Config) ReadPid() (int, error) {
    return config.ReadWorkFileNumber(config.PidFile)
}
