package ui

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"os"
	"strings"
)

// Color palette
var (
	colorError   = lipgloss.Color("#FF6B6B")
	colorSuccess = lipgloss.Color("#51CF66")
	colorWarning = lipgloss.Color("#FFD43B")
	colorInfo    = lipgloss.Color("#74C0FC")
	colorCode    = lipgloss.Color("#FFA8A8")
	colorCodeBg  = lipgloss.Color("#2D1B1B")
	colorTable   = lipgloss.Color("#8BE9FD")
)

// Style templates
var (
	baseStyle = lipgloss.NewStyle().
			Bold(true).
			Padding(0, 1).
			Border(lipgloss.RoundedBorder())

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Padding(0, 2).
			MarginBottom(1)

	codeStyle = lipgloss.NewStyle().
			Foreground(colorCode).
			Background(colorCodeBg).
			Padding(0, 1).
			Italic(true)

	// Table styles
	tableStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorTable).
			Padding(0, 1).
			MarginBottom(1)

	tableHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(colorTable).
				Padding(0, 1)

	tableCellStyle = lipgloss.NewStyle().
			Foreground(colorTable).
			Padding(0, 1)
)

// MessageType represents different types of messages
type MessageType int

const (
	ErrorType MessageType = iota
	SuccessType
	WarningType
	InfoType
)

// Message represents a structured message with optional details
type Message struct {
	Type    MessageType
	Title   string
	Content string
	Details []string
	Code    string // For highlighting code/input
}

// Error creates an error message
func Error(title, content string) Message {
	return Message{
		Type:    ErrorType,
		Title:   title,
		Content: content,
	}
}

// Success creates a success message
func Success(content string) Message {
	return Message{
		Type:    SuccessType,
		Content: content,
	}
}

// Warning creates a warning message
func Warning(content string) Message {
	return Message{
		Type:    WarningType,
		Content: content,
	}
}

// Info creates an info message
func Info(content string) Message {
	return Message{
		Type:    InfoType,
		Content: content,
	}
}

// Print displays a message with appropriate styling based on its type
func (m Message) Print() {
	var style lipgloss.Style
	var icon string
	var titleBg lipgloss.Color

	switch m.Type {
	case ErrorType:
		style = baseStyle.Foreground(colorError).BorderForeground(colorError)
		titleBg = colorError
		icon = "❌"
	case SuccessType:
		style = baseStyle.Foreground(colorSuccess).BorderForeground(colorSuccess)
		titleBg = colorSuccess
		icon = "✅"
	case WarningType:
		style = baseStyle.Foreground(colorWarning).BorderForeground(colorWarning)
		titleBg = colorWarning
		icon = "⚠️"
	case InfoType:
		style = baseStyle.Foreground(colorInfo).BorderForeground(colorInfo)
		titleBg = colorInfo
		icon = "ℹ️"
	}

	fmt.Println()

	// Print title if provided
	if m.Title != "" {
		fmt.Println(titleStyle.Background(titleBg).Render(fmt.Sprintf("%s %s", icon, m.Title)))
	}

	// Build content
	var content string
	if m.Content != "" {
		content += fmt.Sprintf("🔥 %s\n", m.Content)
	}

	// Add code highlighting if provided
	if m.Code != "" {
		content += fmt.Sprintf("\n📝 Input: %s\n", codeStyle.Render(m.Code))
	}

	// Add details if provided
	if len(m.Details) > 0 {
		content += "\n💡 Details:\n"
		for _, detail := range m.Details {
			content += fmt.Sprintf("   • %s\n", detail)
		}
	}

	if content != "" {
		fmt.Println(style.Render(strings.TrimSpace(content)))
	}

	fmt.Println()
}

// PrintAndExit prints the message and exits the program (for errors)
func (m Message) PrintAndExit() {
	m.Print()
	os.Exit(1)
}

// WithDetails adds details to a message
func (m Message) WithDetails(details ...string) Message {
	m.Details = details
	return m
}

// WithCode adds code highlighting to a message
func (m Message) WithCode(code string) Message {
	m.Code = code
	return m
}
