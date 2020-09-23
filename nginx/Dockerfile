FROM nginx:1.17.8-alpine

RUN rm /etc/nginx/nginx.conf

COPY nginx.conf /etc/nginx/nginx.conf

# RUN adduser -D nginxUser

# RUN touch /var/run/nginx.pid && \
#     chown -R nginxUser:nginxUser /var/cache/nginx/ && \
#     chown -R nginxUser:nginxUser /var/log/nginx/ && \
#     chown -R nginxUser:nginxUser /etc/nginx/nginx.conf && \
#     chown -R nginxUser:nginxUser /var/run/nginx.pid

# USER nginxUser
