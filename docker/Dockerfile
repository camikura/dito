# Oracle NoSQL Database Community Edition
FROM ghcr.io/oracle/nosql:latest-ce

# メタデータ
LABEL maintainer="nosql-client"
LABEL description="Oracle NoSQL Database for development and testing"

# 環境変数の設定
ENV KV_PROXY_PORT=8080

# ヘルスチェック用のcurlをインストール（必要に応じて）
USER root
RUN microdnf install -y curl && microdnf clean all

# オリジナルのユーザーに戻す
USER oracle

# ポート公開
EXPOSE 8080

# データディレクトリ
VOLUME /kvroot

# 起動コマンド（ベースイメージのENTRYPOINTを使用）
