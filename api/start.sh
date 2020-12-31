#!/bin/bash

source venv/bin/activate

N_PROC=$(nproc --all)
N_WORKERS=$(($N_PROC + 1))
N_THREADS=$(($N_PROC + 1))

# exec gunicorn -b :5000 --statsd-host=prometheus-scrapers.monitoring:9125 --log-level=debug --statsd-prefix=api server:app
exec gunicorn -b :5000 --statsd-host=localhost:9125 --preload --reuse-port --worker-class=gevent --workers=$N_WORKERS --threads=$N_THREADS --worker-connections=20 --log-level=debug --statsd-prefix=api server:app
# exec gunicorn -b :5000 --statsd-host=prometheus-scrapers.monitoring:9125 --log-level=debug --statsd-prefix=api --worker-class=gevent --workers=$N_WORKERS --threads=$N_WORKERS server:app
# exec python server.py
# --access-logfile=/etc/api-logs/access.log --error-logfile=/etc/api-logs/error.log
#--worker-connections=1000 