FROM prom/statsd-exporter

ADD statsd.conf /statsd/

ENTRYPOINT ["/bin/statsd_exporter", "--web.telemetry-path=/statsd/metrics", "--statsd.listen-tcp=:9125", "--statsd.mapping-config=/statsd/statsd.conf", "--web.listen-address=0.0.0.0:9091", "--log.level=debug"]