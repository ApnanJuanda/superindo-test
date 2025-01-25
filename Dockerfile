################
# BUILD BINARY #
################
# golang:1.18.2-alpine3.16
FROM golang:1.23.5-alpine3.21 as builder


# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

WORKDIR $GOPATH/superindo
COPY . .

RUN echo $PWD && ls -la

# Fetch dependencies.
# RUN go get -d -v
RUN go mod download
RUN go mod verify

#CMD go build -v
# go build command with the -ldflags="-w -s" option to produce a smaller binary file by stripping debug information and symbol tables.
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o /go/bin/superindo .

#####################
# MAKE SMALL BINARY #
#####################
FROM alpine:3.21

RUN apk update

# Import from builder.
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd

# Copy the executable.
COPY --from=builder /go/bin/superindo /go/bin/superindo

# set env
# docker run -d -p 8080:8080 --name superindo-container --network superindo-network -e MYUSERNAME=root -e MYPASSWORD= -e MYHOST=host.docker.internal -e MYPORT=330
#6 -e MYDATABASE=superindo -e PORT=8080 -e RSPORT=6379 superindo

ENTRYPOINT ["/go/bin/superindo"]
