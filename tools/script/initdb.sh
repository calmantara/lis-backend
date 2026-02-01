#!/bin/sh

set -e

echo "Creating application database $MASTER_DATABASE_NAME"
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
	CREATE DATABASE "$MASTER_DATABASE_NAME";
	GRANT ALL PRIVILEGES ON DATABASE  "$MASTER_DATABASE_NAME" TO "$POSTGRES_USER";
EOSQL

#  add extension 
echo "shared_preload_libraries='pg_cron,pg_partman_bgw'" >> $PGDATA/postgresql.conf
echo "cron.database_name='$MASTER_DATABASE_NAME'" >> $PGDATA/postgresql.conf
echo "pg_partman_bgw.interval = 3600" >> $PGDATA/postgresql.conf
echo "pg_partman_bgw.role = '$POSTGRES_USER'" >> $PGDATA/postgresql.conf
echo "pg_partman_bgw.dbname = '$MASTER_DATABASE_NAME'" >> $PGDATA/postgresql.conf

# change db max connection
echo "max_connections = 1000" >> $PGDATA/postgresql.conf

echo "shared_buffers = 100MB" >> $PGDATA/postgresql.conf
# restart postgres
echo "reload config"
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$MASTER_DATABASE_NAME" <<-EOSQL
    SELECT pg_reload_conf();
EOSQL
su -s /bin/sh postgres <<EOF
    pg_ctl restart
EOF
poweroff
