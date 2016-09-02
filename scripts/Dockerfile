FROM alpine:3.4

# Install pops
RUN apk add --update ca-certificates curl openssl && \
  curl -L -o /tmp/pops_linux_amd64.tar.gz \
    https://github.com/whence/pops/releases/download/1.1.2/pops_linux_amd64.tar.gz && \
  openssl dgst -sha256 /tmp/pops_linux_amd64.tar.gz \
    | grep 'b8215952e2b187b2616eddc1d68e1f5ee65ed57b7708ff91f395059c0234870e' \
    || (echo 'shasum mismatch' && false) && \
  tar -xzf /tmp/pops_linux_amd64.tar.gz -C /tmp && \
  rm /tmp/pops_linux_amd64.tar.gz && \
  mv /tmp/pops_linux_amd64/pops /usr/local/bin && \
  rm -rf /tmp/pops_linux_amd64 && \
  apk del ca-certificates curl openssl && \
  rm -rf /tmp/* /var/cache/apk/*

CMD ["pops", "--help"]
