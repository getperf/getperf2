Consul 周辺ツール調査
====

リファレンス調査
----------------

デプロイの検証

https://gendosu.jp/archives/2840

Stretcher 開発者の記事。一番詳しい

https://techblog.kayac.com/10_stretcher.html

capistranoと連携したアプリデプロイ事例

https://tech.drecom.co.jp/examples-consul-and-stretcher-push-to-pull/

図による構成説明。判りやすい

http://tech.feedforce.jp/stretcher-consul-capistrano.html

Stretcherチュートリアル。Railsアプリのデプロイ。※ 実機検証計画

https://cross-black777.hatenablog.com/entry/2016/01/23/190211

Counsul 運用系機能色々

https://github.com/hashicorp-japan/consul-workshop/blob/master/contents/utilities.md

A Practical Guide to HashiCorp Consul — Part 1。モノリスサーバ
、マイクロサービス、サービスメッシュの説明

https://medium.com/velotio-perspectives/a-practical-guide-to-hashicorp-consul-part-1-5ee778a7fcf4

MongoDB 3台, Django 2台, Fabio(ロードバランサー) 1台の構築例

https://medium.com/velotio-perspectives/a-practical-guide-to-hashicorp-consul-part-2-3c0ebc0351e8

HashiCorp Consul Demo。上記サンプル。 Docker Compose で立ち上げる
※ 実機検証計画

https://github.com/pranavcode/consul-demo

使用法概要
----------

マニュフェスト記入例

s3://example/manifest.yml

::

   commands:
     pre:
       - echo 'staring deploy'
       - echo `hostname`
       - curl -X PUT -d 'deploy' http://localhost:8500/v1/kv/deploy/{`hostname`}
     post:
       - echo 'deploy done'
       - 'supervisorctl restart nginx'
       - curl -X DELETE http://localhost:8500/v1/kv/deploy/{`hostname`}
     success:
       - echo 'deploy success'
     failure:
       - echo 'deploy failed!!'
       - cat >> /root/failure.log
       - curl -X PUT -d 'failure' http://localhost:8500/v1/kv/deployFailure/{`hostname`}

Consuleイベントに Strecher 登録

::

   consul watch -type event -name deploy /path/to/stretcher

イベントの発行。manifest.yml を送る

::

   consul event -name deploy s3://example/manifest.yml

tarとrsyncが必要。Windows Server だと？

その他のツール
--------------

fireap

https://keyamb.hatenablog.com/entry/fireap-release-v0.1.0

tuggle

curl -XPUT -H"Content-Type: application/gzip" --data-binary @test.gz localhost:8900/test.gz

Stretcher 調査
--------------

Conusl 試行の環境を tmux 上で起動

192.168.0.70 で tmux 起動

::

   # Server
    ostrich  192.168.10.1               # tmux #0
   # spring boot consul agent 3台
    192.168.0.15 redmine CentOS 6.10   # #1:java, #2:agent1
    192.168.0.17 getperf CentOS 6.9    # #3:java, #4:agent2
    192.168.0.20 centos7 CentOS 7.3.1611 # #5:java, #6:agent3

ostrich のtmux #1

::

   # nginx consul template
    192.168.0.70 centos75 CentOS 7.5   
    
    consul members
   Node     Address            Status  Type    Build  Protocol  DC   Segment
   getperf  192.168.10.1:8301  alive   server  1.7.1  2         dc1  <all>
   agent1   192.168.0.20:8301  alive   client  1.7.1  2         dc1  <default>
   agent2   192.168.0.17:8301  alive   client  1.7.1  2         dc1  <default>
   agent3   192.168.0.15:8301  alive   client  1.7.1  2         dc1  <default>
   win16    192.168.0.27:8301  alive   client  1.7.1  2         dc1  <default>

