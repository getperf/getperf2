Consul テンプレート
===================

リファレンス
------------

さくらインターネット前佛さんの記事

::

   # Gihyo
   https://gihyo.jp/admin/feature/01/serf-consul/0008

   # Slide share
   https://www.slideshare.net/zembutsu/consul-template-beat-tanpopo-on-sashimi

* ヘルスチェック時、ファイル更新時、キーバリューの値が変化した場合にトリガー発火
* テンプレートを基に設定ファイルを自動生成。エージェントに配布
* 後処理のコマンドを自動実行(サービス再起動など)

Consul Template チュートリアル

::

   https://learn.hashicorp.com/consul/developer-configuration/consul-template

試行
----

Consul Template チュートリアルを基に試行

インストール

::

   wget https://releases.hashicorp.com/consul-template/0.24.1/consul-template_0.24.1_linux_amd64.tgz

Consul KV ユースケース
~~~~~~~~~~~~~~~~~~~~~~

テンプレートファイル作成

::

   vi find_address.tpl
   {{ key "/hashicorp/street_address" }}

事前に consul を devモードで起動

::

   ./consul agent -dev

consul-template を起動

::

   consul-template -template "find_address.tpl:hashicorp_address.txt"

::

   consul kv put hashicorp/street_address "102 2nd St"

consul-templateを実行した同じディレクトリ下に、hashicorp_address.txt ができる

キーバリュー値を変えると動的に、.txt ファイルも更新される

以下のWARNが発生する。Vaultとの連携が必要か要調査

::

   vault.token: failed to renew: 
   Put https://127.0.0.1:8200/v1/auth/token/renew-self: 
   dial tcp 127.0.0.1:8200: connect: connection refused

サービスディスカバリー ユースケース
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

テンプレートファイル作成

::

   vi all-services.tpl
   {{range services}}# {{.Name}}{{range service .Name}}
   {{.Address}}{{end}}

   {{end}}

consul-template起動

::

   consul-template -template="all-services.tpl:all-services.txt" -once

他のユースケースの試行を進める

::

   # nginx のロードバランサー構成設定事例
   https://kazuhira-r.hatenablog.com/entry/20170611/1497189791

