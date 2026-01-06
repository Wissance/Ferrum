import os
import csv
import argparse
import datetime
import typing

default_config = "config_docker_w_redis.json"
default_client_secret = "00000000000000000000000000000000"

realms_n_arg = "--realms"
user_per_realm_arg = "--realm_users"
output_dir_arg = "--output"

parser=argparse.ArgumentParser()
parser.add_argument(realms_n_arg, help="Number of realms", required=True, type=int)
parser.add_argument(user_per_realm_arg, help="Number of users per realm", required=True, type=int)
parser.add_argument(output_dir_arg, help="Generating file", required=True, type=str)

# example - ./gen_users_for_perm_tests --realms=100 --realm_users=100 --output=../tools/r100_u100_demo_data.sh
def main():
    print ("********* Ferrum perf test data generation started ********")
    try:
        args = parser.parse_args()
        output_file = args.output
        realms_number = args.realms
        users_per_realm_number = args.realm_users
        with open(output_file, 'w+', encoding='utf-8') as f :
            for r in range(1, realms_number+1):
                realm, realm_cmd = make_create_realm_cmd(r, default_config)
                f.write(realm_cmd)
                f.write("\n")
                #todo(umv) :create client
                client, client_cmd = make_create_client_cmd(realm, r, default_config, default_client_secret)
                f.write(client_cmd)
                f.write("\n")
                #todo(umv) :create users

        print ("********* Ferrum perf test data generation successfully finished ********")
    except Exception as e:
        print (str.format("An error occurred during data generation, error: {0}", str(e)))
        print ("********* Ferrum perf test data generation failed ********")

def make_create_realm_cmd(realm_number:int, app_config:str) -> tuple[str, str]:
    realm_create_template = "./ferrum-admin --config={0} --resource=realm --operation=create --value=--value=\'{{\"name\": \"{1}\", \"user_federation_services\":[], \"token_expiration\": 600, \"refresh_expiration\": 300 }}\'"
    realm = str.format("realm_{0}", realm_number)
    realm_create_cmd = realm_create_template.format(app_config, realm)
    return realm, realm_create_cmd

def make_create_client_cmd(realm:str, client_number:int, app_config:str, client_secret:str) -> tuple[str, str]:
    client_create_template = "./ferrum-admin --config={0} --resource=client --operation=create --value=\'{{\"id\": \"{1}\", \"name\": \"{2}\", \"type\": \"confidential\", \"auth\": {{\"type\": 1, \"value\": \"{3}\"}} }}' --params=\"{4}\""
    client_id = str.format("00000000-0000-0000-0000-{0:012d}", client_number)
    client = str.format("client_{0}", client_number)
    client_create_cmd =  client_create_template.format(app_config, client_id, client, client_secret, realm)
    return client, client_create_cmd

if __name__ == "__main__":
    main()