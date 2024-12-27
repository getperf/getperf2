再現調査
リグレッション

# 再現調査

TestSoapCallSendData

go test ./agent/ --run TestSoapCallSendData -v


2023/05/15 09:06:48 ResponseReturn: 1
Result : OK

TestSoapCallSendMessage

go test ./agent/ --run TestSoapCallSendMessage -v

メッセージ出力 

 tail -f tomcat.log
09:09:02.772 [ajp-nio-0:0:0:0:0:0:0:1-58009-exec-5] INFO  com.getperf.perf.EventManager - [site1,host1] this is a test3

# レグレッション

両者の違い。

```golang
func TestSoapCallSendMessage(t *testing.T) {

    httpReq, err := soapSenderData.MakeSoapRequest(req)

func TestSoapCallSendData(t *testing.T) {

    httpReq, err := soapSenderData.MakeSoapRequestWithAttachment(req, "../arc_host1__Linux_20230506_0800.zip")

```

プロトでの実行

```bash
go run test5_sendData/main.go

```

結果　ns:return>OK</ns:return>　となるが、tomcat.log が出ない

```xml
MEDIA:multipart/related,PARAMS:map[boundary:MIMEBoundary_33dfad477082d08ea7a28d6f25d1c1e27de66d44f14312c5 start:<0.03dfad477082d08ea7a28d6f25d1c1e27de66d44f14312c5@apache.org> type:text/xml]
MEDIA:text/xml,PARAMS:map[charset:UTF-8]
Part : "<?xml version='1.0' encoding='UTF-8'?><soapenv:Envelope xmlns:soapenv=\"http://schemas.xmlsoap.org/soap/envelope/\"><soapenv:Header xmlns:wsa=\"http://www.w3.org/2005/08/addressing\"><wsa:Action>urn:sendDataResponse</wsa:Action><wsa:RelatesTo>urn:uuid:254de93f-816d-4fc2-b424-de340f173c8b</wsa:RelatesTo></soapenv:Header><soapenv:Body><ns:sendDataResponse xmlns:ns=\"http://perf.getperf.com\"><ns:return>OK</ns:return></ns:sendDataResponse></soapenv:Body></soapenv:Envelope>"
XML DECODE:&{<nil>}, <nil>
2023/05/15 09:17:31 &{<nil>}
```

ls -ltr ~/getperf/t/staging_data/site1/arc_host1__Linux_20230506_0800.zip

-rw-r----- 1 psadmin psadmin 17248  5月 15 09:17 arc_host1__Linux_20230506_0800.zip

一旦消してみる

rm ~/getperf/t/staging_data/site1/arc_host1__Linux_20230506_0800.zip

soap_test.go では保存されないが、 main.go では保存される。

main.go を調査

wsa:To　を http ではなく、https を指定している。

    <wsa:To>https://192.168.231.160:58443/axis2/services/GetperfService.GetperfServiceHttpsSoap11Endpoint/</wsa:To>

指定を変えたら保存できるようになった。


type SoapSender struct {
    ServerIP         string
    MessageID        string
    ToURL            string <=== 廃止
    ServiceURL       string
    AttachedFilePath string
    Timeout          int
}
