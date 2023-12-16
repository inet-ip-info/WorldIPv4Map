# 国別IPv4アドレスリスト

## 概要
このプロジェクトでは、複数の地域インターネットレジストリ（RIR）が管理するIPアドレスデータを定期的に収集し、CIDR表記やサブネットマスク表記に変換して1つのファイルにまとめています。

## 詳細
RIR統計ファイルは、地域インターネットレジストリ（[Regional Internet Registry, RIR](https://en.wikipedia.org/wiki/Regional_Internet_registry)）によって提供されるデータファイルで地域毎に公開されています。
具体的には[こちら]((https://github.com/inet-ip-info/WorldIPv4Map/blob/8981e2c07987fc15be3f005c008b4ec1b960a72b/main.go#L12-L16)に記載しているURLの一覧から取得可能です。



[こちら記載しているURL](https://github.com/inet-ip-info/WorldIPv4Map/blob/8981e2c07987fc15be3f005c008b4ec1b960a72b/main.go#L12-L16)で地域毎で公開されています。

これらのファイルには、国別に割り当てられたIPアドレスの範囲と数が記載されており、それによって各国のインターネットリソース使用状況を確認できます。しかし、元のデータは「開始アドレスからのアドレス数」という形式で、直接的な使用には適していません。このプロジェクトでは、これらのデータをLinuxコマンドなどで扱いやすいCIDR表記（例: 192.168.0.0/24）やサブネットマスク表記（例: 192.168.0.0/255.255.255.0）に変換し、保存しています。

## 利用方法
以下のURLからデータをダウンロードできます。

### すべての国
- CIDR表記: [ダウンロード(all-ipv4cidr.tsv.gz)](https://github.com/inet-ip-info/WorldIPv4Map/releases/latest/download/all-ipv4cidr.tsv.gz)
- サブネットマスク表記: [ダウンロード(all-ipv4mask.tsv.gz)](https://github.com/inet-ip-info/WorldIPv4Map/releases/latest/download/all-ipv4mask.tsv.gz)

### 日本
- CIDR表記: [ダウンロード(jp-ipv4cidr.txt.gz)](https://github.com/inet-ip-info/WorldIPv4Map/releases/latest/download/jp-ipv4cidr.txt.gz)
- サブネットマスク表記: [ダウンロード(jp-ipv4mask.txt.gz)](https://github.com/inet-ip-info/WorldIPv4Map/releases/latest/download/jp-ipv4mask.txt.gz)

## Forkの歓迎
このリポジトリは、皆さんの貢献と協力を歓迎します。プロジェクトをより良くするため、また、この重要な情報リソースの継続的な更新と保守を確実にするため、Forkしての参加を奨励しています。多くの方々がこのリストを独自に管理・更新することで、データの正確性とアップデートの速度が向上し、コミュニティ全体に利益をもたらすことができると考えています。自分だけではなく、皆さんの手によってこのプロジェクトが継続し続けることを願っています。


## 参考サイト
このコードを書く際に、[世界の国別 IPv4 アドレス割り当てリスト](http://nami.jp/ipv4bycc/)を大いに参考にさせていただきました。このサイトでは、IPアドレスリストの変換仕様が詳しく記載されており、コーディングの大きなヒントになりました。

このサイトの作者様には、アイデアを授けてくださったことに心より感謝しています。