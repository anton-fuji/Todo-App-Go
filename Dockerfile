# Go 1.22のベースイメージを使用
# 注意 go.mod記載のversionと合わせる
# 軽量にするためalpineを使用
FROM golang:1.23.1-alpine

# コンテナ内の作業ディレクトリを/appに設定
WORKDIR /app

# ホストマシンからtodo-app内の全ファイルを、コンテナ内の作業ディレクトリ(/app)にコピー
COPY . ./

# go.modファイルに記載された依存パッケージをdownload
RUN go mod download

# main.goファイルをコンパイルして、実行ファイルを作成
RUN go build -o main /app/main.go

# コンテナの8080ポートを公開
EXPOSE 8080

# コンテナ起動後に、実行ファイルを実行 (JSON形式で指定)
CMD ["/app/main"]