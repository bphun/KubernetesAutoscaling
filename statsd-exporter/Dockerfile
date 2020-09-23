FROM prom/statsd-exporter

ADD statsd.conf /statsd/

ENTRYPOINT ["/bin/statsd_exporter", "--statsd.mapping-config=/statsd/statsd.conf", "--web.listen-address=:9091", "--log.level=debug"]