package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/xml"
	"fmt"
	"html/template"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
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
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:wsa="http://www.w3.org/2005/08/addressing" xmlns:ns0="http://perf.getperf.com">
    <soapenv:Header>
        <wsa:Action>urn:getLatestBuild</wsa:Action>
        <wsa:MessageID>urn:uuid:0d4cec05-e391-4958-a423-923f1af0c364</wsa:MessageID>
        <wsa:To>http://192.168.231.160:57000/axis2/services/GetperfService.GetperfServiceHttpsSoap11Endpoint/</wsa:To>
    </soapenv:Header>
    <soapenv:Body>
        <ns0:getLatestBuild xmlns:ns0="http://perf.getperf.com" xmlns="http://rmi.java/xsd">
            <ns0:moduleTag>Linux</ns0:moduleTag>
            <ns0:majorVer>2</ns0:majorVer>
        </ns0:getLatestBuild xmlns:ns0="http://perf.getperf.com">
    </soapenv:Body>
</soapenv:Envelope>
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
	r, err := http.NewRequest(
		http.MethodPost,
		"https://192.168.231.160:57443/axis2/services/GetperfService.GetperfServiceHttpsSoap11Endpoint/",
		doc,
	)
	r.Header.Add("Content-Type", "text/xml; charset=UTF-8")
	// r.Header.Add("Accept", "text/xml")
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
	pem, err := os.ReadFile("./network/ca.crt")
	if err != nil {
		log.Fatal(err)
	}

	caCertPool, err := x509.SystemCertPool()
	if err != nil {
		log.Fatal(err)
	}

	if !caCertPool.AppendCertsFromPEM(pem) {
		log.Fatal("failed to add ca cert")
	}

	defaultTransport, ok := http.DefaultTransport.(*http.Transport)
	if !ok {
		log.Fatal("invalid default transport")
	}

	transport := defaultTransport.Clone()

	transport.TLSClientConfig = &tls.Config{
		RootCAs:    caCertPool,
		ServerName: "192.168.231.160",
	}
	return transport, nil
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
			fmt.Printf("READ %v, ERROR %v\n", p, err)
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			mediaType, params, err := mime.ParseMediaType(p.Header.Get("Content-Type"))
			fmt.Printf("MEDIA:%v,PARAMS:%v\n", mediaType, params)
			slurp, err := io.ReadAll(p)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Part : %q\n", slurp)
			if strings.HasSuffix(mediaType, "/xml") {
				err = xml.Unmarshal(slurp, &r)
			}
			fmt.Printf("XML DECODE:%v, %v\n", r, err)
		}
	}
	fmt.Printf("BODY\n")
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
