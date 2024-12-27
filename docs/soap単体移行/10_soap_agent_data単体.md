# soap agent data 単体

ReserveSender(c, filename )
SendData(c, filename)
SendMessage(c, severity, message )
DownloadCertificate(c, timestamp )   
ReserveFileSender(c, onOff, waitSec )    
SendZipData(c, filename)
DownloadConfigFilePM(c, filename )

# ReserveSender

レグレッション


go test ./agent/ --run TestReserveSender  -v
response:Invarid role　→ OK になった

# SendData

レグレッション

go test ./agent/ --run TestSendData  -v

→ OK

go test ./agent/ --run TestSoapCallSendData  -v

=== RUN   TestSoapCallSendData
    soap_test.go:176: response Invarid role
--- PASS: TestSoapCallSendData (0.68s)

Invarid role になる

    static final String role = nvl(System.getProperty("GETPERF_WS_ROLE"), "admin");

; Admin management web service url.
URL_CM = https://192.168.133.128:57443/axis2/services/GetperfService

; Data management web service url.
URL_PM = https://192.168.133.128:58443/axis2/services/GetperfService


./bin/setenv.sh:GETPERF_WS_ROLE=data
./bin/setenv.sh:JAVA_OPTS="$JAVA_OPTS -DGETPERF_WS_ROLE=$GETPERF_WS_ROLE"

soap_test.go 修正

    soapSenderData, err = NewSoapSender("192.168.133.128", 58443)


# SendMessage(c, severity, message )

レグレッション

go test ./agent/ --run TestSendMessage  -v


# DownloadCertificate(c, timestamp )

レグレッション

go test ./agent/ --run TestDownloadCertificate  -v


# ReserveFileSender(c, onOff, waitSec )

web サービスコードにはなかったため保留とする。 ReserveSender のみ

# SendZipData(c, filename)

レグレッション

go test ./agent/ --run TestSendZipData  -v
これもなさそう fileSender のみ

# DownloadConfigFilePM(c, filename )

これもなさそう downloadCertificate のみ
