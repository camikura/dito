package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/oracle/nosql-go-sdk/nosqldb"
	"github.com/oracle/nosql-go-sdk/nosqldb/nosqlerr"

	"github.com/camikura/dito/internal/ui"
)

// 画面の種類
type screen int

const (
	screenSelection screen = iota
	screenOnPremiseConfig
	screenCloudConfig
	screenTableList
)

// 接続状態
type connectionStatus int

const (
	statusDisconnected connectionStatus = iota
	statusConnecting
	statusConnected
	statusError
)

// 接続結果メッセージ
type connectionResultMsg struct {
	err      error
	version  string
	client   *nosqldb.Client
	endpoint string
	isTest   bool // trueの場合はテスト接続（画面遷移しない）
}

// テーブル一覧取得結果メッセージ
type tableListResultMsg struct {
	tables []string
	err    error
}

// テーブル詳細取得結果メッセージ
type tableDetailsResultMsg struct {
	tableName string
	schema    *nosqldb.TableResult
	indexes   []nosqldb.IndexInfo
	err       error
}

// テーブルデータ取得結果メッセージ
type tableDataResultMsg struct {
	tableName       string
	rows            []map[string]interface{}
	lastPKValues    map[string]interface{} // 最後の行のPRIMARY KEY値（カーソルとして使用）
	hasMore         bool                   // さらにデータがあるかどうか
	err             error
	isAppend        bool   // 既存データに追加するかどうか
	sql             string // デバッグ用: 実行したSQL
	displaySQL      string // 表示用: LIMIT句なしのSQL
}

// On-Premise接続設定
type onPremiseConfig struct {
	endpoint      string
	port          string
	secure        bool
	focus         int // フォーカス中のフィールド
	status        connectionStatus
	errorMsg      string
	serverVersion string
	cursorPos     int // テキスト入力のカーソル位置
}

// Cloud接続設定
type cloudConfig struct {
	region        string
	compartment   string
	authMethod    int // 0: OCI Config Profile, 1: Instance Principal, 2: Resource Principal
	configFile    string
	focus         int // フォーカス中のフィールド
	status        connectionStatus
	errorMsg      string
	serverVersion string
	cursorPos     int // テキスト入力のカーソル位置
}

// 右ペインの表示モード
type rightPaneMode int

const (
	rightPaneModeSchema rightPaneMode = iota
	rightPaneModeList   // データ一覧表示
	rightPaneModeDetail // レコード表示
)

// モデル定義
type model struct {
	screen          screen
	choices         []string
	cursor          int
	selected        map[int]struct{}
	onPremiseConfig onPremiseConfig
	cloudConfig     cloudConfig
	width           int
	height          int
	// テーブル一覧画面用
	nosqlClient     *nosqldb.Client
	tables          []string
	selectedTable   int
	endpoint        string // 接続先エンドポイント（ステータス表示用）
	tableDetails    map[string]*tableDetailsResultMsg
	loadingDetails  bool
	// データ表示用
	rightPaneMode      rightPaneMode
	tableData          map[string]*tableDataResultMsg
	dataOffset         int // データの取得オフセット（無限スクロール用）
	fetchSize          int // 一度に取得するデータ数
	loadingData        bool
	selectedDataRow    int // データビューモードで選択中の行（全体の絶対位置）
	viewportOffset     int // 表示開始位置
	viewportSize       int // 一度に画面に表示する行数
	horizontalOffset   int // 横スクロールのオフセット（カラム単位、0始まり）
}

// 初期化関数
func initialModel() model {
	return model{
		screen:   screenSelection,
		choices:  []string{"Oracle NoSQL Cloud Service", "On-Premise"},
		selected: make(map[int]struct{}),
		onPremiseConfig: onPremiseConfig{
			endpoint:  "localhost",
			port:      "8080",
			secure:    false,
			focus:     0,
			status:    statusDisconnected,
			cursorPos: 9, // "localhost"の末尾
		},
		cloudConfig: cloudConfig{
			region:      "us-ashburn-1",
			compartment: "",
			authMethod:  0, // OCI Config Profile
			configFile:  "DEFAULT",
			focus:       0,
			status:      statusDisconnected,
			cursorPos:   12, // "us-ashburn-1"の末尾
		},
		tableDetails:  make(map[string]*tableDetailsResultMsg),
		rightPaneMode: rightPaneModeSchema,
		tableData:     make(map[string]*tableDataResultMsg),
		dataOffset:    0,
		fetchSize:     100, // 一度に100行取得（無限スクロール）
		viewportSize:  10,
	}
}

// Initメソッド
func (m model) Init() tea.Cmd {
	return nil
}

// テーブル一覧を取得するCommand
func fetchTables(client *nosqldb.Client) tea.Cmd {
	return func() tea.Msg {
		req := &nosqldb.ListTablesRequest{}
		result, err := client.ListTables(req)
		if err != nil {
			return tableListResultMsg{err: err}
		}

		// システムテーブル（SYS$*）をフィルタリング
		var userTables []string
		for _, table := range result.Tables {
			if !strings.HasPrefix(table, "SYS$") {
				userTables = append(userTables, table)
			}
		}

		return tableListResultMsg{tables: userTables, err: nil}
	}
}

// テーブル詳細を取得するCommand
func fetchTableDetails(client *nosqldb.Client, tableName string) tea.Cmd {
	return func() tea.Msg {
		// テーブル情報を取得
		tableReq := &nosqldb.GetTableRequest{
			TableName: tableName,
		}
		tableResult, err := client.GetTable(tableReq)
		if err != nil {
			return tableDetailsResultMsg{tableName: tableName, err: err}
		}

		// インデックス情報を取得
		indexReq := &nosqldb.GetIndexesRequest{
			TableName: tableName,
		}
		indexResult, err := client.GetIndexes(indexReq)
		if err != nil {
			// インデックス取得エラーは無視して、スキーマ情報だけ返す
			return tableDetailsResultMsg{tableName: tableName, schema: tableResult, indexes: nil, err: nil}
		}

		return tableDetailsResultMsg{tableName: tableName, schema: tableResult, indexes: indexResult.Indexes, err: nil}
	}
}

// テーブルデータを取得するCommand（PRIMARY KEYでソート、初回取得）
func fetchTableData(client *nosqldb.Client, tableName string, limit int, primaryKeys []string) tea.Cmd {
	return fetchTableDataWithCursor(client, tableName, limit, primaryKeys, nil, false)
}

// テーブルデータを追加取得するCommand（PRIMARY KEYカーソル使用）
func fetchMoreTableData(client *nosqldb.Client, tableName string, limit int, primaryKeys []string, lastPKValues map[string]interface{}) tea.Cmd {
	return fetchTableDataWithCursor(client, tableName, limit, primaryKeys, lastPKValues, true)
}

// テーブルデータを取得する内部関数（PRIMARY KEYカーソル対応）
func fetchTableDataWithCursor(client *nosqldb.Client, tableName string, limit int, primaryKeys []string, lastPKValues map[string]interface{}, isAppend bool) tea.Cmd {
	return func() tea.Msg {
		// PRIMARY KEY順に明示的にソート
		var orderByClause string
		if len(primaryKeys) > 0 {
			orderByClause = " ORDER BY " + strings.Join(primaryKeys, ", ")
		}

		// WHERE句を構築（PRIMARY KEYカーソルがある場合）
		var whereClause string
		if lastPKValues != nil && len(lastPKValues) > 0 && len(primaryKeys) > 0 {
			// 複合PRIMARY KEYの場合の条件を構築
			// 例: WHERE pk1 > ? OR (pk1 = ? AND pk2 > ?) OR (pk1 = ? AND pk2 = ? AND pk3 > ?)
			var conditions []string
			for i := 0; i < len(primaryKeys); i++ {
				var cond string
				if i == 0 {
					// 最初のキー: pk1 > ?
					val := lastPKValues[primaryKeys[i]]
					cond = fmt.Sprintf("%s > %s", primaryKeys[i], formatValue(val))
				} else {
					// それ以降: (pk1 = ? AND pk2 = ? AND ... AND pkN > ?)
					var parts []string
					for j := 0; j < i; j++ {
						val := lastPKValues[primaryKeys[j]]
						parts = append(parts, fmt.Sprintf("%s = %s", primaryKeys[j], formatValue(val)))
					}
					val := lastPKValues[primaryKeys[i]]
					parts = append(parts, fmt.Sprintf("%s > %s", primaryKeys[i], formatValue(val)))
					cond = "(" + strings.Join(parts, " AND ") + ")"
				}
				conditions = append(conditions, cond)
			}
			whereClause = " WHERE " + strings.Join(conditions, " OR ")
		}

		statement := fmt.Sprintf("SELECT * FROM %s%s%s LIMIT %d", tableName, whereClause, orderByClause, limit)

		// 表示用SQL（LIMIT句なし）
		displayStatement := fmt.Sprintf("SELECT * FROM %s%s%s", tableName, whereClause, orderByClause)

		prepReq := &nosqldb.PrepareRequest{
			Statement: statement,
		}
		prepResult, err := client.Prepare(prepReq)
		if err != nil {
			return tableDataResultMsg{tableName: tableName, err: err, isAppend: isAppend, sql: statement, displaySQL: displayStatement}
		}

		queryReq := &nosqldb.QueryRequest{
			PreparedStatement: &prepResult.PreparedStatement,
		}

		// すべての結果を取得（SDKの内部ページネーションを使用）
		var rows []map[string]interface{}
		for {
			queryResult, err := client.Query(queryReq)
			if err != nil {
				return tableDataResultMsg{tableName: tableName, err: err, isAppend: isAppend, sql: statement, displaySQL: displayStatement}
			}

			// 結果を取得
			results, err := queryResult.GetResults()
			if err != nil {
				return tableDataResultMsg{tableName: tableName, err: err, isAppend: isAppend, sql: statement, displaySQL: displayStatement}
			}

			for _, result := range results {
				rows = append(rows, result.Map())
			}

			// 継続トークンがなければ終了
			if queryReq.IsDone() {
				break
			}
		}

		// 最後の行のPRIMARY KEY値を保存
		var newLastPKValues map[string]interface{}
		if len(rows) > 0 && len(primaryKeys) > 0 {
			lastRow := rows[len(rows)-1]
			newLastPKValues = make(map[string]interface{})
			for _, pk := range primaryKeys {
				if val, exists := lastRow[pk]; exists {
					newLastPKValues[pk] = val
				}
			}
		}

		// 次のページがあるかチェック
		// 取得した行数がlimitと同じなら、まだデータがある可能性がある
		hasMore := len(rows) == limit

		return tableDataResultMsg{
			tableName:    tableName,
			rows:         rows,
			lastPKValues: newLastPKValues,
			hasMore:      hasMore,
			err:          nil,
			isAppend:     isAppend,
			sql:          statement,
			displaySQL:   displayStatement,
		}
	}
}

// 値をSQL文字列にフォーマット
func formatValue(val interface{}) string {
	switch v := val.(type) {
	case string:
		// 文字列はシングルクォートで囲む
		return fmt.Sprintf("'%s'", strings.ReplaceAll(v, "'", "''"))
	case int, int32, int64, float32, float64:
		return fmt.Sprintf("%v", v)
	default:
		return fmt.Sprintf("'%v'", v)
	}
}

// 接続を試みるCommand
func connectToNoSQL(endpoint, port string, isTest bool) tea.Cmd {
	return func() tea.Msg {
		// 接続設定
		endpointURL := fmt.Sprintf("http://%s:%s", endpoint, port)
		cfg := nosqldb.Config{
			Mode:     "onprem",
			Endpoint: endpointURL,
		}

		// クライアント作成
		client, err := nosqldb.NewClient(cfg)
		if err != nil {
			return connectionResultMsg{err: err, isTest: isTest}
		}

		// 簡単なテスト（テーブル一覧取得）
		req := &nosqldb.ListTablesRequest{}
		_, err = client.ListTables(req)
		if err != nil {
			client.Close()
			// エラーの詳細を取得
			if nosqlErr, ok := err.(*nosqlerr.Error); ok {
				return connectionResultMsg{err: fmt.Errorf("NoSQL Error: %s", nosqlErr.Error()), isTest: isTest}
			}
			return connectionResultMsg{err: err, isTest: isTest}
		}

		// テスト接続の場合はクライアントをクローズ
		if isTest {
			client.Close()
			return connectionResultMsg{
				version: "Connected",
				err:     nil,
				isTest:  true,
			}
		}

		// 接続成功 - クライアントを返す（closeしない）
		return connectionResultMsg{
			version:  "Connected",
			err:      nil,
			client:   client,
			endpoint: fmt.Sprintf("%s:%s", endpoint, port),
			isTest:   false,
		}
	}
}

// Updateメソッド
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// ビューポートサイズを画面の高さから計算
		// 右ペインの高さ (m.height - 8) からヘッダー等を引く
		rightPaneHeight := m.height - 8
		m.viewportSize = rightPaneHeight - 4 // SQLエリア+ボーダー: 2行 + カラムヘッダー+セパレーター: 2行
		if m.viewportSize < 1 {
			m.viewportSize = 1
		}
		return m, nil
	case tea.KeyMsg:
		switch m.screen {
		case screenSelection:
			return m.updateSelection(msg)
		case screenOnPremiseConfig:
			return m.updateOnPremiseConfig(msg)
		case screenCloudConfig:
			return m.updateCloudConfig(msg)
		case screenTableList:
			return m.updateTableList(msg)
		}
	case connectionResultMsg:
		// 接続結果を処理
		if msg.err != nil {
			m.onPremiseConfig.status = statusError
			m.onPremiseConfig.errorMsg = msg.err.Error()
		} else {
			m.onPremiseConfig.status = statusConnected
			m.onPremiseConfig.serverVersion = msg.version
			m.onPremiseConfig.errorMsg = ""

			// テスト接続でない場合のみテーブル一覧を取得して画面遷移
			if !msg.isTest {
				// クライアントとエンドポイントを保存
				m.nosqlClient = msg.client
				m.endpoint = msg.endpoint
				// テーブル一覧を取得
				return m, fetchTables(msg.client)
			}
		}
		return m, nil
	case tableListResultMsg:
		// テーブル一覧取得結果を処理
		if msg.err != nil {
			m.onPremiseConfig.status = statusError
			m.onPremiseConfig.errorMsg = fmt.Sprintf("Failed to fetch tables: %v", msg.err)
		} else {
			m.tables = msg.tables
			m.selectedTable = 0
			// テーブル一覧画面に遷移
			m.screen = screenTableList
			// 最初のテーブルの詳細を取得
			if len(m.tables) > 0 {
				return m, fetchTableDetails(m.nosqlClient, m.tables[0])
			}
		}
		return m, nil
	case tableDetailsResultMsg:
		// テーブル詳細取得結果を処理
		if msg.err == nil {
			m.tableDetails[msg.tableName] = &msg
		}
		m.loadingDetails = false

		// グリッドビューモードで、このテーブルのデータがまだ取得されていない場合は取得
		if m.rightPaneMode == rightPaneModeList && msg.err == nil {
			if _, exists := m.tableData[msg.tableName]; !exists {
				m.loadingData = true
				primaryKeys := parsePrimaryKeysFromDDL(msg.schema.DDL)
				return m, fetchTableData(m.nosqlClient, msg.tableName, m.fetchSize, primaryKeys)
			}
		}
		return m, nil
	case tableDataResultMsg:
		// テーブルデータ取得結果を処理
		if msg.err == nil {
			if msg.isAppend {
				// 既存データに追加（SQLは更新しない）
				if existingData, exists := m.tableData[msg.tableName]; exists {
					existingData.rows = append(existingData.rows, msg.rows...)
					existingData.lastPKValues = msg.lastPKValues
					existingData.hasMore = msg.hasMore
				}
			} else {
				// 新規データとして設定
				m.tableData[msg.tableName] = &msg
			}
		} else {
			// エラーの場合もデータを保存（エラーメッセージとSQLを表示するため）
			m.tableData[msg.tableName] = &msg
		}
		m.loadingData = false
		return m, nil
	}

	return m, nil
}

// エディション選択画面のUpdate
func (m model) updateSelection(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "up", "shift+tab":
		m.cursor--
		if m.cursor < 0 {
			m.cursor = len(m.choices) - 1
		}
	case "down", "tab":
		m.cursor = (m.cursor + 1) % len(m.choices)
	case "enter":
		// 0: Cloud, 1: On-Premise
		switch m.cursor {
		case 0:
			// Cloud: 接続設定画面に遷移
			m.screen = screenCloudConfig
			return m, nil
		case 1:
			// On-Premise: 接続設定画面に遷移
			m.screen = screenOnPremiseConfig
			return m, nil
		}
	}
	return m, nil
}

// Cloud接続設定画面のUpdate
func (m model) updateCloudConfig(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		// エディション選択画面に戻る
		m.screen = screenSelection
		return m, nil
	case "tab":
		// 次のフィールドへ
		m.cloudConfig.focus = (m.cloudConfig.focus + 1) % 8
		// テキスト入力フィールドの場合、カーソルを末尾に
		if m.cloudConfig.focus == 0 {
			m.cloudConfig.cursorPos = len(m.cloudConfig.region)
		} else if m.cloudConfig.focus == 1 {
			m.cloudConfig.cursorPos = len(m.cloudConfig.compartment)
		} else if m.cloudConfig.focus == 5 {
			m.cloudConfig.cursorPos = len(m.cloudConfig.configFile)
		}
		return m, nil
	case "shift+tab":
		// 前のフィールドへ
		m.cloudConfig.focus--
		if m.cloudConfig.focus < 0 {
			m.cloudConfig.focus = 7
		}
		// テキスト入力フィールドの場合、カーソルを末尾に
		if m.cloudConfig.focus == 0 {
			m.cloudConfig.cursorPos = len(m.cloudConfig.region)
		} else if m.cloudConfig.focus == 1 {
			m.cloudConfig.cursorPos = len(m.cloudConfig.compartment)
		} else if m.cloudConfig.focus == 5 {
			m.cloudConfig.cursorPos = len(m.cloudConfig.configFile)
		}
		return m, nil
	case "enter":
		// ボタンが選択されている場合
		if m.cloudConfig.focus == 6 {
			// 接続テスト - TODO: Cloud接続実装
			m.cloudConfig.status = statusConnecting
			m.cloudConfig.errorMsg = ""
			return m, nil
		} else if m.cloudConfig.focus == 7 {
			// 接続する - TODO: Cloud接続実装
			m.cloudConfig.status = statusConnecting
			m.cloudConfig.errorMsg = ""
			return m, nil
		}
		return m, nil
	case " ":
		// ラジオボタンの選択
		if m.cloudConfig.focus >= 2 && m.cloudConfig.focus <= 4 {
			m.cloudConfig.authMethod = m.cloudConfig.focus - 2
		}
		return m, nil
	case "left":
		// カーソルを左に移動
		if m.cloudConfig.cursorPos > 0 {
			m.cloudConfig.cursorPos--
		}
		return m, nil
	case "right":
		// カーソルを右に移動
		var maxPos int
		if m.cloudConfig.focus == 0 {
			maxPos = len(m.cloudConfig.region)
		} else if m.cloudConfig.focus == 1 {
			maxPos = len(m.cloudConfig.compartment)
		} else if m.cloudConfig.focus == 5 {
			maxPos = len(m.cloudConfig.configFile)
		}
		if m.cloudConfig.cursorPos < maxPos {
			m.cloudConfig.cursorPos++
		}
		return m, nil
	case "backspace":
		// テキストフィールドの入力削除
		if m.cloudConfig.focus == 0 && m.cloudConfig.cursorPos > 0 {
			m.cloudConfig.region = m.cloudConfig.region[:m.cloudConfig.cursorPos-1] + m.cloudConfig.region[m.cloudConfig.cursorPos:]
			m.cloudConfig.cursorPos--
		} else if m.cloudConfig.focus == 1 && m.cloudConfig.cursorPos > 0 {
			m.cloudConfig.compartment = m.cloudConfig.compartment[:m.cloudConfig.cursorPos-1] + m.cloudConfig.compartment[m.cloudConfig.cursorPos:]
			m.cloudConfig.cursorPos--
		} else if m.cloudConfig.focus == 5 && m.cloudConfig.cursorPos > 0 {
			m.cloudConfig.configFile = m.cloudConfig.configFile[:m.cloudConfig.cursorPos-1] + m.cloudConfig.configFile[m.cloudConfig.cursorPos:]
			m.cloudConfig.cursorPos--
		}
		return m, nil
	default:
		// テキスト入力
		if len(msg.String()) == 1 {
			if m.cloudConfig.focus == 0 {
				m.cloudConfig.region = m.cloudConfig.region[:m.cloudConfig.cursorPos] + msg.String() + m.cloudConfig.region[m.cloudConfig.cursorPos:]
				m.cloudConfig.cursorPos++
			} else if m.cloudConfig.focus == 1 {
				m.cloudConfig.compartment = m.cloudConfig.compartment[:m.cloudConfig.cursorPos] + msg.String() + m.cloudConfig.compartment[m.cloudConfig.cursorPos:]
				m.cloudConfig.cursorPos++
			} else if m.cloudConfig.focus == 5 {
				m.cloudConfig.configFile = m.cloudConfig.configFile[:m.cloudConfig.cursorPos] + msg.String() + m.cloudConfig.configFile[m.cloudConfig.cursorPos:]
				m.cloudConfig.cursorPos++
			}
		}
		return m, nil
	}
}

// On-Premise接続設定画面のUpdate
func (m model) updateOnPremiseConfig(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		// エディション選択画面に戻る
		m.screen = screenSelection
		return m, nil
	case "tab":
		// 次のフィールドへ
		m.onPremiseConfig.focus = (m.onPremiseConfig.focus + 1) % 5
		// テキスト入力フィールドの場合、カーソルを末尾に
		if m.onPremiseConfig.focus == 0 {
			m.onPremiseConfig.cursorPos = len(m.onPremiseConfig.endpoint)
		} else if m.onPremiseConfig.focus == 1 {
			m.onPremiseConfig.cursorPos = len(m.onPremiseConfig.port)
		}
		return m, nil
	case "shift+tab":
		// 前のフィールドへ
		m.onPremiseConfig.focus--
		if m.onPremiseConfig.focus < 0 {
			m.onPremiseConfig.focus = 4
		}
		// テキスト入力フィールドの場合、カーソルを末尾に
		if m.onPremiseConfig.focus == 0 {
			m.onPremiseConfig.cursorPos = len(m.onPremiseConfig.endpoint)
		} else if m.onPremiseConfig.focus == 1 {
			m.onPremiseConfig.cursorPos = len(m.onPremiseConfig.port)
		}
		return m, nil
	case "enter":
		// ボタンが選択されている場合
		if m.onPremiseConfig.focus == 3 {
			// 接続テスト（テスト接続なので画面遷移しない）
			m.onPremiseConfig.status = statusConnecting
			m.onPremiseConfig.errorMsg = ""
			return m, connectToNoSQL(m.onPremiseConfig.endpoint, m.onPremiseConfig.port, true)
		} else if m.onPremiseConfig.focus == 4 {
			// 接続する（実接続なのでテーブル一覧画面に遷移）
			m.onPremiseConfig.status = statusConnecting
			m.onPremiseConfig.errorMsg = ""
			return m, connectToNoSQL(m.onPremiseConfig.endpoint, m.onPremiseConfig.port, false)
		}
		return m, nil
	case " ":
		// セキュアチェックボックスのトグル
		if m.onPremiseConfig.focus == 2 {
			m.onPremiseConfig.secure = !m.onPremiseConfig.secure
		}
		return m, nil
	case "left":
		// カーソルを左に移動
		if m.onPremiseConfig.cursorPos > 0 {
			m.onPremiseConfig.cursorPos--
		}
		return m, nil
	case "right":
		// カーソルを右に移動
		var maxPos int
		if m.onPremiseConfig.focus == 0 {
			maxPos = len(m.onPremiseConfig.endpoint)
		} else if m.onPremiseConfig.focus == 1 {
			maxPos = len(m.onPremiseConfig.port)
		}
		if m.onPremiseConfig.cursorPos < maxPos {
			m.onPremiseConfig.cursorPos++
		}
		return m, nil
	case "backspace":
		// テキストフィールドの入力削除
		if m.onPremiseConfig.focus == 0 && m.onPremiseConfig.cursorPos > 0 {
			m.onPremiseConfig.endpoint = m.onPremiseConfig.endpoint[:m.onPremiseConfig.cursorPos-1] + m.onPremiseConfig.endpoint[m.onPremiseConfig.cursorPos:]
			m.onPremiseConfig.cursorPos--
		} else if m.onPremiseConfig.focus == 1 && m.onPremiseConfig.cursorPos > 0 {
			m.onPremiseConfig.port = m.onPremiseConfig.port[:m.onPremiseConfig.cursorPos-1] + m.onPremiseConfig.port[m.onPremiseConfig.cursorPos:]
			m.onPremiseConfig.cursorPos--
		}
		return m, nil
	default:
		// テキスト入力
		if len(msg.String()) == 1 {
			if m.onPremiseConfig.focus == 0 {
				m.onPremiseConfig.endpoint = m.onPremiseConfig.endpoint[:m.onPremiseConfig.cursorPos] + msg.String() + m.onPremiseConfig.endpoint[m.onPremiseConfig.cursorPos:]
				m.onPremiseConfig.cursorPos++
			} else if m.onPremiseConfig.focus == 1 {
				m.onPremiseConfig.port = m.onPremiseConfig.port[:m.onPremiseConfig.cursorPos] + msg.String() + m.onPremiseConfig.port[m.onPremiseConfig.cursorPos:]
				m.onPremiseConfig.cursorPos++
			}
		}
		return m, nil
	}
}

// テーブル一覧画面のUpdate
func (m model) updateTableList(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	oldSelection := m.selectedTable

	switch msg.String() {
	case "ctrl+c", "q":
		// クライアントをクローズしてから終了
		if m.nosqlClient != nil {
			m.nosqlClient.Close()
		}
		return m, tea.Quit
	case "up", "k":
		if m.rightPaneMode == rightPaneModeList || m.rightPaneMode == rightPaneModeDetail {
			// グリッドビュー/レコードビュー: データ行を選択
			if m.selectedDataRow > 0 {
				m.selectedDataRow--
				// ビューポートを調整（選択行がビューポートの上端より上になった場合）
				if m.selectedDataRow < m.viewportOffset {
					m.viewportOffset = m.selectedDataRow
				}
			}
		} else {
			// スキーマビューモード: テーブルを選択
			if m.selectedTable > 0 {
				m.selectedTable--
			}
		}
	case "down", "j":
		if m.rightPaneMode == rightPaneModeList || m.rightPaneMode == rightPaneModeDetail {
			// グリッドビュー/レコードビュー: データ行を選択
			tableName := m.tables[m.selectedTable]
			if data, exists := m.tableData[tableName]; exists && data.err == nil {
				totalRows := len(data.rows)
				if m.selectedDataRow < totalRows-1 {
					m.selectedDataRow++
					// ビューポートを調整（選択行がビューポートの下端より下になった場合）
					if m.selectedDataRow >= m.viewportOffset+m.viewportSize {
						m.viewportOffset = m.selectedDataRow - m.viewportSize + 1
					}

					// 残り10行以内まで来たら、さらにデータがある場合は追加取得
					remainingRows := totalRows - m.selectedDataRow - 1
					if remainingRows <= 10 && data.hasMore && !m.loadingData {
						m.loadingData = true
						// PRIMARY KEYを取得
						var primaryKeys []string
						if details, exists := m.tableDetails[tableName]; exists && details.schema != nil && details.schema.DDL != "" {
							primaryKeys = parsePrimaryKeysFromDDL(details.schema.DDL)
						}
						return m, fetchMoreTableData(m.nosqlClient, tableName, m.fetchSize, primaryKeys, data.lastPKValues)
					}
				}
			}
		} else {
			// スキーマビューモード: テーブルを選択
			if m.selectedTable < len(m.tables)-1 {
				m.selectedTable++
			}
		}
	case "h", "left":
		// データビューモード: 左にスクロール
		if m.rightPaneMode == rightPaneModeList {
			if m.horizontalOffset > 0 {
				m.horizontalOffset--
			}
		}
	case "l", "right":
		// データビューモード: 右にスクロール
		if m.rightPaneMode == rightPaneModeList {
			tableName := m.tables[m.selectedTable]
			// カラム数を取得
			var totalColumns int
			if details, exists := m.tableDetails[tableName]; exists && details.schema != nil && details.schema.DDL != "" {
				primaryKeys := parsePrimaryKeysFromDDL(details.schema.DDL)
				columns := parseColumnsFromDDL(details.schema.DDL, primaryKeys)
				totalColumns = len(columns)
			} else if data, exists := m.tableData[tableName]; exists && len(data.rows) > 0 {
				totalColumns = len(data.rows[0])
			}
			// 最後のカラムまでスクロールできるが、少なくとも1カラムは表示する
			if m.horizontalOffset < totalColumns-1 {
				m.horizontalOffset++
			}
		}
	case "esc", "u":
		if m.rightPaneMode == rightPaneModeDetail {
			// レコードビュー → グリッドビュー
			m.rightPaneMode = rightPaneModeList
			return m, nil
		} else if m.rightPaneMode == rightPaneModeList {
			// グリッドビュー → スキーマビュー
			m.rightPaneMode = rightPaneModeSchema
			m.horizontalOffset = 0 // 横スクロールをリセット
			return m, nil
		}
		// スキーマビュー → 接続設定画面に戻る
		m.screen = screenOnPremiseConfig
		return m, nil
	case "enter", "o":
		if m.rightPaneMode == rightPaneModeSchema {
			// スキーマビュー → グリッドビュー
			m.rightPaneMode = rightPaneModeList
			m.selectedDataRow = 0    // 行選択をリセット
			m.viewportOffset = 0     // ビューポートをリセット
			m.horizontalOffset = 0   // 横スクロールをリセット
			// データ表示モードに切り替えたとき、データとテーブル詳細を取得
			if len(m.tables) > 0 {
				tableName := m.tables[m.selectedTable]

				// テーブル詳細がまだ取得されていない場合は取得
				var cmds []tea.Cmd
				if _, exists := m.tableDetails[tableName]; !exists {
					m.loadingDetails = true
					cmds = append(cmds, fetchTableDetails(m.nosqlClient, tableName))
				}

				// データがまだ取得されていない場合は取得
				if _, exists := m.tableData[tableName]; !exists {
					m.loadingData = true
					// PRIMARY KEYを取得（テーブル詳細があれば）
					var primaryKeys []string
					if details, exists := m.tableDetails[tableName]; exists && details.schema != nil && details.schema.DDL != "" {
						primaryKeys = parsePrimaryKeysFromDDL(details.schema.DDL)
					}
					cmds = append(cmds, fetchTableData(m.nosqlClient, tableName, m.fetchSize, primaryKeys))
				}

				if len(cmds) > 0 {
					return m, tea.Batch(cmds...)
				}
			}
		} else if m.rightPaneMode == rightPaneModeList {
			// グリッドビュー → レコードビュー
			m.rightPaneMode = rightPaneModeDetail
		}
	}

	// テーブル選択が変わった場合、詳細を取得（スキーマビューモードのみ）
	if m.rightPaneMode == rightPaneModeSchema && oldSelection != m.selectedTable && len(m.tables) > 0 {
		tableName := m.tables[m.selectedTable]
		// まだ取得していないテーブルの場合のみ取得
		if _, exists := m.tableDetails[tableName]; !exists {
			m.loadingDetails = true
			return m, fetchTableDetails(m.nosqlClient, tableName)
		}
	}

	return m, nil
}

// Viewメソッド
func (m model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	// 共通スタイル
	statusBarStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Padding(0, 1).
		Width(m.width - 2)

	// メインコンテンツ
	var content string
	switch m.screen {
	case screenSelection:
		content = m.viewSelectionContent()
	case screenOnPremiseConfig:
		content = m.viewOnPremiseConfigContent()
	case screenCloudConfig:
		content = m.viewCloudConfigContent()
	case screenTableList:
		return m.viewTableList() // テーブル一覧は独自レイアウト
	default:
		content = "Unknown screen"
	}

	// コンテンツを左寄せ
	contentHeight := m.height - 7 // タイトル行、空行、セパレーター×3、ステータスエリア、フッターを除く
	contentStyle := lipgloss.NewStyle().
		Width(m.width - 2).
		Height(contentHeight).
		AlignVertical(lipgloss.Top).
		AlignHorizontal(lipgloss.Left).
		PaddingLeft(1)

	leftAlignedContent := contentStyle.Render(content)

	// セパレーター
	separator := ui.Separator(m.width - 2)

	// ステータス表示エリア（1行）
	var statusMessage string
	statusStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FF00"))
	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF0000"))

	if m.screen == screenOnPremiseConfig {
		switch m.onPremiseConfig.status {
		case statusConnecting:
			statusMessage = statusStyle.Render("Connecting...")
		case statusConnected:
			msg := "Connected"
			if m.onPremiseConfig.serverVersion != "" {
				msg = m.onPremiseConfig.serverVersion
			}
			statusMessage = statusStyle.Render(msg)
		case statusError:
			msg := "Connection failed"
			if m.onPremiseConfig.errorMsg != "" {
				errMsg := m.onPremiseConfig.errorMsg
				maxWidth := m.width - 10
				if len(errMsg) > maxWidth {
					errMsg = errMsg[:maxWidth] + "..."
				}
				msg = errMsg
			}
			statusMessage = errorStyle.Render(msg)
		}
	} else if m.screen == screenCloudConfig {
		switch m.cloudConfig.status {
		case statusConnecting:
			statusMessage = statusStyle.Render("Connecting...")
		case statusConnected:
			msg := "Connected"
			if m.cloudConfig.serverVersion != "" {
				msg = m.cloudConfig.serverVersion
			}
			statusMessage = statusStyle.Render(msg)
		case statusError:
			msg := "Connection failed"
			if m.cloudConfig.errorMsg != "" {
				errMsg := m.cloudConfig.errorMsg
				maxWidth := m.width - 10
				if len(errMsg) > maxWidth {
					errMsg = errMsg[:maxWidth] + "..."
				}
				msg = errMsg
			}
			statusMessage = errorStyle.Render(msg)
		}
	}

	statusAreaStyle := lipgloss.NewStyle().
		Padding(0, 1).
		Width(m.width - 2)
	statusArea := statusAreaStyle.Render(statusMessage)

	// フッター（ヘルプテキスト）
	var helpText string
	switch m.screen {
	case screenSelection:
		helpText = "Tab/Shift+Tab or ↑/↓: Navigate  Enter: Select  q: Quit"
	case screenOnPremiseConfig:
		helpText = "Tab/Shift+Tab: Navigate  Space: Toggle  Enter: Execute  Esc: Back  Ctrl+C: Quit"
	case screenCloudConfig:
		helpText = "Tab/Shift+Tab: Navigate  Space: Toggle  Enter: Execute  Esc: Back  Ctrl+C: Quit"
	}
	footer := statusBarStyle.Render(helpText)

	// 全体を組み立て（手動でボーダーを描画）
	borderStyleColor := lipgloss.NewStyle().Foreground(lipgloss.Color("#00D9FF"))

	// 上部ボーダー: ╭── Dito ─────...╮
	title := " Dito "
	// 全体の幅 = m.width
	// "╭──" = 3文字, title = 6文字, "╮" = 1文字
	// 残りの "─" = m.width - 3 - 6 - 1 = m.width - 10
	topBorder := borderStyleColor.Render("╭──" + title + strings.Repeat("─", m.width-10) + "╮")

	// 左右のボーダー文字
	leftBorder := borderStyleColor.Render("│")
	rightBorder := borderStyleColor.Render("│")

	// コンテンツの各行にボーダーを追加
	var result strings.Builder
	result.WriteString(topBorder + "\n")

	// タイトル行の下に空行を追加
	emptyLine := strings.Repeat(" ", m.width-2)
	result.WriteString(leftBorder + emptyLine + rightBorder + "\n")

	// コンテンツを行ごとに分割してボーダーを追加
	lines := []string{
		leftAlignedContent,
		separator,
		statusArea,
		separator,
		footer,
	}

	for _, line := range lines {
		for _, l := range strings.Split(line, "\n") {
			if l != "" {
				result.WriteString(leftBorder + l + rightBorder + "\n")
			}
		}
	}

	// 下部ボーダー: ╰─────...╯
	bottomBorder := borderStyleColor.Render("╰" + strings.Repeat("─", m.width-2) + "╯")
	result.WriteString(bottomBorder)

	return result.String()
}

// エディション選択画面のコンテンツ
func (m model) viewSelectionContent() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		MarginBottom(1)

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF"))

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00D9FF")).
		Bold(true)

	var s strings.Builder
	s.WriteString(titleStyle.Render("Select Connection") + "\n")

	for i, choice := range m.choices {
		if m.cursor == i {
			s.WriteString(selectedStyle.Render(" > " + choice) + "\n")
		} else {
			s.WriteString(normalStyle.Render("   " + choice) + "\n")
		}
	}

	return s.String()
}

// On-Premise接続設定画面のコンテンツ
func (m model) viewOnPremiseConfigContent() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		MarginBottom(1)

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Width(11).
		Align(lipgloss.Left)

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF"))

	focusedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00D9FF")).
		Bold(true)

	var s strings.Builder

	s.WriteString(titleStyle.Render("On-Premise Connection") + "\n")

	// Endpoint
	endpointField := ui.TextField(m.onPremiseConfig.endpoint, 25, m.onPremiseConfig.focus == 0, m.onPremiseConfig.cursorPos)
	if m.onPremiseConfig.focus == 0 {
		s.WriteString(" " + labelStyle.Render("Endpoint:") + " " + focusedStyle.Render(endpointField) + "\n")
	} else {
		s.WriteString(" " + labelStyle.Render("Endpoint:") + " " + normalStyle.Render(endpointField) + "\n")
	}

	// Port
	portField := ui.TextField(m.onPremiseConfig.port, 8, m.onPremiseConfig.focus == 1, m.onPremiseConfig.cursorPos)
	if m.onPremiseConfig.focus == 1 {
		s.WriteString(" " + labelStyle.Render("Port:") + " " + focusedStyle.Render(portField) + "\n")
	} else {
		s.WriteString(" " + labelStyle.Render("Port:") + " " + normalStyle.Render(portField) + "\n")
	}

	// Secure checkbox
	secureText := ui.Checkbox("HTTPS/TLS", m.onPremiseConfig.secure, m.onPremiseConfig.focus == 2)
	s.WriteString(" " + labelStyle.Render("Secure:") + " " + secureText + "\n\n")

	// ボタン（縦配置）
	s.WriteString(" " + ui.Button("Test Connection", m.onPremiseConfig.focus == 3) + "\n")
	s.WriteString(" " + ui.Button("Connect", m.onPremiseConfig.focus == 4) + "\n")

	return s.String()
}

// Cloud接続設定画面のコンテンツ
func (m model) viewCloudConfigContent() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		MarginBottom(1)

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Width(15).
		Align(lipgloss.Left)

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF"))

	focusedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00D9FF")).
		Bold(true)

	var s strings.Builder

	s.WriteString(titleStyle.Render("Cloud Connection") + "\n")

	// Region
	regionField := ui.TextField(m.cloudConfig.region, 25, m.cloudConfig.focus == 0, m.cloudConfig.cursorPos)
	if m.cloudConfig.focus == 0 {
		s.WriteString(" " + labelStyle.Render("Region:") + " " + focusedStyle.Render(regionField) + "\n")
	} else {
		s.WriteString(" " + labelStyle.Render("Region:") + " " + normalStyle.Render(regionField) + "\n")
	}

	// Compartment
	compartmentField := ui.TextField(m.cloudConfig.compartment, 25, m.cloudConfig.focus == 1, m.cloudConfig.cursorPos)
	if m.cloudConfig.focus == 1 {
		s.WriteString(" " + labelStyle.Render("Compartment:") + " " + focusedStyle.Render(compartmentField) + "\n\n")
	} else {
		s.WriteString(" " + labelStyle.Render("Compartment:") + " " + normalStyle.Render(compartmentField) + "\n\n")
	}

	// Auth Method (ラジオボタン)
	s.WriteString(" " + labelStyle.Render("Auth Method:") + "\n")

	authMethods := []string{"OCI Config Profile (default)", "Instance Principal", "Resource Principal"}
	for i, method := range authMethods {
		focus := 2 + i
		s.WriteString(" " + ui.RadioButton(method, m.cloudConfig.authMethod == i, m.cloudConfig.focus == focus) + "\n")
	}
	s.WriteString("\n")

	// Config File
	configFileField := ui.TextField(m.cloudConfig.configFile, 25, m.cloudConfig.focus == 5, m.cloudConfig.cursorPos)
	if m.cloudConfig.focus == 5 {
		s.WriteString(" " + labelStyle.Render("Config File:") + " " + focusedStyle.Render(configFileField) + "\n\n")
	} else {
		s.WriteString(" " + labelStyle.Render("Config File:") + " " + normalStyle.Render(configFileField) + "\n\n")
	}

	// ボタン
	s.WriteString(" " + ui.Button("Test Connection", m.cloudConfig.focus == 6) + "\n")
	s.WriteString(" " + ui.Button("Connect", m.cloudConfig.focus == 7) + "\n")

	return s.String()
}

// カラム情報の構造体
type columnInfo struct {
	name string
	typ  string
}

// DDL文字列からカラム情報を抽出
// DDLからPRIMARY KEYのカラム名を抽出
func parsePrimaryKeysFromDDL(ddl string) []string {
	var primaryKeys []string

	// PRIMARY KEY(col1, col2, ...) の部分を探す
	upperDDL := strings.ToUpper(ddl)
	pkIndex := strings.Index(upperDDL, "PRIMARY KEY")
	if pkIndex == -1 {
		return primaryKeys
	}

	// PRIMARY KEYの後の括弧内を取得
	pkPart := ddl[pkIndex:]
	start := strings.Index(pkPart, "(")
	end := strings.LastIndex(pkPart, ")") // 最後の括弧を取得
	if start == -1 || end == -1 || start >= end {
		return primaryKeys
	}

	// カラム名リストを抽出
	keysPart := pkPart[start+1 : end]

	// SHARD()構文を処理
	// PRIMARY KEY(SHARD(id), name) のような形式に対応
	keysPart = strings.ReplaceAll(keysPart, "SHARD(", "")
	keysPart = strings.ReplaceAll(keysPart, "shard(", "")
	keysPart = strings.ReplaceAll(keysPart, ")", "")

	keys := strings.Split(keysPart, ",")
	for _, key := range keys {
		key = strings.TrimSpace(key)
		if key != "" {
			primaryKeys = append(primaryKeys, key)
		}
	}

	return primaryKeys
}

func parseColumnsFromDDL(ddl string, primaryKeys []string) []columnInfo {
	var columns []columnInfo

	// PRIMARY KEYのマップを作成（高速検索用）
	pkMap := make(map[string]bool)
	for _, pk := range primaryKeys {
		pkMap[pk] = true
	}

	// CREATE TABLE ... ( ... ) からカラム定義部分を抽出
	start := strings.Index(ddl, "(")
	end := strings.LastIndex(ddl, ")")
	if start == -1 || end == -1 || start >= end {
		return columns
	}

	columnDefs := ddl[start+1 : end]

	// PRIMARY KEY定義を除外
	lines := strings.Split(columnDefs, ",")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// PRIMARY KEY行をスキップ
		if strings.HasPrefix(strings.ToUpper(line), "PRIMARY KEY") {
			continue
		}

		// カラム名と型を抽出
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			name := parts[0]
			typ := parts[1]

			// PRIMARY KEYかどうかを判定
			if pkMap[name] {
				typ += " (Primary Key)"
			}

			columns = append(columns, columnInfo{name: name, typ: typ})
		}
	}

	return columns
}

// テーブル一覧画面のView
func (m model) viewTableList() string {
	// 2ペインレイアウト
	leftPaneWidth := 30 // 固定幅
	// rightPaneWidth = (borderの内側の幅) - (leftPaneWidth + leftPaneBorderRight)
	// = (m.width - 2) - (30 + 1) = m.width - 33
	rightPaneWidth := m.width - leftPaneWidth - 3

	// ヘッダー
	// borderStyleの内側の幅 m.width - 2 に合わせる
	// 右寄せで接続サーバ情報を表示
	rightText := "Connected to " + m.endpoint

	// 使用可能な幅（パディング分を引く）
	availableWidth := m.width - 4
	spaceBefore := availableWidth - len(rightText)
	if spaceBefore < 0 {
		spaceBefore = 0
	}

	headerContent := strings.Repeat(" ", spaceBefore) + rightText

	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Padding(0, 1).
		Width(m.width - 2)
	header := headerStyle.Render(headerContent)

	// 左ペイン: テーブルリスト
	// SelectableListを使用
	tableList := ui.SelectableList{
		Title:         fmt.Sprintf("Tables (%d)", len(m.tables)),
		Items:         m.tables,
		SelectedIndex: m.selectedTable,
		Focused:       m.rightPaneMode == rightPaneModeSchema, // スキーマビューモードの時のみフォーカス
	}
	leftPaneContent := tableList.Render()

	// ボーダー色の決定
	var borderColor string
	if m.rightPaneMode == rightPaneModeList || m.rightPaneMode == rightPaneModeDetail {
		borderColor = "#666666"
	} else {
		borderColor = "#555555"
	}
	leftPaneStyle := lipgloss.NewStyle().
		Width(leftPaneWidth).
		Height(m.height - 8). // タイトル行、ヘッダー、セパレーター×3、ステータス、フッター、ボーダー×2を除く
		BorderStyle(lipgloss.NormalBorder()).
		BorderRight(true).
		BorderForeground(lipgloss.Color(borderColor)).
		Padding(0, 1)
	leftPane := leftPaneStyle.Render(leftPaneContent)

	// 右ペイン: テーブル詳細またはデータ表示
	rightPaneContent := ""
	if len(m.tables) > 0 && m.selectedTable < len(m.tables) {
		selectedTableName := m.tables[m.selectedTable]

		// モードに応じてヘッダーを変更
		if m.rightPaneMode == rightPaneModeList || m.rightPaneMode == rightPaneModeDetail {
			// グリッドビュー/レコードビューモード: SQLエリアを表示
			if data, exists := m.tableData[selectedTableName]; exists {
				// SQLエリアのスタイル（背景なし）
				sqlStyle := lipgloss.NewStyle().
					Foreground(lipgloss.Color("#CCCCCC"))

				// SQLとセパレーターを手動で組み立て
				sqlText := sqlStyle.Render(data.displaySQL)
				separator := ui.Separator(rightPaneWidth - 2)

				rightPaneContent = sqlText + "\n" + separator
			}
		}

		if m.rightPaneMode == rightPaneModeSchema {
			// Schema表示モード
			// ヘッダーを表示
			rightPaneContent += fmt.Sprintf("Table:    %s", selectedTableName) + "\n"

			// 親子関係の判定
			if strings.Contains(selectedTableName, ".") {
				// 子テーブル
				parts := strings.Split(selectedTableName, ".")
				rightPaneContent += fmt.Sprintf("Parent:   %s\n", parts[0])
				rightPaneContent += "Children: -\n"
			} else {
				// 親テーブル - 子テーブルをカウント
				rightPaneContent += "Parent:   -\n"
				childCount := 0
				var childNames []string
				prefix := selectedTableName + "."
				for _, t := range m.tables {
					if strings.HasPrefix(t, prefix) {
						childCount++
						childNames = append(childNames, strings.TrimPrefix(t, prefix))
					}
				}
				if childCount > 0 {
					rightPaneContent += fmt.Sprintf("Children: %s\n", strings.Join(childNames, ", "))
				} else {
					rightPaneContent += "Children: -\n"
				}
			}

			// カラム情報とインデックス情報を表示
			if details, exists := m.tableDetails[selectedTableName]; exists && details != nil {
				// カラム情報（DDL文字列から抽出）
				rightPaneContent += "\nColumns:\n"
				if details.schema != nil && details.schema.DDL != "" {
					// DDLからカラム情報を抽出
					primaryKeys := parsePrimaryKeysFromDDL(details.schema.DDL)
					columns := parseColumnsFromDDL(details.schema.DDL, primaryKeys)
					if len(columns) > 0 {
						for _, col := range columns {
							rightPaneContent += fmt.Sprintf("  %-20s %s\n", col.name, col.typ)
						}
					} else {
						rightPaneContent += "  (No column information available)\n"
					}
				} else if details.schema != nil && details.schema.Schema != "" {
					rightPaneContent += "  " + details.schema.Schema + "\n"
				} else {
					rightPaneContent += "  (No column information available)\n"
				}

				// インデックス情報
				rightPaneContent += "\nIndexes:\n"
				if len(details.indexes) > 0 {
					for _, index := range details.indexes {
						fields := strings.Join(index.FieldNames, ", ")
						rightPaneContent += fmt.Sprintf("  %-20s (%s)\n", index.IndexName, fields)
					}
				} else {
					rightPaneContent += "  (none)\n"
				}
			} else if m.loadingDetails {
				rightPaneContent += "\nColumns:\n"
				rightPaneContent += "  Loading...\n"
				rightPaneContent += "\nIndexes:\n"
				rightPaneContent += "  Loading...\n"
			} else {
				rightPaneContent += "\nColumns:\n"
				rightPaneContent += "  (Schema information will be displayed here)\n"
				rightPaneContent += "\nIndexes:\n"
				rightPaneContent += "  (Index information will be displayed here)\n"
			}
		} else if m.rightPaneMode == rightPaneModeList {
			// グリッドビューモード
			rightPaneHeight := m.height - 8

			// データの取得状態を確認
			data, exists := m.tableData[selectedTableName]
			if m.loadingData {
				rightPaneContent += "Loading data..."
			} else if !exists || data == nil {
				rightPaneContent += "No data available"
			} else if data.err != nil {
				rightPaneContent += fmt.Sprintf("Error: %v\n\nSQL:\n%s", data.err, data.sql)
			} else if len(data.rows) == 0 {
				rightPaneContent += fmt.Sprintf("No data found\n\nSQL:\n%s", data.sql)
			} else {
				// カラム名をDDL定義順で取得
				var columnNames []string
				if details, exists := m.tableDetails[selectedTableName]; exists && details.schema != nil && details.schema.DDL != "" {
					primaryKeys := parsePrimaryKeysFromDDL(details.schema.DDL)
					columns := parseColumnsFromDDL(details.schema.DDL, primaryKeys)
					for _, col := range columns {
						columnNames = append(columnNames, col.name)
					}
				} else if len(data.rows) > 0 {
					// DDLが取得できない場合は、データから取得（順序は不定）
					for key := range data.rows[0] {
						columnNames = append(columnNames, key)
					}
				}

				// ui.DataGridを使用してレンダリング
				grid := ui.DataGrid{
					Rows:             data.rows,
					Columns:          columnNames,
					SelectedRow:      m.selectedDataRow,
					HorizontalOffset: m.horizontalOffset,
					ViewportOffset:   m.viewportOffset,
				}
				rightPaneContent += grid.Render(rightPaneWidth, rightPaneHeight)
			}
		} else if m.rightPaneMode == rightPaneModeDetail {
			// レコードビューモード
			// データの取得状態を確認
			data, exists := m.tableData[selectedTableName]
			if m.loadingData {
				rightPaneContent += "Loading data..."
			} else if !exists || data == nil {
				rightPaneContent += "No data available"
			} else if data.err != nil {
				rightPaneContent += fmt.Sprintf("Error: %v", data.err)
			} else if len(data.rows) == 0 {
				rightPaneContent += "No data found"
			} else if m.selectedDataRow < 0 || m.selectedDataRow >= len(data.rows) {
				rightPaneContent += "Invalid row selection"
			} else {
				// 選択された行を取得
				selectedRow := data.rows[m.selectedDataRow]

				// カラム名をDDL定義順で取得
				var columnNames []string
				if details, exists := m.tableDetails[selectedTableName]; exists && details.schema != nil && details.schema.DDL != "" {
					primaryKeys := parsePrimaryKeysFromDDL(details.schema.DDL)
					columns := parseColumnsFromDDL(details.schema.DDL, primaryKeys)
					for _, col := range columns {
						columnNames = append(columnNames, col.name)
					}
				} else if len(selectedRow) > 0 {
					// DDLが取得できない場合は、データから取得（順序は不定）
					for key := range selectedRow {
						columnNames = append(columnNames, key)
					}
				}

				// ui.VerticalTableを使用してレンダリング
				verticalTable := ui.VerticalTable{
					Data: selectedRow,
					Keys: columnNames,
				}
				rightPaneContent += verticalTable.Render()
			}
		}
	}

	rightPaneStyle := lipgloss.NewStyle().
		Width(rightPaneWidth).
		Height(m.height - 8).
		Padding(0, 1)
	// 末尾の空行を削除
	rightPaneContent = strings.TrimSuffix(rightPaneContent, "\n")
	rightPane := rightPaneStyle.Render(rightPaneContent)

	// 2ペインを横に並べる
	panes := lipgloss.JoinHorizontal(lipgloss.Top, leftPane, rightPane)

	// ステータスバー
	statusBarStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#CCCCCC")).
		Padding(0, 1).
		Width(m.width - 2)
	var status string
	if len(m.tables) > 0 {
		selectedTableName := m.tables[m.selectedTable]
		if m.rightPaneMode == rightPaneModeList || m.rightPaneMode == rightPaneModeDetail {
			// グリッドビュー/レコードビューモード: テーブル名と行数を表示
			if data, exists := m.tableData[selectedTableName]; exists {
				if data.err != nil {
					// エラーが発生した場合は赤色で表示
					errorStyle := lipgloss.NewStyle().
						Foreground(lipgloss.Color("#FF0000")).
						Padding(0, 1)
					status = errorStyle.Render(fmt.Sprintf("Error: %v", data.err))
				} else if len(data.rows) > 0 {
					totalRows := len(data.rows)
					// データがまだある場合は "+" を追加
					moreIndicator := ""
					if data.hasMore {
						moreIndicator = "+"
					}
					// テーブル名と行数のみ表示
					status = statusBarStyle.Render(fmt.Sprintf("Table: %s (%d%s rows)", selectedTableName, totalRows, moreIndicator))
				} else {
					status = statusBarStyle.Render(fmt.Sprintf("Table: %s (0 rows)", selectedTableName))
				}
			} else if m.loadingData {
				status = statusBarStyle.Render(fmt.Sprintf("Table: %s (loading...)", selectedTableName))
			} else {
				status = statusBarStyle.Render(fmt.Sprintf("Table: %s", selectedTableName))
			}
		} else {
			// スキーマビューモード: テーブル名のみ表示
			status = statusBarStyle.Render(fmt.Sprintf("Table: %s", selectedTableName))
		}
	} else {
		status = statusBarStyle.Render("")
	}

	// フッター
	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Padding(0, 1).
		Width(m.width - 2)
	var footer string
	if m.rightPaneMode == rightPaneModeList {
		footer = footerStyle.Render("j/k: Scroll  h/l: Scroll Left/Right  o: Detail  u: Back  q: Quit")
	} else if m.rightPaneMode == rightPaneModeDetail {
		footer = footerStyle.Render("j/k: Scroll  u: Back to List  q: Quit")
	} else {
		footer = footerStyle.Render("j/k: Navigate  o: View Data  u: Back  q: Quit")
	}

	// セパレーター
	topSeparator := ui.Separator(m.width - 2)
	statusSeparator := ui.Separator(m.width - 2)

	// 全体を組み立て
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		topSeparator,
		panes,
		statusSeparator,
		status,
		statusSeparator,
		footer,
	)

	// 手動でボーダーを描画
	borderStyleColor := lipgloss.NewStyle().Foreground(lipgloss.Color("#00D9FF"))

	// 上部ボーダー: ╭── Dito ─────...╮
	title := " Dito "
	// 全体の幅 = m.width
	// "╭──" = 3文字, title = 6文字, "╮" = 1文字
	// 残りの "─" = m.width - 3 - 6 - 1 = m.width - 10
	topBorder := borderStyleColor.Render("╭──" + title + strings.Repeat("─", m.width-10) + "╮")

	// 左右のボーダー文字
	leftBorder := borderStyleColor.Render("│")
	rightBorder := borderStyleColor.Render("│")

	// コンテンツの各行にボーダーを追加
	var result strings.Builder
	result.WriteString(topBorder + "\n")

	for _, line := range strings.Split(content, "\n") {
		result.WriteString(leftBorder + line + rightBorder + "\n")
	}

	// 下部ボーダー: ╰─────...╯
	bottomBorder := borderStyleColor.Render("╰" + strings.Repeat("─", m.width-2) + "╯")
	result.WriteString(bottomBorder)

	return result.String()
}

func main() {
	p := tea.NewProgram(
		initialModel(),
		tea.WithAltScreen(),       // 全画面モード
		tea.WithMouseCellMotion(), // マウスサポート（オプション）
	)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
