##################################
# STEP 1 build executable binary #
##################################
FROM golang:1.12.3-alpine3.9 AS build-env

# Install some dependencies needed to build the project
RUN apk add bash ca-certificates git gcc g++ libc-dev

# All these steps will be cached
RUN mkdir /build
WORKDIR /build

# Copying the .mod and .sum files before the rest of the code
# supposedly improves the caching behavior of Docker
# See: https://medium.com/@petomalina/using-go-mod-download-to-speed-up-golang-docker-builds-707591336888
# And: https://medium.com/@pierreprinetti/the-go-1-11-dockerfile-a3218319d191
COPY go.mod . 
COPY go.sum .  

# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download
# COPY the source code as the last step
COPY . .

# And compile the project 
# CGO_ENABLED and installsuffix are part of the scheme to get better caching on builds
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/appentrypoint

##############################
# STEP 2 build a small image #
##############################
FROM alpine:3.9
WORKDIR /app
RUN apk update && apk add ca-certificates
COPY --from=build-env /build/. /app/.
COPY --from=build-env /go/bin/appentrypoint /app/appentrypoint
ENTRYPOINT ["/app/appentrypoint"]
