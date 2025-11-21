# UI Mockup

## 実装済みUI

### ステップ1: エディション選択（全画面モード・マージン追加版）

```
╭─────────────────────────────────────────────────────────────────╮
│  dito - Oracle NoSQL Database TUI Client                        │
│                                                                 │
│  Select Connection                                              │
│   > Oracle NoSQL Cloud Service                                  │
│     On-Premise                                                  │
│                                                                 │
│                                                                 │
│                                                                 │
│                                                                 │
│                                                                 │
│                                                                 │
│                                                                 │
│                                                                 │
│                                                                 │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
├─────────────────────────────────────────────────────────────────┤
│  Tab/Shift+Tab or ↑/↓: Navigate  Enter: Select  q: Quit        │
╰─────────────────────────────────────────────────────────────────╯
```

### ステップ2-a: Cloud接続設定

```
╭─────────────────────────────────────────────────────────────────╮
│  dito - Oracle NoSQL Database TUI Client                        │
│                                                                 │
│  Cloud Connection                                               │
│  Region:        [ us-ashburn-1                  ]               │
│  Compartment:   [ ocid1.compartment.oc1..aaaaaa ]               │
│                                                                 │
│  Auth Method:                                                   │
│   > (*) OCI Config Profile (default)                            │
│     ( ) Instance Principal                                      │
│     ( ) Resource Principal                                      │
│                                                                 │
│  Config File:   [ DEFAULT                       ]               │
│                                                                 │
│   > Test Connection                                             │
│     Connect                                                     │
│                                                                 │
│                                                                 │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
├─────────────────────────────────────────────────────────────────┤
│  Tab/Shift+Tab: Navigate  Space: Toggle  Enter: Execute  ...   │
╰─────────────────────────────────────────────────────────────────╯
```

### ステップ2-b: On-Premise接続設定（vimライク左寄せ版）

```
╭─────────────────────────────────────────────────────────────────╮
│  dito - Oracle NoSQL Database TUI Client                        │
│                                                                 │
│  On-Premise Connection                                          │
│  Endpoint: [ localhost              ]                           │
│  Port:     [ 8080   ]                                           │
│  Secure:   [ ] HTTPS/TLS                                        │
│                                                                 │
│   > Test Connection                                             │
│     Connect                                                     │
│                                                                 │
│                                                                 │
│                                                                 │
│                                                                 │
│                                                                 │
├─────────────────────────────────────────────────────────────────┤
│ Connected                                                       │
├─────────────────────────────────────────────────────────────────┤
│  Tab/Shift+Tab: Navigate  Space: Toggle  Enter: Execute  ...   │
╰─────────────────────────────────────────────────────────────────╯
```

### ステップ3: テーブル一覧（2ペイン構成）

接続成功後、左ペインにテーブル一覧、右ペインに選択したテーブルの詳細を表示します。

```
╭─────────────────────────────────────────────────────────────────────────────────────╮
│  dito - Oracle NoSQL Database TUI Client                                            │
├─────────────────────────────────┬───────────────────────────────────────────────────┤
│  Tables (7)                     │  Table: users                                     │
│                                 │                                                   │
│  > users                        │  Parent:        -                                 │
│    users.addresses              │  Children:      2 (addresses, phones)             │
│    users.phones                 │                                                   │
│    products                     │  Columns:                                         │
│    products.reviews             │    id           INTEGER (Primary Key)             │
│    orders                       │    name         STRING                            │
│    orders.items                 │    email        STRING                            │
│                                 │    age          INTEGER                           │
│                                 │    created_at   TIMESTAMP(3)                      │
│                                 │                                                   │
│                                 │  Indexes:                                         │
│                                 │    email_idx    (email)                           │
│                                 │                                                   │
│                                 │                                                   │
│                                 │                                                   │
│                                 │                                                   │
│                                 │                                                   │
├─────────────────────────────────┴───────────────────────────────────────────────────┤
│ Connected: localhost:8080                                                           │
├─────────────────────────────────────────────────────────────────────────────────────┤
│  j/k or ↑/↓: Navigate  Tab: Toggle View  Esc: Back  q: Quit                       │
╰─────────────────────────────────────────────────────────────────────────────────────╯
```

**特徴:**
- 左ペイン: テーブル一覧（固定幅30文字）
- 右ペイン: 選択したテーブルの詳細（残りの幅）
- 親子テーブルの階層を表現（インデントなし、ドット記法）
- システムテーブル（`SYS$*`）は非表示
- テーブル数を表示
- 選択中のテーブルは `>` でハイライト

**親子テーブルの表示例:**

子テーブル選択時：
```
╭─────────────────────────────────────────────────────────────────────────────────────╮
│  dito - Oracle NoSQL Database TUI Client                                            │
├─────────────────────────────────┬───────────────────────────────────────────────────┤
│  Tables (7)                     │  Table: users.addresses                           │
│                                 │                                                   │
│    users                        │  Parent:        users                             │
│  > users.addresses              │  Children:      -                                 │
│    users.phones                 │                                                   │
│    products                     │  Columns:                                         │
│    products.reviews             │    id           INTEGER (From parent)             │
│    orders                       │    address_id   INTEGER (Primary Key)             │
│    orders.items                 │    type         STRING                            │
│                                 │    street       STRING                            │
│                                 │    city         STRING                            │
│                                 │    state        STRING                            │
│                                 │    postal_code  STRING                            │
│                                 │    country      STRING                            │
│                                 │    is_primary   BOOLEAN                           │
│                                 │                                                   │
│                                 │  Indexes:                                         │
│                                 │    (none)                                         │
│                                 │                                                   │
│                                 │                                                   │
├─────────────────────────────────┴───────────────────────────────────────────────────┤
│ Connected: localhost:8080                                                           │
├─────────────────────────────────────────────────────────────────────────────────────┤
│  j/k or ↑/↓: Navigate  Tab: Toggle View  Esc: Back  q: Quit                       │
╰─────────────────────────────────────────────────────────────────────────────────────╯
```

### ステップ4: データ表示（右ペイン切り替え）

テーブル一覧画面で`Enter`キーを押すと、右ペインが「データ表示」に切り替わります。`Esc`キーで「概要表示」に戻ります。

#### データ表示モード（Enter押下後）

```
╭─────────────────────────────────────────────────────────────────────────────────────╮
│  dito - Oracle NoSQL Database TUI Client                                            │
├─────────────────────────────────┬───────────────────────────────────────────────────┤
│  Tables (7)                     │  Table: users (150 rows) [Data View]             │
│                                 │                                                   │
│  > users                        │  id  name          email              age  creat…│
│    users.addresses              │  ─…  ────          ─────              ───  ─────…│
│    users.phones                 │  1   Alice Smith   alice@example.com  28   2024-…│
│    products                     │  2   Bob Johnson   bob@example.com    35   2024-…│
│    products.reviews             │  3   Charlie Bro…  charlie@exampl…    42   2024-…│
│    orders                       │  4   Diana Princ…  diana@example.c…   31   2024-…│
│    orders.items                 │  5   Eve Wilson    eve@example.com    29   2024-…│
│                                 │  6   Frank Mille…  frank@example.c…   45   2024-…│
│                                 │  7   Grace Lee     grace@example.com  33   2024-…│
│                                 │  8   Henry Davis   henry@example.com  38   2024-…│
│                                 │  9   Iris Chen     iris@example.com   27   2024-…│
│                                 │  10  Jack Thomps…  jack@example.com   41   2024-…│
│                                 │                                                   │
├─────────────────────────────────┴───────────────────────────────────────────────────┤
│ Showing 1-10 of 150                                                                 │
├─────────────────────────────────────────────────────────────────────────────────────┤
│  j/k or ↑/↓: Navigate  n/p: Page  Esc: Back to Schema  q: Quit                    │
╰─────────────────────────────────────────────────────────────────────────────────────╯
```

**特徴:**
- 2ペインレイアウトを維持（左：テーブル一覧、右：データ表示）
- `Enter`キーで右ペインをデータ表示に切り替え
- `Esc`キーで右ペインを概要表示に戻す
- **データビューモード時:**
  - 左ペインはグレーアウトされ、インアクティブ状態
  - `j`/`k`（または `↑`/`↓`）でデータ行を選択
  - 選択中の行はシアン色 + 太字 + `>` マーカー
  - テーブルを切り替えるには `Esc` でスキーマビューに戻る必要がある
- データ表示時は行番号なし、コンパクトな表示
- カラム幅は自動調整、長いデータは省略（`...`）
- ページング対応（デフォルト10行/右ペイン）
- `n` / `p`キーで次/前ページ（ページ移動時は先頭行を選択）
- ステータスエリアに表示範囲を表示

#### 概要表示モード（Esc押下で戻る）

```
╭─────────────────────────────────────────────────────────────────────────────────────╮
│  dito - Oracle NoSQL Database TUI Client                                            │
├─────────────────────────────────┬───────────────────────────────────────────────────┤
│  Tables (7)                     │  Table: users [Schema View]                       │
│                                 │                                                   │
│  > users                        │  Parent:        -                                 │
│    users.addresses              │  Children:      2 (addresses, phones)             │
│    users.phones                 │                                                   │
│    products                     │  Columns:                                         │
│    products.reviews             │    id           INTEGER (Primary Key)             │
│    orders                       │    name         STRING                            │
│    orders.items                 │    email        STRING                            │
│                                 │    age          INTEGER                           │
│                                 │    created_at   TIMESTAMP(3)                      │
│                                 │                                                   │
│                                 │  Indexes:                                         │
│                                 │    email_idx    (email)                           │
│                                 │                                                   │
├─────────────────────────────────┴───────────────────────────────────────────────────┤
│ Connected: localhost:8080                                                           │
├─────────────────────────────────────────────────────────────────────────────────────┤
│  j/k or ↑/↓: Navigate  Enter: View Data  Esc: Back  q: Quit                       │
╰─────────────────────────────────────────────────────────────────────────────────────╯
```

**操作フロー:**
1. **スキーマビューモード**: テーブル選択（j/k で移動）
2. `Enter` → 右ペインがデータ表示に切り替わる、左ペインがグレーアウト
3. **データビューモード**: データ行選択（j/k で移動）
4. `n`/`p` → ページ移動
5. `Esc` → スキーマビューモードに戻る（左ペインが通常色に戻る）
6. もう一度 `Esc` → 接続設定画面に戻る

#### 空のテーブルのデータ表示

```
╭─────────────────────────────────────────────────────────────────────────────────────╮
│  dito - Oracle NoSQL Database TUI Client                                            │
├─────────────────────────────────┬───────────────────────────────────────────────────┤
│  Tables (7)                     │  Table: empty_table (0 rows) [Data View]         │
│                                 │                                                   │
│  > users                        │  id  name  email                                  │
│    users.addresses              │  ──  ────  ─────                                  │
│    users.phones                 │                                                   │
│    products                     │            No data found                          │
│    products.reviews             │                                                   │
│    orders                       │                                                   │
│    orders.items                 │                                                   │
│                                 │                                                   │
│                                 │                                                   │
├─────────────────────────────────┴───────────────────────────────────────────────────┤
│ No rows                                                                             │
├─────────────────────────────────────────────────────────────────────────────────────┤
│  j/k or ↑/↓: Navigate  Esc: Back to Schema  q: Quit                               │
╰─────────────────────────────────────────────────────────────────────────────────────╯
```

**将来実装予定:**
- 左ペイン折りたたみ機能（例: `Ctrl+B`）でデータ表示領域を拡大
