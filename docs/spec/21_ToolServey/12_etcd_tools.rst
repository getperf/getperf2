etcd ツール調査
===============

* リファレンス調査
* 検証
* アーキテクチャ検討

リファレンス調査
----------------

Tool 調査
~~~~~~~~~

https://etcd.io/docs/v3.3.12/integrations/

flannel
~~~~~~~

* CoreOSのプロダクトである
* VXLAN ルーター、ルーティングテーブルは etcd で管理

go-etcd
~~~~~~~~

http://ttsubo.hatenablog.com/entry/2015/12/27/170005

etcd の Go ライブラリと API コードサンプル。
キーバリュー値監視のコードサンプル。チャネルを引数に指定

::

   client.Watch("/state", 0, true, etcdResponseChan, nil)

Readmeのサンプルを実行。connect error。etcd を個別に起動する必要がある。

CentOS、 etcd セットアップ事例
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

https://computingforgeeks.com/setup-etcd-cluster-on-centos-debian-ubuntu/

Kubernetes関連
~~~~~~~~~~~~~~

https://etcd.io/docs/v3.4.0/

https://kubernetes.io/docs/tasks/administer-cluster/configure-upgrade-etcd/#backing-up-an-etcd-cluster

Kubernetes とのインテグレーション事例が多く、
etcd と周辺ツールのみのデプロイ事例は少ない。

etcd データのバックアップ運用を推奨している。

VXLAN 管理などITインフラ全般の管理ツールとして利用されている。
比較的規模の大きなシステム向けで、あまりカジュアルに使える印象はない



