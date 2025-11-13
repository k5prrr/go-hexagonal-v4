# TradeIN 2.0

## Установка
1. cp .env.example .env
2. Tg token без вебхуков
3. cd docker && podman compose -f db.yml up -d
4. sudo go mod tidy
5. clear && go run ./cmd/app/main.go

## Для теста
https://github.com/swagger-api/swagger-ui/releases

http://localhost:8080/api/v2/increment
http://localhost:8080/api/v2/doc/
http://localhost:8080/




http://localhost:8080/api/v2/registration/linkcheckphone
http://localhost:45010/?pgsql=db&username=postgres_user&db=postgres_db&ns=public&select=auth_codes



https://h7team.ru/installGO/

