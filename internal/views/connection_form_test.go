package views

import (
	"strings"
	"testing"
)

func TestRenderOnPremiseForm(t *testing.T) {
	tests := []struct {
		name     string
		model    OnPremiseFormModel
		contains []string // strings that should be in the output
	}{
		{
			name: "basic form with endpoint focused",
			model: OnPremiseFormModel{
				Endpoint:  "localhost",
				Port:      "8080",
				Secure:    false,
				Focus:     0,
				CursorPos: 9,
			},
			contains: []string{"On-Premise Connection", "Endpoint:", "Port:", "Secure:", "Test Connection", "Connect"},
		},
		{
			name: "form with port focused",
			model: OnPremiseFormModel{
				Endpoint:  "example.com",
				Port:      "443",
				Secure:    true,
				Focus:     1,
				CursorPos: 3,
			},
			contains: []string{"On-Premise Connection", "Endpoint:", "Port:", "Secure:", "HTTPS/TLS"},
		},
		{
			name: "form with secure checkbox focused",
			model: OnPremiseFormModel{
				Endpoint:  "db.example.com",
				Port:      "8080",
				Secure:    false,
				Focus:     2,
				CursorPos: 0,
			},
			contains: []string{"On-Premise Connection", "Secure:", "HTTPS/TLS"},
		},
		{
			name: "form with test button focused",
			model: OnPremiseFormModel{
				Endpoint:  "localhost",
				Port:      "8080",
				Secure:    false,
				Focus:     3,
				CursorPos: 0,
			},
			contains: []string{"Test Connection"},
		},
		{
			name: "form with connect button focused",
			model: OnPremiseFormModel{
				Endpoint:  "localhost",
				Port:      "8080",
				Secure:    false,
				Focus:     4,
				CursorPos: 0,
			},
			contains: []string{"Connect"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderOnPremiseForm(tt.model)

			for _, substr := range tt.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("RenderOnPremiseForm() = %q, should contain %q", result, substr)
				}
			}

			if result == "" {
				t.Error("RenderOnPremiseForm() should not return empty string")
			}
		})
	}
}

func TestRenderCloudForm(t *testing.T) {
	tests := []struct {
		name     string
		model    CloudFormModel
		contains []string // strings that should be in the output
	}{
		{
			name: "basic form with region focused",
			model: CloudFormModel{
				Region:      "us-ashburn-1",
				Compartment: "ocid1.compartment.oc1..example",
				AuthMethod:  0,
				ConfigFile:  "~/.oci/config",
				Focus:       0,
				CursorPos:   12,
			},
			contains: []string{"Cloud Connection", "Region:", "Compartment:", "Auth Method:", "Config File:", "Test Connection", "Connect"},
		},
		{
			name: "form with compartment focused",
			model: CloudFormModel{
				Region:      "us-phoenix-1",
				Compartment: "my-compartment",
				AuthMethod:  0,
				ConfigFile:  "~/.oci/config",
				Focus:       1,
				CursorPos:   5,
			},
			contains: []string{"Cloud Connection", "Compartment:"},
		},
		{
			name: "form with first auth method focused",
			model: CloudFormModel{
				Region:      "us-ashburn-1",
				Compartment: "ocid1.compartment.oc1..example",
				AuthMethod:  0,
				ConfigFile:  "~/.oci/config",
				Focus:       2,
				CursorPos:   0,
			},
			contains: []string{"OCI Config Profile (default)", "Instance Principal", "Resource Principal"},
		},
		{
			name: "form with instance principal selected",
			model: CloudFormModel{
				Region:      "us-ashburn-1",
				Compartment: "ocid1.compartment.oc1..example",
				AuthMethod:  1,
				ConfigFile:  "",
				Focus:       3,
				CursorPos:   0,
			},
			contains: []string{"Instance Principal"},
		},
		{
			name: "form with resource principal selected",
			model: CloudFormModel{
				Region:      "us-ashburn-1",
				Compartment: "ocid1.compartment.oc1..example",
				AuthMethod:  2,
				ConfigFile:  "",
				Focus:       4,
				CursorPos:   0,
			},
			contains: []string{"Resource Principal"},
		},
		{
			name: "form with config file focused",
			model: CloudFormModel{
				Region:      "us-ashburn-1",
				Compartment: "ocid1.compartment.oc1..example",
				AuthMethod:  0,
				ConfigFile:  "/custom/path/config",
				Focus:       5,
				CursorPos:   10,
			},
			contains: []string{"Config File:"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderCloudForm(tt.model)

			for _, substr := range tt.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("RenderCloudForm() = %q, should contain %q", result, substr)
				}
			}

			if result == "" {
				t.Error("RenderCloudForm() should not return empty string")
			}
		})
	}
}
