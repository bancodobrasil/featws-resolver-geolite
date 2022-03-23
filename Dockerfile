FROM golang:1.18-alpine AS BUILD

WORKDIR /opt

# download db if arg is not empty
ARG MAXMIND_LICENSE_KEY
ADD /scripts/download-databases.sh /opt/
RUN /opt/download-databases.sh $MAXMIND_LICENSE_KEY

WORKDIR /app

COPY go.* /app/
RUN go mod download

COPY . /app/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -ldflags="-w -s" -o /resolver-geolite

FROM alpine:3.15

COPY --from=BUILD /app/resolver-geolite /bin/

CMD [ "resolver-geolite serve" ]