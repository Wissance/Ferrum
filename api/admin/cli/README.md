go build -o ./api/admin/cli/ferrum-admin.exe ./api/admin/cli


client:    
./api/admin/cli/ferrum-admin.exe --resource=client --operation=get --resource_id=test-service-app-client   
./api/admin/cli/ferrum-admin.exe --resource=client --operation=get --resource_id=test-service-app-client --params=myApp   

./api/admin/cli/ferrum-admin.exe --resource=client --operation=delete --resource_id=test-service-app-client   
./api/admin/cli/ferrum-admin.exe --resource=client --operation=delete --resource_id=test-service-app-client --params=myApp   

./api/admin/cli/ferrum-admin.exe --resource=client --operation=create --value='{"id": "d4dc483d-7d0d-4d2e-a0a0-2d34b55e6666", "name": "training", "type": "confidential", "auth": {"type": 1, "value": "fb6Z4RsOadVycQoeQiN57xpu8w8wTEST"}}'     
./api/admin/cli/ferrum-admin.exe --resource=client --operation=create --params=myApp   

./api/admin/ferrum-admin.exe --resource=client --operation=update --resource_id=training --value='{"id": "d4dc483d-7d0d-4d2e-a0a0-2d34b55e6666", "name": "training", "type": "confidential", "auth": {"type": 1, "value": "fb6Z4RsOadVycQoeQiN57xpu8w8wTEST"}}'     


user:   
./api/admin/cli/ferrum-admin.exe --resource=user --operation=get --resource_id=admin   
./api/admin/cli/ferrum-admin.exe --resource=user --operation=get --resource_id=admin --params=myApp   

./api/admin/cli/ferrum-admin.exe --resource=user --operation=delete --resource_id=admin
./api/admin/cli/ferrum-admin.exe --resource=user --operation=delete --resource_id=admin --params=myApp   

./api/admin/cli/ferrum-admin.exe --resource=user --operation=create --value='{"info": {"sub": "667ff6a7-3f6b-449b-a217-6fc5d9actest", "email_verified": false, "roles": ["admin"], "name": "firstTestName lastTestName", "preferred_username": "testuser", "given_name": "firstTestName", "family_name": "lastTestName"}, "credentials": {"password": "1s2d3f4g90xs"}}' 
./api/admin/cli/ferrum-admin.exe --resource=user --operation=create --params=myApp    

./api/admin/cli/ferrum-admin.exe --resource=user --operation=update --resource_id=testuser --value='{"info": {"sub": "667ff6a7-3f6b-449b-a217-6fc5d9actest", "email_verified": false, "roles": ["admin"], "name": "firstTestName lastTestName", "preferred_username": "testuser", "given_name": "firstTestName", "family_name": "lastTestName"}, "credentials": {"password": "1s2d3f4g90xs"}}' 

realm:   
./api/admin/cli/ferrum-admin.exe --resource=realm --operation=get --resource_id=myApp   

./api/admin/cli/ferrum-admin.exe --resource=realm --operation=delete --resource_id=myApp   

./api/admin/cli/ferrum-admin.exe --resource=realm --operation=create --value='{"name": "testRalm", "token_expiration": 600, "refresh_expiration": 300, "clients": [{"id": "d4dc483d-7d0d-4d2e-a0a0-2d34b55e1111", "name": "trainingFirst", "type": "confidential", "auth": {"type": 1, "value": "fb6Z4RsOadVycQoeQiN57xpu8w8wTEST"}}, {"id": "d4dc483d-7d0d-4d2e-a0a0-2d34b55e2222", "name": "trainingSecond", "type": "confidential", "auth": {"type": 1, "value": "fb6Z4RsOadVycQoeQiN57xpu8w8wTEST"}}], "users": [{"info": {"sub": "667ff6a7-3f6b-449b-a217-111111actest", "email_verified": false, "roles": ["admin"], "name": "firstTestName lastTestName", "preferred_username": "testFirst", "given_name": "firstTestName", "family_name": "lastTestName"}, "credentials": {"password": "1s2d3f4g90xs"}}, {"info": {"sub": "667ff6a7-3f6b-449b-a217-222222actest", "email_verified": false, "roles": ["admin"], "name": "firstTestName lastTestName", "preferred_username": "testSecond", "given_name": "firstTestName", "family_name": "lastTestName"}, "credentials": {"password": "1s2d3f4g90xs"}}]}'

./api/admin/cli/ferrum-admin.exe --resource=realm --operation=update --resource_id=testRalm --value='{"name": "testRalm", "token_expiration": 600, "refresh_expiration": 300}'








./api/admin/ferrum-admin.exe --params=testRalm --resource=user --operation=get --resource_id=testFirst   
./api/admin/ferrum-admin.exe --params=testRalm --resource=user --operation=get --resource_id=testSecond   

./api/admin/ferrum-admin.exe --operation=create --resource=realm --resource_id=IRP --value='{"name":"irp_rt"}'   
./api/admin/ferrum-admin.exe --operation=create --resource=realm --resource_id=IRP --value='{"name":"irp_rt"}'   