FROM prom/prometheus

ADD prometheus.yml /etc/prometheus/

ENTRYPOINT [ "/bin/prometheus" ]
CMD        [ "--config.file=/etc/prometheus/prometheus.yml", \
             "--storage.tsdb.path=/prometheus", \
             "--web.console.libraries=/usr/share/prometheus/console_libraries", \
             "--web.console.templates=/usr/share/prometheus/consoles", \
             "--web.external-url=http://192.168.1.64/prometheus/" ]
            #  "--web.route-prefix=/prometheus" ]