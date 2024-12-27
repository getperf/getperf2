package agent

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Response struct {
	Header interface{} `xml:"soapenv:Header" json:"Header"`
}

type SoapSender struct {
	ServerIP         string
	MessageID        string
	ServiceURL       string
	Transport        *http.Transport
	AttachedFilePath string
	Timeout          int
}

var endPointSuffix = "/axis2/services/GetperfService.GetperfServiceHttpsSoap11Endpoint/"

func NewSoapSender(serverIp string, port int) (*SoapSender, error) {
	u, err := uuid.NewRandom()
	if err != nil {
		return nil, errors.Wrap(err, "initialize soap sender")
	}
	serviceURL := fmt.Sprintf("https://%s:%d", serverIp, port)
	soapSender := SoapSender{
		ServerIP:   serverIp,
		MessageID:  u.String(),
		ServiceURL: serviceURL + endPointSuffix,
		Timeout:    10,
	}
	return &soapSender, nil
}

func (sender *SoapSender) WithTimeout(timeout int) *SoapSender {
	sender.Timeout = timeout
	return sender
}

func (sender *SoapSender) WithAttachedFilePath(filePath string) *SoapSender {
	sender.AttachedFilePath = filePath
	return sender
}

func (sender *SoapSender) MakeSoapRequestMsg(action string, requests map[string]string) (string, error) {
	requestMsgTemplate, ok := SoapRequestMsgTemplates[action]
	if !ok {
		return "", errors.New("soap action template not found")
	}
	temp, err := template.New("InputRequest").Parse(requestMsgTemplate)
	if err != nil {
		return "", errors.Wrap(err, "soap action template create")
	}

	// インスタンスから共通定義の MessageID,ToURL を取得してセット
	requests["MessageID"] = sender.MessageID
	requests["ToURL"] = sender.ServiceURL

	doc := &bytes.Buffer{}
	err = temp.Execute(doc, requests)
	if err != nil {
		return "", errors.Wrap(err, "execute soap action template create")
	}
	return doc.String(), nil
}

func (soapSender *SoapSender) MakeSoapRequest(requestMsg string) (*http.Request, error) {
	body := new(bytes.Buffer)
	body.Write([]byte(requestMsg))
	r, err := http.NewRequest(
		http.MethodPost,
		soapSender.ServiceURL,
		body,
	)
	r.Header.Add("Content-Type", "text/xml; charset=UTF-8")
	// r.Header.Add("Accept", "text/xml")
	if err != nil {
		return nil, errors.Wrap(err, "make soap request")
	}
	return r, nil
}

func (soapSender *SoapSender) MakeSoapRequestWithAttachment(requestMsg string, filePath string) (*http.Request, error) {
	body := new(bytes.Buffer)
	mpWriter := multipart.NewWriter(body)
	{
		part := make(textproto.MIMEHeader)
		part.Set("Content-Type", "text/xml; charset=UTF-8")
		writer, err := mpWriter.CreatePart(part)
		if err != nil {
			return nil, errors.Wrap(err, "create xml part for make soap request")
		}
		writer.Write([]byte(requestMsg))
	}
	{
		file, err := os.Open(filePath)
		if err != nil {
			return nil, errors.Wrapf(err, "open soap attachment file")
		}
		defer file.Close()

		filename := filepath.Base(file.Name())
		log.Debugf("FILENAME:%v\n", filename)
		part := make(textproto.MIMEHeader)
		part.Set("Content-Type", "application/octet-stream")
		part.Set("content-transfer-encoding", "binary")
		part.Set("Content-Disposition", `form-data; name="file"; filename="`+filename+`"`)
		part.Set("content-id", "<"+filename+">")
		writer, err := mpWriter.CreatePart(part)
		log.Debugf("MIME PART:%v\n", part)
		if err != nil {
			return nil, errors.Wrap(err, "create zip part for make soap request")
		}
		io.Copy(writer, file)
	}
	mpWriter.Close()

	r, err := http.NewRequest(
		http.MethodPost,
		soapSender.ServiceURL,
		body,
	)
	if err != nil {
		log.Printf("Error making a request. %s ", err.Error())
		return nil, errors.Wrap(err, "create new request for soap with attachement")
	}

	content_type := fmt.Sprintf(
		"multipart/related;boundary=%s;type=\"text/xml\";charset=UTF-8;",
		mpWriter.Boundary(),
	)
	log.Debugf("CONTENT_TYPE:%s\n", content_type)
	r.Header.Add("Content-Type", content_type)

	return r, nil
}

func (soapSender *SoapSender) getResponseReturn(xmlMsg string) (string, error) {
	rep := regexp.MustCompile(`<ns:return>(.*)</ns:return>`)
	result := rep.FindAllStringSubmatch(xmlMsg, -1)
	log.Debugf("response length: %v\n", len(result))
	if len(result) == 0 {
		return "", errors.New("error to get <ns:return>{val}</ns:return> from xml : " + xmlMsg)
	}
	return result[0][1], nil
}

func (soapSender *SoapSender) soapCall(req *http.Request) (string, error) {
	client := &http.Client{
		Timeout: time.Duration(soapSender.Timeout) * time.Second,
	}
	if soapSender.Transport != nil {
		client.Transport = soapSender.Transport
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "soap call")
	}

	defer resp.Body.Close()
	result := ""
	mediaType, params, _ := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	log.Debugf("MEDIA:%v,PARAMS:%v\n", mediaType, params)
	if strings.HasPrefix(mediaType, "multipart/") {
		mr := multipart.NewReader(resp.Body, params["boundary"])
		for {
			p, err := mr.NextPart()
			log.Debugf("READ %v, ERROR %v\n", p, err)
			if err == io.EOF {
				break
			}
			if err != nil {
				return "", errors.Wrap(err, "parse multipart after soap call")
			}
			mediaType, params, _ := mime.ParseMediaType(p.Header.Get("Content-Type"))
			log.Debugf("MEDIA:%v,PARAMS:%v\n", mediaType, params)
			if strings.HasSuffix(mediaType, "/xml") {
				// slurp, err := io.ReadAll(p)
				slurp, err := ioutil.ReadAll(p)
				if err != nil {
					return "", errors.Wrap(err, "read xml part after soap call")
				}
				result, err = soapSender.getResponseReturn(string(slurp))
				if err != nil {
					return "", errors.Wrap(err, "parse xml part after soap call")
				}
				// err = xml.Unmarshal(slurp, &r)
				log.Debugf("Result : %v\n", result)
			}
			if strings.HasSuffix(mediaType, "/octet-stream") {
				tmpfile, err := os.Create(soapSender.AttachedFilePath)
				if err != nil {
					return "", errors.Wrap(err, "save zip file after soap call")
				}
				defer tmpfile.Close()
				_, err = io.Copy(tmpfile, p)
				if err != nil {
					return "", errors.Wrap(err, "write zip file after soap call")
				}
			}

		}
	}
	return result, nil
}
