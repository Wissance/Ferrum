#!/bin/sh

#TODO(UMV): make Smart insert of initial data
./create_wissance_demo_users_docker.sh
./ferrum --config /app/config_docker_w_redis.json
