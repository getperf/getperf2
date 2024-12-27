package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/xml"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// リクエストの型
type Request struct {
	Key      string
	ID       string
	Password string
}

// 結果
type Response struct {
	Header interface{} `xml:"soapenv:Header" json:"Header"`
}

// PostするXML
type Post struct {
	XMLName     xml.Name `xml:"Request"`
	Credentials struct {
		ID       string `xml:"id"`
		Password string `xml:"password"`
	} `xml:"Credentials"`
	Identity struct {
		Key string `xml:"key"`
	}
}

// テンプレート
var getTemplate = `
<soap-env:Envelope xmlns:soap-env="http://schemas.xmlsoap.org/soap/envelope/">
  <soap-env:Header xmlns:wsa="http://www.w3.org/2005/08/addressing">
    <wsa:Action>urn:sendData</wsa:Action>
    <wsa:MessageID>urn:uuid:254de93f-816d-4fc2-b424-de340f173c8b</wsa:MessageID>
    <wsa:To>https://192.168.231.160:58443/axis2/services/GetperfService.GetperfServiceHttpsSoap11Endpoint/</wsa:To>
  </soap-env:Header>
  <soap-env:Body>
    <ns0:sendData xmlns:ns0="http://perf.getperf.com">
      <ns0:siteKey>site1</ns0:siteKey>
      <ns0:filename>arc_host1__Linux_20230506_0800.zip</ns0:filename>
    </ns0:sendData>
  </soap-env:Body>
</soap-env:Envelope>
`

func generateSOAPRequest(req *Request) (*http.Request, error) {
	// テンプレートを使ってXMLを作成
	temp, err := template.New("InputRequest").Parse(getTemplate)

	if err != nil {
		log.Printf("Error while marshling object. %s ", err.Error())
		return nil, err
	}

	doc := &bytes.Buffer{}
	err = temp.Execute(doc, req)
	if err != nil {
		log.Printf("template.Execute error. %s ", err.Error())
		return nil, err
	}
	// err = uploadZip()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	body := new(bytes.Buffer)
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

		// part := make(textproto.MIMEHeader)
		// writer, err := mpWriter.CreatePart(part)
		// part, err := mpWriter.CreateFormFile("file", filepath.Base(file.Name()))
		// if err != nil {
		// 	log.Printf("create zip part %s", err.Error())
		// 	return nil, err
		// }
		// io.Copy(part, file)
	}
	mpWriter.Close()
	r, err := http.NewRequest(
		http.MethodPost,
		"https://192.168.231.160:58443/axis2/services/GetperfService.GetperfServiceHttpsSoap11Endpoint/",
		body,
		// doc,
	)
	// r.Header.Add("Content-Type", "text/xml; charset=UTF-8")
	fmt.Printf("CONTENT_TYPE:%s\n", mpWriter.FormDataContentType())
	fmt.Printf("BOUDNDARY:%s\n", mpWriter.Boundary())
	content_type := fmt.Sprintf(
		"multipart/related;boundary=%s;type=\"text/xml\";charset=UTF-8;",
		mpWriter.Boundary(),
	)
	// content_type := fmt.Sprintf(
	// 	"multipart/related;boundary=%s;type=\"text/xml\";start=\"%s\";start-info=\"text/xml\";charset=UTF-8;",
	// 	mpWriter.Boundary(),
	// 	content_id_start,
	// )
	r.Header.Add("Content-Type", content_type)
	// r.Header.Add("Content-Type", mpWriter.FormDataContentType())

	fmt.Printf("HEADER:\n%s\n\n", r.Header.Get("Content-Type"))
	// fmt.Printf("BODY:\n%s\n\n", body)

	if err != nil {
		log.Printf("Error making a request. %s ", err.Error())
		return nil, err
	}

	return r, nil
}

// リクエストの内容
func populateRequest() *Request {
	req := Request{}
	req.Key = "12345678"
	req.ID = "SENOUE"
	req.Password = "Password"
	return &req
}

func get_transport_with_ssl() (*http.Transport, error) {
	caCertFile := string("./network/ca.crt")
	clientCertFile := string("./network/client.crt")
	clientKeyFile := string("./network/client.key")

	cert, err := tls.LoadX509KeyPair(clientCertFile, clientKeyFile)
	if err != nil {
		log.Fatalf("Error creating x509 keypair from client cert file %s and client key file %s", clientCertFile, clientKeyFile)
	}

	caCert, err := ioutil.ReadFile(caCertFile)
	if err != nil {
		log.Fatalf("Error opening cert file %s, Error: %s", caCertFile, err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	defaultTransport, ok := http.DefaultTransport.(*http.Transport)
	if !ok {
		log.Fatal("invalid default transport")
	}

	transport := defaultTransport.Clone()

	transport.TLSClientConfig = &tls.Config{
		RootCAs:      caCertPool,
		Certificates: []tls.Certificate{cert},
		ServerName:   "192.168.231.160",
	}
	return transport, nil
}

// getperf2 の zip 添付実装

func uploadZip() error {
	zipname := "arc_host1__Linux_20230506_0800.zip"
	zippath := "./arc_host1__Linux_20230506_0800.zip"
	file, err := os.Open(zippath)
	if err != nil {
		log.Printf("open error %s ", err.Error())
		return err
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
	log.Printf("CONTENTTYPE: %s", contentType)

	// var w http.ResponseWriter
	// w.Header().Set("Content-Disposition",
	// 	fmt.Sprintf(`form-data; name="%s"; filename="%s"`, fieldname, zipname))
	// w.Header().Set("Content-Type", contentType)

	// fwで作ったパートにファイルのデータを書き込む
	_, err = io.Copy(fw, file)
	if err != nil {
		log.Printf("write MIME multipart %s", err.Error())
		return err
	}

	// リクエストのContent-Typeヘッダに使う値を取得する（バウンダリを含む）
	err = mw.Close()
	if err != nil {
		log.Printf("close MIME multipart %s", err.Error())
		return err
	}
	// w.WriteHeader(200)

	// w.Write(body.Bytes())
	log.Printf("download : %s", zipname)
	// log.Print(body)
	return nil
}

func soapCall(req *http.Request) (*Response, error) {
	transport, err := get_transport_with_ssl()
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{
		Timeout:   1500 * time.Millisecond,
		Transport: transport,
	}
	log.Print("httpReq:\n")
	log.Println(req)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	return nil, err
	// }
	defer resp.Body.Close()
	r := &Response{}
	fmt.Printf("RESP\n")
	// fmt.Println(body)
	mediaType, params, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	fmt.Printf("MEDIA:%v,PARAMS:%v\n", mediaType, params)
	if strings.HasPrefix(mediaType, "multipart/") {
		mr := multipart.NewReader(resp.Body, params["boundary"])
		for {
			p, err := mr.NextPart()
			// fmt.Printf("READ %v, ERROR %v\n", p, err)
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			mediaType, params, err := mime.ParseMediaType(p.Header.Get("Content-Type"))
			fmt.Printf("MEDIA:%v,PARAMS:%v\n", mediaType, params)
			if strings.HasSuffix(mediaType, "/xml") {
				slurp, err := io.ReadAll(p)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf("Part : %q\n", slurp)
				err = xml.Unmarshal(slurp, &r)
				fmt.Printf("XML DECODE:%v, %v\n", r, err)
			}
			if strings.HasSuffix(mediaType, "/octet-stream") {
				tmpfile, err := os.Create("./" + "resp.dat")
				if err != nil {
					log.Fatal(err)
				}
				defer tmpfile.Close()
				// slurp, err := io.ReadAll(p)
				// if err != nil {
				// 	log.Fatal(err)
				// }
				_, err = io.Copy(tmpfile, p)
				// _, err = tmpfile.Write(slurp)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}
	// fmt.Printf("BODY\n")
	// fmt.Printf(string(body))
	// err = xml.Unmarshal(body, &r)

	// if err != nil {
	// 	return nil, err
	// }
	return r, nil
}

func main() {
	req := populateRequest()
	httpReq, err := generateSOAPRequest(req)
	if err != nil {
		log.Println("Some problem occurred in request generation")
	}
	response, err := soapCall(httpReq)
	if err != nil {
		log.Printf("soap call error %s", err)
	}
	log.Print(response)
}
