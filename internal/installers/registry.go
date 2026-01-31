package installers

// DefaultRegistry returns a registry with all available installers
func DefaultRegistry() *Registry {
	r := NewRegistry()

	// Containers
	r.Register(NewDockerInstaller())

	// Runtime & Languages
	r.Register(NewNodeJSInstaller())
	r.Register(NewGolangInstaller())

	// Web Servers
	r.Register(NewLAMPInstaller())
	r.Register(NewLEMPInstaller())

	// Databases
	r.Register(NewMySQLInstaller())
	r.Register(NewPostgreSQLInstaller())
	r.Register(NewMongoDBInstaller())
	r.Register(NewRedisInstaller())

	// Automation
	r.Register(NewN8NInstaller())

	// Frameworks
	r.Register(NewLaravelKitInstaller())
	r.Register(NewNextJSKitInstaller())

	// Security
	r.Register(NewUFWInstaller())
	r.Register(NewCertbotInstaller())
	r.Register(NewFail2banInstaller())

	// Monitoring
	r.Register(NewPrometheusInstaller())
	r.Register(NewGrafanaInstaller())
	r.Register(NewNetdataInstaller())

	// Infrastructure
	r.Register(NewNginxProxyManagerInstaller())
	r.Register(NewTraefikInstaller())
	r.Register(NewMinIOInstaller())

	// VPN
	r.Register(NewWireGuardInstaller())

	// CI/CD
	r.Register(NewJenkinsInstaller())
	r.Register(NewGitLabRunnerInstaller())
	r.Register(NewGitHubActionsRunnerInstaller())

	// DNS & Network
	r.Register(NewPiholeInstaller())

	// CMS & Blog
	r.Register(NewWordPressInstaller())
	r.Register(NewGhostInstaller())

	// Backup
	r.Register(NewResticInstaller())

	// AI CLI Tools
	r.Register(NewClaudeCLIInstaller())
	r.Register(NewGeminiCLIInstaller())
	r.Register(NewCodexCLIInstaller())
	r.Register(NewAiderInstaller())
	r.Register(NewKiloCodeInstaller())
	r.Register(NewContinueInstaller())

	// MCP Servers
	r.Register(NewContext7MCPInstaller())
	r.Register(NewPlaywrightMCPInstaller())
	r.Register(NewGitHubMCPInstaller())
	r.Register(NewFilesystemMCPInstaller())
	r.Register(NewPostgresMCPInstaller())
	r.Register(NewBraveSearchMCPInstaller())
	r.Register(NewMemoryMCPInstaller())
	r.Register(NewSequentialThinkingMCPInstaller())

	return r
}

// CategoryInstallers returns a grouped map of installers by category
func CategoryInstallers(r *Registry) map[Category][]Installer {
	result := make(map[Category][]Installer)
	for _, inst := range r.All() {
		cat := inst.Category()
		result[cat] = append(result[cat], inst)
	}
	return result
}
