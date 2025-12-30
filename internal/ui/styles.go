// Package ui provides shared styling components for the CertWatch Agent CLI.
// All commands should use this package for consistent visual appearance.
package ui

import (
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// Theme colors - CertWatch brand colors
var (
	ColorPrimary   = lipgloss.Color("#0EA5E9") // Sky blue
	ColorSuccess   = lipgloss.Color("#22C55E") // Green
	ColorWarning   = lipgloss.Color("#F59E0B") // Amber
	ColorError     = lipgloss.Color("#EF4444") // Red
	ColorMuted     = lipgloss.Color("#6B7280") // Gray
	ColorHighlight = lipgloss.Color("#A855F7") // Purple
	ColorDark      = lipgloss.Color("#1F2937") // Dark gray
	ColorLight     = lipgloss.Color("#F9FAFB") // Light gray
)

// Component styles
var (
	// TitleStyle for section headers
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			MarginBottom(1)

	// SuccessStyle for success messages
	SuccessStyle = lipgloss.NewStyle().
			Foreground(ColorSuccess).
			Bold(true)

	// ErrorStyle for error messages
	ErrorStyle = lipgloss.NewStyle().
			Foreground(ColorError).
			Bold(true)

	// WarningStyle for warning messages
	WarningStyle = lipgloss.NewStyle().
			Foreground(ColorWarning)

	// MutedStyle for secondary text
	MutedStyle = lipgloss.NewStyle().
			Foreground(ColorMuted)

	// CodeStyle for command/code display
	CodeStyle = lipgloss.NewStyle().
			Background(ColorDark).
			Foreground(ColorLight).
			Padding(0, 1)

	// BoxStyle for summary sections
	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorPrimary).
			Padding(1, 2).
			MarginTop(1)

	// WarningBoxStyle for warning boxes
	WarningBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorWarning).
			Padding(1, 2).
			MarginTop(1)

	// HeaderStyle for the main header
	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			Background(lipgloss.Color("#1E3A5F")).
			Padding(0, 2).
			MarginBottom(1)

	// SectionStyle for section dividers
	SectionStyle = lipgloss.NewStyle().
			Foreground(ColorMuted).
			MarginTop(1).
			MarginBottom(1)

	// LabelStyle for key-value pair labels (fixed width for alignment)
	LabelStyle = lipgloss.NewStyle().
			Foreground(ColorMuted).
			Width(14)

	// ValueStyle for key-value pair values
	ValueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF"))
)

// Prefixes for messages
const (
	SuccessPrefix = "✓ "
	ErrorPrefix   = "✗ "
	WarningPrefix = "! "
	InfoPrefix    = "→ "
)

// DefaultAPIEndpoint is the default API endpoint that should not be displayed
const DefaultAPIEndpoint = "https://api.certwatch.app"

// CreateTheme returns a custom huh theme matching CertWatch branding.
func CreateTheme() *huh.Theme {
	t := huh.ThemeBase()

	// Customize focused state
	t.Focused.Title = t.Focused.Title.Foreground(ColorPrimary)
	t.Focused.Description = t.Focused.Description.Foreground(ColorMuted)
	t.Focused.SelectSelector = t.Focused.SelectSelector.Foreground(ColorHighlight)
	t.Focused.SelectedOption = t.Focused.SelectedOption.Foreground(ColorPrimary)
	t.Focused.TextInput.Cursor = t.Focused.TextInput.Cursor.Foreground(ColorPrimary)

	// Customize blurred state
	t.Blurred.Title = t.Blurred.Title.Foreground(ColorMuted)

	return t
}

// RenderAppHeader renders the main " CertWatch Agent " banner.
func RenderAppHeader() string {
	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorLight).
		Background(ColorPrimary).
		Padding(0, 2).
		Render(" CertWatch Agent ")

	return header
}

// RenderCommandHeader renders a " CertWatch Agent - {context} " banner.
func RenderCommandHeader(context string) string {
	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorLight).
		Background(ColorPrimary).
		Padding(0, 2).
		Render(" CertWatch Agent - " + context + " ")

	return header
}

// RenderSection renders a section divider with title.
func RenderSection(title string) string {
	line := "─────────────────────────────────────────"
	maxLen := 40 - len(title)
	if maxLen < 0 {
		maxLen = 0
	}
	return SectionStyle.Render("─── " + title + " " + line[:maxLen])
}

// RenderSuccess renders a success message with green checkmark.
func RenderSuccess(msg string) string {
	return SuccessStyle.Render(SuccessPrefix + msg)
}

// RenderError renders an error message with red X.
func RenderError(msg string) string {
	return ErrorStyle.Render(ErrorPrefix + msg)
}

// RenderWarning renders a warning message with amber exclamation.
func RenderWarning(msg string) string {
	return WarningStyle.Render(WarningPrefix + msg)
}

// RenderInfo renders an info message with muted arrow.
func RenderInfo(msg string) string {
	return MutedStyle.Render(InfoPrefix + msg)
}

// RenderCode renders a code/command with dark background.
func RenderCode(code string) string {
	return CodeStyle.Render(code)
}

// RenderKeyValue renders a formatted "  Label  Value" line with alignment.
func RenderKeyValue(label, value string) string {
	return "  " + LabelStyle.Render(label) + ValueStyle.Render(value)
}

// RenderKeyValueList renders multiple key-value pairs.
func RenderKeyValueList(items [][2]string) string {
	lines := make([]string, 0, len(items))
	for _, item := range items {
		lines = append(lines, RenderKeyValue(item[0], item[1]))
	}
	return strings.Join(lines, "\n")
}

// RenderWarningBox renders a warning box with multiple lines.
func RenderWarningBox(title string, lines []string) string {
	content := WarningStyle.Render(WarningPrefix+title) + "\n"
	for _, line := range lines {
		content += "\n" + line
	}
	return WarningBoxStyle.Render(content)
}

// TruncateID truncates a UUID for display (shows first 12 chars + ...)
func TruncateID(id string) string {
	if len(id) > 12 {
		return id[:12] + "..."
	}
	return id
}
