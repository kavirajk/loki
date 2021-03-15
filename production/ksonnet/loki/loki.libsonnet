local consul = import "consul/consul.libsonnet";

(import 'ksonnet-util/kausal.libsonnet') +
(import 'jaeger-agent-mixin/jaeger.libsonnet') +
(import 'images.libsonnet') +
(import 'common.libsonnet') +
(import 'config.libsonnet') +
(import 'overrides.libsonnet') +

consul + {
  // Without consul-sidekick
  consul_deployment:
    deployment.new('consul', $._config.consul_replicas, [
      $.consul_container,
      $.consul_statsd_exporter,
      $.consul_exporter,
    ]) +
    $.util.configMapVolumeMount($.consul_config_map, '/etc/config') +
    $.util.antiAffinity,
}
// Loki services
(import 'distributor.libsonnet') +
(import 'ingester.libsonnet') +
(import 'querier.libsonnet') +
(import 'table-manager.libsonnet') +
(import 'query-frontend.libsonnet') +
(import 'ruler.libsonnet') +

// Supporting services
(import 'memcached.libsonnet') +

// WAL support
(import 'wal.libsonnet') +

// BoltDB Shipper support. This should be the last one to get imported.
(import 'boltdb_shipper.libsonnet')
