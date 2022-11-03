using System;
using System.Threading.Tasks;
using Microsoft.Extensions.Logging;
using Wissance.Authorization.Config;
using Wissance.Authorization.Data;
using Wissance.Authorization.OpenId;

namespace Wissance.Auth.FerrumChecker
{
    class Program
    {
        static void Main(string[] args)
        {
            KeyCloakServerConfig ferrumCfg = new KeyCloakServerConfig("http://127.0.0.1:8182", "myapp",  KeyCloakClientType.Confidential,
                "test-service-app-client", "fb6Z4RsOadVycQoeQiN57xpu8w8wplYz");
            string defUserName = "admin";
            string defPassword = "1s2d3f4g90xs";
            string scope = "profile";

            IOpenIdAuthenticator authenticator = new KeyCloakOpenIdAuthenticator(ferrumCfg, new LoggerFactory());
            Task<TokenInfo> authenticateTask = authenticator.AuthenticateAsync(defUserName,defPassword, scope);
            authenticateTask.Wait();
            TokenInfo token = authenticateTask.Result;
            if (token != null && !string.IsNullOrEmpty(token.AccessToken))
            {
                Console.WriteLine($"Successful authentication: {token.Session}");
            }
            else
            {
                Console.WriteLine("Authentication failed");
            }

            string wait = Console.ReadLine();
        }
    }
}
