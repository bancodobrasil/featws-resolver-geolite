FROM golang:1.18-alpine AS BUILD

RUN apk add curl

WORKDIR /opt

# download db if arg is not empty
ARG FEATWS_GEOLITE_TOKEN
COPY scripts/download-databases.sh .
RUN ./download-databases.sh $FEATWS_GEOLITE_TOKEN

# build app
WORKDIR /app

COPY go.* /app/
RUN go mod download -x

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -ldflags="-w -s" -o resolver

# pack binary
FROM alpine:3.15

ENV LOG_LEVEL=error
ENV SERVER_PORT=7000

COPY --from=BUILD /opt/ /opt/
COPY --from=BUILD /app/resolver /bin/
COPY scripts/startup.sh /

EXPOSE 7000

CMD [ "/startup.sh" ]