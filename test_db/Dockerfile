FROM mysql:latest

ENV MYSQL_ROOT_PASSWORD 123
ENV MYSQL_DATABASE mobile_data
ENV MYSQL_USER treigerm
ENV MYSQL_PASSWORD 123

ADD setupBig.sql /docker-entrypoint-initdb.d 

EXPOSE 3306
