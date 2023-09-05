############################
# STEP 1 build executable binary
############################
FROM golang:1.20.3-alpine3.17 as builder

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

WORKDIR /app
COPY . .

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux go build -o krt-files-downloader ./cmd/krt-files-downloader

############################
# STEP 2 build a small image
############################
FROM scratch

# Copy the static executable.
COPY --from=builder /app/krt-files-downloader /app/krt-files-downloader

# Run the hello binary.
ENTRYPOINT ["/app/krt-files-downloader"]
