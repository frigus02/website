#!/bin/bash

set -eo pipefail

token=$(openssl rand -base64 32)

cat >~/cgi-bin/update-website <<__EOF__
#!/bin/bash

echo "Content-Type: text/plain"
echo ""

if [ "$QUERY_STRING" = "token=$token" ]; then
	echo "Deployment started"

	nohup ~/projects/website/deploy/trigger.sh "$USER" >~/projects/website-update.log &
else
	echo "Meh"
fi
__EOF__

chmod 755 ~/cgi-bin/update-website

echo "Created website deployment script"
echo "/cgi-bin/update-website?token=$token"
