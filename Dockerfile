FROM ubuntu:focal
ARG TARGETOS
ENV TARGETOS=${TARGETOS:-linux}
ARG TARGETARCH
ENV TARGETARCH=${TARGETARCH:-amd64}
ENTRYPOINT ["/pixlet"]
RUN apt update && apt install -y \
  ca-certificates \
  libwebp-dev \
  && apt clean
COPY ./build/${TARGETOS}_${TARGETARCH}/pixlet /pixlet
