# Runtimes & Languages

Install programming language runtimes and development environments.

## ⚡ Available Installers

### Node.js

**JavaScript runtime built on Chrome's V8 engine**

| Property | Value |
| -------- | ----- |
| Category | Runtime |
| Supported OS | Linux, macOS, Windows |
| Package Manager | npm, pnpm, yarn, bun |

**Installation via NVM (Recommended):**
```bash
# Install NVM
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash

# Reload shell
source ~/.bashrc

# Install latest LTS
nvm install --lts

# Verify
node --version
npm --version
```

**Install Package Managers:**
```bash
# pnpm
npm install -g pnpm

# yarn
npm install -g yarn

# bun
curl -fsSL https://bun.sh/install | bash
```

---

### Go (Golang)

**Fast, statically typed, compiled programming language**

| Property | Value |
| -------- | ----- |
| Category | Runtime |
| Supported OS | Linux, macOS, Windows |
| Package Manager | go mod |

**Installation:**
```bash
# Download latest
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz

# Extract
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz

# Add to PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Verify
go version
```

**Or via package manager:**
```bash
# Ubuntu (may not be latest)
sudo apt install -y golang-go

# macOS
brew install go
```

---

### Python

**Versatile programming language for scripting, web, data science**

| Property | Value |
| -------- | ----- |
| Category | Runtime |
| Supported OS | Linux, macOS, Windows |
| Package Manager | pip, pipenv, poetry |

**Installation:**
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install -y python3 python3-pip python3-venv

# CentOS/RHEL
sudo yum install -y python3 python3-pip

# macOS
brew install python
```

**Virtual Environments:**
```bash
# Create venv
python3 -m venv myenv

# Activate
source myenv/bin/activate

# Install packages
pip install -r requirements.txt
```

**Install pipenv:**
```bash
pip install pipenv
pipenv install
pipenv shell
```

---

### PHP

**Server-side scripting language for web development**

| Property | Value |
| -------- | ----- |
| Category | Runtime |
| Supported OS | Linux, macOS |
| Package Manager | Composer |

**Installation:**
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install -y php php-cli php-fpm php-mysql php-curl php-json php-mbstring php-xml php-zip

# CentOS/RHEL
sudo yum install -y php php-cli php-fpm php-mysqlnd

# macOS
brew install php
```

**Install Composer:**
```bash
curl -sS https://getcomposer.org/installer | php
sudo mv composer.phar /usr/local/bin/composer
composer --version
```

---

## 📊 Version Management

| Language | Version Manager |
| -------- | --------------- |
| Node.js | nvm, fnm, volta |
| Python | pyenv |
| Go | gvm |
| PHP | phpbrew |

## 💡 Tips

### Check Installed Versions
```bash
node --version
go version
python3 --version
php --version
```

### Update PATH
Add these to `~/.bashrc` or `~/.zshrc`:
```bash
# Node.js (NVM)
export NVM_DIR="$HOME/.nvm"
[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"

# Go
export PATH=$PATH:/usr/local/go/bin
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
```

---

← [[Web Servers]] | [[Databases]] →
