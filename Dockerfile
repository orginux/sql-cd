FROM golang:1.18.2-bullseye AS builder
LABEL maintainer="orginux"
WORKDIR ${GOPATH}/src/sql-cd
COPY . .
RUN CGO_ENABLED=0 GOOS=linux \
        go build -o ${GOPATH}/bin/sql-cd

FROM golang:1.18.2-bullseye
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/sql-cd /go/bin/sql-cd
# COPY github /tmp/key
ENTRYPOINT ["/go/bin/sql-cd"]
