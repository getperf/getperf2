変更調査
レグレッション
プロトタイプ

# 変更調査

soapCall そのままの記述でOK

ファイル保存の指定件古澤です。お疲れ様です。

```golang
    r, err := http.NewRequest(
        http.MethodPost,
        "https://192.168.231.160:58443/axis2/services/GetperfService.GetperfServiceHttpsSoap11Endpoint/",
        body,
        // doc,
    )
```


# タイムアウトパラメータ化

type SoapSender struct {
    ServerIP  string
    MessageID string
    ToURL     string
    Timeout   int
}

テスト

go test ./agent/ --run TestNewSoapSender -v

ok

ファイル受信パラメータ化

type SoapSender struct {
    ServerIP         string
    MessageID        string
    ToURL            string
    ServiceURL       string
    AttachedFilePath string
    Timeout          int
}

# エラーハンドラ調査

            return nil, errors.Wrap(err, "failed to load X509KeyPair")
        return errors.Wrap(err, "prepare listen port")

var errTimedOut = errors.New("command timed out")

