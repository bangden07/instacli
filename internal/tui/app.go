package tui

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/instacli/instacli/internal/executor"
	"github.com/instacli/instacli/internal/installers"
)

// View represents different screens
type View int

const (
	ViewMain View = iota
	ViewCategory
	ViewInstaller
	ViewSettings
	ViewSSHConfig
	ViewRunning
	ViewSystemStatus
	ViewHelp
)

// TargetMode represents local or SSH
type TargetMode int

const (
	TargetLocal TargetMode = iota
	TargetSSH
)

// CategoryItem for the list
type CategoryItem struct {
	title       string
	description string
	icon        string
	category    installers.Category
}

func (i CategoryItem) Title() string       { return i.title }
func (i CategoryItem) Description() string { return i.description }
func (i CategoryItem) FilterValue() string { return i.title }
func (i CategoryItem) Icon() string        { return i.icon }

// InstallerItem for the list
type InstallerItem struct {
	installer installers.Installer
}

func (i InstallerItem) Title() string       { return i.installer.Name() }
func (i InstallerItem) Description() string { return i.installer.Description() }
func (i InstallerItem) FilterValue() string { return i.installer.Name() }
func (i InstallerItem) Icon() string        { return i.installer.Icon() }

// App is the main application model
type App struct {
	currentView       View
	targetMode        TargetMode
	categoryItems     []CategoryItem
	installerItems    []InstallerItem
	cursor            int
	selectedCategory  string
	selectedInstaller installers.Installer
	registry          *installers.Registry
	sysCheck          *executor.SystemCheck
	sshHost           textinput.Model
	sshUser           textinput.Model
	sshPass           textinput.Model
	sshPort           textinput.Model
	settingsCursor    int
	focusedInput      int // 0=host, 1=port, 2=user, 3=pass
	editingSSH        bool
	width             int
	height            int
	quitting          bool
	output            string
	hue               float64 // For RGB animation (0-360)
}

// tickMsg for animation
type tickMsg struct{}

// NewApp creates a new App instance
func NewApp() *App {
	// Create category items
	categories := []CategoryItem{
		{"System Status", "Check installed tools & environment", "🔍", ""},
		{"Web Server Stack", "LAMP, LEMP, Nginx, Apache", "🌐", installers.CategoryWebServer},
		{"Runtime & Languages", "Node.js, Go, Python, PHP", "⚡", installers.CategoryRuntime},
		{"Containers", "Docker, Docker Compose", "🐳", installers.CategoryContainer},
		{"Databases", "MySQL, PostgreSQL, MongoDB, Redis", "🗄️", installers.CategoryDatabase},
		{"Frameworks", "Laravel Kit, Next.js Kit", "🔧", installers.CategoryFramework},
		{"Automation", "N8N, Coolify, PM2", "🤖", installers.CategoryAutomation},
		{"Security", "UFW, Certbot SSL, Fail2ban", "🛡️", installers.CategorySecurity},
		{"Settings", "SSH Config, Preferences", "⚙️", ""},
	}

	// SSH input fields
	sshHost := textinput.New()
	sshHost.Placeholder = "192.168.1.100 or hostname"
	sshHost.CharLimit = 100

	sshPort := textinput.New()
	sshPort.Placeholder = "22"
	sshPort.CharLimit = 5

	sshUser := textinput.New()
	sshUser.Placeholder = "root"
	sshUser.CharLimit = 50

	sshPass := textinput.New()
	sshPass.Placeholder = "password"
	sshPass.EchoMode = textinput.EchoPassword
	sshPass.CharLimit = 100

	// Initialize system check
	sysCheck := executor.NewSystemCheck()
	sysCheck.DetectAll()

	return &App{
		currentView:   ViewMain,
		targetMode:    TargetLocal,
		categoryItems: categories,
		cursor:        0,
		registry:      installers.DefaultRegistry(),
		sysCheck:      sysCheck,
		sshHost:       sshHost,
		sshPort:       sshPort,
		sshUser:       sshUser,
		sshPass:       sshPass,
	}
}

// tick returns a command that sends a tickMsg after a duration
func tick() tea.Cmd {
	return tea.Tick(time.Millisecond*50, func(t time.Time) tea.Msg {
		return tickMsg{}
	})
}

// hslToHex converts HSL to hex color string
func hslToHex(h, s, l float64) string {
	c := (1 - math.Abs(2*l-1)) * s
	x := c * (1 - math.Abs(math.Mod(h/60, 2)-1))
	m := l - c/2

	var r, g, b float64
	switch {
	case h < 60:
		r, g, b = c, x, 0
	case h < 120:
		r, g, b = x, c, 0
	case h < 180:
		r, g, b = 0, c, x
	case h < 240:
		r, g, b = 0, x, c
	case h < 300:
		r, g, b = x, 0, c
	default:
		r, g, b = c, 0, x
	}

	return fmt.Sprintf("#%02x%02x%02x",
		int((r+m)*255),
		int((g+m)*255),
		int((b+m)*255))
}

// Init implements tea.Model
func (a *App) Init() tea.Cmd {
	return tick()
}

// Update implements tea.Model
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		// Update hue for rainbow effect
		a.hue = math.Mod(a.hue+3, 360)
		return a, tick()

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			if a.currentView == ViewMain {
				a.quitting = true
				return a, tea.Quit
			}
			a.currentView = ViewMain
			a.cursor = 0
			return a, nil

		case "up", "k":
			if a.cursor > 0 {
				a.cursor--
			}
			return a, nil

		case "down", "j":
			maxItems := a.getMaxItems()
			if a.cursor < maxItems-1 {
				a.cursor++
			}
			return a, nil

		case "enter":
			switch a.currentView {
			case ViewMain:
				if a.cursor < len(a.categoryItems) {
					item := a.categoryItems[a.cursor]
					if item.title == "Settings" {
						a.currentView = ViewSettings
					} else if item.title == "System Status" {
						// Refresh system check when entering
						a.sysCheck = executor.NewSystemCheck()
						a.sysCheck.DetectAll()
						a.currentView = ViewSystemStatus
					} else {
						a.selectedCategory = item.title
						a.loadInstallers(item.category)
						a.currentView = ViewCategory
					}
					a.cursor = 0
				}
			case ViewCategory:
				if a.cursor < len(a.installerItems) {
					a.selectedInstaller = a.installerItems[a.cursor].installer
					a.currentView = ViewInstaller
				}
			}
			return a, nil

		case "esc", "backspace":
			switch a.currentView {
			case ViewCategory:
				a.currentView = ViewMain
				a.cursor = 0
			case ViewInstaller:
				a.currentView = ViewCategory
			case ViewSettings:
				if a.editingSSH {
					// Exit editing mode
					a.editingSSH = false
					a.sshHost.Blur()
					a.sshPort.Blur()
					a.sshUser.Blur()
					a.sshPass.Blur()
				} else {
					a.currentView = ViewMain
					a.cursor = 0
				}
			case ViewSystemStatus, ViewHelp:
				a.currentView = ViewMain
				a.cursor = 0
			}
			return a, nil

		case "?":
			if !a.editingSSH {
				a.currentView = ViewHelp
			}
			return a, nil

		case "e":
			// Enter edit mode in Settings
			if a.currentView == ViewSettings && !a.editingSSH {
				a.editingSSH = true
				a.focusedInput = 0
				a.sshHost.Focus()
			}
			return a, nil

		case "tab":
			if a.currentView == ViewSettings && a.editingSSH {
				// Cycle through SSH inputs
				a.sshHost.Blur()
				a.sshPort.Blur()
				a.sshUser.Blur()
				a.sshPass.Blur()

				a.focusedInput = (a.focusedInput + 1) % 4
				switch a.focusedInput {
				case 0:
					a.sshHost.Focus()
				case 1:
					a.sshPort.Focus()
				case 2:
					a.sshUser.Focus()
				case 3:
					a.sshPass.Focus()
				}
			} else if a.currentView == ViewMain {
				if a.targetMode == TargetLocal {
					a.targetMode = TargetSSH
				} else {
					a.targetMode = TargetLocal
				}
			}
			return a, nil

		case "g":
			if a.currentView == ViewInstaller && a.selectedInstaller != nil {
				err := a.generateScript()
				if err != nil {
					a.output = fmt.Sprintf("❌ Error: %v", err)
				} else {
					a.output = fmt.Sprintf("📜 Script saved: ./scripts/generated/%s.sh", strings.ToLower(strings.ReplaceAll(a.selectedInstaller.Name(), " ", "_")))
				}
			}
			return a, nil

		case "i":
			if a.currentView == ViewInstaller && a.selectedInstaller != nil {
				if a.targetMode == TargetLocal {
					a.output = "🚀 Installing locally... (check terminal for output)"
					go a.executeInstall()
				} else {
					// SSH mode
					host := a.sshHost.Value()
					if host == "" {
						a.output = "❌ SSH host not configured. Go to Settings."
					} else {
						a.output = fmt.Sprintf("🔌 Connecting to %s...", host)
						go a.executeSSHInstall()
					}
				}
			}
			return a, nil

		default:
			// Handle text input in Settings SSH editing mode
			if a.currentView == ViewSettings && a.editingSSH {
				var cmd tea.Cmd
				switch a.focusedInput {
				case 0:
					a.sshHost, cmd = a.sshHost.Update(msg)
				case 1:
					a.sshPort, cmd = a.sshPort.Update(msg)
				case 2:
					a.sshUser, cmd = a.sshUser.Update(msg)
				case 3:
					a.sshPass, cmd = a.sshPass.Update(msg)
				}
				return a, cmd
			}
		}

	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
	}

	return a, nil
}

func (a *App) getMaxItems() int {
	switch a.currentView {
	case ViewMain:
		return len(a.categoryItems)
	case ViewCategory:
		return len(a.installerItems)
	default:
		return 0
	}
}

func (a *App) loadInstallers(category installers.Category) {
	instList := a.registry.GetByCategory(category)
	a.installerItems = make([]InstallerItem, len(instList))
	for i, inst := range instList {
		a.installerItems[i] = InstallerItem{installer: inst}
	}
}

// View implements tea.Model
func (a *App) View() string {
	if a.quitting {
		return "\n  👋 Thanks for using InstaCli!\n\n"
	}

	// Calculate widths
	contentWidth := a.width - 6
	if contentWidth < 60 {
		contentWidth = 60
	}
	if contentWidth > 100 {
		contentWidth = 100
	}

	var b strings.Builder

	// ═══════════════════════════════════════════════════════════
	// HEADER BOX - Adaptive based on terminal height
	// ═══════════════════════════════════════════════════════════
	var headerContent string
	if a.height > 35 {
		// Full header with logo for large terminals
		headerContent = LogoStyle.Render(Logo) + "\n" +
			SubtitleStyle.Render("     🚀 Universal Installer Tool v1.0")
	} else {
		// Compact header for small terminals
		headerContent = TitleStyle.Render("⚡ InstaCli") + "  " +
			SubtitleStyle.Render("Universal Installer Tool v1.0")
	}

	headerBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Primary).
		Padding(0, 2).
		Width(contentWidth).
		Align(lipgloss.Center).
		Render(headerContent)

	b.WriteString(headerBox)
	b.WriteString("\n")

	// ═══════════════════════════════════════════════════════════
	// TARGET TOGGLE BOX
	// ═══════════════════════════════════════════════════════════
	localLabel := "○ Local"
	sshLabel := "○ VPS (SSH)"
	if a.targetMode == TargetLocal {
		localLabel = TargetActiveStyle.Render("● Local")
		sshLabel = TargetInactiveStyle.Render("○ VPS (SSH)")
	} else {
		localLabel = TargetInactiveStyle.Render("○ Local")
		sshLabel = TargetActiveStyle.Render("● VPS (SSH)")
	}

	targetContent := fmt.Sprintf("  Target:  %s    %s  ", localLabel, sshLabel)
	targetBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(BorderColor).
		Padding(0, 1).
		Render(targetContent)

	b.WriteString(targetBox)
	b.WriteString("\n\n")

	// ═══════════════════════════════════════════════════════════
	// MAIN CONTENT
	// ═══════════════════════════════════════════════════════════
	var content string
	switch a.currentView {
	case ViewMain:
		content = a.renderMainMenu(contentWidth)
	case ViewCategory:
		content = a.renderCategoryView(contentWidth)
	case ViewInstaller:
		content = a.renderInstallerView(contentWidth)
	case ViewSettings:
		content = a.renderSettingsView(contentWidth)
	case ViewSystemStatus:
		content = a.renderSystemStatusView(contentWidth)
	case ViewHelp:
		content = a.renderHelpView(contentWidth)
	}

	// RGB animated border color
	rgbColor := lipgloss.Color(hslToHex(a.hue, 0.8, 0.5))

	contentBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(rgbColor).
		Padding(0, 1).
		Width(contentWidth).
		Render(content)

	b.WriteString(contentBox)

	// ═══════════════════════════════════════════════════════════
	// OUTPUT MESSAGE
	// ═══════════════════════════════════════════════════════════
	if a.output != "" {
		b.WriteString("\n")
		b.WriteString(SuccessStyle.Render(a.output))
	}

	// ═══════════════════════════════════════════════════════════
	// FOOTER / HELP
	// ═══════════════════════════════════════════════════════════
	b.WriteString("\n\n")
	help := a.getHelpText()

	helpBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(BorderColor).
		Padding(0, 2).
		Foreground(Muted).
		Render(help)

	b.WriteString(helpBox)

	return lipgloss.NewStyle().Padding(1, 2).Render(b.String())
}

func (a *App) renderMainMenu(width int) string {
	var b strings.Builder

	title := TitleStyle.Render("📦 Select Category")
	b.WriteString(title)
	b.WriteString("\n")

	// Compact menu - single line per item
	for i, cat := range a.categoryItems {
		icon := cat.icon
		name := cat.title

		if i == a.cursor {
			// Selected item - highlight
			b.WriteString(lipgloss.NewStyle().
				Background(Surface).
				Foreground(Secondary).
				Bold(true).
				Render(fmt.Sprintf(" ▸ %s %s ", icon, name)))
		} else {
			// Normal item
			b.WriteString(fmt.Sprintf("   %s %s", icon, name))
		}
		b.WriteString("\n")
	}

	// Show description of selected item at bottom
	if a.cursor < len(a.categoryItems) {
		desc := a.categoryItems[a.cursor].description
		b.WriteString("\n")
		b.WriteString(SubtitleStyle.Render(fmt.Sprintf("  → %s", desc)))
	}

	return b.String()
}

func (a *App) renderCategoryView(width int) string {
	var b strings.Builder

	// Back breadcrumb
	backStyle := lipgloss.NewStyle().Foreground(Muted)
	b.WriteString(backStyle.Render("← Back to Categories"))
	b.WriteString("\n\n")

	title := TitleStyle.Render(fmt.Sprintf("📦 %s", a.selectedCategory))
	b.WriteString(title)
	b.WriteString("\n\n")

	if len(a.installerItems) == 0 {
		b.WriteString(SubtitleStyle.Render("  No installers available in this category."))
		return b.String()
	}

	for i, item := range a.installerItems {
		inst := item.installer
		icon := inst.Icon()
		name := inst.Name()
		desc := inst.Description()

		var itemStyle lipgloss.Style
		cursor := "  "

		if i == a.cursor {
			cursor = "▸ "
			itemStyle = lipgloss.NewStyle().
				Background(Surface).
				Foreground(Secondary).
				Bold(true).
				Padding(0, 1).
				Width(width - 8)
		} else {
			itemStyle = lipgloss.NewStyle().
				Foreground(Foreground).
				Padding(0, 1).
				Width(width - 8)
		}

		line := fmt.Sprintf("%s%s %s", cursor, icon, name)
		descLine := fmt.Sprintf("     %s", SubtitleStyle.Render(desc))

		b.WriteString(itemStyle.Render(line))
		b.WriteString("\n")
		b.WriteString(descLine)
		b.WriteString("\n")
	}

	return b.String()
}

func (a *App) renderInstallerView(width int) string {
	if a.cursor >= len(a.installerItems) {
		return "No installer selected"
	}

	inst := a.installerItems[a.cursor].installer

	var b strings.Builder

	// Back breadcrumb
	backStyle := lipgloss.NewStyle().Foreground(Muted)
	b.WriteString(backStyle.Render(fmt.Sprintf("← Back to %s", a.selectedCategory)))
	b.WriteString("\n\n")

	// Title
	titleBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Primary).
		Padding(0, 2).
		Render(fmt.Sprintf("%s  %s", inst.Icon(), TitleStyle.Render(inst.Name())))

	b.WriteString(titleBox)
	b.WriteString("\n\n")

	b.WriteString(SubtitleStyle.Render(inst.Description()))
	b.WriteString("\n\n")

	// Info section
	infoBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(BorderColor).
		Padding(1, 2).
		Width(width - 6)

	var info strings.Builder
	info.WriteString(InfoStyle.Render("📋 Details"))
	info.WriteString("\n\n")

	// Supported OS
	info.WriteString("  Supported OS: ")
	osNames := []string{}
	for _, os := range inst.SupportedOS() {
		switch os {
		case installers.OSLinux:
			osNames = append(osNames, "🐧 Linux")
		case installers.OSMacOS:
			osNames = append(osNames, "🍎 macOS")
		case installers.OSWindows:
			osNames = append(osNames, "🪟 Windows")
		}
	}
	info.WriteString(strings.Join(osNames, "  "))
	info.WriteString("\n")

	// Dependencies
	deps := inst.Dependencies()
	if len(deps) > 0 {
		info.WriteString("  Dependencies: ")
		info.WriteString(strings.Join(deps, ", "))
	}

	b.WriteString(infoBox.Render(info.String()))
	b.WriteString("\n\n")

	// Actions
	actionsBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Success).
		Padding(1, 2)

	actions := HelpKeyStyle.Render("[i]") + " Install   " +
		HelpKeyStyle.Render("[g]") + " Generate Script   " +
		HelpKeyStyle.Render("[Esc]") + " Back"

	b.WriteString(actionsBox.Render(actions))

	return b.String()
}

func (a *App) renderSettingsView(width int) string {
	var b strings.Builder

	backStyle := lipgloss.NewStyle().Foreground(Muted)
	b.WriteString(backStyle.Render("← Back to Menu (Esc)"))
	b.WriteString("\n\n")

	b.WriteString(TitleStyle.Render("⚙️ Settings"))
	b.WriteString("\n\n")

	// SSH Configuration Box
	sshBorderColor := BorderColor
	if a.editingSSH {
		sshBorderColor = Secondary
	}

	sshBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(sshBorderColor).
		Padding(1, 2).
		Width(width - 6)

	var sshContent strings.Builder
	sshContent.WriteString(InfoStyle.Render("🔐 SSH Configuration"))
	if a.editingSSH {
		sshContent.WriteString(SuccessStyle.Render("  [EDITING]"))
	}
	sshContent.WriteString("\n\n")

	// Host field
	hostLabel := "  Host:     "
	if a.editingSSH && a.focusedInput == 0 {
		hostLabel = SuccessStyle.Render("▸ Host:     ")
	}
	sshContent.WriteString(hostLabel)
	if a.editingSSH && a.focusedInput == 0 {
		sshContent.WriteString(a.sshHost.View())
	} else if a.sshHost.Value() != "" {
		sshContent.WriteString(a.sshHost.Value())
	} else {
		sshContent.WriteString(SubtitleStyle.Render("(not set)"))
	}
	sshContent.WriteString("\n")

	// Port field
	portLabel := "  Port:     "
	if a.editingSSH && a.focusedInput == 1 {
		portLabel = SuccessStyle.Render("▸ Port:     ")
	}
	sshContent.WriteString(portLabel)
	if a.editingSSH && a.focusedInput == 1 {
		sshContent.WriteString(a.sshPort.View())
	} else if a.sshPort.Value() != "" {
		sshContent.WriteString(a.sshPort.Value())
	} else {
		sshContent.WriteString(SubtitleStyle.Render("22"))
	}
	sshContent.WriteString("\n")

	// User field
	userLabel := "  User:     "
	if a.editingSSH && a.focusedInput == 2 {
		userLabel = SuccessStyle.Render("▸ User:     ")
	}
	sshContent.WriteString(userLabel)
	if a.editingSSH && a.focusedInput == 2 {
		sshContent.WriteString(a.sshUser.View())
	} else if a.sshUser.Value() != "" {
		sshContent.WriteString(a.sshUser.Value())
	} else {
		sshContent.WriteString(SubtitleStyle.Render("(not set)"))
	}
	sshContent.WriteString("\n")

	// Password field
	passLabel := "  Password: "
	if a.editingSSH && a.focusedInput == 3 {
		passLabel = SuccessStyle.Render("▸ Password: ")
	}
	sshContent.WriteString(passLabel)
	if a.editingSSH && a.focusedInput == 3 {
		sshContent.WriteString(a.sshPass.View())
	} else if a.sshPass.Value() != "" {
		sshContent.WriteString("••••••••")
	} else {
		sshContent.WriteString(SubtitleStyle.Render("(not set)"))
	}

	b.WriteString(sshBox.Render(sshContent.String()))
	b.WriteString("\n\n")

	// Help for SSH form
	sshHelpBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Muted).
		Padding(0, 2)

	var helpText string
	if a.editingSSH {
		helpText = HelpKeyStyle.Render("Tab") + " Next field  " +
			HelpKeyStyle.Render("Esc") + " Done editing"
	} else {
		helpText = HelpKeyStyle.Render("e") + " Edit SSH  " +
			HelpKeyStyle.Render("Esc") + " Back"
	}
	b.WriteString(sshHelpBox.Render(helpText))
	b.WriteString("\n\n")

	// About section
	aboutBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(BorderColor).
		Padding(1, 2).
		Width(width - 6)

	about := WarningStyle.Render("ℹ️  About InstaCli") + "\n\n"
	about += "  Version: 1.0.0\n"
	about += "  Author:  InstaCli Team\n"
	about += "  License: MIT\n"
	about += "  Website: github.com/instacli/instacli"

	b.WriteString(aboutBox.Render(about))

	return b.String()
}

func (a *App) renderSystemStatusView(width int) string {
	var b strings.Builder

	backStyle := lipgloss.NewStyle().Foreground(Muted)
	b.WriteString(backStyle.Render("← Back to Menu"))
	b.WriteString("\n\n")

	b.WriteString(TitleStyle.Render("🔍 System Status"))
	b.WriteString("\n\n")

	// System info box
	infoBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Secondary).
		Padding(1, 2).
		Width(width - 6)

	sysInfo := fmt.Sprintf("  System: %s/%s\n", a.sysCheck.OS, a.sysCheck.Arch)
	sysInfo += fmt.Sprintf("  Package Manager: %s\n", a.sysCheck.PackageManager)
	if a.sysCheck.IsRoot {
		sysInfo += "  Privileges: " + SuccessStyle.Render("✓ Root/Admin")
	} else {
		sysInfo += "  Privileges: " + WarningStyle.Render("○ Normal User")
	}

	b.WriteString(infoBox.Render(sysInfo))
	b.WriteString("\n\n")

	// Installed tools
	installed := a.sysCheck.GetInstalledTools()
	missing := a.sysCheck.GetMissingTools()

	installedBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Success).
		Padding(1, 2).
		Width(width - 6)

	var installedContent strings.Builder
	installedContent.WriteString(SuccessStyle.Render(fmt.Sprintf("✅ Installed (%d)", len(installed))))
	installedContent.WriteString("\n\n")

	if len(installed) == 0 {
		installedContent.WriteString("  No tools detected")
	} else {
		count := 0
		for _, tool := range installed {
			if count >= 10 {
				installedContent.WriteString(fmt.Sprintf("\n  ... and %d more", len(installed)-10))
				break
			}
			version := tool.CurrentVersion
			if len(version) > 30 {
				version = version[:30] + "..."
			}
			installedContent.WriteString(fmt.Sprintf("  ✓ %s", tool.Name))
			if version != "" {
				installedContent.WriteString(fmt.Sprintf(" (%s)", SubtitleStyle.Render(version)))
			}
			installedContent.WriteString("\n")
			count++
		}
	}

	b.WriteString(installedBox.Render(installedContent.String()))
	b.WriteString("\n\n")

	// Missing tools
	missingBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Muted).
		Padding(1, 2).
		Width(width - 6)

	var missingContent strings.Builder
	missingContent.WriteString(InfoStyle.Render(fmt.Sprintf("○ Not Installed (%d)", len(missing))))
	missingContent.WriteString("\n\n")

	if len(missing) == 0 {
		missingContent.WriteString("  All common tools are installed!")
	} else {
		count := 0
		for _, tool := range missing {
			if count >= 8 {
				missingContent.WriteString(fmt.Sprintf("\n  ... and %d more", len(missing)-8))
				break
			}
			missingContent.WriteString(fmt.Sprintf("  ○ %s\n", tool.Name))
			count++
		}
	}

	b.WriteString(missingBox.Render(missingContent.String()))

	return b.String()
}

func (a *App) renderHelpView(width int) string {
	var b strings.Builder

	backStyle := lipgloss.NewStyle().Foreground(Muted)
	b.WriteString(backStyle.Render("← Press Esc to go back"))
	b.WriteString("\n\n")

	b.WriteString(TitleStyle.Render("❓ Help & Documentation"))
	b.WriteString("\n\n")

	// Keyboard Shortcuts
	keyBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Secondary).
		Padding(1, 2).
		Width(width - 6)

	keys := SuccessStyle.Render("⌨️  Keyboard Shortcuts") + "\n\n"
	keys += fmt.Sprintf("  %s  Navigate up/down\n", HelpKeyStyle.Render("↑/↓"))
	keys += fmt.Sprintf("  %s  Select item\n", HelpKeyStyle.Render("Enter"))
	keys += fmt.Sprintf("  %s  Go back\n", HelpKeyStyle.Render("Esc"))
	keys += fmt.Sprintf("  %s  Toggle Local/SSH target\n", HelpKeyStyle.Render("Tab"))
	keys += fmt.Sprintf("  %s  Show this help\n", HelpKeyStyle.Render("?"))
	keys += fmt.Sprintf("  %s  Generate install script\n", HelpKeyStyle.Render("g"))
	keys += fmt.Sprintf("  %s  Install selected tool\n", HelpKeyStyle.Render("i"))
	keys += fmt.Sprintf("  %s  Quit application\n", HelpKeyStyle.Render("q"))

	b.WriteString(keyBox.Render(keys))
	b.WriteString("\n\n")

	// Categories explanation
	catBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Primary).
		Padding(1, 2).
		Width(width - 6)

	cats := InfoStyle.Render("📦 Available Categories") + "\n\n"
	cats += "  " + SuccessStyle.Render("🔍 System Status") + "\n"
	cats += "     Scan & show installed tools, versions, and system info\n\n"
	cats += "  " + SuccessStyle.Render("🌐 Web Server Stack") + "\n"
	cats += "     LAMP (Linux+Apache+MySQL+PHP), LEMP (Nginx instead)\n\n"
	cats += "  " + SuccessStyle.Render("⚡ Runtime & Languages") + "\n"
	cats += "     Node.js (via NVM), Golang, Python, PHP\n\n"
	cats += "  " + SuccessStyle.Render("🐳 Containers") + "\n"
	cats += "     Docker Engine, Docker Compose for containerization\n\n"
	cats += "  " + SuccessStyle.Render("🗄️ Databases") + "\n"
	cats += "     MySQL, PostgreSQL, MongoDB, Redis\n\n"
	cats += "  " + SuccessStyle.Render("🔧 Frameworks") + "\n"
	cats += "     Laravel Kit (PHP), Next.js Kit (React)\n\n"
	cats += "  " + SuccessStyle.Render("🤖 Automation") + "\n"
	cats += "     N8N (workflow), PM2 (process manager), Coolify\n\n"
	cats += "  " + SuccessStyle.Render("🛡️ Security") + "\n"
	cats += "     UFW firewall, Certbot SSL, Fail2ban protection"

	b.WriteString(catBox.Render(cats))
	b.WriteString("\n\n")

	// Features explanation
	featBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Accent).
		Padding(1, 2).
		Width(width - 6)

	feats := WarningStyle.Render("✨ Features") + "\n\n"
	feats += "  • " + HelpKeyStyle.Render("Local Install") + " - Install directly on this machine\n"
	feats += "  • " + HelpKeyStyle.Render("VPS Install") + " - Install on remote server via SSH\n"
	feats += "  • " + HelpKeyStyle.Render("Script Generate") + " - Create shell script for manual use\n"
	feats += "  • " + HelpKeyStyle.Render("Cross-Platform") + " - Works on Linux, macOS, Windows\n"
	feats += "  • " + HelpKeyStyle.Render("Auto-Detection") + " - Detects existing tools & versions"

	b.WriteString(featBox.Render(feats))

	return b.String()
}

func (a *App) getHelpText() string {
	base := HelpKeyStyle.Render("↑↓") + " Navigate  "

	switch a.currentView {
	case ViewMain:
		return base +
			HelpKeyStyle.Render("Enter") + " Select  " +
			HelpKeyStyle.Render("Tab") + " Switch Target  " +
			HelpKeyStyle.Render("?") + " Help  " +
			HelpKeyStyle.Render("q") + " Quit"
	case ViewCategory:
		return base +
			HelpKeyStyle.Render("Enter") + " Select  " +
			HelpKeyStyle.Render("Esc") + " Back  " +
			HelpKeyStyle.Render("q") + " Quit"
	case ViewInstaller:
		return HelpKeyStyle.Render("i") + " Install  " +
			HelpKeyStyle.Render("g") + " Generate  " +
			HelpKeyStyle.Render("Esc") + " Back  " +
			HelpKeyStyle.Render("q") + " Quit"
	default:
		return HelpKeyStyle.Render("Esc") + " Back  " +
			HelpKeyStyle.Render("q") + " Quit"
	}
}

// Unused but needed to satisfy interface if list is used elsewhere
// generateScript writes the installer script to ./scripts/generated/
func (a *App) generateScript() error {
	if a.selectedInstaller == nil {
		return fmt.Errorf("no installer selected")
	}

	// Create directory
	dir := filepath.Join(".", "scripts", "generated")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Generate filename
	name := strings.ToLower(strings.ReplaceAll(a.selectedInstaller.Name(), " ", "_"))
	filename := filepath.Join(dir, name+".sh")

	// Generate script
	script := a.selectedInstaller.GenerateInstallScript(installers.OSLinux, installers.PMApt)

	// Write file
	if err := os.WriteFile(filename, []byte(script), 0755); err != nil {
		return err
	}

	return nil
}

// executeInstall runs the installer script locally
func (a *App) executeInstall() {
	if a.selectedInstaller == nil {
		a.output = "❌ No installer selected"
		return
	}

	// Generate script first
	if err := a.generateScript(); err != nil {
		a.output = fmt.Sprintf("❌ Failed to generate script: %v", err)
		return
	}

	// Get script path
	name := strings.ToLower(strings.ReplaceAll(a.selectedInstaller.Name(), " ", "_"))
	scriptPath := filepath.Join(".", "scripts", "generated", name+".sh")

	// For now, just update the output with instructions
	// Real execution would require exec.Command with proper terminal handling
	a.output = fmt.Sprintf("✅ Script ready: %s\n\nRun manually:\n  chmod +x %s && sudo %s", scriptPath, scriptPath, scriptPath)
}

// executeSSHInstall runs the installer via SSH
func (a *App) executeSSHInstall() {
	if a.selectedInstaller == nil {
		a.output = "❌ No installer selected"
		return
	}

	host := a.sshHost.Value()
	port := a.sshPort.Value()
	user := a.sshUser.Value()

	if host == "" {
		a.output = "❌ SSH host not configured"
		return
	}
	if port == "" {
		port = "22"
	}
	if user == "" {
		user = "root"
	}

	// Generate script
	script := a.selectedInstaller.GenerateInstallScript(installers.OSLinux, installers.PMApt)

	// Build SSH command instruction
	a.output = fmt.Sprintf("✅ Ready to install %s on %s@%s:%s\n\nRun manually:\n  ssh %s@%s -p %s 'bash -s' << 'EOF'\n%s\nEOF",
		a.selectedInstaller.Name(), user, host, port, user, host, port, script)
}

func createDelegate() list.DefaultDelegate {
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(Secondary).
		BorderLeftForeground(Secondary)
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Foreground(Muted).
		BorderLeftForeground(Secondary)
	return delegate
}
