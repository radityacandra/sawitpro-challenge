# Dockerfile definition for Backend application service.

# From which image we want to build. This is basically our environment.
FROM golang:1.20-alpine as Build

RUN apk add make
RUN go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest

# This will copy all the files in our repo to the inside the container at workdir location.
RUN mkdir /var/app
WORKDIR /var/app
COPY . /var/app/
RUN make generated api.yml

# Build our binary at root location.
RUN go build -o /main ./cmd/main.go

####################################################################
# This is the actual image that we will be using in production.
FROM alpine:latest

# We need to copy the binary from the build image to the production image.
COPY --from=Build /main .

# This is the port that our application will be listening on.
EXPOSE 1323

# This is the command that will be executed when the container is started.
ENTRYPOINT ["./main"]