# これは何か？
LINE経由でNotionにサイトの情報を保存しておくことができるLINEBotです。

# ビルド方法
```
GOOS=linux go build main.go
```

# lambdaにソースをアップする方法
appディレクトリ内のファイルを1つのzipファイルに圧縮してlambdaにアップします。

# 参考記事
- NotionAPIを使用するために必要になる「インテグレーション」について
    - [【NotionAPI 】インテグレーションが表示されない時の対処法](https://zenn.dev/syfut/articles/4906816e6e9118)
    - [私のインテグレーション](https://www.notion.so/my-integrations)

- データベースとNotionAPIのインテグレーションを接続する際に参考にした記事
    - [共有と権限設定](https://www.notion.so/ja-jp/help/sharing-and-permissions)

- linebotをgolangで書く際のドキュメント
    - [line/line-bot-sdk-go](https://github.com/line/line-bot-sdk-go)

# その他開発メモ
- パッケージをインストールする  
`go get -u "github.com/line/line-bot-sdk-go/linebot"`
