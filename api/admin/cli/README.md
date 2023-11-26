go build -o ./api/admin/cli/ferrum-admin.exe ./api/admin/cli

./api/admin/cli/ferrum-admin.exe --resource=realm --operation=create --value='{"name": "myApp", "token_expiration": 600, "refresh_expiration": 300}'    

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
