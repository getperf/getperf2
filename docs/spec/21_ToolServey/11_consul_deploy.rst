Consul template によるアプリデプロイ方法調査
============================================

リファレンス
------------

::

   https://kazuhira-r.hatenablog.com/entry/20170611/1497189791

環境
----

::

   # Consul Server 1台
    ostrich  192.168.10.1
   # spring boot consul agent 3台
    192.168.0.15 redmine CentOS 6.10
    192.168.0.17 getperf CentOS 6.9
    192.168.0.20 centos7 CentOS 7.3.1611
   # nginx consul template
    192.168.0.70 centos75 CentOS 7.5

Spring Boot CLIでSpring Bootアプリケーションを作成

Consul Server
-------------

::

   sdk install springboot
   spring --version

   vi server.groovy
   # 記事のソース編集
   spring jar server.jar server.groovy

暫くすると、server.jar ができる

同ノードで Consul Server 起動

::

   consul agent -bootstrap -server -client=`hostname -i` \
   -bind=192.168.10.1 \
   -enable-script-checks -config-dir=./consul.d \
   -data-dir=/home/psadmin/work/go/consul/template

Consul Agent
------------

jar ファイルを配布

::

   scp /usr/local/bin/consul psadmin@192.168.0.15:~/work/consul/
   scp /usr/local/bin/consul psadmin@192.168.0.17:~/work/consul/
   scp /usr/local/bin/consul psadmin@192.168.0.20:~/work/consul/

各ノードで consol.json 作成

::

   vi ~/work/consul/client.json

   {
     "data_dir": "/home/psadmin/work/consul",
     "start_join": [
         "192.168.10.1"
     ]
   }

サービスの定義

::

   vi ~/work/consul/service-app.json

   {
     "service": {
       "name": "appbackend",
       "tags": ["app.backend"],
       "address": "",
       "port": 8080,
       "checks": [
         {
           "args": ["curl","http://localhost:8080/health"],
           "interval": "10s"
         }
       ]
     }
   }

各ノード起動

::

   ./consul agent -config-dir . -node agent1 -bind 192.168.0.20 -enable-script-checks
   ./consul agent -config-dir . -node agent2 -bind 192.168.0.17 -enable-script-checks
   ./consul agent -config-dir . -node agent3 -bind 192.168.0.15 -enable-script-checks

Spring Bootアプリケーションは最初の1台のみ起動し、残りはあとから順次追加

::

   scp server.jar psadmin@192.168.0.20:~/work/consul/
   scp server.jar psadmin@192.168.0.17:~/work/consul/
   scp server.jar psadmin@192.168.0.15:~/work/consul/

agent1 のみ spring 起動

::

   java server.jar

nginx with Consul Template
--------------------------

nginx インストール

::

   sudo vi /etc/yum.repos.d/nginx.repo
   [nginx]
   name=nginx repo
   baseurl=http://nginx.org/packages/centos/7/$basearch/
   gpgcheck=0
   enabled=1

::

   sudo yum install nginx
   sudo systemctl enable nginx
   sudo systemctl start nginx

ブラウザから起動確認

::

   http://192.168.0.70/

consul-template ダウンロード

::

   wget https://releases.hashicorp.com/consul-template/0.24.1/consul-template_0.24.1_linux_amd64.tgz
   tar xvf consul-template_0.24.1_linux_amd64.tgz


nginx テンプレートの適用例を参考にする

::

   https://github.com/hashicorp/consul-template/blob/master/examples/nginx.md

nginx テンプレートファイル作成

::

   vi nginx.ctmpl

   upstream appbackend {
     {{ range service "appbackend" }} {{ $name := .Name }} {{ $service := service .Name }}
     server {{ .Address }}:{{ .Port }};
     {{ end }}
   }

   server {
       listen 8080;
       server_name localhost;
       location / {
           proxy_pass http://appbackend;
       }
   }

consul-template 起動

::

   sudo ./consul-template -consul-addr 192.168.10.1:8500 \
   -template "nginx.ctmpl:/etc/nginx/conf.d/balancer.conf:service nginx reload"

SELinux が有効になっている場合は以下で許可設定

::

   setsebool -P httpd_can_network_connect 1

各ノードでラウンドロビンされて、出力結果が毎回変わる

::

   curl http://192.168.0.70:8080/hello

nginx 設定ファイルが更新される

::

   more /etc/nginx/conf.d/balancer.conf
