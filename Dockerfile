FROM centos:8

COPY ./app /home/app

ENV POSTGRES_USERNAME=postgres \ 
    POSTGRES_DATABASE=stock \
    POSTGRES_PASSWORD=postgres \
    POSTGRES_PORT=5432 \ 
    POSTGRES_HOST=127.0.0.1 \
    TZ=Asia/Shanghai

EXPOSE 8080

WORKDIR /home

ENTRYPOINT ["./app"]