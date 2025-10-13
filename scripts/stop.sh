# sh scripts/start.sh
clear
docker stop $(docker ps -aq)
docker rm $(docker ps -aq)
