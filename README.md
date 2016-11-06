# Shorty Challenge

Requirements have been moved to [requirements.md](./REQUIREMENTS.md)

## Running instructions:

  1. You need to have Docker installed - [docker installation instructions](https://docs.docker.com/engine/installation)
  2. You need to have Docker Compose installed - [docker compose installation instructions](https://docs.docker.com/compose/install)
  3. In the root directory of the project run `docker-compose -p impraise up`
  4. The service should be accessible on localhost:8080

## Running the tests

  1. Make sure you went through the previous step, and the container is running. Run `docker ps` in terminal. You should see a container named `impraise_shorty_1`
  2. Execute `docker exec -it impraise_shorty_1 go test -v -cover -timeout=5m -goblin.timeout=5m`
  3. ???
  4. Profit?

## Known issues / Reflections

  1. Running the tests on my machine, they finish in ~1 second. When running the tests inside the docker container, the running time spikes to ~3 minutes !?
     * I'm either doing something wrong (very likely)
     * Docker is a bad boy (less likely, but still a plausible cause)
  2. On Redirect, the API returns a 302 Found status code. However the `gorequest` library I have used to make the requests, seems to get a 200 OK :-(
  3. I'm pretty satisfied with the test coverage: 90%. Seems to have covered pretty much all the cases.

## Benchmarks

  ```
    shorty git:development ❯ ab -c 20 -n 2000 http://127.0.0.1:8080/Wy_G8D/stats

    This is ApacheBench, Version 2.3 <$Revision: 1748469 $>
    Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
    Licensed to The Apache Software Foundation, http://www.apache.org/

    Benchmarking 127.0.0.1 (be patient)
    Completed 200 requests
    Completed 400 requests
    Completed 600 requests
    Completed 800 requests
    Completed 1000 requests
    Completed 1200 requests
    Completed 1400 requests
    Completed 1600 requests
    Completed 1800 requests
    Completed 2000 requests
    Finished 2000 requests

    Server Software:
    Server Hostname:        127.0.0.1
    Server Port:            8080

    Document Path:          /Wy_G8D/stats
    Document Length:        90 bytes

    Concurrency Level:      20
    Time taken for tests:   2.739 seconds
    Complete requests:      2000
    Failed requests:        0
    Total transferred:      426000 bytes
    HTML transferred:       180000 bytes
    Requests per second:    730.08 [#/sec] (mean)
    Time per request:       27.394 [ms] (mean)
    Time per request:       1.370 [ms] (mean, across all concurrent requests)
    Transfer rate:          151.86 [Kbytes/sec] received

    Connection Times (ms)
                  min  mean[+/-sd] median   max
    Connect:        0    0   0.3      0       4
    Processing:     1   27  49.6      8     307
    Waiting:        1   27  49.4      8     306
    Total:          1   27  49.6      9     307

    Percentage of the requests served within a certain time (ms)
      50%      9
      66%     18
      75%     28
      80%     37
      90%     71
      95%    119
      98%    242
      99%    279
     100%    307 (longest request)
 ```
