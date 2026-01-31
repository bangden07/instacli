# MCP Servers

InstaCli includes installers for Model Context Protocol (MCP) servers that enhance AI capabilities.

## What is MCP?

Model Context Protocol (MCP) is a standard for connecting AI systems with external tools, data sources, and services. MCP servers provide capabilities like:

- File system access
- Database queries
- API integrations
- Browser automation
- Search functionality

## Available MCP Servers

| Server | Purpose | Requirements |
|--------|---------|--------------|
| **Context7** | Documentation lookup for programming libraries | Node.js, npm |
| **Playwright** | Browser automation and testing | Node.js, npm |
| **GitHub** | Repository integration and management | Node.js, npm |
| **Filesystem** | File system access for AI | Node.js, npm |
| **PostgreSQL** | Database access and queries | Node.js, npm |
| **Brave Search** | Web search integration | Node.js, npm |
| **Memory** | Persistent memory for AI | Node.js, npm |
| **Sequential Thinking** | Step-by-step reasoning | Node.js, npm |

## Usage

1. Navigate to **MCP Servers** category in InstaCli
2. Select the MCP server you want to install
3. Press `i` to install or `g` to generate script

## Server Details

### Context7

Lookup documentation for programming libraries and frameworks.

```json
{
  "mcpServers": {
    "context7": {
      "command": "npx",
      "args": ["-y", "@upstash/context7-mcp"]
    }
  }
}
```

### Playwright

Browser automation for testing and web scraping.

```json
{
  "mcpServers": {
    "playwright": {
      "command": "npx",
      "args": ["-y", "@anthropic-ai/mcp-server-playwright"]
    }
  }
}
```

### GitHub

Repository management and integration.

```json
{
  "mcpServers": {
    "github": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-github"],
      "env": {
        "GITHUB_TOKEN": "your-github-token"
      }
    }
  }
}
```

### Filesystem

File system access for AI.

```json
{
  "mcpServers": {
    "filesystem": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-filesystem", "/path/to/allowed/directory"]
    }
  }
}
```

### PostgreSQL

Database access for AI queries.

```json
{
  "mcpServers": {
    "postgres": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-postgres"],
      "env": {
        "DATABASE_URL": "postgresql://user:pass@localhost/db"
      }
    }
  }
}
```

### Brave Search

Web search integration.

```json
{
  "mcpServers": {
    "brave-search": {
      "command": "npx",
      "args": ["-y", "@anthropic-ai/mcp-server-brave-search"],
      "env": {
        "BRAVE_API_KEY": "your-brave-api-key"
      }
    }
  }
}
```

### Memory

Persistent memory for AI sessions.

```json
{
  "mcpServers": {
    "memory": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-memory"]
    }
  }
}
```

### Sequential Thinking

Step-by-step reasoning for complex problems.

```json
{
  "mcpServers": {
    "sequential-thinking": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-sequential-thinking"]
    }
  }
}
```

## Configuration

MCP servers are typically configured in your AI client's settings:

- **Claude Desktop**: `~/Library/Application Support/Claude/claude_desktop_config.json` (macOS)
- **VS Code Copilot**: Settings > MCP Servers
- **Other clients**: Check your AI tool's documentation

## See Also

- [[AI-CLI]] - AI CLI tools
- [[Installation]] - Getting started with InstaCli
