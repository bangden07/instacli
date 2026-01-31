package installers

// ============================================================
// MySQL Installer
// ============================================================
type MySQLInstaller struct {
	BaseInstaller
}

func NewMySQLInstaller() *MySQLInstaller {
	return &MySQLInstaller{
		BaseInstaller: BaseInstaller{
			name:        "MySQL",
			description: "MySQL Database Server",
			category:    CategoryDatabase,
			icon:        "🐬",
			supportedOS: []OS{OSLinux, OSMacOS},
		},
	}
}

func (i *MySQLInstaller) RequiredPackageManagers() []PackageManager {
	return []PackageManager{PMApt, PMYum, PMBrew}
}

func (i *MySQLInstaller) Dependencies() []string {
	return []string{}
}

func (i *MySQLInstaller) Install(executor Executor) error {
	_, err := executor.Run(i.GenerateInstallScript(executor.GetOS(), executor.GetPackageManager()))
	return err
}

func (i *MySQLInstaller) Uninstall(executor Executor) error {
	_, err := executor.Run(i.GenerateUninstallScript(executor.GetOS(), executor.GetPackageManager()))
	return err
}

func (i *MySQLInstaller) IsInstalled(executor Executor) (bool, error) {
	_, err := executor.Run("mysql --version")
	return err == nil, nil
}

func (i *MySQLInstaller) GenerateInstallScript(os OS, pm PackageManager) string {
	return `#!/bin/bash
set -e
echo "📦 Installing MySQL..."
if [ -f /etc/debian_version ]; then
    sudo apt-get update && sudo apt-get install -y mysql-server mysql-client
    sudo systemctl enable mysql && sudo systemctl start mysql
elif [ -f /etc/redhat-release ]; then
    sudo yum install -y mysql-server mysql
    sudo systemctl enable mysqld && sudo systemctl start mysqld
elif command -v brew &> /dev/null; then
    brew install mysql && brew services start mysql
fi
echo "✅ MySQL installed!"
mysql --version`
}

func (i *MySQLInstaller) GenerateUninstallScript(os OS, pm PackageManager) string {
	return `#!/bin/bash
sudo apt-get remove -y mysql-server mysql-client || sudo yum remove -y mysql-server mysql || brew uninstall mysql`
}

// ============================================================
// PostgreSQL Installer
// ============================================================
type PostgreSQLInstaller struct {
	BaseInstaller
}

func NewPostgreSQLInstaller() *PostgreSQLInstaller {
	return &PostgreSQLInstaller{
		BaseInstaller: BaseInstaller{
			name:        "PostgreSQL",
			description: "PostgreSQL Database Server",
			category:    CategoryDatabase,
			icon:        "🐘",
			supportedOS: []OS{OSLinux, OSMacOS},
		},
	}
}

func (i *PostgreSQLInstaller) RequiredPackageManagers() []PackageManager {
	return []PackageManager{PMApt, PMYum, PMBrew}
}
func (i *PostgreSQLInstaller) Dependencies() []string { return []string{} }
func (i *PostgreSQLInstaller) Install(executor Executor) error {
	_, err := executor.Run(i.GenerateInstallScript(executor.GetOS(), executor.GetPackageManager()))
	return err
}
func (i *PostgreSQLInstaller) Uninstall(executor Executor) error {
	_, err := executor.Run(i.GenerateUninstallScript(executor.GetOS(), executor.GetPackageManager()))
	return err
}
func (i *PostgreSQLInstaller) IsInstalled(executor Executor) (bool, error) {
	_, err := executor.Run("psql --version")
	return err == nil, nil
}

func (i *PostgreSQLInstaller) GenerateInstallScript(os OS, pm PackageManager) string {
	return `#!/bin/bash
set -e
echo "📦 Installing PostgreSQL..."
if [ -f /etc/debian_version ]; then
    sudo apt-get update && sudo apt-get install -y postgresql postgresql-contrib
    sudo systemctl enable postgresql && sudo systemctl start postgresql
elif [ -f /etc/redhat-release ]; then
    sudo yum install -y postgresql-server postgresql-contrib
    sudo postgresql-setup initdb && sudo systemctl enable postgresql && sudo systemctl start postgresql
elif command -v brew &> /dev/null; then
    brew install postgresql && brew services start postgresql
fi
echo "✅ PostgreSQL installed!"
psql --version`
}

func (i *PostgreSQLInstaller) GenerateUninstallScript(os OS, pm PackageManager) string {
	return `#!/bin/bash
sudo apt-get remove -y postgresql postgresql-contrib || sudo yum remove -y postgresql-server || brew uninstall postgresql`
}

// ============================================================
// MongoDB Installer
// ============================================================
type MongoDBInstaller struct {
	BaseInstaller
}

func NewMongoDBInstaller() *MongoDBInstaller {
	return &MongoDBInstaller{
		BaseInstaller: BaseInstaller{
			name:        "MongoDB",
			description: "MongoDB NoSQL Database",
			category:    CategoryDatabase,
			icon:        "🍃",
			supportedOS: []OS{OSLinux, OSMacOS},
		},
	}
}

func (i *MongoDBInstaller) RequiredPackageManagers() []PackageManager {
	return []PackageManager{PMApt, PMYum, PMBrew}
}
func (i *MongoDBInstaller) Dependencies() []string { return []string{} }
func (i *MongoDBInstaller) Install(executor Executor) error {
	_, err := executor.Run(i.GenerateInstallScript(executor.GetOS(), executor.GetPackageManager()))
	return err
}
func (i *MongoDBInstaller) Uninstall(executor Executor) error {
	_, err := executor.Run(i.GenerateUninstallScript(executor.GetOS(), executor.GetPackageManager()))
	return err
}
func (i *MongoDBInstaller) IsInstalled(executor Executor) (bool, error) {
	_, err := executor.Run("mongod --version")
	return err == nil, nil
}

func (i *MongoDBInstaller) GenerateInstallScript(os OS, pm PackageManager) string {
	return `#!/bin/bash
set -e
echo "📦 Installing MongoDB..."
if [ -f /etc/debian_version ]; then
    curl -fsSL https://pgp.mongodb.com/server-7.0.asc | sudo gpg -o /usr/share/keyrings/mongodb-server-7.0.gpg --dearmor
    echo "deb [arch=amd64 signed-by=/usr/share/keyrings/mongodb-server-7.0.gpg] https://repo.mongodb.org/apt/ubuntu jammy/mongodb-org/7.0 multiverse" | sudo tee /etc/apt/sources.list.d/mongodb-org-7.0.list
    sudo apt-get update && sudo apt-get install -y mongodb-org
    sudo systemctl enable mongod && sudo systemctl start mongod
elif command -v brew &> /dev/null; then
    brew tap mongodb/brew && brew install mongodb-community && brew services start mongodb-community
fi
echo "✅ MongoDB installed!"
mongod --version`
}

func (i *MongoDBInstaller) GenerateUninstallScript(os OS, pm PackageManager) string {
	return `#!/bin/bash
sudo apt-get remove -y mongodb-org || brew uninstall mongodb-community`
}

// ============================================================
// Redis Installer
// ============================================================
type RedisInstaller struct {
	BaseInstaller
}

func NewRedisInstaller() *RedisInstaller {
	return &RedisInstaller{
		BaseInstaller: BaseInstaller{
			name:        "Redis",
			description: "Redis In-Memory Data Store",
			category:    CategoryDatabase,
			icon:        "🔴",
			supportedOS: []OS{OSLinux, OSMacOS},
		},
	}
}

func (i *RedisInstaller) RequiredPackageManagers() []PackageManager {
	return []PackageManager{PMApt, PMYum, PMBrew}
}
func (i *RedisInstaller) Dependencies() []string { return []string{} }
func (i *RedisInstaller) Install(executor Executor) error {
	_, err := executor.Run(i.GenerateInstallScript(executor.GetOS(), executor.GetPackageManager()))
	return err
}
func (i *RedisInstaller) Uninstall(executor Executor) error {
	_, err := executor.Run(i.GenerateUninstallScript(executor.GetOS(), executor.GetPackageManager()))
	return err
}
func (i *RedisInstaller) IsInstalled(executor Executor) (bool, error) {
	_, err := executor.Run("redis-server --version")
	return err == nil, nil
}

func (i *RedisInstaller) GenerateInstallScript(os OS, pm PackageManager) string {
	return `#!/bin/bash
set -e
echo "📦 Installing Redis..."
if [ -f /etc/debian_version ]; then
    sudo apt-get update && sudo apt-get install -y redis-server
    sudo systemctl enable redis-server && sudo systemctl start redis-server
elif [ -f /etc/redhat-release ]; then
    sudo yum install -y redis && sudo systemctl enable redis && sudo systemctl start redis
elif command -v brew &> /dev/null; then
    brew install redis && brew services start redis
fi
echo "✅ Redis installed!"
redis-server --version`
}

func (i *RedisInstaller) GenerateUninstallScript(os OS, pm PackageManager) string {
	return `#!/bin/bash
sudo apt-get remove -y redis-server || sudo yum remove -y redis || brew uninstall redis`
}
