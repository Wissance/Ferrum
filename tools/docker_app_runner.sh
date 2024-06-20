#!/bin/sh
# BE careful if you got error "./docker_app_runner.sh: cannot execute: required file not found" it 100% means that you MUST REPLACE ALL LINE ENDINGS \r\n -> \n
#TODO(UMV): make Smart insert of initial data
./create_wissance_demo_users_docker.sh || true
./ferrum --config /app/config_docker_w_redis.json
