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
