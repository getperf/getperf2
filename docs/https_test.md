
https通信テスト
===============

1. 環境の準備
2. 再現テスト方法

疎通テスト
----------

centos80.getperf 192.168.0.5 を検証用環境とする

以下のコマンドで、検査用webサービスを起動する。

    go run ./cmd/nettest/nettest.go -u https://0.0.0.0:59001 -t testdata/tls/network/tls.yaml

出力メッセージ

    INFO[0000] url : http://0.0.0.0:59001, dir : testdata/webservice/backup, config : testdata/tls/network/tls.yaml

tls.yaml は以下を参照

    tlsConfig :
      # Certificate and key files for server to use to authenticate to client
      tlsCertPath : testdata/tls/network/server.crt
      tlsKeyPath : testdata/tls/network/server.key

      # Server policy for client authentication. Maps to ClientAuth Policies
      clientAuth : RequestClientCert

      # CA certificate for client certificate authentication to the server
      clientCAs : testdata/tls/network/ca.crt


クライアントから疎通テスト

    cd testdata/tlsn/network
    wget --no-proxy --output-document=/tmp/01.txt --ca-certificate=./ca.crt --certificate=./client.pem https://centos80.getperf:59001/files/


実行すると認証エラーが発生

    CA証明書 './ca.crt' をロードしました
    centos80.getperf (centos80.getperf) をDNSに問いあわせています... 192.168.0.5
    centos80.getperf (centos80.getperf)|192.168.0.5|:59001 に接続しています... 接続 しました。
    エラー: `centos80.getperf' の証明書は信用されません。
    エラー: `centos80.getperf' の証明書の発行者が不明です。


証明書チェックを無効にすると通信できる

    wget --no-proxy --no-check-certificate --output-document=/tmp/01.txt --ca-certificate=./ca.crt --certificate=./client.pem https://centos80.getperf:59001/files/

ca をルートではなく中間証明書にすると疎通できた


    wget --no-proxy  --output-document=/tmp/01.txt --ca-certificate=./inter/ca.crt --certificate=./client.pem https://centos80.getperf:59001/files/

IP指定だと、CNが一致しないエラー発生

    wget --no-proxy  --output-document=/tmp/01.txt --ca-certificate=./inter/ca.crt --certificate=./client.pem https://192.168.0.5:59001/files/

    証明書の所有者の名前とホスト名 `192.168.0.5' が一致しません

ca をルートではなく中間証明書にすると疎通できた。/etc/hosts にcentos80を登録

    wget --no-proxy  --output-document=/tmp/01.txt --ca-certificate=./inter/ca.crt --certificate=./client.pem https://centos80:59001/files/

NoClientCert以外を試す。RequireAndVerifyClientCertにすると、
client.pm を指定しない場合に送信失敗、再試行になる

まとめ

* サーバ証明書の CN は ドメイン名を除いたホスト名にする
  * 名前解決できない場合は /etc/hosts に IP アドレスを登録
* クライアント認証にはルート証明書ではなく、中間証明書を指定する必要がある
  （中間認証局がある場合）
* https.getTLSConfig(YAMLファイル) でパースしている
  * YAML ではなく、Config から読み込むようにする
    * getTLSConfigFromAgent(agent.Config)を追加

変更方針
--------

* クライアント
  * gconf get コマンド追加
    gconf get -o {output-dir} -f {https://agent:59001} -n {tls.toml}
  * 動作
    * GET /store で zip ファイルリスト検索
    * GET /zip/{file}.zip で順に zip ダウンロード
    * ダウンロードした zip を output-dir 下に解凍
  * 解凍先のパス変換はせずにそのまま解凍する。
    * 呼び出し側で変換をする
* サーバ
  * webservice.go 内変更
    * (c *Config)GetTSLConfig()( *https.TSLConfig) 追加
    * https/nettest.go の処理を移行

SSL　疎通検証
-------------

wget のラッパーではなく。クライアントツールを使う。gconf

    cd getperf2/testdata/tls/network
    go run client1.go
    2020/06/08 06:13:58 [{"keys":{"host":"hogehoge","stat_name":"Windows"},"ZipFile":"arc_hogehoge__Windows_20200520_1000.zip"}]

クライアントキーとクライアント証明書をペアにしたセットでも認証できる

    go run client1.go -cert ./client.pem -key ./client.pem

以下で疎通を確認。

    go run ./cmd/gconf/ get -f https://centos80:59001 \
    --ca ./testdata/tls/network/inter/ca.crt \
    --cert ./testdata/tls/network/client.pem

サーバ側変更
------------

getperf.ini 変更。既定は https で有効にする

    more ~/ptune/network/server/server.ini
    ; --------- Archive file shareing service --------------------------------
    ; The presence or absence of shareing service
    WEB_SERVICE_ENABLE = false

    ; Describe in url format. Set http / https, allowed IP, port
    ;WEB_SERVICE_URL = https://0.0.0.0:59443
    WEB_SERVICE_URL = http://0.0.0.0:59000

クライアント証明書で疎通確認

    ./gconf get -f https://centos80:59443 \
    --ca ./testdata/ptune/network/server/ca.crt \
    --cert ./testdata/ptune/network/client.pem
