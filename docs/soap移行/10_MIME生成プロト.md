クライアント認証プロト
変更調査
レグレッション
プロトタイプ

# クライアント認証プロト

現状調査

```golang

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

```

サーバ証明書の設定をして、send message のテスト実行

リクエストテンプレート

    soap_test.go:40: check request template: InputRequest:13: function "message" not defined


レグレッション

go test ./agent/ -run TestSoapCallSendMessage -v

11:53:15.444 [ajp-nio-0:0:0:0:0:0:0:1-58009-exec-14] INFO  com.getperf.perf.EventManager - [site1,host1] this is a test3

動くようになった


