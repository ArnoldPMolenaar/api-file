FROM golang:1.23-alpine AS builder

LABEL authors="Arnold Molenaar <arnold.molenaar@webmi.nl> (https://arnoldmolenaar.nl/)"

# Install libvips
RUN apk update &&\
    apk add --update --no-cache gcc g++ vips vips-dev

# Move to working directory (/build).
WORKDIR /build

# Copy and download dependency using go mod.
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy the code into the container.
COPY . .

# Set necessary environment variables needed for our image and build the API.
ENV GOOS=linux GOARCH=amd64
RUN go build -a -installsuffix cgo -ldflags="-s -w" -o api .

FROM scratch

# Copy binary and config files from /build to root folder of scratch container.
COPY --from=builder ["/build/api", "/build/.env", "/"]

EXPOSE 5000

# Command to run when starting the container.
ENTRYPOINT ["/api"]