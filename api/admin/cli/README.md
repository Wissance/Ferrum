### 1. Ferrum CLI Admin console

`CLI Admin` is an administrative console that allows to manage all `Ferrum` CLI has a following main peculiarities:

1. It is separate console executable utility
2. Shares same codebase
3. Use `Ferrum` config file (provides as an argument)

Admin CLI could be build as follows:

```ps1
go build -o ferrum-admin.exe ./api/admin/cli
```

### 2. Ferrum CLI Admin operations

All Admin CLI operation have the same scheme as follows:
`{admin_cli_executable} --resource={resorce_name} --operation={operation_type} [additional_arguments]`
where:
* `{admin_cli_executable}` is a name of executable file
* `{resource_name}` - `realm`, `client` or `user`
* `{operation_type}` is an operation to perform over resource (see operation description below)
* `[additional_arguments]` a set of additional `--key=value` pairs i.e. resource id (for get), or value (for create and|or update)

#### 2.1 Operations

`CLI` allows to perform standard CRUD operation via console (`create`, `read`, `update`, `delete`) and some additional
operations:

* reset_password - reset password to random
* change_password - changes password to provided

##### 2.1.1 Standard CRUD operations

**Create operation** on realm looks as follows:
```ps1
.ferrum-admin.exe --resource=realm --operation=create --value='{"name": "WissanceFerrumDemo", "token_expiration": 600, "refresh_expiration": 300}'
```
value is using for providing data to operations, above example creates new `Realm` with name `WissanceFerrumDemo`

But we know that `Realm` also has `Clients` and `Users`, therefore we should provide Realm id = name to newly creating Clients and Users via `--params={realm_name}`:

Creation of new client looks as follows:
```ps1
./ferrum-admin.exe --resource=client --operation=create --value='{"id": "d4dc483d-7d0d-4d2e-a0a0-2d34b55e6666", "name": "WissanceWebDemo", "type": "confidential", "auth": {"type": 1, "value": "fb6Z4RsOadVycQoeQiN57xpu8w8wTEST"}}' --params=WissanceFerrumDemo
```

And User (like in create in `--params` realm name should be passed):
```ps1
./ferrum-admin.exe --resource=user --operation=create --value='{"info": {"sub": "667ff6a7-3f6b-449b-a217-6fc5d9actest", "email_verified": true, "roles": ["admin"], "name": "M.V.Ushakov", "preferred_username": "umv", "given_name": "Michael", "family_name": "Ushakov"}, "credentials": {"password": "1s2d3f4g90xs"}}' --params=WissanceFerrumDemo
```

**Update operation** requires `--resource_id` parameter = name of resource. I.e. for a client it looks as follows:
```ps1
./ferrum-admin.exe --resource=client --operation=update --resource_id=WissanceWebDemo --value='{"id": "d4dc483d-7d0d-4d2e-a0a0-2d34b55e6666", "name": "WissanceWebDemo", "type": "confidential", "auth": {"type": 2, "value": "fb6Z4RsOadVycQoeQiN57xpu8w8wTEST"}}' --params=WissanceFerrumDemo
```

**Update operation** over realm:
```ps1
./ferrum-admin.exe --resource=realm --operation=update --resource_id=WissanceFerrumDemo --value='{"name": "WissanceFerrumDemo", "token_expiration": 2400, "refresh_expiration": 1200}'
```

##### 2.1.2 Additional operations

```
go build -o ./api/admin/cli/ferrum-admin.exe ./api/admin/cli

client:
./api/admin/cli/ferrum-admin.exe --resource=client --operation=create --value='{"id": "d4dc483d-7d0d-4d2e-a0a0-2d34b55e6666", "name": "clientFromCreate", "type": "confidential", "auth": {"type": 1, "value": "fb6Z4RsOadVycQoeQiN57xpu8w8wTEST"}}' --params=myApp

./api/admin/cli/ferrum-admin.exe --resource=client --operation=get --resource_id=clientFromCreate --params=myApp

./api/admin/cli/ferrum-admin.exe --resource=client --operation=update --resource_id=clientFromCreate --value='{"id": "d4dc483d-7d0d-4d2e-a0a0-2d34b55e6666", "name": "clientFromCreate", "type": "confidential", "auth": {"type": 2, "value": "fb6Z4RsOadVycQoeQiN57xpu8w8wTEST"}}' --params=myApp

./api/admin/cli/ferrum-admin.exe --resource=client --operation=update --resource_id=clientFromCreate --value='{"id": "d4dc483d-7d0d-4d2e-a0a0-2d34b55e6666", "name": "RenameClientFromCreate", "type": "confidential", "auth": {"type": 1, "value": "fb6Z4RsOadVycQoeQiN57xpu8w8wTEST"}}' --params=myApp

./api/admin/cli/ferrum-admin.exe --resource=client --operation=delete --resource_id=RenameClientFromCreate --params=myApp


user:
./api/admin/cli/ferrum-admin.exe --resource=user --operation=create --value='{"info": {"sub": "667ff6a7-3f6b-449b-a217-6fc5d9actest", "email_verified": false, "roles": ["admin"], "name": "firstTestName lastTestName", "preferred_username": "userFromCreate", "given_name": "firstTestName", "family_name": "lastTestName"}, "credentials": {"password": "1s2d3f4g90xs"}}' --params=myApp

./api/admin/cli/ferrum-admin.exe --resource=user --operation=get --resource_id=userFromCreate --params=myApp

./api/admin/cli/ferrum-admin.exe --resource=user --operation=update --resource_id=userFromCreate --value='{"info": {"sub": "667ff6a7-3f6b-449b-a217-6fc5d9actest", "email_verified": true, "roles": ["admin"], "name": "firstTestName lastTestName", "preferred_username": "userFromCreate", "given_name": "firstTestName", "family_name": "lastTestName"}, "credentials": {"password": "1s2d3f4g90xs"}}' --params=myApp

./api/admin/cli/ferrum-admin.exe --resource=user --operation=update --resource_id=userFromCreate --value='{"info": {"sub": "667ff6a7-3f6b-449b-a217-6fc5d9actest", "email_verified": false, "roles": ["admin"], "name": "firstTestName lastTestName", "preferred_username": "RenameUserFromCreate", "given_name": "firstTestName", "family_name": "lastTestName"}, "credentials": {"password": "1s2d3f4g90xs"}}' --params=myApp

./api/admin/cli/ferrum-admin.exe --resource=user --operation=reset_password --resource_id=RenameUserFromCreate --params=myApp

./api/admin/cli/ferrum-admin.exe --resource=user --operation=change_password --resource_id=RenameUserFromCreate --value='newPassword' --params=myApp

./api/admin/cli/ferrum-admin.exe --resource=user --operation=delete --resource_id=RenameUserFromCreate --params=myApp


realm:
./api/admin/cli/ferrum-admin.exe --resource=realm --operation=create --value='{"name": "testRealm", "token_expiration": 600, "refresh_expiration": 300, "clients": [{"id": "d4dc483d-7d0d-4d2e-a0a0-2d34b55e1111", "name": "clientTestOne", "type": "confidential", "auth": {"type": 1, "value": "fb6Z4RsOadVycQoeQiN57xpu8w8wTEST"}}, {"id": "d4dc483d-7d0d-4d2e-a0a0-2d34b55e2222", "name": "clientTestTwo", "type": "confidential", "auth": {"type": 1, "value": "fb6Z4RsOadVycQoeQiN57xpu8w8wTEST"}}], "users": [{"info": {"sub": "667ff6a7-3f6b-449b-a217-111111actest", "email_verified": false, "roles": ["admin"], "name": "firstTestName lastTestName", "preferred_username": "userTestOne", "given_name": "firstTestName", "family_name": "lastTestName"}, "credentials": {"password": "1s2d3f4g90xs"}}, {"info": {"sub": "667ff6a7-3f6b-449b-a217-222222actest", "email_verified": false, "roles": ["admin"], "name": "firstTestName lastTestName", "preferred_username": "userTestTwo", "given_name": "firstTestName", "family_name": "lastTestName"}, "credentials": {"password": "1s2d3f4g90xs"}}]}'

./api/admin/cli/ferrum-admin.exe --resource=realm --operation=get --resource_id=testRealm

./api/admin/cli/ferrum-admin.exe --resource=realm --operation=update --resource_id=testRealm --value='{"name": "testRealm", "token_expiration": 700, "refresh_expiration": 300}'

./api/admin/cli/ferrum-admin.exe --resource=realm --operation=update --resource_id=testRealm --value='{"name": "RenameTestRealm", "token_expiration": 600, "refresh_expiration": 300}'

./api/admin/cli/ferrum-admin.exe --resource=realm --operation=delete --resource_id=RenameTestRealm

```
