変更調査
レグレッション
プロト

# 変更調査

既存コード調査

go run test5_sendData/main.go

テンプレート sendData

mpWriter

```golang
    mpWriter := multipart.NewWriter(body)
    // content_id_start := "<0.09BC7F4BE2E4D3EF1B@apache.org>"
    {
        part := make(textproto.MIMEHeader)
        part.Set("Content-Type", "text/xml; charset=UTF-8")
        // part.Set("content-id", content_id_start)
        writer, err := mpWriter.CreatePart(part)
        if err != nil {
            log.Printf("create xml part %s", err.Error())
            return nil, err
        }
        writer.Write(doc.Bytes())
    }
    {
        filePath := "./arc_host1__Linux_20230506_0800.zip"
        file, _ := os.Open(filePath)
        if err != nil {
            log.Printf("create zip part %s", err.Error())
            return nil, err
        }
        defer file.Close()

        part := make(textproto.MIMEHeader)
        part.Set("Content-Type", "application/octet-stream")
        part.Set("content-transfer-encoding", "binary")
        part.Set("Content-Disposition", `form-data; name="file"; filename="arc_host1__Linux_20230506_0800.zip"`)
        part.Set("content-id", "<"+filepath.Base(file.Name())+">")
        writer, err := mpWriter.CreatePart(part)
        if err != nil {
            log.Printf("create zip part %s", err.Error())
            return nil, err
        }
        io.Copy(writer, file)
    }
    mpWriter.Close()

```

ヘッダ作成

```golang
    fmt.Printf("CONTENT_TYPE:%s\n", mpWriter.FormDataContentType())
    fmt.Printf("BOUDNDARY:%s\n", mpWriter.Boundary())
    content_type := fmt.Sprintf(
        "multipart/related;boundary=%s;type=\"text/xml\";charset=UTF-8;",
        mpWriter.Boundary(),
    )
    // content_type := fmt.Sprintf(
    //  "multipart/related;boundary=%s;type=\"text/xml\";start=\"%s\";start-info=\"text/xml\";charset=UTF-8;",
    //  mpWriter.Boundary(),
    //  content_id_start,
    // )
    r.Header.Add("Content-Type", content_type)

```

テストコード作成

go test ./agent/ -run TestSoapCallSendData -v

ヘッダ作成

プロトタイプ

org.apache.axis2.AxisFault: NO Attachment ID arc_host1__Linux_20230506_0800.zip

ok になった

tomcat.log に出力されない

./webapps/axis2/WEB-INF/conf/axis2.xml

13:42:25,779 |-INFO in ch.qos.logback.core.FileAppender[LOGFILE] - File property is set to [/usr/local/tomcat-data/logs/tomcat.log]

ls -ltr ~/getperf/t/staging_data/site1/
ls -ltr ~/getperf/t/staging_data/site1/arc_host1__Linux_20230506_0800.zip

送信されていない。保留とする



