#!/bin/bash
set -e

echo "Начало инициализации"

# Ожидаем, пока мастер станет доступен
until pg_isready -h postgres-master -p 5432 -U replicator; do
  sleep 1
done

echo "Удаляем старые данные (если они есть)"
rm -rf "$PGDATA"/*


echo "Делаем pg_basebackup с мастера"
PGPASSWORD=pass pg_basebackup \
  -h postgres-master -p 5432 -U replicator \
  -D /var/lib/postgresql/data \
  -P -R -X stream -C -S postgres_slave_3

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

primary_conninfo = 'host=postgres-master port=5432 user=replicator password=pass application_name=postgres_slave_3'
hot_standby = on
EOF

# Даём права postgres на данные
chown -R postgres:postgres /var/lib/postgresql/data
chmod -R 700 /var/lib/postgresql/data

# Перезагружаем конфигурацию без перезапуска сервера
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_REPLICA_USER" --dbname "$POSTGRES_DB" -c "SELECT pg_reload_conf();"
echo "Завершение инициализации"
