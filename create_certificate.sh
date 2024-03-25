#!/bin/bash

set -e

echo "Generating certificate..."

openssl req -x509 \
	-newkey rsa:4096 \
	-nodes \
	-keyout /var/lib/postgresql/data/private.pem \
	-out /var/lib/postgresql/data/cert.pem \
	-subj "/C=BR/ST=SP/L=SJBV/O=Rinha/CN=www.rinhabackend2024.com" \
	-days 365

chmod 600 /var/lib/postgresql/data/private.pem
chown postgres /var/lib/postgresql/data/private.pem

echo "Enabling SSL..."

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<EOF

ALTER SYSTEM SET ssl = 'on';
ALTER SYSTEM SET ssl_cert_file = '/var/lib/postgresql/data/cert.pem';
ALTER SYSTEM SET ssl_key_file = '/var/lib/postgresql/data/private.pem';

SELECT pg_reload_conf();

EOF
