soapメッセージ単体

# 変更調査

```xml
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
```

リクエストパラメータ

soap.go

```golang
// リクエストの型
type Request struct {
    MessageID string
    Action    string
    To        string
    ModuleTag string
    MajorID   int
}

// リクエスト引数
type RequestArgs map[string]string
```

単体テスト

make_request(string, string[string])
make_request_message(string, string[string])

func MakeRequest(action string, reqestHeader *RequestHeder) (*Request, error)

soapSender := NewSoapSender(url)

    messege id と、 to を作成する

soapRequest := soapSender.NewRequest(action, args)
soapMessage := soapSender.NewSoapMessage(soapRequest)

r, err := http.NewRequest(
    http.MethodPost,
    "https://192.168.231.160:57443/axis2/services/GetperfService.GetperfServiceHttpsSoap11Endpoint/",
    soapMessage,
)

# soapSenderコンストラクタテスト

func TestSoapSender(t *testing.T) {
    url := "http://192.168.231.160:57000/axis2/services/GetperfService.GetperfServiceHttpsSoap11Endpoint/"
    soapSender, err := NewSoapSender(url)
    if err != nil {
        t.Errorf("check request %s", err)
    }
    t.Logf(soapSender)
}


go test ./agent/ -run TestSoapSender -v

# soapReqestコンストラクタテスト

type SoapRequest struct {
    Action    string
    Args   map[string]string
}

func NewSoapRequest(action string, args map[string]string) *SoapRequest {
    return &SoapRequest{
        Action: action,
        Args:   args,
    }
}

