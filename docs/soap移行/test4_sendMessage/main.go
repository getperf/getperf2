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
<soap-env:Envelope xmlns:soap-env="http://schemas.xmlsoap.org/soap/envelope/">
  <soap-env:Header xmlns:wsa="http://www.w3.org/2005/08/addressing">
    <wsa:Action>urn:sendMessage</wsa:Action>
    <wsa:MessageID>urn:uuid:254de93f-816d-4fc2-b424-de340f173c8b</wsa:MessageID>
    <wsa:To>https://192.168.231.160:58443/axis2/services/GetperfService.GetperfServiceHttpsSoap11Endpoint/</wsa:To>
  </soap-env:Header>
  <soap-env:Body>
    <ns0:sendMessage xmlns:ns0="http://perf.getperf.com">
      <ns0:siteKey>site1</ns0:siteKey>
      <ns0:hostname>host1</ns0:hostname>
      <ns0:severity>1</ns0:severity>
      <ns0:message>this is a test</ns0:message>
    </ns0:sendMessage>
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
	r, err := http.NewRequest(
		http.MethodPost,
		"https://192.168.231.160:58443/axis2/services/GetperfService.GetperfServiceHttpsSoap11Endpoint/",
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

	// caCertPool := x509.NewCertPool()
	// caCertPool.AppendCertsFromPEM(caCert)

	// pem, err := os.ReadFile(caCertFile)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// caCertPool, err := x509.SystemCertPool()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// if !caCertPool.AppendCertsFromPEM(pem) {
	// 	log.Fatal("failed to add ca cert")
	// }

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

func soapCall(req *http.Request) (*Response, error) {
	transport, err := get_transport_with_ssl()
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{
		Timeout:   1500 * time.Millisecond,
		Transport: transport,
	}
	// log.Print("httpReq:\n")
	// log.Println(req)
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
				// fmt.Printf("Part : %q\n", slurp)
				err = xml.Unmarshal(slurp, &r)
				// fmt.Printf("XML DECODE:%v, %v\n", r, err)
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
