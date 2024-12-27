# Cacti エージェントのテストについて

以下でリグレッションテストを実行します。

    go test ./agent/ -v

./testdata/ptune-base-3.0/のテストデータを環境に合わせて作成する必要があります。

# Getconfig の WinRM テスト

WinRM ライブラリ [masterzen's Go WinRM](https://github.com/masterzen/winrm) の接続には NTLM 認証標準設定にBasic認証を有効化する設定が必要となります。
本制約から現在、ライブラリの使用を保留しており、テストコードの一部がエラーになります。
また、将来 WinRM 認証をバイパスする方式の移行を計画しており、合わせてテストも
修正する予定です。

# https の疎通テストについて

実際にWebサービスに接続してテストを行います。そのため、Getperf Webサービスの環境が必要です。

また、以下ディレクトリ下に、クライアント証明書ファイル一式を配布する必要があります。

    ./testdata/ptune-base-3.0/network

以下のssldadmin.pl コマンドで作成できます。

    ssladmin.pl client_cert --sitekey {サイト} --agent {ホスト}

実行すると以下のnetwork ディレクトリが生成され、本ディレクトリをコピーしてください。

    ls /etc/getperf/ssl/client/site1/{ホスト名}/network 

実行例：

    ssladmin.pl client_cert --sitekey site1 --agent alma8
    cp -r /etc/getperf/ssl/client/site1/alma8/network/* ./testdata/ptune-base-3.0/network/

https 疎通テストの関数名は TestSoapから始まりまるため、以下のキーワード指定で
https 疎通テストに絞り込んだテストができます。

go test --run TestSoap ./agent/ -v

# パラメータファイルの解析テストについて

パラメータファイルのパーサーの解析方法を CSV 形式に変更したため、従来のテスト関数
ではエラーが発生しています。

agetnt/param_test.go の TestParseWorkerLineCommand() に
変更したコードのテスト修正が必要です。

