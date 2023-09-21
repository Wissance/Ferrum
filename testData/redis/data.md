This is text contains set of small JSONs (sample data) 4 testing the app.
We are using `ferrum_1` as a Redis namespace (prefix before every key) 

1. Realms, without clients:
```json
{
    "ferrum_1_myApp": {
         "name": "myapp",
         "token_expiration": 600,
         "refresh_expiration": 300,
         "clients": [
         ]
    },
    "ferrum_1_testApp": {
         "name": "testApp",
         "token_expiration": 6000,
         "refresh_expiration": 3000,
         "clients": [
         ]
    }
}
```
`myApp` and `testApp` are ACTUAL Realms names, but they must be stored in `Redis` with keys `ferrum_1_myApp` and `ferrum_1_testApp` respectively, we left client 
blank, clients to realm relation is setting in a separate onject

2. Clients
```json
{

}
```
3. Clients to Realm binding

4. Users

5. Users to Realms binding

{
    "realms": [
        {
            "name": "myapp",
            "token_expiration": 330,
            "refresh_expiration": 200,
            "clients": [
                {
                    "id": "d4dc483d-7d0d-4d2e-a0a0-2d34b55e5a14",
                    "name": "test-service-app-client",
                    "type": "confidential",
                    "auth": {
                        "type": 1,
                        "value": "fb6Z4RsOadVycQoeQiN57xpu8w8wplYz"
                    }
                }
            ],
            "users": [
                {
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
                {
                    "info": {
                        "sub": "8be91328-0f85-408f-966a-fd9a04ce94d9",
                        "email_verified": false,
                        "roles": [
                            "1stfloor",
                            "manager"
                        ],
                        "name": "ivan ivanov",
                        "preferred_username": "vano",
                        "given_name": "ivan",
                        "family_name": "ivanov"
                    },
                    "credentials": {
                        "password": "qwerty_user"
                    }
                }
            ]
        }
    ]
}
