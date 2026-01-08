import http from 'k6/http';
import { check } from "k6";
export let options = {
    stages: [
        // Ramp-up from 1 to 10 VUs in 15s
        { duration: "15s", target: 10 },
        // Stay at rest on 10 VUs for 30s
        { duration: "30s", target: 10 },
        // Ramp-down from 10 to 0 VUs for 15s
        { duration: "15s", target: 0 }
    ],
    thresholds: {
        http_req_duration: ['p(95)<200'], // 95% of requests must complete in less than 500ms for the test to pass
        http_req_failed: ['rate<0.01'],    // Test fails if more than 1% of requests fail
    },
};

/* This function is entrypoint to short Smoke testing over Ferrum running in performance test mode
 * Purposes of Smoke tests:
 * 1. Run parallel working of minimal amount of clients (up to 10) to see whether exists some
 *    issues with parallel call or not
 * 2. Test all methods that are going to be used in average_load and stress testing
 * Questions:
 * 1. We are going to test multiple variants with distribution users on realms, but now we have
 *    only r100_u100_demo.data.sh (100 realms, each realm has 100 users). How we are going to pass here
 *    what variant we are actually using ? By env vars ?
 * 2. How to pass Host && Port for call Ferrum WebAPI (seems 127.0.0.1) is not suitable here
 *
 * */
export default function () {
    const iterationNum = 10
    const userPassword = "P@55W0rD"
    const clientSecret = "00000000000000000000000000000000"
    let ferrumBaseUrl = "http://10.50.40.3:8182";
    // 1. Get random realm (1-100) && user (1-100)
    let randomRealm = getRandomInt(1, 100)
    let realm = "realm_" + randomRealm
    let randomUserRel = getRandomInt(1, 100)
    let absUser = randomRealm * 100 + randomUserRel
    let user = "u" + absUser
    let clientId = "client_" + randomRealm
    // 2. Get Access Token
    let getTokenResponse = getAccessToken(ferrumBaseUrl, realm, clientId, clientSecret, user, userPassword)
    // check status
    check(getTokenResponse, {
        'status is 200': (r) => r.status === 200,
    });
    const responseBody = JSON.parse(getTokenResponse.body);
    let accessToken = responseBody.access_token
    // 3. Iterations over userInfo
    for (let i = 0; i < iterationNum; i++) {
        let pause = getRandomInt(10, 40)
        sleep(pause)
    }

};

/* This function get access token from Ferrum
 * baseUrl is part protocol://host:port
 * */
function getAccessToken(baseUrl, realm, clientId, secret, username, password) {
    const url = baseUrl+"/auth/realms/" + realm + "/protocol/openid-connect/token"; // A service that echoes the request

    var payload = {
        "client_id": clientId,
        "client_secret" : secret,
        "username": username,
        "password": password,
        "grant_type": "password",
        "scope": "profile"
    };

    // k6 automatically handles the encoding and sets the header,
    // but you can explicitly set it for clarity if needed.
    const params = {
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
        },
    };

    return http.post(url, payload, params);
}

/*
 *
 * */
function getUserInfo(baseUrl, realm, accessToken) {

}

function getRandomInt(min, max) {
    min = Math.ceil(min);
    max = Math.floor(max);
    return Math.floor(Math.random() * (max - min + 1)) + min;
}