etcd 調査
=========

* etcd 基本機能のキーバリュー値のセットを試行
* クラスター構成のデプロイ機能の調査はまだ不十分
   * 他のソフトと組み合わせが必要な模様
   * 設定ファイルの配布/再ロード/変更イベントの検知

インストール
------------

ダウンロード

::

   ETCD_VER=v3.3.11

   DOWNLOAD_URL=https://storage.googleapis.com/etcd

   curl -L ${DOWNLOAD_URL}/${ETCD_VER}/etcd-${ETCD_VER}-linux-amd64.tar.gz -o /tmp/etcd-${ETCD_VER}-linux-amd64.tar.gz

/usr/local/binにetcd, etcdctlをコピー

::

   tar xzf /tmp/etcd-v3.3.11-linux-amd64.tar.gz -C /tmp/
   ls -l /tmp/etcd-v3.3.11-linux-amd64/
   sudo cp /tmp/etcd-v3.3.11-linux-amd64/etcd* /usr/local/bin/

   etcd --version
   etcdctl --version

シングルノードで起動

::

   mkdir /tmp/etcd-data-dir
   cd /tmp/etcd-data-dir/
   etcd

キーバリュー値セット

::

   etcdctl mk test_key "test value"
   etcdctl ls
   etcdctl get test_key
   etcdctl set test_key "new value"
   etcdctl get test_key
   etcdctl rm test_key
   etcdctl ls

ディレクトリ構成のキーバリュー値のセット

::

   etcdctl mkdir /hoge/fuga
   etcdctl mk /hoge/key2 "test value2"
   etcdctl mk /hoge/fuga/key3 "test value3"
   etcdctl ls -p --recursive
   etcdctl get /hoge/key2
   etcdctl get /hoge/fuga/key3

   etcdctl rm --recursive /hoge
   etcdctl ls -p --recursive

etcd クラスタ作成
-----------------

3台のノードでetcdを起動。事前に etcd を各ノードにインストール

1号機

::

   etcd \
    --name etcd1 \
    --initial-advertise-peer-urls http://192.168.0.17:2380 \
    --listen-peer-urls http://192.168.0.17:2380 \
    --listen-client-urls http://192.168.0.17:2379,http://127.0.0.1:2379 \
    --advertise-client-urls http://192.168.0.17:2379 \
    --initial-cluster etcd1=http://192.168.0.17:2380,etcd2=http://192.168.0.15:2380,etcd3=http://192.168.0.20:2380

2号機

::

   etcd \
    --name etcd2 \
    --initial-advertise-peer-urls http://192.168.0.15:2380 \
    --listen-peer-urls http://192.168.0.15:2380 \
    --listen-client-urls http://192.168.0.15:2379,http://127.0.0.1:2379 \
    --advertise-client-urls http://192.168.0.15:2379 \
    --initial-cluster etcd1=http://192.168.0.17:2380,etcd2=http://192.168.0.15:2380,etcd3=http://192.168.0.20:2380

3号機

::

   etcd \
    --name etcd3 \
    --initial-advertise-peer-urls http://192.168.0.20:2380 \
    --listen-peer-urls http://192.168.0.20:2380 \
    --listen-client-urls http://192.168.0.20:2379,http://127.0.0.1:2379 \
    --advertise-client-urls http://192.168.0.20:2379 \
    --initial-cluster etcd1=http://192.168.0.17:2380,etcd2=http://192.168.0.15:2380,etcd3=http://192.168.0.20:2380

クラスター構成でのキーバリュー値セット

::

   etcdctl mk key1 "test value1"
   etcdctl ls
   etcdctl get key1
