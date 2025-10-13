Это для быстрого старта проекта
https://github.com/k5prrr/go-hexagonal-v1

Ещё примеры
https://github.com/GolangLessons/sso


make up




go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest




# OLD ---

# Install
1. micro .gitignore
2. sh ./scripts/installModules.sh
3. sh scripts/start.sh

# Start
in main dir
sh scripts/start.sh
http://localhost:8081/testSpeed



# Исправлю
Пример мидлваре
Приведите имена пакетов к нижнему регистру.
Пишите сообщения и логи на английском для универсальности.
Добавьте тесты: В проекте отсутствуют тесты (папка tests пуста).
Используйте линтеры: Добавьте настройки для golangci-lint в .golangci.yml.
Документация: Добавьте описание API в api/openapi.yaml и примеры использования в README_API.md.


Hexagonal architecture! - Сейчас в моде

https://www.youtube.com/watch?v=3YTLDYG5MnQ

https://www.youtube.com/watch?v=0Fhsgmz-Gig

https://www.youtube.com/watch?app=desktop&v=yyrvXnXLnU8

covrom


sudo apt install podman-compose
sudo nano /etc/containers/registries.conf
[registries.search]
registries = ["docker.io"]

[registries.insecure]
registries = []


mkdir -p ./pg_data ./migrations

sudo podman stop $(sudo podman ps -aq) && sudo podman rm $(sudo podman ps -aq)
podman-compose up -d
podman ps -a
