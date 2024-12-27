https://go.dev/dl/

sudo -E
wget https://go.dev/dl/go1.20.5.linux-amd64.tar.gz
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.20.5.linux-amd64.tar.gz

go version
go version go1.20.5 linux/amd64

go build cmd/getperf/getperf.go

agent/schedule.go:      schedule.WebServiceEnable = false
