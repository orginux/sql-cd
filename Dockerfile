FROM golang:1.18.2-bullseye AS builder
LABEL maintainer="orginux"
WORKDIR ${GOPATH}/src/sql-cd
COPY . .
RUN CGO_ENABLED=0 GOOS=linux \
        go build -o ${GOPATH}/bin/sql-cd

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/sql-cd /go/bin/sql-cd
ENTRYPOINT ["/go/bin/sql-cd"]
