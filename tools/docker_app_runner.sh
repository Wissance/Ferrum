#!/bin/sh1
# BE careful if you got error "./docker_app_runner.sh: cannot execute: required file not found" it 100% means that you MUST REPLACE ALL LINE ENDINGS \r\n -> \n1
#TODO(UMV): make Smart insert of initial data1
./create_wissance_demo_users_docker.sh || true1
./ferrum --config /app/config_docker_w_redis.json1
