package app

import (
	"testing"

	"github.com/camikura/dito/internal/db"
)

func TestInitialModel(t *testing.T) {
	model := InitialModel()

	// Check initial screen
	if model.Screen != ScreenSelection {
		t.Errorf("InitialModel() Screen = %v, want %v", model.Screen, ScreenSelection)
	}

	// Check choices
	expectedChoices := []string{"Oracle NoSQL Cloud Service", "On-Premise"}
	if len(model.Choices) != len(expectedChoices) {
		t.Errorf("InitialModel() Choices length = %d, want %d", len(model.Choices), len(expectedChoices))
	}
	for i, choice := range expectedChoices {
		if model.Choices[i] != choice {
			t.Errorf("InitialModel() Choices[%d] = %v, want %v", i, model.Choices[i], choice)
		}
	}

	// Check selected map is initialized
	if model.Selected == nil {
		t.Error("InitialModel() Selected map should be initialized")
	}

	// Check OnPremiseConfig defaults
	if model.OnPremiseConfig.Endpoint != "localhost" {
		t.Errorf("InitialModel() OnPremiseConfig.Endpoint = %v, want localhost", model.OnPremiseConfig.Endpoint)
	}
	if model.OnPremiseConfig.Port != "8080" {
		t.Errorf("InitialModel() OnPremiseConfig.Port = %v, want 8080", model.OnPremiseConfig.Port)
	}
	if model.OnPremiseConfig.Secure != false {
		t.Error("InitialModel() OnPremiseConfig.Secure should be false")
	}
	if model.OnPremiseConfig.Status != StatusDisconnected {
		t.Errorf("InitialModel() OnPremiseConfig.Status = %v, want %v", model.OnPremiseConfig.Status, StatusDisconnected)
	}

	// Check CloudConfig defaults
	if model.CloudConfig.Region != "us-ashburn-1" {
		t.Errorf("InitialModel() CloudConfig.Region = %v, want us-ashburn-1", model.CloudConfig.Region)
	}
	if model.CloudConfig.Compartment != "" {
		t.Error("InitialModel() CloudConfig.Compartment should be empty")
	}
	if model.CloudConfig.AuthMethod != 0 {
		t.Errorf("InitialModel() CloudConfig.AuthMethod = %d, want 0", model.CloudConfig.AuthMethod)
	}
	if model.CloudConfig.ConfigFile != "DEFAULT" {
		t.Errorf("InitialModel() CloudConfig.ConfigFile = %v, want DEFAULT", model.CloudConfig.ConfigFile)
	}
	if model.CloudConfig.Status != StatusDisconnected {
		t.Errorf("InitialModel() CloudConfig.Status = %v, want %v", model.CloudConfig.Status, StatusDisconnected)
	}

	// Check table-related maps are initialized
	if model.TableDetails == nil {
		t.Error("InitialModel() TableDetails map should be initialized")
	}
	if model.TableData == nil {
		t.Error("InitialModel() TableData map should be initialized")
	}

	// Check RightPaneMode default
	if model.RightPaneMode != RightPaneModeSchema {
		t.Errorf("InitialModel() RightPaneMode = %v, want %v", model.RightPaneMode, RightPaneModeSchema)
	}

	// Check fetch size
	if model.FetchSize != 100 {
		t.Errorf("InitialModel() FetchSize = %d, want 100", model.FetchSize)
	}

	// Check viewport size
	if model.ViewportSize != 10 {
		t.Errorf("InitialModel() ViewportSize = %d, want 10", model.ViewportSize)
	}
}

func TestToTableListViewModel(t *testing.T) {
	tests := []struct {
		name  string
		model Model
	}{
		{
			name: "basic model conversion",
			model: Model{
				Width:            100,
				Height:           30,
				Endpoint:         "localhost:8080",
				Tables:           []string{"users", "products"},
				SelectedTable:    1,
				RightPaneMode:    RightPaneModeList,
				TableData:        map[string]*db.TableDataResult{},
				TableDetails:     map[string]*db.TableDetailsResult{},
				LoadingDetails:   false,
				LoadingData:      true,
				SelectedDataRow:  5,
				HorizontalOffset: 2,
				ViewportOffset:   3,
			},
		},
		{
			name: "model with nil maps",
			model: Model{
				Width:            80,
				Height:           25,
				Endpoint:         "example.com:8080",
				Tables:           []string{"table1"},
				SelectedTable:    0,
				RightPaneMode:    RightPaneModeSchema,
				TableData:        nil,
				TableDetails:     nil,
				LoadingDetails:   true,
				LoadingData:      false,
				SelectedDataRow:  0,
				HorizontalOffset: 0,
				ViewportOffset:   0,
			},
		},
		{
			name: "model in detail view mode",
			model: Model{
				Width:            120,
				Height:           40,
				Endpoint:         "cloud.oracle.com",
				Tables:           []string{"customers", "orders", "products"},
				SelectedTable:    2,
				RightPaneMode:    RightPaneModeDetail,
				TableData:        map[string]*db.TableDataResult{},
				TableDetails:     map[string]*db.TableDetailsResult{},
				LoadingDetails:   false,
				LoadingData:      false,
				SelectedDataRow:  10,
				HorizontalOffset: 0,
				ViewportOffset:   5,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vm := tt.model.ToTableListViewModel()

			// Check all fields are correctly copied
			if vm.Width != tt.model.Width {
				t.Errorf("ToTableListViewModel() Width = %v, want %v", vm.Width, tt.model.Width)
			}
			if vm.Height != tt.model.Height {
				t.Errorf("ToTableListViewModel() Height = %v, want %v", vm.Height, tt.model.Height)
			}
			if vm.Endpoint != tt.model.Endpoint {
				t.Errorf("ToTableListViewModel() Endpoint = %v, want %v", vm.Endpoint, tt.model.Endpoint)
			}
			if len(vm.Tables) != len(tt.model.Tables) {
				t.Errorf("ToTableListViewModel() Tables length = %d, want %d", len(vm.Tables), len(tt.model.Tables))
			}
			for i := range tt.model.Tables {
				if vm.Tables[i] != tt.model.Tables[i] {
					t.Errorf("ToTableListViewModel() Tables[%d] = %v, want %v", i, vm.Tables[i], tt.model.Tables[i])
				}
			}
			if vm.SelectedTable != tt.model.SelectedTable {
				t.Errorf("ToTableListViewModel() SelectedTable = %v, want %v", vm.SelectedTable, tt.model.SelectedTable)
			}
			if vm.RightPaneMode != tt.model.RightPaneMode {
				t.Errorf("ToTableListViewModel() RightPaneMode = %v, want %v", vm.RightPaneMode, tt.model.RightPaneMode)
			}
			if vm.LoadingDetails != tt.model.LoadingDetails {
				t.Errorf("ToTableListViewModel() LoadingDetails = %v, want %v", vm.LoadingDetails, tt.model.LoadingDetails)
			}
			if vm.LoadingData != tt.model.LoadingData {
				t.Errorf("ToTableListViewModel() LoadingData = %v, want %v", vm.LoadingData, tt.model.LoadingData)
			}
			if vm.SelectedDataRow != tt.model.SelectedDataRow {
				t.Errorf("ToTableListViewModel() SelectedDataRow = %v, want %v", vm.SelectedDataRow, tt.model.SelectedDataRow)
			}
			if vm.HorizontalOffset != tt.model.HorizontalOffset {
				t.Errorf("ToTableListViewModel() HorizontalOffset = %v, want %v", vm.HorizontalOffset, tt.model.HorizontalOffset)
			}
			if vm.ViewportOffset != tt.model.ViewportOffset {
				t.Errorf("ToTableListViewModel() ViewportOffset = %v, want %v", vm.ViewportOffset, tt.model.ViewportOffset)
			}

			// Check that maps are set (we can't directly compare map references in Go)
			if tt.model.TableData != nil && vm.TableData == nil {
				t.Error("ToTableListViewModel() TableData should not be nil when source is not nil")
			}
			if tt.model.TableDetails != nil && vm.TableDetails == nil {
				t.Error("ToTableListViewModel() TableDetails should not be nil when source is not nil")
			}
		})
	}
}
