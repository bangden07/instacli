# SSH Remote Installation

Install tools on remote servers directly from InstaCli!

## 🌐 Overview

InstaCli supports remote installation via SSH, allowing you to:
- Configure SSH connection details
- Install tools on remote servers
- Generate scripts for remote execution

## ⚙️ Configure SSH

### Step 1: Open Settings
Navigate to **"Settings"** in the main menu.

### Step 2: Enter Edit Mode
Press `e` to enter edit mode.

### Step 3: Fill in Details

| Field | Description | Example |
| ----- | ----------- | ------- |
| Host | Server IP or hostname | `192.168.1.100` |
| Port | SSH port (default: 22) | `22` |
| User | SSH username | `root` |
| Password | SSH password | `****` |

### Step 4: Save
Press `Enter` to save, or `Esc` to cancel.

## 🔄 Switch to SSH Mode

Press `Tab` on the main menu to toggle between:
- **🖥️ Local** - Install on current machine
- **🌐 SSH** - Install on remote server

The current mode is shown in the header:
```
╭────────────────────────────────────╮
│  Target: 🌐 SSH (root@192.168.1.100)  │
╰────────────────────────────────────╯
```

## 📦 Install Remotely

1. Switch to SSH mode (Tab)
2. Select a category
3. Select an installer
4. Press `i` to install

InstaCli will generate and show the SSH command:
```bash
ssh root@192.168.1.100 -p 22 'bash -s' << 'EOF'
#!/bin/bash
# Docker installation script
curl -fsSL https://get.docker.com | sh
EOF
```

## 🔐 Security Best Practices

### Use SSH Keys (Recommended)
Instead of password authentication, use SSH keys:

1. Generate a key pair:
   ```bash
   ssh-keygen -t ed25519 -C "your_email@example.com"
   ```

2. Copy to server:
   ```bash
   ssh-copy-id user@server
   ```

3. Leave password field empty in InstaCli

### Use Non-Root User
Create a dedicated user with sudo access:
```bash
adduser deployer
usermod -aG sudo deployer
```

### Limit SSH Access
Configure firewall to limit SSH access:
```bash
ufw allow from 192.168.1.0/24 to any port 22
```

## 🔧 Troubleshooting

### Connection Refused
- Check if SSH is running: `systemctl status sshd`
- Check firewall: `ufw status`
- Verify port is correct

### Authentication Failed
- Verify username and password
- Check if password auth is enabled in `/etc/ssh/sshd_config`
- Try with SSH key instead

### Permission Denied
- Ensure user has sudo privileges
- Check `/etc/sudoers` for NOPASSWD if needed

## 📝 Manual Execution

If automatic SSH execution fails, you can:

1. Press `g` to generate script
2. Copy the script manually:
   ```bash
   scp ./scripts/generated/docker.sh root@server:/tmp/
   ssh root@server 'chmod +x /tmp/docker.sh && /tmp/docker.sh'
   ```

---

**Next:** [[System Status]] →
