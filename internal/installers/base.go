package installers

// OS represents the operating system
type OS string

const (
	OSLinux   OS = "linux"
	OSMacOS   OS = "darwin"
	OSWindows OS = "windows"
)

// PackageManager represents the system package manager
type PackageManager string

const (
	PMNone   PackageManager = "none"
	PMApt    PackageManager = "apt"
	PMYum    PackageManager = "yum"
	PMDnf    PackageManager = "dnf"
	PMBrew   PackageManager = "brew"
	PMChoco  PackageManager = "choco"
	PMPacman PackageManager = "pacman"
)

// Category represents installer categories
type Category string

const (
	CategoryWebServer  Category = "Web Server Stack"
	CategoryRuntime    Category = "Runtime & Languages"
	CategoryContainer  Category = "Containers"
	CategoryDatabase   Category = "Databases"
	CategoryFramework  Category = "Frameworks"
	CategoryAutomation Category = "Automation"
	CategorySecurity   Category = "Security"
)

// Installer is the interface all installers must implement
type Installer interface {
	// Basic info
	Name() string
	Description() string
	Category() Category
	Icon() string

	// Compatibility
	SupportedOS() []OS
	RequiredPackageManagers() []PackageManager

	// Dependencies
	Dependencies() []string

	// Actions
	Install(executor Executor) error
	Uninstall(executor Executor) error
	IsInstalled(executor Executor) (bool, error)

	// Script generation
	GenerateInstallScript(os OS, pm PackageManager) string
	GenerateUninstallScript(os OS, pm PackageManager) string
}

// Executor is the interface for running commands
type Executor interface {
	// Run executes a command and returns output
	Run(cmd string) (string, error)

	// RunWithProgress runs a command with progress callback
	RunWithProgress(cmd string, onOutput func(line string)) error

	// GetOS returns the target OS
	GetOS() OS

	// GetPackageManager returns the detected package manager
	GetPackageManager() PackageManager

	// IsRoot checks if running with elevated privileges
	IsRoot() bool
}

// BaseInstaller provides common functionality
type BaseInstaller struct {
	name        string
	description string
	category    Category
	icon        string
	supportedOS []OS
}

func NewBaseInstaller(name, description string, category Category, icon string, os []OS) BaseInstaller {
	return BaseInstaller{
		name:        name,
		description: description,
		category:    category,
		icon:        icon,
		supportedOS: os,
	}
}

func (b BaseInstaller) Name() string        { return b.name }
func (b BaseInstaller) Description() string { return b.description }
func (b BaseInstaller) Category() Category  { return b.category }
func (b BaseInstaller) Icon() string        { return b.icon }
func (b BaseInstaller) SupportedOS() []OS   { return b.supportedOS }

// Registry holds all available installers
type Registry struct {
	installers map[string]Installer
}

func NewRegistry() *Registry {
	return &Registry{
		installers: make(map[string]Installer),
	}
}

func (r *Registry) Register(installer Installer) {
	r.installers[installer.Name()] = installer
}

func (r *Registry) Get(name string) (Installer, bool) {
	i, ok := r.installers[name]
	return i, ok
}

func (r *Registry) GetByCategory(category Category) []Installer {
	var result []Installer
	for _, i := range r.installers {
		if i.Category() == category {
			result = append(result, i)
		}
	}
	return result
}

func (r *Registry) All() []Installer {
	result := make([]Installer, 0, len(r.installers))
	for _, i := range r.installers {
		result = append(result, i)
	}
	return result
}
