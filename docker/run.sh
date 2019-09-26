#!/bin/bash
docker run \
	-dit \
	--rm \
	--name website \
	-p 8080:80 \
	-p 8443:443 \
	-v "$PWD/build":/usr/local/apache2/htdocs/ \
	-v "$PWD/docker/setup":/setup/ \
	httpd \
	/setup/start.sh
