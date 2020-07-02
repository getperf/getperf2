Consul イベント調査
===================

* リファレンス調査
* 検証

リファレンス調査
----------------

Stretcher。 Consul 周辺ツール。デプロイツール

::

   https://github.com/fujiwara/stretcher

Consul Event説明

::

   https://www.consul.io/docs/commands/event.html

イベントのペイロードは非常に小さくする必要がある(100バイト未満)
指定するイベントが大きすぎると、エラーが返されます

Consul Examples

::

   https://github.com/JoergM/consul-examples

検証
-----

Consul Examplesを試す

::

   git clone https://github.com/JoergM/consul-examples

基本例

* Simple DNS Setup
* Using the http-api
* Service Checks
* Locks, Watches and Events

複雑な例。 Consul-Template と HAProxy を使用したサービスディスカバリー

::

   https://github.com/JoergM/consul-examples/tree/master/showcase

DNS
~~~~

::

   consul agent -config-file consul.json
   dig @127.0.0.1 -p 8600 moi.node.consul
   dig @127.0.0.1 -p 8600 consul.service.consul
   dig @127.0.0.1 -p 8600 consul.service.consul SRV


http-api
~~~~~~~~

consul.json にサービス追加

::

   "service": {
       "name": "example",
       "tags": ["special"],
       "address":"127.0.0.1",
       "port":3456
   }

curl で APIアクセス

curl localhost:8500/v1/catalog/service/example | jq "."

サービス追加と再ロード

::

   vi service.json

::

   {
     "ID": "example2",
     "Name": "example2",
     "Tags": [
       "specific",
       "v1"
     ],
     "Address": "127.0.0.1",
     "Port": 8000
   }

::

   curl -X PUT -d @service.json localhost:8500/v1/agent/service/register
   curl localhost:8500/v1/agent/services

キーバリュー値セット

::

   curl -X PUT -d @- localhost:8500/v1/kv/example <<< FooValue
   curl -s localhost:8500/v1/kv/example?raw

サービスチェック
~~~~~~~~~~~~~~~~

チェックスクリプトを カレントにコピー

::

   cp consul-examples/checks/check_* .

::

   vi consul.json

オリジナルの記述が古く、以下例の通り PARGS 形式で記述する必要がある

::

      "checks": [
        {
          "args": ["./check_green.sh"],
          "interval": "15s"
        }

APIでの確認

::

   # サービスチェック
   curl localhost:8500/v1/health/service/example_green
   curl localhost:8500/v1/health/service/example_falied
   # gron を使った検索
   curl localhost:8500/v1/health/service/example_green|gron|grep Status
   curl localhost:8500/v1/health/service/example_failed|gron|grep Status

イベント
~~~~~~~~

ロック動作。複数の実行で並列度を指定して外部コマンドを実行できる

端末A

::

   consul lock LockA sleep 15

端末B

::

   consul lock LockA echo released

5並列に制限して実行する場合

ウォッチ
~~~~~~~~

変化があった場合にコマンドを実行する

キー値セット

::

   curl -X PUT -d @- localhost:8500/v1/kv/application/online <<< true


consul watch 実行

::

   consul watch -type keyprefix application/online cat 

以下のエラーがでるため、-type key 指定で監視

::

   consul watch -type keyprefix application/online
   Must specify a single prefix to watch

::

   consul watch -type key -key application/online cat

別端末で実行

::

   curl -X PUT -d @- localhost:8500/v1/kv/application/online <<< false

consul watch の表示

::

   {"Key":"application/online","CreateIndex":1149,"ModifyIndex":1149,"LockIndex":0,"Flags":0,"Value":"dHJ1ZQ==","Session":""}
   {"Key":"application/online","CreateIndex":1149,"ModifyIndex":1182,"LockIndex":0,"Flags":0,"Value":"ZmFsc2U=","Session":""}

イベント
~~~~~~~~

イベント監視

::

   consul watch -type event -name exampleEvent cat

別端末で実行

::

   consul event -name exampleEvent "Event Payload"

イベントIDが出力される

::

   [{"ID":"89914255-8e59-b392-884b-c38d5365717d","Name":"exampleEvent","Payload":"RXZlbnQgUGF5bG9hZA==","NodeFilter":"","ServiceFilter":"","TagFilter":"","Version":1,"LTime":10}]

ショーケース(複雑な例)
~~~~~~~~~~~~~~~~~~~~~~

前回調査した以下事例に類似した設定例。
haproxy を用いてロードバランサーを構成している。
自動構成ツールは、Puppetを使用。 
Vagrant 1.6.3, VirtualBox 4.3 で検証環境構築。
古いので動かない可能性あり。

::

   https://kazuhira-r.hatenablog.com/entry/20170611/1497189791
