FROM golang:1.15 AS build

ADD . /opt/app
WORKDIR /opt/app
RUN go build ./cmd/project_main.go

FROM ubuntu:20.04

MAINTAINER Felix1Green

RUN apt-get -y update && apt-get install -y tzdata

ENV TZ=Russia/Moscow
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

ENV PGVER 12
RUN apt-get -y update && apt-get install -y postgresql-$PGVER

USER postgres

RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER main_user WITH SUPERUSER PASSWORD 'some_password';" &&\
    createdb -O main_user ForumDatabase &&\
    /etc/init.d/postgresql stop

EXPOSE 5432

VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

USER root

WORKDIR /usr/src/app

COPY . .
COPY --from=build /opt/app/project_main .

EXPOSE 5000
ENV PGPASSWORD some_password
CMD service postgresql start &&  psql -h localhost -d ForumDatabase -U main_user -p 5432 -a -q -f ./init/init_db.sql && ./project_main