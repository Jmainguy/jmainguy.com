#!/bin/bash
docker run -p 8005:80 -d --restart always \
    -v /opt/jmainguy.com/public:/usr/share/nginx/html \
    --name jmainguy.com nginx
