# Adding New Installers

Learn how to add new installers to InstaCli.

## 📁 File Structure

Installers are defined in `internal/installers/`:

```
internal/installers/
├── base.go          # Base interface and types
├── registry.go      # Installer registry
├── webserver.go     # Web server installers
├── runtime.go       # Runtime installers
├── database.go      # Database installers
├── monitoring.go    # Monitoring installers
├── infrastructure.go
├── cicd.go          # CI/CD installers
└── additional.go    # Additional installers
```

## 🔧 Step 1: Define the Installer

Create a new struct that embeds `BaseInstaller`:

```go
// MyToolInstaller installs MyTool
type MyToolInstaller struct {
    BaseInstaller
}

func NewMyToolInstaller() *MyToolInstaller {
    return &MyToolInstaller{
        BaseInstaller: BaseInstaller{
            name:        "MyTool",
            description: "Description of MyTool",
            category:    CategoryAutomation, // Choose appropriate category
            icon:        "🔧",
            supportedOS: []OS{OSLinux, OSMacOS},
        },
    }
}
```

## 📝 Step 2: Implement Methods

### GenerateInstallScript

```go
func (i *MyToolInstaller) GenerateInstallScript(os OS, pm PackageManager) string {
    var script strings.Builder
    script.WriteString("#!/bin/bash\nset -e\n\n")
    script.WriteString("echo '🔧 Installing MyTool...'\n\n")

    switch pm {
    case PMApt:
        script.WriteString("sudo apt-get update\n")
        script.WriteString("sudo apt-get install -y mytool\n")
    case PMYum:
        script.WriteString("sudo yum install -y mytool\n")
    case PMBrew:
        script.WriteString("brew install mytool\n")
    }

    script.WriteString("\necho '✅ MyTool installed!'\n")
    return script.String()
}
```

### GenerateUninstallScript

```go
func (i *MyToolInstaller) GenerateUninstallScript(os OS, pm PackageManager) string {
    var script strings.Builder
    script.WriteString("#!/bin/bash\nset -e\n\n")
    script.WriteString("echo '🗑️ Uninstalling MyTool...'\n\n")

    switch pm {
    case PMApt:
        script.WriteString("sudo apt-get remove -y mytool\n")
    case PMYum:
        script.WriteString("sudo yum remove -y mytool\n")
    case PMBrew:
        script.WriteString("brew uninstall mytool\n")
    }

    script.WriteString("\necho '✅ MyTool uninstalled!'\n")
    return script.String()
}
```

### Install and Uninstall

```go
func (i *MyToolInstaller) Install(exec executor.Executor) error {
    return exec.Run(i.GenerateInstallScript(OSLinux, PMApt))
}

func (i *MyToolInstaller) Uninstall(exec executor.Executor) error {
    return exec.Run(i.GenerateUninstallScript(OSLinux, PMApt))
}
```

### IsInstalled

```go
func (i *MyToolInstaller) IsInstalled() bool {
    _, err := exec.LookPath("mytool")
    return err == nil
}
```

## 📋 Step 3: Register the Installer

In `registry.go`, add to `DefaultRegistry()`:

```go
func DefaultRegistry() *Registry {
    r := NewRegistry()

    // ... existing installers ...

    // Automation
    r.Register(NewMyToolInstaller())  // Add your installer

    return r
}
```

## 🏷️ Available Categories

| Constant | Display Name |
| -------- | ------------ |
| `CategoryWebServer` | Web Server Stack |
| `CategoryRuntime` | Runtime & Languages |
| `CategoryContainer` | Containers |
| `CategoryDatabase` | Databases |
| `CategoryFramework` | Frameworks |
| `CategoryAutomation` | Automation |
| `CategorySecurity` | Security |
| `CategoryMonitoring` | Monitoring |
| `CategoryInfrastructure` | Infrastructure |
| `CategoryCICD` | CI/CD |
| `CategoryDNS` | DNS & Network |
| `CategoryCMS` | CMS & Blog |
| `CategoryBackup` | Backup |
| `CategoryVPN` | VPN |

## ✅ Step 4: Build and Test

```bash
go build -o instacli ./cmd/instacli
./instacli
```

Verify your installer appears in the correct category.

## 💡 Best Practices

1. **Use latest stable versions** - Fetch from official sources
2. **Handle all package managers** - apt, yum, brew
3. **Add systemd services** - For background services
4. **Include verification** - Check if installation succeeded
5. **Clean up on uninstall** - Remove config files if appropriate

---

**Next:** [[Contributing]] →
