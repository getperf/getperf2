Goファイル転送
==============

リファレンス調査
----------------

p2p ファイル転送。ftrans

https://qiita.com/Tsuzu/items/d6bef7e531c87168e517

tuggle

https://github.com/fujiwara/tuggle

go製 curl コマンド

https://github.com/astaxie/bat

httpie

https://github.com/jakubroztocil/httpie

grpcurl

https://github.com/fullstorydev/grpcurl

Tuggule
-------

1ファイル tuggle.go に 38 メソッド。6 構造体。
解読に困難。
コンサルクラスター内で各ノードにファイルを配布する仕組み。
ファイル DELETE,PUT,LIST APIを提供。
サービス起動して curl コマンドで受信。
consul/api で、API からサービス、ウォッチ、イベントの設定をしている

Consul 立ち上げ

consul agent -server -data-dir /tmp/consul -bootstrap-expect 1

Tuggle 立ち上げ。Consul エージェントに接続し、8900ポートでリッスン

tuggle -data-dir /tmp/tuggle

ファイルPUT。各ノードの/tmp/tuggleに保存される

curl -XPUT -H"Content-Type: application/gzip" --data-binary @test.gz localhost:8900/test.gz

LIST検索

curl -s localhost:8900 | jq .
 

bat
---

Beego Framework MVCタイプのフレームワーク

ファイルサイズは 8MB弱

::

   ls -l ~/go/bin/bat
   7985816  3月 23 06:09 2020 /home/psadmin/go/bin/bat

使用法

::

   bat PUT localhost X-API-Token:123 name=John -p

アップロード

::

   bat localhost < /tmp/zabbix_agentd.log

Download a file and save it via redirected output:

::

   bat example.org/file > file

grpc
----

ファイルサイズ 20MB強

::

    ls -l ~/go/bin/grpcurl
   22482715  3月 23 06:27 2020 /home/psadmin/go/bin/grpcurl

grpcurl -help
