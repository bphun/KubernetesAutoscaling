#!/bin/bash

source venv/bin/activate

N_PROC=$(nproc --all)
N_WORKERS=$(($N_PROC * 2 + 1))

exec gunicorn -b :5000 --statsd-host=statsd:9125 --statsd-prefix=api --worker-class=gevent --workers=$N_WORKERS --threads=$N_WORKERS --log-level=debug --access-logfile=/etc/api-logs/access.log --error-logfile=/etc/api-logs/error.log server:app
# exec python server.py
#--worker-connections=1000 