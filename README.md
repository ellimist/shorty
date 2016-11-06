# Shorty Challenge

Requirements can be found in [requirements.md](./REQUIREMENTS.md)

## Running instructions:

  1. You need to have Docker installed - [docker installation instructions](https://docs.docker.com/engine/installation)
  2. You need to have Docker Compose installed - [docker compose installation instructions](https://docs.docker.com/compose/install)
  3. In the root directory of the project run `docker-compose -p impraise up`
  4. The service is accessible on localhost:8080

## Running the tests

  1. Make sure you went through the previous step, and the container is running. Run `docker ps` in terminal. You should see a container named `impraise_shorty_1`
  2. Execute `go test -v -cover`
  3. ???
  4. Profit?
