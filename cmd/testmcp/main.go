package main

import (
	"fmt"

	"github.com/instacli/instacli/internal/installers"
)

func main() {
	// Create Context7 installer and generate Linux script
	ctx7 := &installers.Context7MCPInstaller{}
	script := ctx7.GenerateInstallScript(installers.OSLinux, installers.PMNone)
	fmt.Println("=== GENERATED SCRIPT FOR LINUX ===")
	fmt.Println(script)
}
