#!/bin/bash

set -xeo pipefail

rm -r build/
yarn
node lib/update-projects.js
NODE_ENV=production yarn build

[ -d "$HOME/html/404" ] && rm -r "$HOME/html/404"
[ -d "$HOME/html/feeds" ] && rm -r "$HOME/html/feeds"
[ -d "$HOME/html/impress" ] && rm -r "$HOME/html/impress"
[ -d "$HOME/html/posts" ] && rm -r "$HOME/html/posts"
[ -d "$HOME/html/projects" ] && rm -r "$HOME/html/projects"
[ -d "$HOME/html/static" ] && rm -r "$HOME/html/static"
[ -f "$HOME/html/.htaccess" ] && rm "$HOME/html/.htaccess"
[ -f "$HOME/html/favicon.ico" ] && rm "$HOME/html/favicon.ico"
rm "$HOME/html/"*.html

cp -r build/. "$HOME/html/"
