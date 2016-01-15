FROM alpine:3.3

# Install pops
RUN apk add --update ca-certificates curl openssl && \
  curl -L -o /tmp/pops_linux_amd64.tar.gz \
    https://github.com/MYOB-Technology/pops/releases/download/1.1.2/pops_linux_amd64.tar.gz && \
  openssl dgst -sha256 /tmp/pops_linux_amd64.tar.gz \
    | grep 'e3b50690420082d2eeff8d62670b0a3da7d7988f13491077094e5b76ec3eaa31' \
    || (echo 'shasum mismatch' && false) && \
  tar -xzf /tmp/pops_linux_amd64.tar.gz -C /tmp && \
  rm /tmp/pops_linux_amd64.tar.gz && \
  mv /tmp/pops_linux_amd64/pops /usr/local/bin && \
  rm -rf /tmp/pops_linux_amd64 && \
  apk del ca-certificates curl openssl && \
  rm -rf /tmp/* /var/cache/apk/*

CMD ["pops", "--help"]
