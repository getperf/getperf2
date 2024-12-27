環境セットアップ
リグレッション

# 環境セットアップ

テストコード


./agent/soap_test.go


http://10.45.50.210:57000/axis2/

ssh psadmin@10.45.50.210

more ~/getperf/config/site/IS_SITE1.json
{
        "site_key":   "IS_SITE1",
        "access_key": "53f798920a1919a8242b9c26caa5742b5939f882",
        "home":       "/home/psadmin/work/IS_SITE1",


~/ptune/bin/getperfctl setup --url=https://10.45.50.210:57443/

疎通確認は出来そう

環境変数化

    soapSender, err = NewSoapSender("192.168.133.128", 57443)

var TEST_WEB_SERVICE_HOST = "192.168.133.128"

go test 
go test -run TestMakeTransportWithServerAuthSSL ./agent/

# github.com/getperf/getperf2/agent [github.com/getperf/getperf2/agent.test]

# ioutil.ReadALL 変更

agent/soap.go:196:19: undefined: io.ReadAll

ioutil.ReaAll

func readAll(r io.Reader, capacity int64) (b []byte, err error) {

 content, err := ioutil.ReadAll(r)

go test -run TestMakeTransportWithServerAuthSSL ./agent/ -v

# 疎通テスト


grep "3\.0" */*.go
agent/https_test.go:    config := NewConfig("../testdata/ptune-base-3.0", NewConfigEnv())
agent/soap_ssl_test.go: config := NewConfig("../testdata/ptune-base-3.0", NewConfigEnv())
agent/soap_ssl_test.go: config := NewConfig("../testdata/ptune-3.0", NewConfigEnv())
agent/soap_test.go:     // config := NewConfig("../testdata/ptune-base-3.0", NewConfigEnv())
agent/soap_test.go:     config := NewConfig("../testdata/ptune-base-3.0", NewConfigEnv())
agent/soap_test.go:     config := NewConfig("../testdata/ptune-3.0", NewConfigEnv())

cp ~/ptune/network/ca.crt testdata/ptune-base-3.0/network/
cp ~/ptune/network/getperf_ws.ini testdata/ptune-base-3.0/network/

 cp -r ~/ptune/network/ ./testdata/ptune-3.0/


GODEBUG=x509ignoreCN=0 go test -run TestSoapCallGetLatestBuild ./agent/ -v

タイムアウト発生

10.45.50.210 に変える



unset http_proxy HTTP_PROXY https_proxy HTTPS_PROXY

 x509: certificate relies on legacy Common Name field, use SANs or temporarily enable Common Name matching with GODEBUG=x509ignoreCN=0

export GODEBUG=x509ignoreCN=1

GODEBUG=x509ignoreCN=0 go test -run TestSoapCallGetLatestBuild ./agent/ -v

