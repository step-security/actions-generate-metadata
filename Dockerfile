# Copyright IBM Corp. 2021, 2025
# Copyright 2026 StepSecurity
# SPDX-License-Identifier: MPL-2.0

FROM golang:1.25@sha256:699337d620559a59b4a2bb298ad59611e535d2ee755a34cf2d2a98f37578dc80 AS build
LABEL maintainer="step-security security@stepsecurity.io"

# Copy all the action files into the container
WORKDIR /go/src/action
COPY action /go/src/action

# Enable Go modules
ENV GO111MODULE=on
RUN go mod download

# Compile the action
RUN CGO_ENABLED=0 go build -o /action -ldflags="-s -w" .

FROM alpine:latest@sha256:e7c4abb69531cb09e2a2bbb56fad3367ab694865c49df898c1c683185cc4376c
RUN apk --update upgrade && apk add --no-cache ca-certificates git make bash

COPY --from=build /action /
# Specify the container's entrypoint as the action
ENTRYPOINT ["/action"]
