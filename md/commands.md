## установка компилятора go
brew install go  

## проверить версию go 
go version

## инициализация проекта
go mod init url-shortener

## запуск проекта
go run main.go. (или " go run . "  если у нас несколько файлов go )

## прочитать наш Dockerfile и собрать из него готовый шаблон.
docker build -t shortener .

## запускаем контейнер с пробросом портов
docker run -p 5001:8080 shortener

curl -X POST http://localhost:5001/api/shorten \
     -H "Content-Type: application/json" \
     -d '{"url": "https://developer.mozilla.org/ru/"}'


## Запуск одним кликом. Флаг -d (detached) означает "запусти в 
## фоновом режиме". Твой терминал останется свободным.
docker compose up -d

## запуск с пересбором билда
docker compose up -d --build

## вырубить наш сервер с базой
docker compose down

## очистить докер от безымянных висячих образов
docker image prune

## установка драйвера для PostgreSQL
go get github.com/lib/pq

## Как научить Go читать .env?
go get github.com/joho/godotenv

## Учим Go запускать миграции
go get -u github.com/golang-migrate/migrate/v4

## для горячей перезагрузки используем утилиту Air
go install github.com/air-verse/air@latest

## Добавляем путь в настройки твоего терминала (.zshrc) т.к. команду air без этого не будет видна.
## и после перезапускаем терминал
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.zshrc

## команда создаст файл .air.toml (настройки утилиты для Air).
air init

## теперь запуск проекта можно осуществлять через air простой командой :
air

