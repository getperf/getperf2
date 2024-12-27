現状調査
リグレッション
TODO作成

# 現状調査

構造体なくてもよい

soapSender.ReserveSender() error 


type SoapSender struct {
    ServerIP         string
    MessageID        string
    ServiceURL       string
    Transport        *http.Transport
    AttachedFilePath string
    Timeout          int
}

テストコード調査

go test ./agent/ --run TestSoapGetLatestBuildRequest -v

UML 更新

スケルトン修正

# レグレッション

レグレッションテストコード作成

soap_agent_test.go

 go test ./agent/ --run TestReserveFileSender -v


