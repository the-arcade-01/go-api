FROM mysql:8.0.33

# for running migrations on mysql container db
COPY ./db.sql /docker-entrypoint-initdb.d/