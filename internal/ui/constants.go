package ui

import "time"

// Layout constants
const (
	// LeftPaneContentWidth is the content width for left panes (excluding borders).
	LeftPaneContentWidth = 50

	// LeftPaneBorderWidth is the total border width (left + right).
	LeftPaneBorderWidth = 2

	// MinRightPaneWidth is the minimum width for the right (data) pane.
	MinRightPaneWidth = 10

	// MinWindowHeight is the minimum window height before showing "Window too short".
	MinWindowHeight = 20

	// ConnectionDialogWidth is the width of the connection setup dialog.
	ConnectionDialogWidth = 60

	// DialogSizeRatio is the ratio of dialog size to screen size (4/5 = 80%).
	DialogSizeRatio   = 4
	DialogSizeDivisor = 5

	// LayoutBorderOverhead is the total vertical space taken by borders in the left pane layout.
	// Tables(2) + Schema(2) + SQL(2) = 6
	LayoutBorderOverhead = 6

	// FooterHeight is the height of the footer.
	FooterHeight = 1

	// DataPaneHeaderLines is the number of lines for header + separator in data pane.
	DataPaneHeaderLines = 2

	// DataPaneTitleAndBorderLines is the number of lines for title + borders in data pane.
	DataPaneTitleAndBorderLines = 3

	// MinContentLines is the minimum number of content lines in data pane.
	MinContentLines = 5

	// DefaultConnectionPaneHeight is the default height for connection pane.
	DefaultConnectionPaneHeight = 5

	// PaneBorderHeight is the height of borders for each pane.
	PaneBorderHeight = 2
)

// Data loading constants
const (
	// DefaultFetchSize is the default number of rows to fetch at once.
	DefaultFetchSize = 100

	// FetchMoreThreshold is the number of remaining rows that triggers fetching more data.
	FetchMoreThreshold = 10
)

// Pane height ratio constants
const (
	// PaneHeightTotalParts is the total parts for height distribution (2:2:1 ratio).
	PaneHeightTotalParts = 5

	// PaneHeightTablesParts is the parts for tables pane (2 of 5).
	PaneHeightTablesParts = 2

	// PaneHeightSchemaParts is the parts for schema pane (2 of 5).
	PaneHeightSchemaParts = 2

	// PaneHeightSQLParts is the parts for SQL pane (1 of 5).
	PaneHeightSQLParts = 1
)

// Scroll constants
const (
	// PageScrollAmount is the number of lines to scroll for page up/down.
	PageScrollAmount = 10
)

// Message duration constants
const (
	// CopyMessageDuration is how long to show copy success/failure message.
	CopyMessageDuration = 2 * time.Second

	// QuitConfirmationTimeout is how long quit confirmation remains active.
	QuitConfirmationTimeout = 3 * time.Second
)
