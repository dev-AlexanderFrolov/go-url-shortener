1. Установка WSL 2 (Ubuntu)
Открыть PowerShell от имени администратора.

Выполнить: wsl --install.

После перезагрузки настроить username и password.

Назначить Ubuntu главной: wsl --set-default Ubuntu.

2. Интеграция с Docker Desktop
Docker Settings -> Resources -> WSL Integration.

Включить тумблер для Ubuntu.

Нажать Apply & Restart.

3. Настройка VS Code
Установить расширение WSL (от Microsoft).

Нажать на иконку >< (левый нижний угол) -> Connect to WSL using Distro -> Ubuntu.

4. Настройка SSH для GitHub (в терминале WSL)
Представиться:
git config --global user.name "Your Name"
git config --global user.email "your@email.com"

Создать ключ:
ssh-keygen -t ed25519 -C "your@email.com"
// 3 раза нажать Enter


Скопировать ключ:
cat ~/.ssh/id_ed25519.pub

Добавить скопированную строку в GitHub Settings -> SSH keys.

Установка Go в WSL (Ubuntu)
Выполнять по очереди:

# Скачивание
wget https://go.dev/dl/go1.26.1.linux-amd64.tar.gz

# Установка
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.26.1.linux-amd64.tar.gz

# Настройка путей (один раз)
echo 'export PATH=$PATH:/usr/local/go/bin:~/go/bin' >> ~/.bashrc
source ~/.bashrc

# Очистка мусора
rm go1.26.1.linux-amd64.tar.gz

# Проверка
go version

6. Запуск проекта
Клонировать через SSH: git clone git@github.com:USER/REPO.git.

Перейти в папку: cd repo_name.

Создать .env: cp .env.example .env (и прописать пароли).

Установить зависимости: go mod tidy.

Поднять БД: docker compose up db -d.

Установить Hot Reload: go install github.com/air-verse/air@latest.

Запустить: air.

Важное примечание по Air (для Linux/WSL)
Убедись, что в файле .air.toml НЕТ расширения .exe:

Ini, TOML
cmd = "go build -o ./tmp/main ."
bin = "tmp/main"