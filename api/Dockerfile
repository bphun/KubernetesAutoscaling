FROM python:3.6

RUN adduser api-user

WORKDIR /home/api

COPY requirements.txt requirements.txt
RUN python -m venv venv && \
    venv/bin/pip install -r requirements.txt && \
    venv/bin/pip install "gunicorn>=19.5" gevent && \
    mkdir -p /etc/api-logs

COPY main.py server.py start.sh ./
RUN chmod +x start.sh

ENV FLASK_APP server.py

RUN chown -R api-user:api-user ./ && \
    chown -R api-user:api-user /etc/api-logs
    
USER api-user

EXPOSE 5000
ENTRYPOINT ["./start.sh"]