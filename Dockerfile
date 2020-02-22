FROM golang:latest as builder

WORKDIR /app

COPY go.mod go.sum .

RUN go mod download

COPY . .

RUN groupadd --gid 1000 clozapinum \
&& useradd --uid 1000 clozapinum

RUN go build -o clozapinum .



FROM alpine:latest

COPY --from=builder /etc/passwd /etc/passwd

USER clozapinum

COPY --from=builder /app/clozapinum /app/

EXPOSE 8443

ENTRYPOINT ["/app/clozapinum"] 