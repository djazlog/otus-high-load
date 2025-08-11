#!/bin/bash
set -e

echo "Начало инициализации"

# Ожидаем, пока мастер станет доступен
until pg_isready -h postgres-master -p 5432 -U replicator; do
  sleep 1
done

echo "Удаляем старые данные (если они есть)"
# shellcheck disable=SC2115
rm -rf "$PGDATA"/*


echo "Делаем pg_basebackup с мастера"
pg_basebackup -h postgres-master -p 5432 -D /var/lib/postgresql/data -U replicator -Fp -Xs -v -P --wal-method=stream

# Ждём инициализации данных PostgreSQL
while [ ! -f /var/lib/postgresql/data/postgresql.conf ]; do
  sleep 1
done

echo "Создаём файл standby.signal"
touch /var/lib/postgresql/data/standby.signal

head -n -7 /var/lib/postgresql/data/postgresql.conf > /var/lib/postgresql/data/postgresql.conf.tmp &&
mv /var/lib/postgresql/data/postgresql.conf.tmp /var/lib/postgresql/data/postgresql.conf
# Настраиваем подключение к мастеру
cat >> /var/lib/postgresql/data/postgresql.conf <<EOF

primary_conninfo = 'host=postgres-master port=5432 user=replicator password=pass application_name=postgres_slave'
hot_standby = on
EOF

# Даём права postgres на данные
chown -R postgres:postgres /var/lib/postgresql/data
echo "Завершение инициализации"
