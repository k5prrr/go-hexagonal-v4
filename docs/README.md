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

## План

### регистрация
нажимаем на вход
приходит код и ссылка с кодом перенаправляет на телеграм
в телеграмме если чела нет, то всегда запрос номера, после запроса номера возврат на страницу

не
бред

вводит номер телефона
нажимает вход

Фио
дата рождения
зарегистрироваться
кнопка которая отправляет данные и ведёт в тг
в старт

не, тоже не понятно

---

Так
регистрация
Форма
V
при кнопке телеграм
сохранение входа в куки с пометкой регистрация (Это он уже в пользователях? костыльно)
запрос ссылки и генерация кода с пометкой регистрация
+
А если это всё через телефон? 
да пофиг, пусть возвращается на сайт

Gjckt jnghfdrb yjvthf



---
Человек является пользователем только после отправки номера в тг
Со статусом душа
после такого входа идёт перенаправление на форму заполнения данных


