# cloudtail
Extend loki's promtail to capture Cloud logs

# Goal
End goal is eventually to merge into `promtail` as another target manager (`pubsub`).

# Pre-requisites
Generate an service-account key (json) for `ops-tools` (project) -> `logging-client` (service-account).

This will be used to set `GOOGLE_APPLICATION_CREDENTIALS` later to access GCP APIs.

# How it works

GCP Cloud Logging -> Sink("logs-dev") -> PubSub Destination ("loki-subscription") -> Cloudtail -> Loki

>NOTE: Currently `sink` is created with filter `resource.type=project` to capture project events, it can be changed into any other resources say `resource.type="http_load_balancer"`.

# Usage

## 1. Running `cloudtail`

Make sure `loki`is running somewhere and can be reached via `LOKI_URL`

```bash
$ GOOGLE_APPLICATION_CREDENTIALS="<path-to-json>" LOKI_URL=<loki-url> go run main.go
```

## 2. Generate logs

Generate some logs by loading `grafanalabs-dev` project on GCP console

## 3. Query Loki
```bash
$ logcli query '{source="cloudtail"}'
```

# TODO
- [ ] test it with big stream of logs from GCP
- [ ] test with retry!, timestamp ordering
- [ ] benchmark
- [ ] test with `grafanalabs-global` logs
