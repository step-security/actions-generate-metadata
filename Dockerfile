# Copyright IBM Corp. 2021, 2025
# Copyright 2026 StepSecurity
# SPDX-License-Identifier: MPL-2.0

FROM golang:1.25 AS build
LABEL maintainer="step-security security@stepsecurity.io"

# Copy all the action files into the container
WORKDIR /go/src/action
COPY action /go/src/action

# Enable Go modules
ENV GO111MODULE=on
RUN go mod download

# Compile the action
RUN CGO_ENABLED=0 go build -o /action -ldflags="-s -w" .

FROM alpine:latest
RUN apk --update add ca-certificates
RUN apk add --no-cache git make bash

COPY --from=build /action /
# Specify the container's entrypoint as the action
ENTRYPOINT ["/action"]
