Consul Exec調査
================

検証
----

前項で試した環境で、consul exec コマンドを調査

::

    consul members
   Node     Address            Status  Type    Build  Protocol  DC   Segment
   ostrich  192.168.10.1:8301  alive   server  1.7.1  2         dc1  <all>
   agent1   192.168.0.20:8301  alive   client  1.7.1  2         dc1  <default>
   agent2   192.168.0.17:8301  alive   client  1.7.1  2         dc1  <default>
   agent3   192.168.0.15:8301  alive   client  1.7.1  2         dc1  <default>


::

   consul exec uptime

エージェント側でエラーが発生

::

   agent: ignoring remote exec event, disabled.

disable_remote_execを false にしてエージェント再起動

::

   https://www.consul.io/docs/agent/options.html#disable_remote_exec


Consul 設定ファイルのパラメータで指定すると実行できる

::

   {
     "disable_remote_exec": false
   }

実行オプションの設定ではできなかった

::

   ./consul agent -config-dir . -node agent3 -bind 192.168.0.15 -enable-script-checks -disable-remote-exec false

Windows 試行
------------

ダウンロードサイトから Windows 32bit版をダウンロード

::

   https://www.consul.io/downloads.html

設定ファイルを c:\consul\conf 下に作成

::

   notepad++ c:¥consul¥conf¥consul.json

::

   {
       "server": false,
       "bootstrap": false,
       "disable_remote_exec": false,
       "data_dir": "c:\\consul\\data\\",
       "log_level": "INFO",
       "start_join": [
           "192.168.10.1"
       ],
       "node_name":"win16"
   }

単体で起動

::

   .\consul.exe agent -config-dir c:/consul/conf

::

   ./consul exec -node=win16 "powershell Write-Host hello"
   ./consul exec -node=win16 "powershell Get-WmiObject -Class Win32_Processor"

サービス起動

::

   sc.exe create "Consul" binPath= "c:\consul\conf\consul.exe agent -config-dir=c:\consul\conf" start= auto
   sc.exe start "Consul"


