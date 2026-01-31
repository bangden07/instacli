# CI/CD

Continuous Integration and Continuous Deployment tools.

## 🚀 Available Installers

### Jenkins

**Open-source automation server**

| Property | Value |
| -------- | ----- |
| Default Port | 8080 |
| Supported OS | Linux |
| Java Required | Yes (JDK 11+) |

**Installation (Ubuntu):**
```bash
# Install Java
sudo apt update
sudo apt install -y openjdk-17-jdk

# Add Jenkins repo
curl -fsSL https://pkg.jenkins.io/debian-stable/jenkins.io-2023.key | sudo tee /usr/share/keyrings/jenkins-keyring.asc > /dev/null

echo deb [signed-by=/usr/share/keyrings/jenkins-keyring.asc] https://pkg.jenkins.io/debian-stable binary/ | sudo tee /etc/apt/sources.list.d/jenkins.list > /dev/null

# Install
sudo apt update
sudo apt install -y jenkins

# Start
sudo systemctl enable jenkins
sudo systemctl start jenkins
```

**Initial Setup:**
```bash
# Get initial password
sudo cat /var/lib/jenkins/secrets/initialAdminPassword
```

1. Access `http://localhost:8080`
2. Enter initial password
3. Install suggested plugins
4. Create admin user

**Docker Installation:**
```yaml
version: '3.8'
services:
  jenkins:
    image: jenkins/jenkins:lts
    ports:
      - "8080:8080"
      - "50000:50000"
    volumes:
      - jenkins_home:/var/jenkins_home

volumes:
  jenkins_home:
```

---

### GitLab Runner

**CI/CD runner for GitLab**

| Property | Value |
| -------- | ----- |
| Supported OS | Linux, macOS, Windows |
| Executors | Shell, Docker, Kubernetes |

**Installation (Linux):**
```bash
# Add repository
curl -L "https://packages.gitlab.com/install/repositories/runner/gitlab-runner/script.deb.sh" | sudo bash

# Install
sudo apt install -y gitlab-runner
```

**Register Runner:**
```bash
sudo gitlab-runner register
# Enter GitLab URL: https://gitlab.com/
# Enter token: (from GitLab project Settings → CI/CD → Runners)
# Enter description: my-runner
# Enter tags: docker,linux
# Enter executor: docker
# Enter default Docker image: alpine:latest
```

**Configuration:**
```toml
# /etc/gitlab-runner/config.toml
concurrent = 4
check_interval = 0

[[runners]]
  name = "my-runner"
  url = "https://gitlab.com/"
  token = "YOUR_TOKEN"
  executor = "docker"
  [runners.docker]
    image = "alpine:latest"
    privileged = false
    volumes = ["/cache"]
```

---

### GitHub Actions Runner

**Self-hosted runner for GitHub Actions**

| Property | Value |
| -------- | ----- |
| Supported OS | Linux, macOS, Windows |
| Config Location | `~/.github-runner` |

**Installation:**
```bash
# Create directory
mkdir actions-runner && cd actions-runner

# Download
curl -o actions-runner-linux-x64-2.311.0.tar.gz -L https://github.com/actions/runner/releases/download/v2.311.0/actions-runner-linux-x64-2.311.0.tar.gz

# Extract
tar xzf ./actions-runner-linux-x64-2.311.0.tar.gz
```

**Configure:**
```bash
# Get token from: GitHub Repo → Settings → Actions → Runners → New self-hosted runner

./config.sh --url https://github.com/YOUR/REPO --token YOUR_TOKEN
```

**Run as Service:**
```bash
sudo ./svc.sh install
sudo ./svc.sh start
```

**Workflow Example:**
```yaml
# .github/workflows/ci.yml
name: CI

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  build:
    runs-on: self-hosted
    steps:
      - uses: actions/checkout@v4
      - name: Build
        run: |
          npm install
          npm run build
```

---

## 📊 Comparison

| Tool | Best For | Hosting |
| ---- | -------- | ------- |
| Jenkins | Complex pipelines, plugins | Self-hosted |
| GitLab Runner | GitLab projects | Self-hosted / GitLab.com |
| GitHub Actions | GitHub repos | GitHub.com / Self-hosted |

## 🔧 Pipeline Examples

### Jenkinsfile
```groovy
pipeline {
    agent any
    stages {
        stage('Build') {
            steps {
                sh 'npm install'
                sh 'npm run build'
            }
        }
        stage('Test') {
            steps {
                sh 'npm test'
            }
        }
        stage('Deploy') {
            steps {
                sh './deploy.sh'
            }
        }
    }
}
```

### .gitlab-ci.yml
```yaml
stages:
  - build
  - test
  - deploy

build:
  stage: build
  script:
    - npm install
    - npm run build

test:
  stage: test
  script:
    - npm test

deploy:
  stage: deploy
  script:
    - ./deploy.sh
  only:
    - main
```

---

← [[Infrastructure]] | [[Security]] →
