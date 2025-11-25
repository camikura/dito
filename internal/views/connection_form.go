package views

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/camikura/dito/internal/ui"
)

// OnPremiseFormModel represents the data needed to render the on-premise connection form.
type OnPremiseFormModel struct {
	Endpoint string
	Port     string
	Secure   bool
	Focus    int // Index of focused field (0: endpoint, 1: port, 2: secure, 3: test button, 4: connect button)
}

// CloudFormModel represents the data needed to render the cloud connection form.
type CloudFormModel struct {
	Region      string
	Compartment string
	AuthMethod  int    // 0: OCI Config Profile, 1: Instance Principal, 2: Resource Principal
	ConfigFile  string
	Focus       int // Index of focused field
}

// RenderOnPremiseForm renders the on-premise connection configuration form.
// This is a pure rendering function that takes model data and returns a string.
func RenderOnPremiseForm(m OnPremiseFormModel) string {
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
	endpointField := ui.TextField(m.Endpoint, 25, m.Focus == 0, 0)
	if m.Focus == 0 {
		s.WriteString(" " + labelStyle.Render("Endpoint:") + " " + focusedStyle.Render(endpointField) + "\n")
	} else {
		s.WriteString(" " + labelStyle.Render("Endpoint:") + " " + normalStyle.Render(endpointField) + "\n")
	}

	// Port
	portField := ui.TextField(m.Port, 8, m.Focus == 1, 0)
	if m.Focus == 1 {
		s.WriteString(" " + labelStyle.Render("Port:") + " " + focusedStyle.Render(portField) + "\n")
	} else {
		s.WriteString(" " + labelStyle.Render("Port:") + " " + normalStyle.Render(portField) + "\n")
	}

	// Secure checkbox
	secureText := ui.Checkbox("HTTPS/TLS", m.Secure, m.Focus == 2)
	s.WriteString(" " + labelStyle.Render("Secure:") + " " + secureText + "\n\n")

	// Buttons (vertical layout)
	s.WriteString(" " + ui.Button("Test Connection", m.Focus == 3) + "\n")
	s.WriteString(" " + ui.Button("Connect", m.Focus == 4) + "\n")

	return s.String()
}

// RenderCloudForm renders the cloud connection configuration form.
// This is a pure rendering function that takes model data and returns a string.
func RenderCloudForm(m CloudFormModel) string {
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
	regionField := ui.TextField(m.Region, 25, m.Focus == 0, 0)
	if m.Focus == 0 {
		s.WriteString(" " + labelStyle.Render("Region:") + " " + focusedStyle.Render(regionField) + "\n")
	} else {
		s.WriteString(" " + labelStyle.Render("Region:") + " " + normalStyle.Render(regionField) + "\n")
	}

	// Compartment
	compartmentField := ui.TextField(m.Compartment, 25, m.Focus == 1, 0)
	if m.Focus == 1 {
		s.WriteString(" " + labelStyle.Render("Compartment:") + " " + focusedStyle.Render(compartmentField) + "\n\n")
	} else {
		s.WriteString(" " + labelStyle.Render("Compartment:") + " " + normalStyle.Render(compartmentField) + "\n\n")
	}

	// Auth Method (radio buttons)
	s.WriteString(" " + labelStyle.Render("Auth Method:") + "\n")

	authMethods := []string{"OCI Config Profile (default)", "Instance Principal", "Resource Principal"}
	for i, method := range authMethods {
		focus := 2 + i
		s.WriteString(" " + ui.RadioButton(method, m.AuthMethod == i, m.Focus == focus) + "\n")
	}
	s.WriteString("\n")

	// Config File
	configFileField := ui.TextField(m.ConfigFile, 25, m.Focus == 5, 0)
	if m.Focus == 5 {
		s.WriteString(" " + labelStyle.Render("Config File:") + " " + focusedStyle.Render(configFileField) + "\n\n")
	} else {
		s.WriteString(" " + labelStyle.Render("Config File:") + " " + normalStyle.Render(configFileField) + "\n\n")
	}

	// Buttons
	s.WriteString(" " + ui.Button("Test Connection", m.Focus == 6) + "\n")
	s.WriteString(" " + ui.Button("Connect", m.Focus == 7) + "\n")

	return s.String()
}
