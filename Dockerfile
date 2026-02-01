FROM golang:1.22-alpine

WORKDIR /app

COPY main /app/main

RUN chmod +x /app/main

EXPOSE 8080

CMD ["./main"]