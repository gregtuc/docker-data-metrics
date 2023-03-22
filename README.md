# docker data metrics

small program to collect data on a docker container and make it available on a webserver

## Try it out:
---

Spin up a docker container using the below command(s) in cmd/terminal:

`docker pull progrium/stress`

`docker run --name testcontainer --rm -it progrium/stress --cpu 1 --io 1 --vm 2 --vm-bytes 128M`

Alternatively, start your own Docker container + name and change the name value in the `.env` file.

## Execution:

---

`cd docker-data-metrics` then `go run .`

## Testing:

---

### Get all Logs

`GET` to `http://localhost:8080/metrics/`

### Get all logs under a certain CPU %

`GET` to `http://localhost:8080/metrics/cpu?under=<cpu_percent>`

### Get all logs over a certain CPU %

`GET` to `http://localhost:8080/metrics/cpu?over=<cpu_percent>`

### Get all logs within a certain CPU % range

`GET` to `http://localhost:8080/metrics/cpu?over=<cpu_percent>&under=<cpu_percent>`

### Get all logs before a certain time

`GET` to `http://localhost:8080/metrics/?before=<unix_timestamp>`

### Get all logs after a certain time

`GET` to `http://localhost:8080/metrics/?after=<unix_timestamp>`

### Get all logs within a certain time range

`GET` to `http://localhost:8080/metrics/?before=<unix_timestamp>&after=<unix_timestamp>`

### Get most recent log (live)

`GET` to `http://localhost:8080/metrics/live`
