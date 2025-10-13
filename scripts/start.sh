# sh scripts/start.sh


go fmt ./...
clear
#docker-compose up --force-recreate -d
#source ./.env
go run cmd/app/main.go
