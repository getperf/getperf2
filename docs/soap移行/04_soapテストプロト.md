# soap テストプロト

以下のテストをagent の下に,実装してみる。
http_test は使用しない。
直接ローカルサーバにアクセスする。

test1_getLatestBuild/main.go


soap.go
soap_test.go

network ディレクトリをコピー


export no_proxy=192.168.231.160


go test ./agent/ -run TestSoapGetLatestBuildRequest -v

```xml
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
        } 0x754c60 772 [] false 192.168.231.160:57443 map[] map[] <nil> map[]   <nil> <nil> <nil> 0xc00002a100}
--- PASS: TestSoapGetLatestBuildRequest (0.00s)
PASS
ok      github.com/getperf/getperf2/agent       0.899s

```

OK になり、 return 値で 0 が返ってくるようになった。

ハードコーディングされている箇所から順に修正を進める。午後から。

docs の下に soaptest を追加。

uuid を生成

go get github.com/google/uuid

