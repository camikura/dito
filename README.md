# dito

**dito（ディト）** は、Oracle NoSQL Database用のモダンなTUI（Text User Interface）クライアントです。

## 特徴

- 🎨 モダンでカラフルなターミナルUI
- ☁️ Oracle NoSQL Cloud ServiceとOn-Premise版の両方に対応
- ⚡ 高速で軽量なGoで実装
- 🔐 セキュアな接続管理

## 開発状況

現在、Phase 1（MVP）を開発中です：

- [x] 要件定義
- [x] UI設計
- [x] 開発環境構築
- [x] TUI Hello World
- [ ] 実装
  - [x] 接続機能（On-Premise）
    - [x] エディション選択画面
    - [x] 接続設定フォーム
    - [x] 接続テスト機能
  - [ ] テーブル一覧表示

## クイックスタート

### ビルド

```bash
# miseを使用する場合
mise run build

# または直接
go build -o dito cmd/dito/main.go
```

### 実行

```bash
# ビルド済みバイナリを実行
./dito

# または miseタスクで実行
mise run run

# 開発モード（ビルドせずに実行）
mise run dev
```

## 開発環境

### 必要要件

- Go 1.21以上
- Docker & Docker Compose
- mise (推奨)

### テスト用データベースの起動

```bash
# データベースを起動
mise run db-start

# ステータス確認
mise run db-status

# ログ確認
mise run db-logs
```

詳しくは [docker/README.md](docker/README.md) を参照してください。

## ドキュメント

- [要件仕様書](docs/REQUIREMENTS_SPEC.md)
- [UIモックアップ](docs/UI_MOCKUP.md)
- [Docker環境](docker/README.md)

## ライセンス

TBD
