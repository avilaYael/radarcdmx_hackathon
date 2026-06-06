############################
# STEP 1 build executable binary
############################

FROM golang:alpine as builder

# Install git.
# Git is required for fetching the dependencies.
RUN apk add --no-cache git

# set directory for all go code
WORKDIR $GOPATH/src/rcapi

# copy all code into working directory
COPY . .

# Fetch dependencies
RUN go mod tidy

# Build the binary.
RUN go build .

############################
# STEP 2 final thin build
############################
FROM alpine:3.22.1
ENV ENV=prod
RUN apk --no-cache add bash curl ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/rcapi .
ENTRYPOINT ["bash", "-c", "/root/rcapi"]
