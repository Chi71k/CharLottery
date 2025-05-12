#!/bin/sh

echo "Ожидаем MongoDB..."
until nc -z mongo 27017; do
  sleep 1
done
echo "MongoDB готов!"

echo "Ожидаем NATS..."
until nc -z nats 4222; do
  sleep 1
done
echo "NATS готов!"

# Запуск приложения
exec "$@"
