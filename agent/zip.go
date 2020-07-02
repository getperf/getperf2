package agent

import (
    "archive/zip"
    "io"
    "io/ioutil"
    "os"
    "path/filepath"
    "strings"

    // "github.com/jhoonb/archivex"
    "github.com/pkg/errors"
    log "github.com/sirupsen/logrus"
)

// func Zip2(zipPath, outDir, logPath string) error {
//     zipArc := archivex.ZipFile{}

//     prev, err := filepath.Abs(".")
//     if err != nil {
//         return fmt.Errorf("zip %s", err)
//     }
//     defer os.Chdir(prev)

//     os.Chdir(outDir)

//     // archivex の Close()処理はオープンファイルのクローズをしていないため、
//     // Create()を使用せずに独自にオープン、クローズ処理を追加

//     // if err := zip.Create(zipPath); err != nil {
//     //     return fmt.Errorf("zip %s", err)
//     // }
//     // defer zip.Close()
//     file, err := os.Create(zipPath)
//     if err != nil {
//         return err
//     }
//     defer file.Close()
//     zipArc.Writer = zip.NewWriter(file)
//     defer zipArc.Writer.Close()

//     if err := zipArc.AddAll(logPath, true); err != nil {
//         return fmt.Errorf("zip %s", err)
//     }
//     return nil
// }

func Zip(zipPath, outDir, logPath string) error {
    destinationFile, err := os.Create(zipPath)
    if err != nil {
        return err
    }
    defer destinationFile.Close()
    myZip := zip.NewWriter(destinationFile)
    defer myZip.Close()
    pathToZip := filepath.Join(outDir, logPath)
    outDir = outDir + string(os.PathSeparator)
    err = filepath.Walk(pathToZip, func(filePath string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if info.IsDir() {
            return nil
        }
        relPath := strings.TrimPrefix(filePath, outDir)

        // Windows でディレクトリパスはファイル名として展開されてしまうため"/"に置換
        // Windows\20200109\141630\ProcessorMemory.csv など
        relPath = strings.Replace(relPath, "\\", "/", -1)
        log.Debug("zip add ", relPath)
        header, _ := zip.FileInfoHeader(info)
        header.Name = relPath

        zipFile, err := myZip.CreateHeader(header)
        if err != nil {
            return err
        }
        fsFile, err := os.Open(filePath)
        if err != nil {
            return err
        }
        _, err = io.Copy(zipFile, fsFile)
        if err != nil {
            return err
        }
        return nil
    })
    if err != nil {
        return err
    }
    return nil
}

func Unzip(src, dest string) error {
    r, err := zip.OpenReader(src)
    if err != nil {
        return errors.Wrap(err, "prepare unzip")
    }
    defer r.Close()
    for _, f := range r.File {
        rc, err := f.Open()
        if err != nil {
            return errors.Wrap(err, "unzip open")
        }
        defer rc.Close()
        if f.FileInfo().IsDir() {
            path := filepath.Join(dest, f.Name)
            os.MkdirAll(path, 0755)
        } else {
            buf := make([]byte, f.UncompressedSize)
            _, err = io.ReadFull(rc, buf)
            if err != nil {
                return errors.Wrap(err, "unzip read")
            }

            path := filepath.Join(dest, f.Name)
            dirname := filepath.Dir(path)
            if err := os.MkdirAll(dirname, 0755); err != nil {
                return errors.Wrap(err, "unzip mkdir")
            }
            if err = ioutil.WriteFile(path, buf, f.Mode()); err != nil {
                return errors.Wrap(err, "unzip file write")
            }
        }
    }
    return nil
}
