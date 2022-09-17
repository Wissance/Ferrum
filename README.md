## Ferrum

### General info

`Ferrum` is `OpenId-Connect` Authorization server written on GO. It behaves like `Keycloak` server
and could:

1. issue new tokens:
2. control user sessions;
3. get useinfo;
4. managed from external code (`Start` and `Stop`) making them an ***ideal candidate*** for using in ***integration
   tests*** for WEB API services that uses `Keycloak` as authorization server;
5. ability to use any user data and attributes (any valid JSON but with some requirements), if you have to
   properly configure your users just add what user have to `data.json` or in memory
6. ability to ***become serious enterprise level Authorization server (GIVE US 100+ STARS on Github
   and we'll make them enterprise)***.

it has `endpoints` SIMILAR to `Keycloak`, at present time we are having following:
1. `POST "/auth/realms/{realm}/protocol/openid-connect/token/"``
2. `GET "/auth/realms/{realm}/protocol/openid-connect/userinfo/"`

### How to use

### Configure user data as you wish

### Use from code