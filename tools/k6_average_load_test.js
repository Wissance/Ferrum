import http from 'k6/http';
import { check } from "k6";
import { sleep } from 'k6';
export let options = {
    stages: [
        // Ramp-up from 1 to 100 VUs in 1m
        { duration: "1m", target: 100 },
        // Stay at rest on 100 VUs for 4m
        { duration: "5m", target: 100 },
        // Ramp-down from 100 to 200 VUs for 2m
        { duration: "7m", target: 200 },
        { duration: "10m", target: 200 },
        { duration: "12m", target: 300 },
        { duration: "14m", target: 400 },
        { duration: "15m", target: 500 },
        { duration: "20m", target: 500 },
        { duration: "25m", target: 0 }
    ],
    thresholds: {
        http_req_duration: ['p(95)<250'], // 95% of requests must complete in less than 500ms for the test to pass
        http_req_failed: ['rate<5'],      // Test fails if more than 5% of requests fail
    },

};

export default function () {
    const numOfIterationsStage1 = 100
    const numOfIterationsStage2 = 50
    // 1. Select Random User
    const userPassword = "P@55W0rD"
    const clientSecret = "00000000000000000000000000000000"
    let ferrumBaseUrl = "http://10.50.40.3:8182";
    // 1. Get random realm (1-100) && user (1-100)
    let randomRealm = getRandomInt(1, 100)
    let realm = "realm_" + randomRealm
    console.log("Using Realm is: " + realm)
    let randomUserRel = getRandomInt(1, 100)
    let absUser = (randomRealm - 1) * 100 + randomUserRel
    let user = "u" + absUser
    console.log("Using User is: " + user)
    let clientId = "client_" + randomRealm
    console.log("Using Client is: " + clientId)
    // 2. Get initial access token
    let getTokenResponse = getAccessToken(ferrumBaseUrl, realm, clientId, clientSecret, user, userPassword)
    // check status
    check(getTokenResponse, {
        'Get access token status is 200': (r) => r.status === 200,
    });
    const responseBody = JSON.parse(getTokenResponse.body);
    let accessToken = responseBody.access_token
    // 3. send up 100 requests userinfo (1-2 sec interval)
    for (let i = 0; i < numOfIterationsStage1; i++) {
        if (i > 0 && i%10 === 0)
        {
            // send refresh token + introspect
        }
        // 3.1 after every 10 request rotate key
        let getUserInfoResponse = getUserInfo(ferrumBaseUrl, realm, accessToken)
        check(getUserInfoResponse, {
            'Get userinfo status is not 500': (r) => r.status !== 500,
        });
        if (getUserInfoResponse.status === 401) {
            getTokenResponse = getAccessToken(ferrumBaseUrl, realm, clientId, clientSecret, user, userPassword)
            // check status
            check(getTokenResponse, {
                'Get access token status is 200': (r) => r.status === 200,
            });
            const responseBody = JSON.parse(getTokenResponse.body);
            accessToken = responseBody.access_token
        }
        getUserInfoResponse = getUserInfo(ferrumBaseUrl, realm, accessToken)
        check(getUserInfoResponse, {
            'Get userinfo status is 200': (r) => r.status !== 500,
        });
        let pause = getRandomInt(1,3)
        sleep(pause)
    }

    // 4 wait 2m
    sleep(120)
    // 5 send up to 50 requests userinfo (2-3 sec interval)
    // 5.1 after every 10 request rotate key - send refresh token + introspect
}

/* This function get access token from Ferrum
 * baseUrl is part protocol://host:port
 * */
function getAccessToken(baseUrl, realm, clientId, secret, username, password) {
    const url = baseUrl+"/auth/realms/" + realm + "/protocol/openid-connect/token";

    var payload = {
        client_id: clientId,
        client_secret : secret,
        username: username,
        password: password,
        grant_type: "password",
        scope: "profile"
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

/* This function get userInfo from Ferrum
 * baseUrl is part protocol://host:port
 * */
function getUserInfo(baseUrl, realm, accessToken) {
    const url = baseUrl+"/auth/realms/" + realm + "/protocol/openid-connect/userinfo";

    let params = {
        headers: {
            "Authorization": "Bearer " + accessToken,
        },
    };

    return http.get(url, params);
}

function refreshAccessToken() {

}

function  introspectToken() {

}

function getRandomInt(min, max) {
    min = Math.ceil(min);
    max = Math.floor(max);
    return Math.floor(Math.random() * (max - min + 1)) + min;
}
