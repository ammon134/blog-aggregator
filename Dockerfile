FROM debian:stable-slim

COPY blog-aggregator bin/blog-aggregator
COPY .env .env

CMD ["/bin/blog-aggregator"]
