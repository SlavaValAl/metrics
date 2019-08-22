# Metrics collector

Simple tool receiving events via http and puts it to file.

## Installing / Getting started

### Building
```sh
$ cd cmds/metrics/ && go build -o metrics
```

### Run
```sh
./metrics
```

### Usage example:
```sh
curl -L -w "\nStatus code: %{http_code}\n" -X POST 127.0.0.1:9000/ga -H "Content-Type: application/json" --data '{"event_type": "test_type", "url": "google.com"}'
```

### Load testing
Run from project root
```sh
ab -p ./tests/metric.txt -T application/json -c 10 -n 10 http://127.0.0.1:9000/ga
```
