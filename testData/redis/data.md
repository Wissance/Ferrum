This is text contains set of small JSONs (sample data) 4 testing the app.
We are using `ferrum_1` as a Redis namespace (prefix before every key) 

1. Realms
```json
{
    "ferrum_1.realm_myApp": {
         "name": "myApp",
         "token_expiration": 600,
         "refresh_expiration": 300,
         "clients": [
         ]
    },
    "ferrum_1.realm_testApp": {
         "name": "testapp",
         "token_expiration": 6000,
         "refresh_expiration": 3000,
         "clients": [
         ]
    }
}
```
`myApp` and `testApp` are ACTUAL Realms names, but they must be stored in `Redis` with keys `ferrum_1.realm_myApp` and `ferrum_1.testApp` respectively, we left client 
blank, clients to realm relation is setting in a separate object

2. Clients
All the clients should have the following key pattern `namespace.{realmName}_client_{clientName}`
```json
{
    "ferrum_1.myApp_client_test-service-app-client": {
        "id": "d4dc483d-7d0d-4d2e-a0a0-2d34b55e5a14",
        "name": "test-service-app-client",
        "type": "confidential",
        "auth": {
            "type": 1,
            "value": "fb6Z4RsOadVycQoeQiN57xpu8w8wplYz"
        }
    },
    "ferrum_1.myApp_client_test-mobile-app-client": {
        "id": "d4dc483d-7d0d-4d2e-a0a0-2d34b55e5199",
        "name": "test-mobile-app-client",
        "type": "confidential",
        "auth": {
            "type": 1,
            "value": "fb6Z4RsOadVycQoeQiN57xpu8w8wplYz"
        }
    },
    "ferrum_1.testApp_client_test-test-app-client": {
        "id": "d4dc483d-7d0d-4d2e-a0a0-2d34b55e5207",
        "name": "test-test-app-client",
        "type": "confidential",
        "auth": {
            "type": 1,
            "value": "fb6Z4RsOadVycQoeQiN57xpu8w8wplYz"
        }
    }
}
```
3. Clients to Realm binding i.e. consider realm `myApp` we should add all clients identifiers+names as array to object with key `ferrum_1.realm_myApp_clients`
```json
{
    "ferrum_1.realm_myApp_clients": [
        {
            "id": "d4dc483d-7d0d-4d2e-a0a0-2d34b55e5a14",
            "name": "test-service-app-client"
        },
        {
            "id": "d4dc483d-7d0d-4d2e-a0a0-2d34b55e5199",
            "name": "test-mobile-app-client"
        }
    ],

    "ferrum_1.realm_testApp_clients": [
        {
            "id": "d4dc483d-7d0d-4d2e-a0a0-2d34b55e5207",
            "name": "test-test-app-client"
        }
    ]
}
```
4. Users itself stores in redis by key with a pattern - `{namespace}.{realmName}_user_{userName}`, if we are having user with `admin` userName. Users could have
different structure, but must meet some common user requirements:
   * must have `info` object on a `JSON` top level with field `preferred_username`  && object `credentials` on the same level as `info`, if user
     authenticates with a password than user has to have a password field with a value. 
```json
{
    "ferrum_1.myApp_user_vano": {
        "info": {
            "sub": "667ff6a7-3f6b-449b-a217-6fc5d9ac0723",
            "email_verified": false,
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
    },
    "ferrum_1.myApp_user_admin": {
        "info": {
            "sub": "667ff6a7-3f6b-449b-a217-6fc5d9ac0724",
            "email_verified": false,
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
}
```
5. Users to Realms binding in Redis object by the following key pattern -  `{namespace}.realm_{realmName}_users`
```json
{
    "ferrum_1.realm_myApp_users": [
        {
            "id": "667ff6a7-3f6b-449b-a217-6fc5d9ac0723",
            "name": "vano"
        },
        {
            "id": "667ff6a7-3f6b-449b-a217-6fc5d9ac0724",
            "name": "admin"
        }
    ]
}
```
