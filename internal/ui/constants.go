package ui

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
