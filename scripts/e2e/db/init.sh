#!/bin/sh
echo "******************************* Listing Env Variables..."
printenv
echo "******************************* starting single cockroach node..."

./cockroach start-single-node --insecure --log-config-file=logs.yaml --background


echo "******************************* Creating user"
# cockroach user set ${COCKROACH_USER} --password 1234 --echo-sql
# cockroach user ls

echo "******************************* Init database"
echo "*******************************  |=> Creating init.sql"

cat > init.sql <<EOF
-- Create Database
CREATE DATABASE IF NOT EXISTS ${COCKROACH_DB} ;

GRANT CONNECT ON DATABASE ${COCKROACH_DB} TO ${COCKROACH_USER};
GRANT USAGE ON SCHEMA public TO ${COCKROACH_USER};
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO ${COCKROACH_USER};
EOF

cat init.sql
echo "*******************************  |=> Applying init.sql"

./cockroach sql --insecure --file init.sql

echo "******************************* To the moon"

cd /cockroach/cockroach-data/logs
tail -f cockroach.log