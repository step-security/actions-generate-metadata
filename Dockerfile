# Copyright IBM Corp. 2021, 2025
# Copyright 2026 StepSecurity
# SPDX-License-Identifier: MPL-2.0

FROM golang:1.27@sha256:7543a96ce82c8e9003cae079ee3e0bc5b7799df8eed2a041e403af0d31fa4e67 AS build
LABEL maintainer="step-security security@stepsecurity.io"

# Copy all the action files into the container
WORKDIR /go/src/action
COPY action /go/src/action

# Enable Go modules
ENV GO111MODULE=on
RUN go mod download

# Compile the action
RUN CGO_ENABLED=0 go build -o /action -ldflags="-s -w" .

FROM alpine:latest@sha256:28bd5fe8b56d1bd048e5babf5b10710ebe0bae67db86916198a6eec434943f8b
RUN apk --update upgrade && apk add --no-cache ca-certificates git make bash

COPY --from=build /action /
# Specify the container's entrypoint as the action
ENTRYPOINT ["/action"]
