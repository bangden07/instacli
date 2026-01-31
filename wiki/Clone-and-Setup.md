# Clone & Setup Feature

**New in v1.2.0** - Automatically clone and setup any Git repository!

## 📥 Overview

The Clone & Setup feature allows you to:
1. Paste any Git repository URL
2. Automatically detect the project type
3. Install required runtime if needed
4. Install all dependencies

## 🎯 How to Use

### Step 1: Open Clone & Setup
Navigate to **"Clone & Setup"** in the main menu.

### Step 2: Enter Repository URL
Paste your Git repository URL:
```
https://github.com/user/my-project.git
```

### Step 3: Press Enter
InstaCli will:
1. Clone the repository
2. Detect project type
3. Install dependencies

## ✨ Supported Project Types

| Icon | Type | Detection | Package Manager |
| ---- | ---- | --------- | --------------- |
| 📦 | Node.js | `package.json` | npm, pnpm, yarn, bun |
| 🐍 | Python | `requirements.txt`, `pyproject.toml` | pip, pipenv |
| 🐹 | Go | `go.mod` | go mod |
| 🐘 | PHP | `composer.json` | composer |
| 💎 | Ruby | `Gemfile` | bundler |
| 🦀 | Rust | `Cargo.toml` | cargo |
| 🐳 | Docker | `docker-compose.yml` | docker compose |

## 🔍 Framework Detection

InstaCli also detects specific frameworks:

### Node.js Frameworks
- **Next.js** - `next.config.js`
- **Nuxt.js** - `nuxt.config.js`
- **Vite** - `vite.config.js`
- **Angular** - `angular.json`
- **Svelte** - `svelte.config.js`

### Python Frameworks
- **Django** - `manage.py`
- **Flask** - `app.py`, `wsgi.py`

### PHP Frameworks
- **Laravel** - `artisan`
- **Symfony** - `bin/console`

## 📦 Package Manager Detection

For Node.js projects, InstaCli detects the preferred package manager:

| Lock File | Package Manager |
| --------- | --------------- |
| `pnpm-lock.yaml` | pnpm |
| `yarn.lock` | yarn |
| `bun.lockb` | bun |
| `package-lock.json` | npm |

## 🚀 Example Workflows

### Next.js Project
```
1. Clone: https://github.com/user/nextjs-app
2. Detected: Next.js (Node.js)
3. Package Manager: pnpm (pnpm-lock.yaml found)
4. Run: pnpm install
5. Ready!
```

### Laravel Project
```
1. Clone: https://github.com/user/laravel-app
2. Detected: Laravel (PHP)
3. Run: composer install
4. Run: cp .env.example .env
5. Run: php artisan key:generate
6. Ready!
```

### Python Django Project
```
1. Clone: https://github.com/user/django-app
2. Detected: Django (Python)
3. Run: python3 -m venv venv
4. Run: source venv/bin/activate
5. Run: pip install -r requirements.txt
6. Ready!
```

## ⚠️ Runtime Installation

If the required runtime is not installed, InstaCli will offer to install it:

- **Node.js** → Installs via NVM
- **Python** → Installs via apt/brew
- **Go** → Downloads from go.dev
- **PHP** → Installs via apt/brew
- **Ruby** → Installs via apt/brew
- **Rust** → Installs via rustup

## 🐳 Docker Support

If `docker-compose.yml` is detected, InstaCli will ask:
```
🐳 Docker Compose detected
Start with Docker Compose? [y/N]
```

---

**Next:** [[SSH Remote Installation]] →
