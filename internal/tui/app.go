package tui

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/instacli/instacli/internal/executor"
	"github.com/instacli/instacli/internal/installers"
	"github.com/instacli/instacli/internal/version"
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
	ViewCloneSetup
	ViewInstalling
	ViewMCPTargetSelect // For selecting MCP installation targets
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
	scrollOffset      int // For scrolling long lists
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
	repoURL           textinput.Model
	repoTargetDir     textinput.Model
	repoStep          int // 0=URL input, 1=target dir input, 2=cloning
	detectedProject   *executor.ProjectInfo
	// Installation progress
	installLog       []string
	installRunning   bool
	installComplete  bool
	installError     error
	installStartTime time.Time
	// MCP target selection
	mcpTargets     []installers.MCPTarget
	mcpSelectedIdx []bool // Multi-select checkboxes
}

// installOutputMsg for receiving install output
type installOutputMsg struct {
	line string
}

// installDoneMsg when install completes
type installDoneMsg struct {
	err error
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
		{"Monitoring", "Prometheus, Grafana, Netdata", "📊", installers.CategoryMonitoring},
		{"Infrastructure", "Nginx Proxy Manager, Traefik, MinIO", "🔀", installers.CategoryInfrastructure},
		{"VPN", "WireGuard VPN", "🔐", installers.CategoryVPN},
		{"CI/CD", "Jenkins, GitLab Runner, GitHub Actions", "🚀", installers.CategoryCICD},
		{"DNS & Network", "Pi-hole ad blocker", "🕳️", installers.CategoryDNS},
		{"CMS & Blog", "WordPress, Ghost", "📝", installers.CategoryCMS},
		{"Backup", "Restic backup", "💾", installers.CategoryBackup},
		{"AI CLI Tools", "Claude, Gemini, Codex, Aider, OpenCode", "🤖", installers.CategoryAICLI},
		{"MCP Servers", "Context7, Playwright, GitHub, Memory", "📦", installers.CategoryMCP},
		{"Clone & Setup", "Auto-setup from Git repository", "📥", ""},
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

	// Repository URL input
	repoURL := textinput.New()
	repoURL.Placeholder = "https://github.com/user/repo.git"
	repoURL.CharLimit = 200
	repoURL.Width = 50

	// Target directory input
	repoTargetDir := textinput.New()
	repoTargetDir.Placeholder = "/var/www/myproject"
	repoTargetDir.CharLimit = 200
	repoTargetDir.Width = 50

	return &App{
		currentView:   ViewMain,
		targetMode:    TargetLocal,
		categoryItems: categories,
		cursor:        0,
		scrollOffset:  0,
		registry:      installers.DefaultRegistry(),
		sysCheck:      sysCheck,
		sshHost:       sshHost,
		sshPort:       sshPort,
		sshUser:       sshUser,
		sshPass:       sshPass,
		repoURL:       repoURL,
		repoTargetDir: repoTargetDir,
		repoStep:      0,
		width:         100, // Default width, will be updated by WindowSizeMsg
		height:        40,  // Default height, will be updated by WindowSizeMsg
	}
}

// tick returns a command that sends a tickMsg after a duration
func tick() tea.Cmd {
	return tea.Tick(time.Millisecond*200, func(t time.Time) tea.Msg {
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

	case installOutputMsg:
		// Append install output line
		a.installLog = append(a.installLog, msg.line)
		// Keep only last 100 lines to prevent memory issues
		if len(a.installLog) > 100 {
			a.installLog = a.installLog[len(a.installLog)-100:]
		}
		return a, nil

	case installDoneMsg:
		a.installRunning = false
		a.installComplete = true
		a.installError = msg.err
		if msg.err != nil {
			a.installLog = append(a.installLog, fmt.Sprintf("❌ Error: %v", msg.err))
		} else {
			a.installLog = append(a.installLog, "✅ Installation completed successfully!")
		}
		return a, nil

	case tea.KeyMsg:
		// PRIORITY: Handle text input first when editing
		// This ensures backspace and other keys work correctly in text fields
		if a.currentView == ViewSettings && a.editingSSH {
			// Only handle navigation keys specially, forward all others to textinput
			switch msg.String() {
			case "esc":
				// Exit editing mode
				a.editingSSH = false
				a.sshHost.Blur()
				a.sshPort.Blur()
				a.sshUser.Blur()
				a.sshPass.Blur()
				return a, nil
			case "tab", "down":
				a.focusedInput = (a.focusedInput + 1) % 4
				a.sshHost.Blur()
				a.sshPort.Blur()
				a.sshUser.Blur()
				a.sshPass.Blur()
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
				return a, nil
			case "shift+tab", "up":
				a.focusedInput = (a.focusedInput - 1 + 4) % 4
				a.sshHost.Blur()
				a.sshPort.Blur()
				a.sshUser.Blur()
				a.sshPass.Blur()
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
				return a, nil
			case "enter":
				// Save and exit editing
				a.editingSSH = false
				a.sshHost.Blur()
				a.sshPort.Blur()
				a.sshUser.Blur()
				a.sshPass.Blur()
				a.output = "✅ SSH configuration saved"
				return a, nil
			default:
				// Forward ALL other keys (including backspace) to textinput
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

		// Handle Clone & Setup text input - 2-step flow
		if a.currentView == ViewCloneSetup {
			switch msg.String() {
			case "esc":
				if a.repoStep > 0 {
					// Go back to previous step
					a.repoStep--
					if a.repoStep == 0 {
						a.repoTargetDir.Blur()
						a.repoURL.Focus()
					}
				} else {
					// Exit to main menu
					a.currentView = ViewMain
					a.cursor = 0
					a.repoURL.Blur()
					a.repoURL.SetValue("")
					a.repoTargetDir.SetValue("")
				}
				return a, nil
			case "enter":
				if a.repoStep == 0 {
					// Step 0: URL entered, move to target dir input
					if a.repoURL.Value() != "" {
						a.repoStep = 1
						a.repoURL.Blur()
						a.repoTargetDir.Focus()
						// Set default target dir based on repo name
						_, repoName, _ := executor.ParseRepoURL(a.repoURL.Value())
						if repoName != "" {
							a.repoTargetDir.SetValue("/var/www/" + repoName)
						}
					}
				} else if a.repoStep == 1 {
					// Step 1: Target dir entered, start cloning
					if a.repoTargetDir.Value() != "" {
						a.repoStep = 2
						a.repoTargetDir.Blur()
						// Switch to install view and start cloning
						a.currentView = ViewInstalling
						go a.executeRepoSetup(a.repoURL.Value(), a.repoTargetDir.Value())
					}
				}
				return a, nil
			case "tab":
				// Toggle between URL and target dir inputs
				if a.repoStep == 0 && a.repoURL.Value() != "" {
					a.repoStep = 1
					a.repoURL.Blur()
					a.repoTargetDir.Focus()
					_, repoName, _ := executor.ParseRepoURL(a.repoURL.Value())
					if repoName != "" && a.repoTargetDir.Value() == "" {
						a.repoTargetDir.SetValue("/var/www/" + repoName)
					}
				} else if a.repoStep == 1 {
					a.repoStep = 0
					a.repoTargetDir.Blur()
					a.repoURL.Focus()
				}
				return a, nil
			default:
				// Forward to active textinput
				var cmd tea.Cmd
				if a.repoStep == 0 {
					a.repoURL, cmd = a.repoURL.Update(msg)
				} else {
					a.repoTargetDir, cmd = a.repoTargetDir.Update(msg)
				}
				return a, cmd
			}
		}

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
			if a.currentView == ViewMain {
				// Grid navigation: move up by row (numCols items)
				numCols := a.getGridCols()
				if a.cursor >= numCols {
					a.cursor -= numCols
				}
			} else if a.cursor > 0 {
				a.cursor--
			}
			return a, nil

		case "down", "j":
			maxItems := a.getMaxItems()
			if a.currentView == ViewMain {
				// Grid navigation: move down by row (numCols items)
				numCols := a.getGridCols()
				if a.cursor+numCols < maxItems {
					a.cursor += numCols
				}
			} else if a.cursor < maxItems-1 {
				a.cursor++
			}
			return a, nil

		case "left", "h":
			if a.currentView == ViewMain && a.cursor > 0 {
				a.cursor--
			}
			return a, nil

		case "right", "l":
			maxItems := a.getMaxItems()
			if a.currentView == ViewMain && a.cursor < maxItems-1 {
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
					} else if item.title == "Clone & Setup" {
						a.repoURL.Focus()
						a.currentView = ViewCloneSetup
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

		case "esc":
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
			case ViewSystemStatus, ViewHelp, ViewCloneSetup:
				a.currentView = ViewMain
				a.cursor = 0
				a.repoURL.Blur()
			case ViewInstalling:
				// Only allow exit if not running
				if !a.installRunning {
					a.currentView = ViewInstaller
					a.installLog = nil
					a.installComplete = false
					a.installError = nil
				}
			}
			return a, nil

		case "backspace":
			// Don't go back if editing text - let the textinput handle backspace
			if a.currentView == ViewSettings && a.editingSSH {
				// Forward to text input (handled in default case)
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
			if a.currentView == ViewCloneSetup {
				var cmd tea.Cmd
				a.repoURL, cmd = a.repoURL.Update(msg)
				return a, cmd
			}
			// Normal back navigation
			switch a.currentView {
			case ViewCategory:
				a.currentView = ViewMain
				a.cursor = 0
				a.scrollOffset = 0
			case ViewInstaller:
				a.currentView = ViewCategory
				a.cursor = 0
				a.scrollOffset = 0
			case ViewSettings:
				a.currentView = ViewMain
				a.cursor = 0
			case ViewSystemStatus, ViewHelp, ViewCloneSetup:
				a.currentView = ViewMain
				a.cursor = 0
			case ViewMCPTargetSelect:
				a.currentView = ViewInstaller
				a.cursor = 0
				a.scrollOffset = 0
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
				// Check if this is an MCP installer
				if a.selectedInstaller.Category() == installers.CategoryMCP {
					// Detect available targets and show selection
					a.mcpTargets = installers.DetectAllMCPTargets()
					if len(a.mcpTargets) == 0 {
						a.output = "❌ No IDE or AI CLI detected. Install an IDE like Cursor, VS Code, or an AI CLI tool first."
						return a, nil
					}
					// Initialize selection (all selected by default)
					a.mcpSelectedIdx = make([]bool, len(a.mcpTargets))
					for i := range a.mcpSelectedIdx {
						a.mcpSelectedIdx[i] = true
					}
					a.cursor = 0
					a.scrollOffset = 0
					a.currentView = ViewMCPTargetSelect
					return a, nil
				}

				// Normal installation flow
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
			} else if a.currentView == ViewMCPTargetSelect {
				// Install MCP to selected targets
				selectedTargets := []installers.MCPTarget{}
				for i, selected := range a.mcpSelectedIdx {
					if selected {
						selectedTargets = append(selectedTargets, a.mcpTargets[i])
					}
				}
				if len(selectedTargets) == 0 {
					a.output = "❌ Please select at least one target. Use SPACE to toggle."
					return a, nil
				}
				// Execute MCP installation with selected targets
				a.output = fmt.Sprintf("🚀 Installing MCP to %d target(s)...", len(selectedTargets))
				go a.executeMCPInstall(selectedTargets)
				a.currentView = ViewInstalling
			}
			return a, nil

		case " ":
			// Space toggles selection in MCP target select view
			if a.currentView == ViewMCPTargetSelect && a.cursor < len(a.mcpSelectedIdx) {
				a.mcpSelectedIdx[a.cursor] = !a.mcpSelectedIdx[a.cursor]
			}
			return a, nil

		case "a":
			// Select/deselect all in MCP target select view
			if a.currentView == ViewMCPTargetSelect && len(a.mcpSelectedIdx) > 0 {
				allSelected := true
				for _, sel := range a.mcpSelectedIdx {
					if !sel {
						allSelected = false
						break
					}
				}
				for i := range a.mcpSelectedIdx {
					a.mcpSelectedIdx[i] = !allSelected
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
			// Handle text input in Clone & Setup view
			if a.currentView == ViewCloneSetup {
				var cmd tea.Cmd
				a.repoURL, cmd = a.repoURL.Update(msg)
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
	case ViewMCPTargetSelect:
		return len(a.mcpTargets)
	default:
		return 0
	}
}

func (a *App) getGridCols() int {
	// Calculate columns based on content width
	contentWidth := a.width - 10
	if contentWidth < 60 {
		contentWidth = 60
	}
	if contentWidth > 100 {
		contentWidth = 100
	}
	itemWidth := 25
	numCols := contentWidth / itemWidth
	if numCols < 2 {
		numCols = 2
	}
	if numCols > 4 {
		numCols = 4
	}
	return numCols
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

	// Calculate content width - responsive based on terminal size
	contentWidth := a.width - 10
	if contentWidth < 60 {
		contentWidth = 60
	}
	if contentWidth > 100 {
		contentWidth = 100
	}

	var b strings.Builder

	// ═══════════════════════════════════════════════════════════
	// HEADER BOX - Premium cyberpunk style
	// ═══════════════════════════════════════════════════════════
	var headerContent string
	versionBadge := lipgloss.NewStyle().
		Foreground(Secondary).
		Bold(true).
		Render(fmt.Sprintf("v%s", version.Version))

	if a.height > 35 {
		// Full header with logo for large terminals
		headerContent = LogoStyle.Render(Logo) + "\n\n" +
			SubtitleStyle.Render("⚡ Universal Installer Tool") + "  " + versionBadge
	} else {
		// Compact header for small terminals
		headerContent = TitleStyle.Render("█ INSTACLI") + "  " +
			SubtitleStyle.Render("Universal Installer") + " " + versionBadge
	}

	headerBox := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
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
		Width(contentWidth).
		Align(lipgloss.Center).
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
	case ViewCloneSetup:
		content = a.renderCloneSetupView(contentWidth)
	case ViewInstalling:
		content = a.renderInstallingView(contentWidth)
	case ViewMCPTargetSelect:
		content = a.renderMCPTargetSelect(contentWidth)
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
		Width(contentWidth).
		Align(lipgloss.Center).
		Render(help)

	b.WriteString(helpBox)

	// Center the entire UI horizontally and vertically
	finalContent := b.String()

	return lipgloss.Place(
		a.width,
		a.height,
		lipgloss.Center,
		lipgloss.Center,
		finalContent,
	)
}

func (a *App) renderMainMenu(width int) string {
	var b strings.Builder

	title := TitleStyle.Render("📦 Select Category")
	b.WriteString(title)
	b.WriteString("\n\n")

	// Calculate columns based on width
	itemWidth := 25 // Width per item
	numCols := width / itemWidth
	if numCols < 2 {
		numCols = 2
	}
	if numCols > 4 {
		numCols = 4
	}

	// Render items in grid
	items := a.categoryItems
	for row := 0; row < (len(items)+numCols-1)/numCols; row++ {
		var rowItems []string
		for col := 0; col < numCols; col++ {
			idx := row*numCols + col
			if idx >= len(items) {
				// Empty cell for alignment
				rowItems = append(rowItems, strings.Repeat(" ", itemWidth))
				continue
			}

			cat := items[idx]
			icon := cat.icon
			name := cat.title

			// Truncate name if needed
			maxNameLen := itemWidth - 6
			if len(name) > maxNameLen {
				name = name[:maxNameLen-2] + ".."
			}

			var cell string
			if idx == a.cursor {
				// Selected item
				cell = lipgloss.NewStyle().
					Background(Surface).
					Foreground(Secondary).
					Bold(true).
					Width(itemWidth).
					Render(fmt.Sprintf("▸ %s %s", icon, name))
			} else {
				// Normal item
				cell = lipgloss.NewStyle().
					Foreground(Foreground).
					Width(itemWidth).
					Render(fmt.Sprintf("  %s %s", icon, name))
			}
			rowItems = append(rowItems, cell)
		}
		b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, rowItems...))
		b.WriteString("\n")
	}

	// Show description of selected item
	if a.cursor < len(a.categoryItems) {
		desc := a.categoryItems[a.cursor].description
		maxDescLen := width - 4
		if len(desc) > maxDescLen {
			desc = desc[:maxDescLen-3] + "..."
		}
		b.WriteString("\n")
		b.WriteString(SubtitleStyle.Render(fmt.Sprintf("→ %s", desc)))
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

	// Calculate visible items based on terminal height
	// Each item takes 2 lines (name + description)
	reservedLines := 20
	if a.height < 35 {
		reservedLines = 15
	}
	maxVisibleItems := (a.height - reservedLines) / 2 // Each item is 2 lines
	if maxVisibleItems < 3 {
		maxVisibleItems = 3
	}
	if maxVisibleItems > len(a.installerItems) {
		maxVisibleItems = len(a.installerItems)
	}

	// Update scroll offset
	if a.cursor < a.scrollOffset {
		a.scrollOffset = a.cursor
	}
	if a.cursor >= a.scrollOffset+maxVisibleItems {
		a.scrollOffset = a.cursor - maxVisibleItems + 1
	}

	// Scroll indicator top
	if a.scrollOffset > 0 {
		b.WriteString(lipgloss.NewStyle().Foreground(Muted).Render("   ▲ scroll up for more\n"))
	}

	// Render visible items
	endIdx := a.scrollOffset + maxVisibleItems
	if endIdx > len(a.installerItems) {
		endIdx = len(a.installerItems)
	}

	for i := a.scrollOffset; i < endIdx; i++ {
		item := a.installerItems[i]
		inst := item.installer
		icon := inst.Icon()
		name := inst.Name()
		desc := inst.Description()

		// Truncate for narrow terminals
		maxNameLen := width - 12
		if maxNameLen < 20 {
			maxNameLen = 20
		}
		if len(name) > maxNameLen {
			name = name[:maxNameLen-3] + "..."
		}
		maxDescLen := width - 10
		if len(desc) > maxDescLen {
			desc = desc[:maxDescLen-3] + "..."
		}

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

	// Scroll indicator bottom
	if endIdx < len(a.installerItems) {
		b.WriteString(lipgloss.NewStyle().Foreground(Muted).Render("   ▼ scroll down for more\n"))
	}

	// Show scroll position indicator
	if len(a.installerItems) > maxVisibleItems {
		scrollInfo := fmt.Sprintf(" [%d/%d]", a.cursor+1, len(a.installerItems))
		b.WriteString(lipgloss.NewStyle().Foreground(Muted).Render(scrollInfo))
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
	about += fmt.Sprintf("  Version: %s\n", version.Version)
	about += "  Author:  InstaCli Team\n"
	about += "  License: MIT\n"
	about += "  Website: github.com/bangden07/instacli"

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

// executeInstall runs the installer script locally with progress view
func (a *App) executeInstall() {
	if a.selectedInstaller == nil {
		a.output = "❌ No installer selected"
		return
	}

	// Switch to installing view
	a.currentView = ViewInstalling
	a.installLog = nil
	a.installRunning = true
	a.installComplete = false
	a.installError = nil
	a.installStartTime = time.Now()

	// Add initial log entries
	a.installLog = append(a.installLog, "🖥 Running local installation...")
	a.installLog = append(a.installLog, fmt.Sprintf("📦 Installing: %s", a.selectedInstaller.Name()))
	a.installLog = append(a.installLog, "")

	// Generate script
	script := a.selectedInstaller.GenerateInstallScript(installers.OSLinux, installers.PMApt)

	// Run in goroutine
	go func() {
		// Create temp script file
		tmpFile, err := os.CreateTemp("", "instacli-*.sh")
		if err != nil {
			a.installLog = append(a.installLog, fmt.Sprintf("❌ Failed to create temp file: %v", err))
			a.installRunning = false
			a.installComplete = true
			a.installError = err
			return
		}
		scriptPath := tmpFile.Name()
		defer os.Remove(scriptPath)

		// Write script
		if _, err := tmpFile.WriteString(script); err != nil {
			a.installLog = append(a.installLog, fmt.Sprintf("❌ Failed to write script: %v", err))
			a.installRunning = false
			a.installComplete = true
			a.installError = err
			return
		}
		tmpFile.Close()

		// Make executable
		os.Chmod(scriptPath, 0755)

		a.installLog = append(a.installLog, "▶ Running installation script...")
		a.installLog = append(a.installLog, "")

		// Run script with bash
		cmd := exec.Command("bash", scriptPath)
		cmd.Env = os.Environ()

		// Capture output
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			a.installLog = append(a.installLog, fmt.Sprintf("❌ Failed to capture output: %v", err))
			a.installRunning = false
			a.installComplete = true
			a.installError = err
			return
		}

		stderr, err := cmd.StderrPipe()
		if err != nil {
			a.installLog = append(a.installLog, fmt.Sprintf("❌ Failed to capture stderr: %v", err))
			a.installRunning = false
			a.installComplete = true
			a.installError = err
			return
		}

		// Start command
		if err := cmd.Start(); err != nil {
			a.installLog = append(a.installLog, fmt.Sprintf("❌ Failed to start script: %v", err))
			a.installRunning = false
			a.installComplete = true
			a.installError = err
			return
		}

		// Read stdout in goroutine
		go func() {
			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				line := scanner.Text()
				a.installLog = append(a.installLog, "  "+line)
				// Keep only last 100 lines
				if len(a.installLog) > 100 {
					a.installLog = a.installLog[len(a.installLog)-100:]
				}
			}
		}()

		// Read stderr
		go func() {
			scanner := bufio.NewScanner(stderr)
			for scanner.Scan() {
				line := scanner.Text()
				a.installLog = append(a.installLog, "  "+line)
			}
		}()

		// Wait for completion
		err = cmd.Wait()

		if err != nil {
			a.installLog = append(a.installLog, "")
			a.installLog = append(a.installLog, fmt.Sprintf("❌ Installation failed: %v", err))
			a.installError = err
		} else {
			a.installLog = append(a.installLog, "")
			a.installLog = append(a.installLog, fmt.Sprintf("✅ Successfully installed %s!", a.selectedInstaller.Name()))
		}

		a.installRunning = false
		a.installComplete = true
	}()
}

// executeSSHInstall runs the installer via SSH with progress view
func (a *App) executeSSHInstall() {
	if a.selectedInstaller == nil {
		a.output = "❌ No installer selected"
		return
	}

	host := a.sshHost.Value()
	portStr := a.sshPort.Value()
	user := a.sshUser.Value()
	pass := a.sshPass.Value()

	if host == "" {
		a.output = "❌ SSH host not configured. Go to Settings."
		return
	}
	if portStr == "" {
		portStr = "22"
	}
	if user == "" {
		user = "root"
	}

	// Parse port
	port := 22
	if portStr != "" {
		fmt.Sscanf(portStr, "%d", &port)
	}

	// Switch to installing view
	a.currentView = ViewInstalling
	a.installLog = nil
	a.installRunning = true
	a.installComplete = false
	a.installError = nil
	a.installStartTime = time.Now()

	// Add initial log entries
	a.installLog = append(a.installLog, fmt.Sprintf("🔌 Connecting to %s@%s:%d...", user, host, port))
	a.installLog = append(a.installLog, fmt.Sprintf("📦 Installing: %s", a.selectedInstaller.Name()))
	a.installLog = append(a.installLog, "")

	// Run install in background
	installerToRun := a.selectedInstaller
	go func() {
		// Create SSH executor
		sshConfig := executor.SSHConfig{
			Host:     host,
			Port:     port,
			User:     user,
			Password: pass,
		}
		sshExec := executor.NewSSHExecutor(sshConfig)

		// Connect
		err := sshExec.Connect()
		if err != nil {
			a.installLog = append(a.installLog, fmt.Sprintf("❌ SSH connection failed: %v", err))
			a.installRunning = false
			a.installComplete = true
			a.installError = err
			return
		}
		defer sshExec.Close()

		a.installLog = append(a.installLog, fmt.Sprintf("✅ Connected to %s!", host))
		a.installLog = append(a.installLog, "")
		a.installLog = append(a.installLog, "▶ Running installation script...")
		a.installLog = append(a.installLog, "")

		// Generate and execute install script
		script := installerToRun.GenerateInstallScript(sshExec.GetOS(), sshExec.GetPackageManager())

		// Run the script with streaming output
		output, err := sshExec.Run(script)

		// Parse output lines
		lines := strings.Split(output, "\n")
		for _, line := range lines {
			if line != "" {
				a.installLog = append(a.installLog, "  "+line)
			}
		}

		if err != nil {
			a.installLog = append(a.installLog, "")
			a.installLog = append(a.installLog, fmt.Sprintf("❌ Installation failed: %v", err))
			a.installError = err
		} else {
			a.installLog = append(a.installLog, "")
			a.installLog = append(a.installLog, fmt.Sprintf("✅ Successfully installed %s!", installerToRun.Name()))
		}

		a.installRunning = false
		a.installComplete = true
	}()
}

// renderInstallingView renders the installation progress view
func (a *App) renderInstallingView(width int) string {
	var b strings.Builder

	// Title with spinner animation
	spinnerFrames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	frameIndex := int(time.Now().UnixNano()/100000000) % len(spinnerFrames)

	var titleText string
	if a.installRunning {
		spinner := spinnerFrames[frameIndex]
		titleText = fmt.Sprintf("%s Installing...", spinner)
	} else if a.installComplete {
		if a.installError != nil {
			titleText = "❌ Installation Failed"
		} else {
			titleText = "✅ Installation Complete"
		}
	}

	b.WriteString(TitleStyle.Render(titleText))
	b.WriteString("\n")

	// Show elapsed time
	elapsed := time.Since(a.installStartTime).Round(time.Second)
	b.WriteString(MutedStyle.Render(fmt.Sprintf("⏱ Elapsed: %s", elapsed)))
	b.WriteString("\n\n")

	// Progress log box
	logHeight := 15 // Show last 15 lines
	startIdx := 0
	if len(a.installLog) > logHeight {
		startIdx = len(a.installLog) - logHeight
	}

	b.WriteString(SubtitleStyle.Render("📋 Installation Log"))
	b.WriteString("\n")

	// Log content
	logContent := strings.Builder{}
	for i := startIdx; i < len(a.installLog); i++ {
		line := a.installLog[i]
		// Color code lines
		if strings.HasPrefix(line, "✅") {
			logContent.WriteString(SuccessStyle.Render(line))
		} else if strings.HasPrefix(line, "❌") {
			logContent.WriteString(ErrorStyle.Render(line))
		} else if strings.HasPrefix(line, "🔌") || strings.HasPrefix(line, "📦") || strings.HasPrefix(line, "▶") {
			logContent.WriteString(InfoStyle.Render(line))
		} else {
			logContent.WriteString(MutedStyle.Render(line))
		}
		logContent.WriteString("\n")
	}

	logBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(BorderColor).
		Padding(0, 1).
		Width(width - 4).
		Height(logHeight + 2).
		Render(logContent.String())

	b.WriteString(logBox)
	b.WriteString("\n\n")

	// Help text
	if a.installRunning {
		b.WriteString(MutedStyle.Render("Installation in progress... Please wait."))
	} else {
		b.WriteString(HelpStyle.Render("[Esc] Back"))
	}

	return b.String()
}

// renderCloneSetupView renders the Clone & Setup view
func (a *App) renderCloneSetupView(width int) string {
	var b strings.Builder

	b.WriteString(TitleStyle.Render("📥 Clone & Setup Repository"))
	b.WriteString("\n\n")

	// Step indicator
	var stepText string
	if a.repoStep == 0 {
		stepText = "[1] Repository URL → [2] Target Directory"
	} else {
		stepText = "[1] Repository URL → [2] Target Directory ✓"
	}
	b.WriteString(HelpStyle.Render(stepText))
	b.WriteString("\n\n")

	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Primary).
		Padding(0, 1).
		Width(width - 10)

	inputStyleInactive := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(0, 1).
		Width(width - 10)

	// Step 1: URL Input
	if a.repoStep == 0 {
		b.WriteString(SubtitleStyle.Render("Step 1: Repository URL"))
		b.WriteString("\n")
		b.WriteString(inputStyle.Render(a.repoURL.View()))
	} else {
		b.WriteString(SuccessStyle.Render("✅ Step 1: " + a.repoURL.Value()))
		b.WriteString("\n")
	}
	b.WriteString("\n\n")

	// Step 2: Target Directory
	if a.repoStep >= 1 {
		b.WriteString(SubtitleStyle.Render("Step 2: Target Directory"))
		b.WriteString("\n")
		b.WriteString(inputStyle.Render(a.repoTargetDir.View()))
		b.WriteString("\n\n")
	} else {
		b.WriteString(HelpStyle.Render("Step 2: Target Directory"))
		b.WriteString("\n")
		b.WriteString(inputStyleInactive.Render("(Enter URL first)"))
		b.WriteString("\n\n")
	}

	// Supported project types
	b.WriteString(SubtitleStyle.Render("✨ Auto-detected project types:"))
	b.WriteString("\n")

	projectTypes := []struct {
		icon string
		name string
		pm   string
	}{
		{"📦", "Node.js", "npm, pnpm, yarn, bun"},
		{"🐍", "Python", "pip, pipenv, venv"},
		{"🐹", "Go", "go mod"},
		{"🐘", "PHP/Laravel", "composer"},
		{"💎", "Ruby", "bundler"},
		{"🦀", "Rust", "cargo"},
		{"🐳", "Docker", "docker compose"},
	}

	for _, pt := range projectTypes {
		b.WriteString(fmt.Sprintf("  %s %s (%s)\n", pt.icon, pt.name, pt.pm))
	}

	b.WriteString("\n")

	// Help text based on step
	if a.repoStep == 0 {
		b.WriteString(HelpStyle.Render("Enter URL and press Enter to continue, Tab to jump to Step 2, Esc to go back"))
	} else {
		b.WriteString(HelpStyle.Render("Press Enter to start cloning, Tab to go back, Esc to cancel"))
	}

	return b.String()
}

// executeRepoSetup clones a repository and runs project setup
func (a *App) executeRepoSetup(repoURL, targetDir string) {
	// Initialize install progress
	a.installLog = nil
	a.installRunning = true
	a.installComplete = false
	a.installError = nil
	a.installStartTime = time.Now()

	a.installLog = append(a.installLog, "📥 Clone & Setup Repository")
	a.installLog = append(a.installLog, "")
	a.installLog = append(a.installLog, fmt.Sprintf("🔗 Repository: %s", repoURL))
	a.installLog = append(a.installLog, fmt.Sprintf("📁 Target: %s", targetDir))
	a.installLog = append(a.installLog, "")

	// Generate clone script
	cloneScript := executor.GenerateCloneScript(repoURL, targetDir)

	// Run the clone script
	a.installLog = append(a.installLog, "⏳ Cloning repository...")

	var cmd *exec.Cmd
	var scriptPath string

	if a.targetMode == TargetSSH && a.sshHost.Value() != "" {
		// Remote execution via SSH
		sshCmd := fmt.Sprintf("ssh -o StrictHostKeyChecking=no -p %s %s@%s 'bash -s'",
			a.sshPort.Value(), a.sshUser.Value(), a.sshHost.Value())
		cmd = exec.Command("bash", "-c", fmt.Sprintf("echo '%s' | %s", cloneScript, sshCmd))
	} else {
		// Local execution
		scriptPath = filepath.Join(os.TempDir(), "repo_clone.sh")
		if err := os.WriteFile(scriptPath, []byte(cloneScript), 0755); err != nil {
			a.installLog = append(a.installLog, fmt.Sprintf("❌ Failed to create script: %v", err))
			a.installRunning = false
			a.installComplete = true
			a.installError = err
			return
		}
		cmd = exec.Command("bash", scriptPath)
	}

	// Capture output
	output, err := cmd.CombinedOutput()
	if err != nil {
		a.installLog = append(a.installLog, "")
		a.installLog = append(a.installLog, "❌ Clone failed:")
		a.installLog = append(a.installLog, string(output))
		a.installLog = append(a.installLog, fmt.Sprintf("Error: %v", err))
		a.installRunning = false
		a.installComplete = true
		a.installError = err
		return
	}

	// Add output to log
	for _, line := range strings.Split(string(output), "\n") {
		if line != "" {
			a.installLog = append(a.installLog, line)
		}
	}

	// Detect project type and run setup
	a.installLog = append(a.installLog, "")
	a.installLog = append(a.installLog, "🔍 Detecting project type...")

	repoSetup := executor.NewRepoSetup(repoURL, targetDir)
	projectInfo, err := repoSetup.DetectProject()
	if err == nil && projectInfo.Type != executor.ProjectUnknown {
		a.detectedProject = projectInfo
		a.installLog = append(a.installLog, fmt.Sprintf("✅ Detected: %s (%s)", projectInfo.Framework, projectInfo.Type))

		if projectInfo.PackageManager != "" {
			a.installLog = append(a.installLog, fmt.Sprintf("📦 Package Manager: %s", projectInfo.PackageManager))
		}

		// Generate and run setup script
		a.installLog = append(a.installLog, "")
		a.installLog = append(a.installLog, "⏳ Running project setup...")

		setupScript := repoSetup.GenerateSetupScript()
		if a.targetMode == TargetSSH && a.sshHost.Value() != "" {
			sshCmd := fmt.Sprintf("ssh -o StrictHostKeyChecking=no -p %s %s@%s 'bash -s'",
				a.sshPort.Value(), a.sshUser.Value(), a.sshHost.Value())
			cmd = exec.Command("bash", "-c", fmt.Sprintf("echo '%s' | %s", setupScript, sshCmd))
		} else {
			scriptPath = filepath.Join(os.TempDir(), "repo_setup.sh")
			os.WriteFile(scriptPath, []byte(setupScript), 0755)
			cmd = exec.Command("bash", scriptPath)
		}

		output, _ = cmd.CombinedOutput()
		for _, line := range strings.Split(string(output), "\n") {
			if line != "" {
				a.installLog = append(a.installLog, line)
			}
		}
	} else {
		a.installLog = append(a.installLog, "ℹ️ Unknown project type - skipping auto-setup")
	}

	// Success message
	a.installLog = append(a.installLog, "")
	a.installLog = append(a.installLog, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	a.installLog = append(a.installLog, "✅ Repository cloned successfully!")
	a.installLog = append(a.installLog, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	a.installLog = append(a.installLog, "")
	a.installLog = append(a.installLog, fmt.Sprintf("📁 Location: %s", targetDir))
	a.installLog = append(a.installLog, "")
	a.installLog = append(a.installLog, "To access your project:")
	a.installLog = append(a.installLog, fmt.Sprintf("   cd %s", targetDir))
	a.installLog = append(a.installLog, "")

	a.installRunning = false
	a.installComplete = true
}

// executeMCPInstall installs MCP server to selected targets
func (a *App) executeMCPInstall(targets []installers.MCPTarget) {
	if a.selectedInstaller == nil {
		a.output = "❌ No installer selected"
		return
	}

	// Switch to installing view
	a.installLog = nil
	a.installRunning = true
	a.installComplete = false
	a.installError = nil
	a.installStartTime = time.Now()

	a.installLog = append(a.installLog, fmt.Sprintf("📦 Installing: %s", a.selectedInstaller.Name()))
	a.installLog = append(a.installLog, "")
	a.installLog = append(a.installLog, "🎯 Configuring for targets:")
	for _, target := range targets {
		a.installLog = append(a.installLog, fmt.Sprintf("   %s %s", target.Icon, target.Name))
	}
	a.installLog = append(a.installLog, "")

	installerToRun := a.selectedInstaller

	go func() {
		// Detect OS
		var osType installers.OS
		switch runtime.GOOS {
		case "darwin":
			osType = installers.OSMacOS
		case "windows":
			osType = installers.OSWindows
		default:
			osType = installers.OSLinux
		}

		// Generate and run the MCP install script
		script := installerToRun.GenerateInstallScript(osType, installers.PMNone)

		a.installLog = append(a.installLog, "📜 Running install script...")
		a.installLog = append(a.installLog, "")

		// Write script to temp file and execute
		var scriptPath string
		var cmd *exec.Cmd

		if runtime.GOOS == "windows" {
			// Windows PowerShell
			scriptPath = filepath.Join(os.TempDir(), "mcp_install.ps1")
			if err := os.WriteFile(scriptPath, []byte(script), 0755); err != nil {
				a.installLog = append(a.installLog, fmt.Sprintf("❌ Failed to create script: %v", err))
				a.installRunning = false
				a.installComplete = true
				a.installError = err
				return
			}
			cmd = exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-File", scriptPath)
		} else {
			// Linux/macOS bash
			scriptPath = filepath.Join(os.TempDir(), "mcp_install.sh")
			if err := os.WriteFile(scriptPath, []byte(script), 0755); err != nil {
				a.installLog = append(a.installLog, fmt.Sprintf("❌ Failed to create script: %v", err))
				a.installRunning = false
				a.installComplete = true
				a.installError = err
				return
			}
			cmd = exec.Command("bash", scriptPath)
		}

		// Get stdout pipe
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			a.installLog = append(a.installLog, fmt.Sprintf("❌ Failed to capture stdout: %v", err))
			a.installRunning = false
			a.installComplete = true
			a.installError = err
			return
		}

		stderr, err := cmd.StderrPipe()
		if err != nil {
			a.installLog = append(a.installLog, fmt.Sprintf("❌ Failed to capture stderr: %v", err))
			a.installRunning = false
			a.installComplete = true
			a.installError = err
			return
		}

		// Start command
		if err := cmd.Start(); err != nil {
			a.installLog = append(a.installLog, fmt.Sprintf("❌ Failed to start script: %v", err))
			a.installRunning = false
			a.installComplete = true
			a.installError = err
			return
		}

		// Read stdout in goroutine
		go func() {
			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				line := scanner.Text()
				a.installLog = append(a.installLog, "  "+line)
				if len(a.installLog) > 100 {
					a.installLog = a.installLog[len(a.installLog)-100:]
				}
			}
		}()

		// Read stderr
		go func() {
			scanner := bufio.NewScanner(stderr)
			for scanner.Scan() {
				line := scanner.Text()
				a.installLog = append(a.installLog, "  "+line)
			}
		}()

		// Wait for completion
		err = cmd.Wait()

		// Clean up
		os.Remove(scriptPath)

		if err != nil {
			a.installLog = append(a.installLog, "")
			a.installLog = append(a.installLog, fmt.Sprintf("❌ Installation failed: %v", err))
			a.installError = err
		} else {
			a.installLog = append(a.installLog, "")
			a.installLog = append(a.installLog, fmt.Sprintf("✅ Successfully installed %s!", a.selectedInstaller.Name()))
			a.installLog = append(a.installLog, "")
			a.installLog = append(a.installLog, "🎉 Restart your IDE/CLI to use the MCP server.")
		}

		a.installRunning = false
		a.installComplete = true
	}()
}

// renderMCPTargetSelect renders the MCP target selection view
func (a *App) renderMCPTargetSelect(_ int) string {
	var b strings.Builder

	// Back breadcrumb
	backStyle := lipgloss.NewStyle().Foreground(Muted)
	b.WriteString(backStyle.Render("← Back (backspace)"))
	b.WriteString("\n\n")

	title := TitleStyle.Render(fmt.Sprintf("📦 Install %s", a.selectedInstaller.Name()))
	b.WriteString(title)
	b.WriteString("\n\n")

	b.WriteString(SubtitleStyle.Render("Select where to install MCP server:"))
	b.WriteString("\n\n")

	if len(a.mcpTargets) == 0 {
		b.WriteString("❌ No IDE or AI CLI detected.\n")
		b.WriteString("Install an IDE (Cursor, VS Code, Windsurf) or\n")
		b.WriteString("an AI CLI tool (Claude, Gemini, OpenCode) first.")
		return b.String()
	}

	// Group targets by type
	ideTargets := []int{}
	cliTargets := []int{}
	for i, t := range a.mcpTargets {
		if t.Type == "ide" {
			ideTargets = append(ideTargets, i)
		} else {
			cliTargets = append(cliTargets, i)
		}
	}

	renderTarget := func(idx int) string {
		target := a.mcpTargets[idx]
		checkbox := "[ ]"
		if a.mcpSelectedIdx[idx] {
			checkbox = "[✓]"
		}

		cursor := "  "
		style := lipgloss.NewStyle().Foreground(Foreground)
		if idx == a.cursor {
			cursor = "▸ "
			style = lipgloss.NewStyle().Background(Surface).Foreground(Secondary).Bold(true)
		}

		line := fmt.Sprintf("%s%s %s %s", cursor, checkbox, target.Icon, target.Name)
		return style.Render(line)
	}

	// Render IDEs section
	if len(ideTargets) > 0 {
		b.WriteString(lipgloss.NewStyle().Foreground(Primary).Bold(true).Render("💻 IDEs"))
		b.WriteString("\n")
		for _, idx := range ideTargets {
			b.WriteString(renderTarget(idx))
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	// Render CLI tools section
	if len(cliTargets) > 0 {
		b.WriteString(lipgloss.NewStyle().Foreground(Primary).Bold(true).Render("🤖 AI CLI Tools"))
		b.WriteString("\n")
		for _, idx := range cliTargets {
			b.WriteString(renderTarget(idx))
			b.WriteString("\n")
		}
	}

	// Selected count
	selectedCount := 0
	for _, sel := range a.mcpSelectedIdx {
		if sel {
			selectedCount++
		}
	}
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("Selected: %d/%d targets", selectedCount, len(a.mcpTargets)))
	b.WriteString("\n\n")

	// Controls
	controls := lipgloss.NewStyle().Foreground(Muted).Render(
		"[SPACE] Toggle • [a] Select/Deselect All • [i] Install • [backspace] Back")
	b.WriteString(controls)

	return b.String()
}
