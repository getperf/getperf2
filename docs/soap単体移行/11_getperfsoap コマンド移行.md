# getperfsoap コマンド移行

変更調査
プロトタイプ

# 変更調査

go ソース

getperf.go

    err := agent.Run(context.Background(), os.Args[1:], os.Stdout, os.Stderr)

agent.go

func (config *Config) Run() error {
    ctx, cancel := MakeContext(0)
    defer cancel()
    return config.RunWithContext(ctx)
}

getperf.go の中でサブコマンドの処理追加

    err := agent.Run(context.Background(), os.Args[1:], os.Stdout, os.Stderr)

# 既存の実行引数体系は変えずに実装できるか調査


agent/admin.go:func (c *Config) RunSetup() error {

agent/getperf2.go:func Run(ctx context.Context, argv []string, stdout, stderr io
.Writer) error {

サブコマンドなしなら、agent.Run を呼び出す

サブコマンドありがら、
    send    agent.RunSender を呼ぶ
    setup   agetn.RunSetup を呼ぶ

make build だと時間が掛かる

time go build  ./cmd/getperf

real    0m0.799s
user    0m0.035s
sys     0m0.439s

# soap_file_exchange.go 実装

```go
func RunSender(ctx context.Context, argv []string, stdout, stderr io.Writer) error {
    var (
        c = flag.String("config", "", "It performs by the specified directory.")
    )
    flag.StringVar(c, "c", "", "")
    flag.Parse()
    log.Info("run sender proto")

    return nil
}
```

./getperf sender a
Args : sender
INFO[0000] run sender proto




