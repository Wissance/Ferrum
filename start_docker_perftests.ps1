try
{
   cp .env.perftests .env
   docker-compose -d -f docker-compose.perf.yml  up --build
}
finally
{
    cp .env.local .env
}