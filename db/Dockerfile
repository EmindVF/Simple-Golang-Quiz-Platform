FROM postgres:alpine

ENV POSTGRES_USER=admin
ENV POSTGRES_PASSWORD=password
ENV POSTGRES_DB=database

COPY postgres_init.sql /docker-entrypoint-initdb.d/