FROM alpine:latest

# Install pops
RUN apk add --update ca-certificates curl openssl && \
  curl -L -o /tmp/pops_linux_amd64.tar.gz \
    https://github.com/MYOB-Technology/pops/releases/download/1.1.1/pops_linux_amd64.tar.gz && \
  openssl dgst -sha256 /tmp/pops_linux_amd64.tar.gz \
    | grep 'f2bbf5fd60132c0f9c742e206e610800698a65576b02461ac70b40a3c227f007' \
    || (echo 'shasum mismatch' && false) && \
  tar -xzf /tmp/pops_linux_amd64.tar.gz -C /tmp && \
  rm /tmp/pops_linux_amd64.tar.gz && \
  mv /tmp/pops_linux_amd64/pops /usr/local/bin && \
  rm -rf /tmp/pops_linux_amd64 && \
  apk del ca-certificates curl openssl && \
  rm -rf /tmp/* /var/cache/apk/*

CMD ["pops", "--help"]
