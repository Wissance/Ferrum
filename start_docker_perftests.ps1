try
{
   cp .env.perftests .env
   docker-compose -f docker-compose.perf.yml  up --build
}
finally
{
    cp .env.local .env
}