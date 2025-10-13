#!/bin/bash
# sudo sh scripts/updateGO.sh

set -e  # Останавливать выполнение при ошибках

echo "Получаем информацию о последней версии Go..."

# Получаем последнюю версию с официального JSON API
VERSION=$(curl -s https://go.dev/dl/?mode=json | head -n 1 | jq -r '.version')

if [ -z "$VERSION" ]; then
  echo "Не удалось получить версию Go"
  exit 1
fi

echo "Последняя версия: $VERSION"

# Формируем имя файла и ссылку
FILENAME="go${VERSION}.linux-amd64.tar.gz"
URL="https://go.dev/dl/${FILENAME}"

# Проверяем, доступен ли файл
if curl --output /dev/null --silent --head --fail "$URL"; then
  echo "Файл найден: $FILENAME"
else
  echo "Ошибка: файл $FILENAME не найден по адресу $URL"
  exit 1
fi

# Удаляем предыдущую установку
echo "Удаляем предыдущую версию Go..."
sudo rm -rf /usr/local/go

# Скачиваем и распаковываем
echo "Скачиваем $FILENAME..."
wget -q "$URL" -O "$FILENAME"

echo "Распаковываем в /usr/local..."
sudo tar -C /usr/local -xzf "$FILENAME"

# Удаляем скаченный архив
rm "$FILENAME"

# Настройка переменных окружения, если ещё не сделано
GO_BIN="/usr/local/go/bin"
if [[ ":$PATH:" != *":$GO_BIN:"* ]]; then
  echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
  echo 'export GOPATH=$HOME/go' >> ~/.bashrc
  echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.bashrc
  source ~/.bashrc
fi

# Проверка результата
echo "Go успешно установлен!"
go version
