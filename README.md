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
  - [x] テーブル一覧表示
    - [x] 2ペインレイアウト（左：テーブルリスト、右：テーブル詳細）
    - [x] 左ペイン固定幅（30文字）
    - [x] システムテーブルのフィルタリング
    - [x] 親子テーブルの関係表示
    - [x] アクティブなテーブルのハイライト表示（シアン + 太字）
    - [x] キーボードナビゲーション（j/k, ↑/↓）
    - [x] テーブル詳細表示
      - [x] カラム一覧（DDL解析による整形表示）
      - [x] インデックス一覧
      - [x] 親子関係の表示

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
# 1. データベースを起動（初回のみテーブル作成）
mise run db-start
mise run db-init

# 2. ditoを実行
./dito

# または miseタスクで実行
mise run run
```

### 使い方

1. **エディション選択**: `On-Premise` を選択して Enter
2. **接続設定**: デフォルト設定（`localhost:8080`）のまま `Connect` を選択して Enter
3. **テーブル一覧**: 接続成功後、テーブル一覧が表示されます
   - `j`/`k` または `↑`/`↓` でテーブルを選択
   - 右ペインに選択したテーブルの詳細が表示されます
   - `Enter` でデータ表示モードに切り替え
   - データ表示モードでは `j`/`k` でスクロール（最大1000行）
   - `Esc` で戻る、`q` で終了

## 開発環境

### 必要要件

- Go 1.21以上
- Docker & Docker Compose
- mise (推奨)

### テスト用データベースの起動

```bash
# データベースを起動
mise run db-start

# テスト用テーブルを作成
mise run db-init

# ステータス確認
mise run db-status

# ログ確認
mise run db-logs
```

詳しくは [docker/README.md](docker/README.md) を参照してください。

## ドキュメント

- [要件仕様書](docs/REQUIREMENTS_SPEC.md)
- [デザインポリシー](docs/DESIGN_POLICY.md)
- [UIモックアップ](docs/UI_MOCKUP.md)
- [Docker環境](docker/README.md)

## ライセンス

TBD
