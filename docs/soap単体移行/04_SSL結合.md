レグレッション
変更調査
プロトタイプ

# レグレッション

現行テスト調査

https_test.go 

```bash
more testdata/ptune/network/License.txt
HOSTNAME=centos80
EXPIRE=20471023
CODE=150b835f65942444f8dc4b955066d7755a5f336b

```


TestGetTSLConfig テスト

```golang
    config := NewConfig("../testdata/ptune", NewConfigEnv())
    config.InitAgent()
    config.ParseConfigFile(config.SvParamFile)

```

TestConfigToTLSConfig テスト

```golang
func TestConfigToTLSConfig(t *testing.T) {
    config := NewConfig("../testdata/ptune", NewConfigEnv())
    config.InitAgent()
    config.ParseConfigFile(config.SvParamFile)
    tlsConfig, err := ConfigToTLSConfig(config)

```


ConfigToTLSConfig メソッド

```golang
func ConfigToTLSConfig(c *Config) (*tls.Config, error)
```

クライアント認証固定にしている

```golang
const ClientAuth = "RequireAndVerifyClientCert"
```

使ってなさそう

```bash
grep ConfigToTLSConfig `find . -name "*.go"`
./agent/https.go:func ConfigToTLSConfig(c *Config) (*tls.Config, error) {
./agent/https.go: server.TLSConfig, err = ConfigToTLSConfig(config)
./agent/https.go:       return ConfigToTLSConfig(config)
./agent/https_test.go:func TestConfigToTLSConfig(t *testing.T) {
./agent/https_test.go:  tlsConfig, err := ConfigToTLSConfig(config)
./agent/nettestmain.go:         server.TLSConfig, err = ConfigToTLSConfig(cfg)
./agent/nettestmain.go:                 return ConfigToTLSConfig(cfg)

```


