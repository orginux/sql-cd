FROM golang:1.18.1-alpine3.15 AS builder
LABEL maintainer="orginux"
RUN apk add --no-cache make git
WORKDIR ${GOPATH}/src/sql-cd
COPY . .
RUN CGO_ENABLED=0 GOOS=linux \
        go build -o ${GOPATH}/bin/sql-cd

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/sql-cd /go/bin/sql-cd
ENTRYPOINT ["/go/bin/sql-cd"]
