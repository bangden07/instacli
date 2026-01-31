# AI CLI Tools

InstaCli includes installers for popular AI coding assistant command-line tools.

## Available AI CLI Tools

| Tool | Description | Requirements |
|------|-------------|--------------|
| **Claude CLI** | Anthropic's Claude Code AI assistant | Node.js, npm |
| **Gemini CLI** | Google's Gemini AI CLI | Node.js, npm |
| **Codex CLI** | OpenAI's Codex CLI | Node.js, npm |
| **Aider** | AI pair programming assistant | Python |
| **Kilo Code** | VS Code AI extension | VS Code |
| **Continue** | Open-source AI code assistant | VS Code |
| **OpenCode CLI** | Open-source AI coding assistant | Node.js, npm |

## Usage

1. Navigate to **AI CLI** category in InstaCli
2. Select the tool you want to install
3. Press `i` to install or `g` to generate script

## Tool Details

### Claude CLI

Anthropic's official Claude Code CLI for interactive AI coding assistance.

```bash
# After installation
claude

# With a prompt
claude "explain this code"
```

### Gemini CLI

Google's Gemini AI CLI for AI-assisted development.

```bash
# After installation
gemini

# Ask a question
gemini "how do I create a REST API"
```

### Codex CLI

OpenAI's Codex-powered CLI for code generation.

```bash
# After installation
codex

# Generate code
codex "write a function to sort an array"
```

### Aider

AI pair programmer that works with your local Git repository.

```bash
# After installation
aider

# Start with a file
aider myfile.py
```

### Kilo Code

AI coding extension for VS Code with multi-model support.

After installation, open VS Code and the extension will be available in the sidebar.

### Continue

Open-source AI code assistant extension for VS Code and JetBrains.

After installation:
1. Open VS Code
2. Look for Continue in the sidebar
3. Configure your preferred AI provider

### OpenCode CLI

Open-source AI coding assistant with terminal interface.

```bash
# After installation
opencode

# Start coding with AI
opencode "create a todo app"
```

## API Keys

Most AI CLI tools require API keys:

| Tool | API Key Source |
|------|---------------|
| Claude CLI | [Anthropic Console](https://console.anthropic.com) |
| Gemini CLI | [Google AI Studio](https://aistudio.google.com) |
| Codex CLI | [OpenAI Platform](https://platform.openai.com) |
| Aider | OpenAI, Anthropic, or local models |
| OpenCode CLI | Various supported providers |

## See Also

- [[MCP-Servers]] - Model Context Protocol servers
- [[Installation]] - Getting started with InstaCli
