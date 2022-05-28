FROM golang:1.18.2-bullseye AS builder
LABEL maintainer="orginux"
WORKDIR ${GOPATH}/src/sql-cd
COPY . .
RUN ssh-keyscan -H github.com > /ssh_known_hosts
RUN CGO_ENABLED=0 GOOS=linux \
        go build -o ${GOPATH}/bin/sql-cd

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/sql-cd /go/bin/sql-cd
COPY --from=builder /ssh_known_hosts /etc/ssh/ssh_known_hosts
COPY github /tmp/key
ENTRYPOINT ["/go/bin/sql-cd"]
