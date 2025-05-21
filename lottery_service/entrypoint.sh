#!/bin/sh
echo "Waiting for MongoDB..."
until nc -z mongo 27017; do
  sleep 1
done
echo "MongoDB is ready!"
echo "Waiting for NATS..."
until nc -z nats 4222; do
  sleep 1
done
echo "NATS is ready!"
# Start application
exec "$@"