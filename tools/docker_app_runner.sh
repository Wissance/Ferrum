#!/bin/sh
# BE careful if you got error "./docker_app_runner.sh: cannot execute: required file not found"
# it 100% means that you MUST REPLACE ALL LINE ENDINGS \r\n -> \n
if [ ! -f "/app_data/inited.txt" ]; then
    /app/tools/${FERRUM_DATA_INIT_SCRIPT} || true
    touch /app_data/inited.txt
fi
/app/ferrum --config /app/config_docker_w_redis.json $FERRUM_ADDITIONAL_OPTS
