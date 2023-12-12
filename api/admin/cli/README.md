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

* `reset_password` - reset password to random value
* `change_password` - changes password to provided

##### 2.1.1 Standard CRUD operations

##### 2.1.1.1 Create operations

Create operation should provide `--value` with resource body, key will be constructed from body. For `client` and `user` creation realm id (name) 
must be provided via `--params`.

Create `realm` example
```ps1
.ferrum-admin.exe --resource=realm --operation=create --value='{"name": "WissanceFerrumDemo", "token_expiration": 600, "refresh_expiration": 300}'
```

Create `client` example:
```ps1
./ferrum-admin.exe --resource=client --operation=create --value='{"id": "d4dc483d-7d0d-4d2e-a0a0-2d34b55e6666", "name": "WissanceWebDemo", "type": "confidential", "auth": {"type": 1, "value": "fb6Z4RsOadVycQoeQiN57xpu8w8wTEST"}}' --params=WissanceFerrumDemo
```

Create `user` example:
```ps1
./ferrum-admin.exe --resource=user --operation=create --value='{"info": {"sub": "667ff6a7-3f6b-449b-a217-6fc5d9actest", "email_verified": true, "roles": ["admin"], "name": "M.V.Ushakov", "preferred_username": "umv", "given_name": "Michael", "family_name": "Ushakov"}, "credentials": {"password": "1s2d3f4g90xs"}}' --params=WissanceFerrumDemo
```
##### 2.1.1.2 Update operations

Update operation fully replace item by key `--resource_id` + `--param={realm_name}` (realm does not requires)
New key content provides via `--value=`. Why we don't provide just a DB key? Answer is there are could be different storage 
and key is often composite, therefore it is more user-friendly to provide separately key and realn

Update `realm` example
```ps1
./ferrum-admin.exe --resource=realm --operation=update --resource_id=WissanceFerrumDemo --value='{"name": "WissanceFerrumDemo", "token_expiration": 2400, "refresh_expiration": 1200}'
```

Update `client` example:
```ps1
./ferrum-admin.exe --resource=client --operation=update --resource_id=WissanceWebDemo --value='{"id": "d4dc483d-7d0d-4d2e-a0a0-2d34b55e6666", "name": "WissanceWebDemo", "type": "confidential", "auth": {"type": 2, "value": "fb6Z4RsOadVycQoeQiN57xpu8w8wTEST"}}' --params=WissanceFerrumDemo
```

Update `user` example:
```ps1
./ferrum-admin.exe --resource=user --operation=update --resource_id=umv --value='{"info": {"sub": "667ff6a7-3f6b-449b-a217-6fc5d9actest", "email_verified": true, "roles": ["admin", "managers"], "name": "M.V.Ushakov", "preferred_username": "umv", "given_name": "Michael", "family_name": "Ushakov"}, "credentials": {"password": "1s2d3f4g90xs"}}' --params=WissanceFerrumDemo
```

Question:
1. What is using for user identification, because it has `preferred_username`, and `given_name` fields. I've not tested this yet but `preferred_username` must be used as `resource_id`. Here and in all `CRUD` operations that are requires identifier. 

##### 2.1.1.3 Get operations

**Get by id operation** requires resource identifier (`resource_id`) and realm name via `--params`.

Get `realm` example:
```ps1
./ferrum-admin.exe --resource=realm --operation=get --resource_id=WissanceFerrumDemo
```

Get `client` example:
```ps1
./ferrum-admin.exe --resource=client --operation=get --resource_id=WissanceWebDemo --params=WissanceFerrumDemo
```

Get `user` example:
```ps1
./ferrum-admin.exe --resource=user --operation=get --resource_id=userFromCreate --params=WissanceFerrumDemo
```
Get user should hide credential section (have to test, not tested yet).

##### 2.1.1.3 Delete operations

Delete operation requires `--resource_id` and `--params` to be provided.

Delete `realm` example:
```ps1
./ferrum-admin.exe --resource=realm --operation=delete --resource_id=WissanceFerrumDemo
```

Delete `client` example:
```ps1
./ferrum-admin.exe --resource=client --operation=delete --resource_id=WissanceWebDemo --params=WissanceFerrumDemo
```

Delete `user` example:
```ps1
./ferrum-admin.exe --resource=user --operation=delete --resource_id=umv --params=WissanceFerrumDemo
```

Questions (todo for work):
1. What happened to clients and users if realm was deleted ? Should be a CASCADE removing.

##### 2.1.2 Additional operations

###### 2.1.2.1 User password reset

Password reset makes set `user` password value to random, new password outputs to console. As for get, update or delete
operation it requires username to be provided via `--resourec_id` and a realm name via `--params`, example:
```ps1
./ferrum-admin.exe --resource=user --operation=reset_password --resource_id=umv --params=WissanceFerrumDemo
```

###### 2.1.2.1 User password change

Password change requires username to be provided via `--resourec_id` and a realm name via `--params. New password
is passing via `--value=`, example:

```ps1
./ferrum-admin.exe --resource=user --operation=change_password --resource_id=umv --value='newPassword' --params=WissanceFerrumDemo
```
