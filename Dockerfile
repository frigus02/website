FROM node:15 AS builder

WORKDIR /opt/app
COPY package.json yarn.lock ./
RUN yarn install --frozen-lockfile
COPY . ./
RUN node lib/update-projects.js && \
    NODE_ENV=production yarn build

FROM nginx:1.19

COPY deploy/docker/default.conf /etc/nginx/conf.d/default.conf
COPY deploy/docker/common-headers.conf /etc/nginx/snippets/common-headers.conf
COPY --from=builder /opt/app/build /usr/share/nginx/html
