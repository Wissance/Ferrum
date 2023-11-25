import json

import redis
from redis.commands.json.path import Path

redis_host = 'redis'
redis_port = 6379
db = 0
username = 'test_user'
password = 'test_password'
client = redis.Redis(host=redis_host, port=redis_port, db=db, username=username, password=password)
try:
    response = client.ping()
except redis.ConnectionError:
    print('Bad connect to redis, host - "{}".'.format(redis_host))

isExistsRealm = client.exists("ferrum_1.realm_myApp")
if isExistsRealm:
    print('The radis has "ferrum_1.realm_myApp". Data not inserted during initialization.')
    exit()


realm_myApp = {
         "name": "myApp",
         "token_expiration": 600,
         "refresh_expiration": 300,
         "clients": [
         ]
}
realm_testApp = {
         "name": "testapp",
         "token_expiration": 6000,
         "refresh_expiration": 3000,
         "clients": [
         ]
}
realm_myApp = json.dumps(realm_myApp)
realm_testApp = json.dumps(realm_testApp)
client.set('ferrum_1.realm_myApp', realm_myApp)
client.set('ferrum_1.realm_testApp', realm_testApp)

client_service = {
        "id": "d4dc483d-7d0d-4d2e-a0a0-2d34b55e5a14",
        "name": "test-service-app-client",
        "type": "confidential",
        "auth": {
            "type": 1,
            "value": "fb6Z4RsOadVycQoeQiN57xpu8w8wplYz"
        }
}
client_mobile = {
        "id": "d4dc483d-7d0d-4d2e-a0a0-2d34b55e5199",
        "name": "test-mobile-app-client",
        "type": "confidential",
        "auth": {
            "type": 1,
            "value": "fb6Z4RsOadVycQoeQiN57xpu8w8wplYz"
        }
}
client_test = {
        "id": "d4dc483d-7d0d-4d2e-a0a0-2d34b55e5207",
        "name": "test-test-app-client",
        "type": "confidential",
        "auth": {
            "type": 1,
            "value": "fb6Z4RsOadVycQoeQiN57xpu8w8wplYz"
        }
}
client_service = json.dumps(client_service)
client_mobile = json.dumps(client_mobile)
client_test = json.dumps(client_test)
client.set('ferrum_1.myApp_client_test-service-app-client', client_service)
client.set('ferrum_1.myApp_client_test-mobile-app-client', client_mobile)
client.set('ferrum_1.testApp_client_test-test-app-client', client_test)

realm_myApp_clients = [
    {
        "id": "d4dc483d-7d0d-4d2e-a0a0-2d34b55e5a14",
        "name": "test-service-app-client"
    },
    {
        "id": "d4dc483d-7d0d-4d2e-a0a0-2d34b55e5199",
        "name": "test-mobile-app-client"
    }
]
realm_testApp_clients = [
    {
        "id": "d4dc483d-7d0d-4d2e-a0a0-2d34b55e5207",
        "name": "test-test-app-client"
    }
]
realm_myApp_clients = [json.dumps(realm_myApp_clients)]
realm_testApp_clients = [json.dumps(realm_testApp_clients)]
client.rpush('ferrum_1.realm_myApp_clients', *realm_myApp_clients)
client.rpush('ferrum_1.realm_testApp_clients', *realm_testApp_clients)

user_vano = {
        "info": {
            "sub": "667ff6a7-3f6b-449b-a217-6fc5d9ac0723",
            "email_verified": False,
            "roles": [
                "admin"
            ],
            "name": "admin sys",
            "preferred_username": "admin",
            "given_name": "admin",
            "family_name": "sys"
        },
        "credentials": {
            "password": "1s2d3f4g90xs"
        }
}
user_admin = {
        "info": {
            "sub": "667ff6a7-3f6b-449b-a217-6fc5d9ac0724",
            "email_verified": False,
            "roles": [
                "admin"
            ],
            "name": "admin sys",
            "preferred_username": "admin",
            "given_name": "admin",
            "family_name": "sys"
        },
        "credentials": {
            "password": "1s2d3f4g90xs"
        }
}
user_vano = json.dumps(user_vano)
user_admin = json.dumps(user_admin)
client.set('ferrum_1.myApp_user_vano', user_vano)
client.set('ferrum_1.myApp_user_admin', user_admin)

realm_myApp_users = [
    {
        "id": "667ff6a7-3f6b-449b-a217-6fc5d9ac0723",
        "name": "vano"
    },
    {
        "id": "667ff6a7-3f6b-449b-a217-6fc5d9ac0724",
        "name": "admin"
    }
]
realm_myApp_users = [json.dumps(realm_myApp_users)]
client.rpush('ferrum_1.realm_myApp_users', *realm_myApp_users)

print('Data is inserted into the radis during initialization.')


