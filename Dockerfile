FROM golang:1.22-alpine

COPY main /app/main

CMD ["sh", "/app/main"]