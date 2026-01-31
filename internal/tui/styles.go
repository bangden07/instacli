package tui

import "github.com/charmbracelet/lipgloss"

// ═══════════════════════════════════════════════════════════════════════════
// 🎨 PREMIUM COLOR PALETTE - Cyberpunk Neon Theme
// ═══════════════════════════════════════════════════════════════════════════
var (
	// Primary gradient colors
	Primary   = lipgloss.Color("#00D9FF") // Electric Cyan
	Secondary = lipgloss.Color("#FF00E5") // Neon Magenta
	Accent    = lipgloss.Color("#FFD000") // Golden Yellow
	Tertiary  = lipgloss.Color("#00FF9F") // Matrix Green

	// Semantic colors
	Success = lipgloss.Color("#00FF88") // Bright Green
	Warning = lipgloss.Color("#FFB800") // Amber
	Danger  = lipgloss.Color("#FF3366") // Hot Red
	Info    = lipgloss.Color("#00BFFF") // Sky Blue

	// Neutral palette
	Muted       = lipgloss.Color("#5C6370") // Grey
	Background  = lipgloss.Color("#0D1117") // Deep Black
	Surface     = lipgloss.Color("#161B22") // Elevated surface
	SurfaceHigh = lipgloss.Color("#21262D") // Higher surface
	Foreground  = lipgloss.Color("#E6EDF3") // Bright white
	ForegroundD = lipgloss.Color("#8B949E") // Dimmed text
	BorderColor = lipgloss.Color("#30363D") // Subtle border
	Highlight   = lipgloss.Color("#1F6FEB") // Blue highlight
)

// ═══════════════════════════════════════════════════════════════════════════
// 🎯 PREMIUM ASCII LOGO - Modern minimalist
// ═══════════════════════════════════════════════════════════════════════════
const Logo = `
 ██╗███╗   ██╗███████╗████████╗ █████╗  ██████╗██╗     ██╗
 ██║████╗  ██║██╔════╝╚══██╔══╝██╔══██╗██╔════╝██║     ██║
 ██║██╔██╗ ██║███████╗   ██║   ███████║██║     ██║     ██║
 ██║██║╚██╗██║╚════██║   ██║   ██╔══██║██║     ██║     ██║
 ██║██║ ╚████║███████║   ██║   ██║  ██║╚██████╗███████╗██║
 ╚═╝╚═╝  ╚═══╝╚══════╝   ╚═╝   ╚═╝  ╚═╝ ╚═════╝╚══════╝╚═╝`

const LogoCompact = `█ INSTACLI`

const LogoMini = `⚡ InstaCli`

var LogoStyle = lipgloss.NewStyle().
	Foreground(Primary).
	Bold(true)

// Gradient effect for logo (alternating colors per line)
var LogoGradientStyle = lipgloss.NewStyle().
	Foreground(Primary)

// ═══════════════════════════════════════════════════════════════════════════
// 📦 MAIN CONTAINER STYLES
// ═══════════════════════════════════════════════════════════════════════════
var (
	// App container with full padding
	AppStyle = lipgloss.NewStyle().
			Padding(1, 2)

	// Premium main box with double border
	MainBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(Primary).
			Padding(1, 2)

	// Glass-like header section
	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Background).
			Background(Primary).
			Padding(0, 3).
			MarginBottom(1).
			Align(lipgloss.Center)

	// Title style with glow effect
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Primary).
			MarginBottom(1)

	// Subtitle with elegant italic
	SubtitleStyle = lipgloss.NewStyle().
			Foreground(ForegroundD).
			Italic(true)

	// Modern content box
	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(BorderColor).
			Padding(1, 2)

	// Highlighted selected box
	SelectedBoxStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(Primary).
				Padding(1, 2)

	// Premium category card
	CategoryCardStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(BorderColor).
				Padding(0, 1).
				MarginBottom(0)

	// Selected category with glow border
	SelectedCategoryStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(Primary).
				Background(SurfaceHigh).
				Padding(0, 1).
				MarginBottom(0)
)

// ═══════════════════════════════════════════════════════════════════════════
// 📋 LIST ITEM STYLES
// ═══════════════════════════════════════════════════════════════════════════
var (
	ItemStyle = lipgloss.NewStyle().
			Foreground(Foreground).
			PaddingLeft(2)

	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(Primary).
				Bold(true).
				PaddingLeft(2)

	// Icon styles
	IconStyle = lipgloss.NewStyle().
			MarginRight(1)

	// Checkbox states
	CheckedStyle = lipgloss.NewStyle().
			Foreground(Success).
			Bold(true)

	UncheckedStyle = lipgloss.NewStyle().
			Foreground(Muted)

	// Disabled state
	DisabledStyle = lipgloss.NewStyle().
			Foreground(Muted).
			Strikethrough(true)
)

// ═══════════════════════════════════════════════════════════════════════════
// 🚦 STATUS STYLES
// ═══════════════════════════════════════════════════════════════════════════
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
			Foreground(Info)

	MutedStyle = lipgloss.NewStyle().
			Foreground(Muted)

	// Animated loading indicator style
	LoadingStyle = lipgloss.NewStyle().
			Foreground(Primary).
			Bold(true)
)

// ═══════════════════════════════════════════════════════════════════════════
// ❓ HELP & FOOTER STYLES
// ═══════════════════════════════════════════════════════════════════════════
var (
	HelpStyle = lipgloss.NewStyle().
			Foreground(Muted)

	HelpKeyStyle = lipgloss.NewStyle().
			Foreground(Primary).
			Bold(true)

	HelpDescStyle = lipgloss.NewStyle().
			Foreground(ForegroundD)

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

// ═══════════════════════════════════════════════════════════════════════════
// 🎯 TARGET TOGGLE STYLES
// ═══════════════════════════════════════════════════════════════════════════
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

// ═══════════════════════════════════════════════════════════════════════════
// 🏷️ BADGE STYLES
// ═══════════════════════════════════════════════════════════════════════════
var (
	// Primary badge
	BadgeStyle = lipgloss.NewStyle().
			Foreground(Background).
			Background(Primary).
			Padding(0, 1).
			Bold(true)

	// Success badge
	BadgeSuccessStyle = lipgloss.NewStyle().
				Foreground(Background).
				Background(Success).
				Padding(0, 1)

	// Warning badge
	BadgeWarningStyle = lipgloss.NewStyle().
				Foreground(Background).
				Background(Warning).
				Padding(0, 1)

	// Danger badge
	BadgeDangerStyle = lipgloss.NewStyle().
				Foreground(Foreground).
				Background(Danger).
				Padding(0, 1)

	// Info badge
	BadgeInfoStyle = lipgloss.NewStyle().
			Foreground(Background).
			Background(Info).
			Padding(0, 1)

	// Version badge
	VersionBadge = lipgloss.NewStyle().
			Foreground(Secondary).
			Background(SurfaceHigh).
			Padding(0, 1).
			Bold(true)
)

// ═══════════════════════════════════════════════════════════════════════════
// 🎨 PROGRESS BAR STYLES
// ═══════════════════════════════════════════════════════════════════════════
var (
	ProgressBarStyle = lipgloss.NewStyle().
				Foreground(Primary)

	ProgressTrackStyle = lipgloss.NewStyle().
				Foreground(BorderColor)

	ProgressTextStyle = lipgloss.NewStyle().
				Foreground(Foreground).
				Bold(true)
)

// ═══════════════════════════════════════════════════════════════════════════
// 🖥️ PANEL STYLES
// ═══════════════════════════════════════════════════════════════════════════
var (
	// Left panel (navigation)
	LeftPanelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(BorderColor).
			Padding(1, 1)

	// Right panel (details)
	RightPanelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Primary).
			Padding(1, 1)

	// Divider
	DividerStyle = lipgloss.NewStyle().
			Foreground(BorderColor)
)

// ═══════════════════════════════════════════════════════════════════════════
// ✨ SPECIAL EFFECT STRINGS
// ═══════════════════════════════════════════════════════════════════════════
var (
	// Arrow indicators
	ArrowRight = "›"
	ArrowLeft  = "‹"
	ArrowUp    = "▲"
	ArrowDown  = "▼"

	// Selection indicators
	BulletFull  = "●"
	BulletEmpty = "○"
	CheckMark   = "✓"
	CrossMark   = "✗"

	// Decorative
	Sparkle  = "✨"
	Star     = "★"
	Diamond  = "◆"
	Triangle = "▸"

	// Progress chars
	ProgressFull    = "█"
	ProgressPartial = "▓"
	ProgressEmpty   = "░"
)
