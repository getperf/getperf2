現状調査、ssl疎通できるようにする
SSLコード分離変更調査
レグレッション

# 現状調査

soap.go

generateSOAPRequest

func (soapSender *SoapSender) generateSOAPRequest(requestMsg string) (*http.Request, error)
func soapCall(req *http.Request) (*Response, error)

soap_test.go

TestSoapGetLatestBuildRequest

go test ./agent/ -run TestSoapGetLatestBuildRequest -v

テストコード作成

go test ./agent/ -run TestSoapCallGetLatestBuild -v

```xml
>}, ERROR <nil>
MEDIA:text/xml,PARAMS:map[charset:UTF-8]
Part : "<?xml version='1.0' encoding='UTF-8'?><soapenv:Envelope xmlns:soapenv=\"http://schemas.xmlsoap.org/soap/envelope/\"><soapenv:Header xmlns:wsa=\"http://www.w3.org/2005/08/addressing\"><wsa:Action>http://www.w3.org/2005/08/addressing/soap/fault</wsa:Action></soapenv:Header><soapenv:Body><soapenv:Fault><faultcode></faultcode><faultstring>com.ctc.wstx.exc.WstxEOFException: Unexpected EOF in prolog\n at [row,col {unknown-source}]: [1,0]</faultstring><detail /></soapenv:Fault></soapenv:Body></soapenv:Envelope>"
XML DECODE:&{<nil>}, <nil>
READ <nil>, ERROR EOF
BODY

```

