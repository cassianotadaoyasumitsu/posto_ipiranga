FROM golang:1.18.1-bullseye as builder
WORKDIR /src

ENV GOPRIVATE=git.wealth-park.com
ARG REPO_ACCESS

RUN git config --global url."https://${REPO_ACCESS}@git.wealth-park.com/".insteadOf "https://git.wealth-park.com/"

COPY go.mod go.sum ./
RUN go mod download

COPY . .
ARG VERSION
RUN make ci-build VERSION="${VERSION}" \
    && mkdir -p /out \
    && mv bin/$(go env GOOS)_$(go env GOARCH)/posto_ipiranga /out/ \
    && chmod +x /out/posto_ipiranga

# https://console.cloud.google.com/gcr/images/distroless/global/static-debian11@sha256:90c1cd58d49840ec035eedc6b7d7a4fb6033ae73489a36c669657ed3bfe8ac41/details
# This is the `debug-nonroot` image which comes with a busybox shell.
# This should be updated regularly.
FROM gcr.io/distroless/static-debian11@sha256:90c1cd58d49840ec035eedc6b7d7a4fb6033ae73489a36c669657ed3bfe8ac41
COPY --from=builder /out/posto_ipiranga /posto_ipiranga
ENTRYPOINT ["/posto_ipiranga"]
