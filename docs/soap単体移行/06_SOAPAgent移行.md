Config レグレッション
変更調査
プロト

# Config レグレッション

soap_agent.go 

type SoapAgent struct {
    Config     *Config
    SopaSender *SoapSender
}

構造体なくてもよい

soapSender.ReserveSender() error 

soap_test.go

```golang
func TestSoapCallGetLatestBuild(t *testing.T) {
    req, err := soapSender.MakeSoapRequestMsg(
        "getLatestBuild",
        map[string]string{
            "moduleTag": "Linux",
            "majorVer":  "2",
        },
    )
    if err != nil {
        t.Errorf("Some problem occurred in request generation")
    }
    httpReq, err := soapSender.MakeSoapRequest(req)
    if err != nil {
        t.Errorf("check request %s", err)
    }
    response, err := soapSender.soapCall(httpReq)
    if err != nil {
        t.Logf("soap call error %s", err)
    }
    t.Logf("response %v", response)
}
```

soap_agent_test.go

soapSender

```golang
func TestReserveSender(t *testing.T) {
    soapSender.ReserveSender(ファイル)
}
```

testdata/ptune/backup/arc_centos80__HW_20200528_132000.zip

soapSender と使用するディレクトリパスを宣言

gpf_soap_admin.c と gpf_soap_agent.c を調査

gpf_soap_agent.c

gpfReserveSender()

archiveDir
schedule->urlPM, "", schedule->siteKey, filename, &fileSize

gpfSendData()
gpfSendMessage()

schedule->siteKey, config->host

gpfDownloadCertificate()


urlCM
urlPM

インスタンスは soapSender にする
引数に Config や、schedule などを追加する


ReserveSender(c, filename )
SendData(c, filename)
SendMessage(c, severity, message )
DownloadCertificate(c, timestamp )   
ReserveFileSender(c, onOff, waitSec )    
SendZipData(c, filename)
DownloadConfigFilePM(c, filename )

GetLatestBuild( c )
RegistAgent( c, setup )
DownloadUpdateModule( c, _build, moduleFile)  
CheckHostStatus( c, setup)

スケルトン作成

