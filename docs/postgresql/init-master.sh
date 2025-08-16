#!/bin/bash
set -e

echo "Начало инициализации"

# Ждём инициализации данных PostgreSQL
while [ ! -f /var/lib/postgresql/data/postgresql.conf ]; do
  sleep 1
done

# Создаём пользователя для репликации
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
    -- Включаем расширение Citus
    CREATE EXTENSION IF NOT EXISTS citus;


    create role replicator with REPLICATION LOGIN ENCRYPTED password 'pass';
    CREATE USER readonly WITH PASSWORD 'pass';
    GRANT CONNECT ON DATABASE otus TO readonly;
    GRANT USAGE ON SCHEMA public TO readonly;
    GRANT SELECT ON ALL TABLES IN SCHEMA public TO readonly;
    ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT ON TABLES TO readonly;
EOSQL

# Добавляем настройки в postgresql.conf
cat >> /var/lib/postgresql/data/postgresql.conf <<EOF
ssl = off
wal_level = replica
max_wal_senders = 10
wal_keep_size = 1GB
hot_standby = on
synchronous_commit = on
synchronous_standby_names = 'FIRST 1 (postgres_slave)'
EOF

# Настройка pg_hba.conf через временный файл
TEMP_HBA=$(mktemp)
cat > "$TEMP_HBA" <<EOF
# TYPE  DATABASE        USER            ADDRESS                 METHOD
local   all             all                                     trust
host    all             all             127.0.0.1/32            md5
host    all             all             ::1/128                 md5
host    replication     replicator      0.0.0.0/0               md5
host    replication     readonly      0.0.0.0/0               md5
host    all             all             0.0.0.0/0               md5
EOF

# Копируем с сохранением прав
cp "$TEMP_HBA" /var/lib/postgresql/data/pg_hba.conf
chown postgres:postgres /var/lib/postgresql/data/pg_hba.conf
chmod 600 /var/lib/postgresql/data/pg_hba.conf
rm "$TEMP_HBA"


# Перезагружаем конфигурацию без перезапуска сервера
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" -c "SELECT pg_reload_conf();"
echo "Завершение инициализации"
