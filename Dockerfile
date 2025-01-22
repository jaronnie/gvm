FROM alpine:latest

LABEL \
  org.opencontainers.image.title="gvm" \
  org.opencontainers.image.description="golang version manage tool" \
  org.opencontainers.image.url="https://github.com/jaronnie/gvm" \
  org.opencontainers.image.documentation="https://github.com/jaronnie/gvm#readme" \
  org.opencontainers.image.source="https://github.com/jaronnie/gvm" \
  org.opencontainers.image.licenses="Apache-2.0" \
  maintainer="jaronnie <jaron@jaronnie.com>"

COPY dist/gvm_linux_amd64_v1/gvm /dist/gvm_linux_amd64/gvm
COPY dist/gvm_linux_arm64_v8.0/gvm /dist/gvm_linux_arm64/gvm

# Select the appropriate binary based on the architecture
RUN ARCH=$(uname -m) && \
    if [ "$ARCH" = "x86_64" ] || [ "$ARCH" = "amd64" ]; then \
      cp /dist/gvm_linux_amd64/gvm /usr/local/bin/gvm; \
    elif [ "$ARCH" = "arm64" ] || [ "$ARCH" = "aarch64" ]; then \
      cp /dist/gvm_linux_arm64/gvm /usr/local/bin/gvm; \
    fi

RUN gvm init sh \
    && cp /root/gvm/.gvmrc /etc/profile.d/.gvmrc.sh