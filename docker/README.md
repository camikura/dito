# Oracle NoSQL Database Docker環境

このディレクトリには、開発・テスト用のOracle NoSQL Database環境が含まれています。

## 必要要件

- Docker
- Docker Compose
- mise (optional)

## 起動方法

### miseを使用する場合（推奨）

プロジェクトルートで以下のコマンドを実行：

```bash
# データベースを起動
mise run db-start

# ログを確認
mise run db-logs

# ステータス確認
mise run db-status

# データベースを停止
mise run db-stop

# データベースを再起動
mise run db-restart

# データベースを完全削除（データも削除）
mise run db-clean
```

### Docker Composeを直接使用する場合

```bash
# データベースを起動
docker compose -f docker/docker-compose.yml up -d

# ログを確認
docker compose -f docker/docker-compose.yml logs -f nosql

# データベースを停止
docker compose -f docker/docker-compose.yml down
```

## 接続情報

起動後、以下の情報で接続できます：

- **Endpoint**: `localhost`
- **Port**: `8080`
- **Protocol**: HTTP
- **Proxy URL**: `http://localhost:8080`

## ヘルスチェック

データベースが起動したか確認：

```bash
curl http://localhost:8080/
```

正常に起動していれば、レスポンスが返ってきます。

## データの永続化

データは `nosql-data` という名前付きボリュームに保存されます。コンテナを削除してもデータは保持されます。

データを完全に削除する場合：

```bash
mise run db-clean
# または
docker compose -f docker/docker-compose.yml down -v
```

## トラブルシューティング

### ポート8080が既に使用されている

`docker/docker-compose.yml` のポートマッピングを変更してください：

```yaml
ports:
  - "8081:8080"  # 例: ホスト側を8081に変更
```

### コンテナが起動しない

ログを確認：

```bash
mise run db-logs
```

### データベースをリセットしたい

```bash
mise run db-clean
mise run db-start
```

## 参考リンク

- [Oracle NoSQL Database Documentation](https://docs.oracle.com/en/database/other-databases/nosql-database/)
- [Oracle NoSQL Database Docker Hub](https://hub.docker.com/r/oracle/nosql)
