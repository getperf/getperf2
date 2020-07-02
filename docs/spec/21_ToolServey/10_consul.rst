Consul と周辺ツール調査
=======================

リファレンス調査
----------------

Strecher
~~~~~~~~

::

   https://tech.drecom.co.jp/examples-consul-and-stretcher-push-to-pull/
   https://github.com/fujiwara/stretcher
   https://cloudpack.media/12992

解決しようとしている課題

* push型デプロイはホスト台数が増減しやすい環境に適さない
* 各種問題を解決するpull型デプロイツールを書いた。Consul と連携

大まかな流れ

* 配布物を全て一つのアーカイブにまとめて、Amazon S3など(http, fileも可)に保存する
* デプロイ手順を定義した manifest (YAML) を作成し、S3などに保存する
* consul event を利用して各ホストにイベント通知を行う
* 各ホストで動作している stretcher agent が、イベントを受けて配布物を取得しデプロイを実行する

Consul
~~~~~~

::

   # Terraformを使ってConsul導入
   http://urapico.hatenablog.com/entry/gp-advent-calendar-2015-18

   # Gihyo 記事。第4回〜8回が Consul 記事
   https://gihyo.jp/admin/feature/01/serf-consul/0004

   # チュートリアル。Slideshare
   https://www.slideshare.net/ssuser07ce9c/consul-58146464

   # Thinkit 記事
   https://thinkit.co.jp/story/2015/08/24/6343

   # Windows でのリモート操作
   https://qiita.com/gamisan9999/items/afeddf024948057faf8a

   # CentOS 1台でインストール、チュートリアル。比較的新しい
   http://kb-instep.com/consul_setup/

* 使用例

   * consul server 起動。クラスターには 3台以上必要
   
   ::

      nohub consul agent -server -bootstrap-expect=3 -data-dir=/var/lib/consul -advertise=192.168.20.21 &
      nohub consul agent -server -bootstrap-expect=3 -data-dir=/var/lib/consul -advertise=192.168.20.22 -join=node01.example.jp &
      nohub consul agent -server -bootstrap-expect=3 -data-dir=/var/lib/consul -advertise=192.168.20.23 -join=node01.example.jp &

   * エージェント起動。台数が増減するものはクライアント(エージェント)で立ち上げる

   ::

      nohub consul agent -data-dir=/var/lib/consul -advertise=192.168.20.41 -join node01.example.jp &
      nohub consul agent -data-dir=/var/lib/consul -advertise=192.168.20.42 -join node01.example.jp &

   * リモート実行

   ::

      consul exec uptime

主な機能

* サービスディスカバリー、HA構成、ラウンドロビン DNS
* ヘルスチェック、イベント受信、アクション
* Key/Value Store、永続データの管理
* 複数DC対応

開発元チュートリアルを今後試行する

::

   https://learn.hashicorp.com/consul

Consul 試行
-----------

インストール

::

   wget https://releases.hashicorp.com/consul/1.7.1/consul_1.7.1_linux_amd64.zip
   unzip consul_1.6.0_linux_amd64.zip
   sudo cp consul /usr/local/bin/
 
設定ファイル編集

::

   mkdir /tmp/consul
   mkdir -p /work/go/consul
   cd /work/go/consul
   vi consul.json
   {
     "server": true,
     "bootstrap_expect": 1,
     "client_addr": "0.0.0.0",
     "bind_addr": "192.168.91.128",
     "node_name": "moi",
     "datacenter": "dc1",
     "data_dir": "/tmp/consul",
     "log_level": "INFO",
     "enable_script_checks": true,
     "enable_syslog": true,
     "check_update_interval": "10s",
     "ui": true,
     "start_join": [
           "192.168.91.128"
     ]
   }

単体での起動 

::

   consul agent -config-dir /home/psadmin/work/go/consul

http://moi:8500

インターフェースと使用するポート

* consulエージェント（RPC）	consulコマンドで参照（メンバー情報のみ）	8100(TCP)
* HTTPインターフェース	curlやwget等で参照	8500(TCP)
* DNSインターフェース	hostやdigなどホスト名の名前解決時に参照	8600(UDP)


メンバーの確認

::

   consul members

ノード情報の取得

::

   curl http://127.0.0.1:8500/v1/catalog/nodes | jq "."

DNS の利用

::

   sudo -E yum -y install bind-utils
   dig @127.0.0.1 -p 8600 moi.node.consul

サービス設定

::

   vi web.json
   {
       "service": {
           "name": "web",
           "tags": ["test"],
           "port": 80
       }
   }

consul 再起動

::

   consul agent -config-dir /home/psadmin/work/go/consul

状態チェック。うまく動作していない模様。継続調査

::

   curl http://127.0.0.1:8500/v1/catalog/service/web | jq "."


DNS確認

::

   dig @127.0.0.1 -p 8600 moi.node.consul
   # A       192.168.91.128
   dig @127.0.0.1 -p 8600 web.service.consul
   # 同様に Aレコードに IP記載
   dig @127.0.0.1 -p 8600 web.service.consul SRV
   # SRV     1 1 80 moi.node.dc1.consul.
   dig @127.0.0.1 -p 8600 test.web.service.consul
   # タグを追加した検索

サービス設定更新

::

   echo '{"service":
     {"name": "web",
       "tags": ["test"],
       "port": 80,
       "check": {
         "args": ["curl", "localhost"],
         "interval": "10s"
       }
     }
   }' > ./web.json

::

   consul reload

Apache を起動／停止してみて、Consul ログを確認

::

   sudo /etc/init.d/httpd stop
   sudo /etc/init.d/httpd start

ヘルスチェックが成功した場合のみサービスドメイン名のAレコードが返る

::

   dig @127.0.0.1 -p 8600 web.service.consul

キーバリューデータストア試行。consul agent を起動している状態で実行

::

   consul kv put redis/config/minconns 1
   consul kv put redis/config/maxconns 25
   consul kv put -flags=42 redis/config/users/admin abcd1234

値の検索

::

   consul kv get redis/config/minconns
   consul kv get -detailed redis/config/users/admin

値の削除

::

   consul kv delete redis/config/minconns
   consul kv delete -recurse redis

UI操作

::

   http://moi:8500/ui

* ACL

   - アクセス制御リスト（ACL）を使用
   - UI、API、CLI、サービス通信、およびエージェント通信を保護

Consul Connect, Intensionについては保留

socat (Multipurpose relay (SOcket CAT)) とは、簡単にソケット通信をサーバー側としてもクライアント側
としても 使えるコマンド
プロキシーなどに用いる

::

   sudo -E yum install --enablerepo=epel socat

複数ノードでの利用チュートリアル

ノード1,2 で起動

以下に日本語のチュートリアル実行ログあり。
Connect 以外は概要の把握ができた

::

   https://kakakakakku.hatenablog.com/entry/2019/03/22/113847


* リモートでコマンドを実行できる。acl 設定が必要
* サービスチェックはインターバル指定で指定したスクリプトを定期実行する
* 判定結果で DNS レコードの登録、失効を行う
* 設定ファイルを更新して consul reload で設定を反映できる
* 設定フィルの書き方、実行オプションの指定方法などの操作性がとても良い
   * ソースコードを参考する
* consul-template を使うと管理ノードの増減に合わせて設定ファイルを自動更新できる
* Connect はまだ用途不明。プロキシー環境で使用する？

