#!/bin/sh
set -e

scriptdir=$(dirname "$0")

# Apache gets grumpy about PID files pre-existing
rm -f /usr/local/apache2/logs/httpd.pid

# Enable SSL
apt-get update
apt-get install -y --no-install-recommends openssl
rm -rf /var/lib/apt/lists/*

openssl req \
	-x509 \
	-newkey rsa:2048 \
	-keyout "/usr/local/apache2/conf/server.key" \
	-out "/usr/local/apache2/conf/server.crt" \
	-days 365 \
	-nodes \
	-config "$scriptdir/openssl.conf"

sed -i \
	-e 's/^#\(Include .*httpd-ssl.conf\)/\1/' \
	-e 's/^#\(LoadModule .*mod_ssl.so\)/\1/' \
	-e 's/^#\(LoadModule .*mod_socache_shmcb.so\)/\1/' \
	conf/httpd.conf

# Enable .htaccess
cat >>conf/httpd.conf <<EOF
<Directory "/usr/local/apache2/htdocs">
    AllowOverride All
</Directory>
EOF

# Enable other required modules
sed -i \
	-e 's/^#\(LoadModule .*mod_rewrite.so\)/\1/' \
	-e 's/^#\(LoadModule .*mod_expires.so\)/\1/' \
	-e 's/^#\(LoadModule .*mod_deflate.so\)/\1/' \
	conf/httpd.conf

exec httpd -DFOREGROUND "$@"
