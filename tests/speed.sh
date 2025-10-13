#!/bin/bash
# sh tests/speed.sh

# sudo apt-get install apache2-utils

# URL для тестирования
URL="http://localhost:8081/testSpeed"

# Количество запросов
REQUESTS=10000

# Параллельные запросы
CONCURRENCY=100

# Запуск Apache Benchmark
ab -n $REQUESTS -c $CONCURRENCY $URL
