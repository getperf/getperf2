Consul ツール調査
=================

Consul Demo
-----------

MongoDB 3台, Django 2台, Fabio(ロードバランサー) 1台の構築例

::

   https://medium.com/velotio-perspectives/a-practical-guide-to-hashicorp-consul-part-2-3c0ebc0351e8

上記サンプル。 Docker Compose で立ち上げる 

::

   git clone https://github.com/pranavcode/consul-demo

twitter アプリ

Stretcher pull型デプロイツール
------------------------------

examples の下のデモを試す。

事前準備
~~~~~~~~

Consul 起動

::

   consul agent -server -data-dir /tmp/consul -bootstrap-expect 1

サンプル main.go コンパイル。 GO111MODULE をオフにしておく。
Go コンパイル、example.tar.gz にアーカイブ、example.yml　作成

::

   export GO111MODULE=off
   ./prepare.sh

prepare.sh 内で Consul デプロイイベント "example_deploy" を登録している

::

   consul event -name "example_deploy" "file://${CWD}/example.yml"

Stretcher 実行
~~~~~~~~~~~~~~

::

   ./exec.sh

exec.sh 内で Consul イベントを発行。stretcher を起動している。

::

   consul watch -type event -name example_deploy stretcher

Consul Agent 側での処理。ログから捕捉

* manifest.yml 受信
* 前処理 sleep 3
* ファイルロード file://.../example.tar.gz
* チェックサム
* tar 解凍
* ファイル同期。rsync 使用。元：/tmp/xxx、先：/.../deployed/
* 後処理、コマンド実行

rsync 実行オプション

* exclude オプションで除外指定、
* --delete でコピー元 にないファイルは削除

::

   [-a --delete -v --exclude *.go --exclude Makefile 
   	/tmp/stretcher_src903720934/ 
   	/home/psadmin/work/go/consul/stretcher_test1/deployed/]

ファイル同期はストラテジーパターンで、rsync 以外も選択できる模様

UML作成,コードレビュー
~~~~~~~~~~~~~~~~~~~~~~

go ソースは 6ファイル、787行と小さく、読みやすい。

外部コマンドで、rsync, tar を使用していることから、Windows には対応していない模様。


