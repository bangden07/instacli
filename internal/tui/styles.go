package tui

import "github.com/charmbracelet/lipgloss"

// Color palette - Modern dark theme
var (
	Primary     = lipgloss.Color("#8B5CF6") // Violet
	Secondary   = lipgloss.Color("#06B6D4") // Cyan
	Accent      = lipgloss.Color("#F472B6") // Pink
	Success     = lipgloss.Color("#22C55E") // Green
	Warning     = lipgloss.Color("#FBBF24") // Amber
	Danger      = lipgloss.Color("#EF4444") // Red
	Muted       = lipgloss.Color("#64748B") // Slate
	Background  = lipgloss.Color("#0F172A") // Dark slate
	Surface     = lipgloss.Color("#1E293B") // Lighter slate
	Foreground  = lipgloss.Color("#F8FAFC") // Almost white
	BorderColor = lipgloss.Color("#334155") // Slate border
	Highlight   = lipgloss.Color("#7C3AED") // Purple highlight
)

// Logo ASCII art - Smaller and cleaner
const Logo = `
  ╦╔╗╔╔═╗╔╦╗╔═╗╔═╗╦  ╦
  ║║║║╚═╗ ║ ╠═╣║  ║  ║
  ╩╝╚╝╚═╝ ╩ ╩ ╩╚═╝╩═╝╩`

var LogoStyle = lipgloss.NewStyle().
	Foreground(Primary).
	Bold(true)

// Main container styles
var (
	// App container
	AppStyle = lipgloss.NewStyle().
			Padding(1, 2)

	// Main box with rounded border
	MainBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Primary).
			Padding(1, 2)

	// Header section
	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Foreground).
			Background(Primary).
			Padding(0, 2).
			MarginBottom(1).
			Align(lipgloss.Center)

	// Title style
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Primary).
			MarginBottom(1)

	// Subtitle
	SubtitleStyle = lipgloss.NewStyle().
			Foreground(Muted).
			Italic(true)

	// Content box
	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(BorderColor).
			Padding(1, 2)

	// Selected/highlighted box
	SelectedBoxStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(Secondary).
				Padding(1, 2)

	// Category card style
	CategoryCardStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(BorderColor).
				Padding(0, 1).
				MarginBottom(0)

	// Selected category card
	SelectedCategoryStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(Primary).
				Background(Surface).
				Padding(0, 1).
				MarginBottom(0)
)

// List item styles
var (
	ItemStyle = lipgloss.NewStyle().
			Foreground(Foreground).
			PaddingLeft(2)

	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(Secondary).
				Bold(true).
				PaddingLeft(2)

	// Icon styles
	IconStyle = lipgloss.NewStyle().
			MarginRight(1)

	// Checked/unchecked
	CheckedStyle = lipgloss.NewStyle().
			Foreground(Success).
			Bold(true)

	UncheckedStyle = lipgloss.NewStyle().
			Foreground(Muted)
)

// Status styles
var (
	SuccessStyle = lipgloss.NewStyle().
			Foreground(Success).
			Bold(true)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(Danger).
			Bold(true)

	WarningStyle = lipgloss.NewStyle().
			Foreground(Warning)

	InfoStyle = lipgloss.NewStyle().
			Foreground(Secondary)

	MutedStyle = lipgloss.NewStyle().
			Foreground(Muted)
)

// Help/footer styles
var (
	HelpStyle = lipgloss.NewStyle().
			Foreground(Muted)

	HelpKeyStyle = lipgloss.NewStyle().
			Foreground(Primary).
			Bold(true)

	FooterStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderTop(true).
			BorderBottom(false).
			BorderLeft(false).
			BorderRight(false).
			BorderForeground(BorderColor).
			Padding(1, 0).
			MarginTop(1)
)

// Target toggle styles
var (
	TargetActiveStyle = lipgloss.NewStyle().
				Foreground(Success).
				Bold(true)

	TargetInactiveStyle = lipgloss.NewStyle().
				Foreground(Muted)

	TargetBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(BorderColor).
			Padding(0, 2)
)

// Badge styles
var (
	BadgeStyle = lipgloss.NewStyle().
			Foreground(Foreground).
			Background(Primary).
			Padding(0, 1).
			Bold(true)

	BadgeSuccessStyle = lipgloss.NewStyle().
				Foreground(Foreground).
				Background(Success).
				Padding(0, 1)
)
