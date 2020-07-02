zabbix-agent2 調査
==================

リファレンス調査
----------------

https://www.zabbix.com/whats_new_4_4

プラグインフレームワーク。
	エクスポーター
	コレクター
	ランナー
	ウォッチャー
チェック間で状態を維持する機能（たとえば、永続的なDB接続の維持）。
柔軟な時間間隔をサポートする組み込みスケジューラー。
バルクデータ転送を使用した効率的なネットワーク使用。

https://tech-lab.sios.jp/archives/18664

Zabbix-Agent2 は Zabbix4.4 では RHEL8 / CentOS8 以降にしか対応していません
 Zabbix5.0 LTS 
 2020年4月リリース予定の Zabbix5.0
 
https://assets.zabbix.com/img/zabconf2019_jp/presentations/10_zabconf2019.pdf

Alexey Pustovalov氏　Zabbix LLC 
Zabbix 5.0より正式サポート予定
• 現時点ではLinux系のOSのみに対応
• Windows向けZabbixエージェント2 ‒ 開発中!
• ドキュメントも近日中にGit上に公開予定

検証
----

git clone https://github.com/zabbix/zabbix.git

./configure --enable-agent2

make

171ファイル

   385   1296  11326 ./internal/agent/serverconnector/serverconnector.go
    68    242   1626 ./internal/agent/serverconnector/activeconnection.go
 22310  78845 592070 合計

モジュール： cmd, conf, internal, pkg, plugins


 