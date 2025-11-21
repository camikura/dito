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
│  j/k or ↑/↓: Navigate  Enter: View Data  Esc: Back  q: Quit                        │
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
│  j/k or ↑/↓: Navigate  Enter: View Data  Esc: Back  q: Quit                        │
╰─────────────────────────────────────────────────────────────────────────────────────╯
```
