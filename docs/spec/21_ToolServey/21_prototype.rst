簡易プロトタイプ作成
=====================

構成
----

* コンフィグレーション
* ファイル転送
   * SCP プロトコル。エンドポイントのマネージャに転送
* テンプレート処理
   * セッション確立、コマンドリスト実行、ファイル転送、セッションクローズ
* ロギング

ロギング
--------

* logrus, go-colorable 使用
* 簡易版とし、ターミナルにログ出力する

cd c:\home
git clone  https://github.com/getperf/getperf2.git

::

   package main

   import (
     "github.com/mattn/go-colorable"
     log "github.com/sirupsen/logrus"
   )

   func main() {
     log.SetFormatter(&log.TextFormatter{ForceColors: true})
     log.SetOutput(colorable.NewColorableStdout())

     log.Info("succeeded")
     log.Warn("not correct")
     log.Error("something error")
     log.Fatal("panic")
   }

::

   cd C:\home\getperf2\misc\prototype
   go run log.go

scp
---

* scp 転送先設定
   * ssh 鍵パス、ユーザ、パスワード
* 指定したパスのファイルをマネージャにscp転送する
* 事前にsshキー、公開鍵を作成して配布する

::

   ssh-keygen -t rsassh-keygen -t rsa

ssh\id_rsa をrsaキーにして公開鍵作成。マネージャ側に公開鍵登録

::

   vi ~/.ssh/authorized_keys


動作確認

::

   scp -i .\ssh\id_rsa scp.go psadmin@192.168.10.1:/tmp
   go run .\scp.go

事前にraa キーを作成する手間がかかる。 
http の場合、サーバ側で http サービスを上げておく必要がある。
簡易版プロトタイプでは scp 方式を採用する。


テンプレート
------------

* windowsinventory.go
* 構造体
  * Inventory
    * TestItems *TestItem[]
  * TestItem
    * Level InvetoryGetLevel
    * TestId
    * Script
* メソッド
  * NewTestItem(level, testId, script)
  * init() TestItems に登録
  * createGetInventoryScript(logDir, level) PowerShellテンプレート作成
  * test3() コマンド実行

Zip圧縮
-------

* mholt/archiver ライブラリ使用
* logDir下のディレクトリ圧縮
* logDir パスフォーマット {ノード}/{プラットフォーム}/{メトリック}
* {ノード}__{プラットフォーム}_{YYYYMMDDHHMI}.zip に圧縮

https://github.com/mholt/archiver


コンフィグレーション
---------------------

* toml ライブラリ使用
