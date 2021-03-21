FROM mongo:latest

ENV PUID=1000
ENV PGID=1000

EXPOSE 27017

ENTRYPOINT ["mongod", "--bind_ip", "0.0.0.0"]