# 国別IPv4アドレスリスト

## 概要
このプロジェクトでは、複数の地域インターネットレジストリ（RIR）が管理するIPアドレスデータを定期的に収集し、CIDR表記やサブネットマスク表記に変換して1つのファイルにまとめています。

## 詳細
RIR統計ファイルは、地域インターネットレジストリ（[Regional Internet Registry, RIR](https://en.wikipedia.org/wiki/Regional_Internet_registry)）によって提供されるデータファイルで、[こちら記載しているURL](https://github.com/inet-ip-info/WorldIPv4Map/blob/8981e2c07987fc15be3f005c008b4ec1b960a72b/main.go#L12-L16)で地域毎で公開されています。これらのファイルには、国別に割り当てられたIPアドレスの範囲と数が記載されており、それによって各国のインターネットリソース使用状況を確認できます。しかし、元のデータは「開始アドレスからのアドレス数」という形式で、直接的な使用には適していません。このプロジェクトでは、これらのデータをLinuxコマンドなどで扱いやすいCIDR表記（例: 192.168.0.0/24）やサブネットマスク表記（例: 192.168.0.0/255.255.255.0）に変換し、保存しています。

## 利用方法
以下のURLからデータをダウンロードできます。

### すべての国
- CIDR表記: [ダウンロード(all-ipv4cidr.tsv.gz)](https://github.com/inet-ip-info/WorldIPv4Map/releases/latest/download/all-ipv4cidr.tsv.gz)
- URL: https://github.com/inet-ip-info/WorldIPv4Map/releases/latest/download/all-ipv4cidr.tsv.gz
- サブネットマスク表記: [ダウンロード(all-ipv4mask.tsv.gz)](https://github.com/inet-ip-info/WorldIPv4Map/releases/latest/download/all-ipv4mask.tsv.gz)
  - URL:https://github.com/inet-ip-info/WorldIPv4Map/releases/latest/download/all-ipv4mask.tsv.gz

### 日本
- CIDR表記: [ダウンロード(jp-ipv4cidr.txt.gz)](https://github.com/inet-ip-info/WorldIPv4Map/releases/latest/download/jp-ipv4cidr.txt.gz)
  - URL: https://github.com/inet-ip-info/WorldIPv4Map/releases/latest/download/jp-ipv4cidr.txt.gz
- サブネットマスク表記: [ダウンロード(jp-ipv4mask.txt.gz)](https://github.com/inet-ip-info/WorldIPv4Map/releases/latest/download/jp-ipv4mask.txt.gz)
  - URL: https://github.com/inet-ip-info/WorldIPv4Map/releases/latest/download/jp-ipv4mask.txt.gz


## その他
http://nami.jp/ipv4bycc/ の情報を参考にさせて頂きました。
