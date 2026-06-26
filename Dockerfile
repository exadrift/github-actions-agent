FROM alpine:3 AS certs
RUN apk --update add ca-certificates

FROM scratch
ARG NAME=_
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ENV SSL_CERT_FILE=/etc/ssl/certs/ca-certificates.crt

COPY bin/${NAME} /${NAME}
ENTRYPOINT [ "/${NAME}" ]
