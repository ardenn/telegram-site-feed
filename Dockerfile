############################
# STEP 1 build executable binary
############################
FROM golang:alpine as builder

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

# Create a non-root user.
ENV USER=appuser
ENV UID=10001

# See https://stackoverflow.com/a/55757473/12429735RUN 
RUN adduser \    
    --disabled-password \    
    --gecos "" \    
    --home "/nonexistent" \    
    --shell "/sbin/nologin" \    
    --no-create-home \    
    --uid "${UID}" \    
    "${USER}"

WORKDIR $GOPATH/src/sitebot/
COPY . .

# Fetch dependencies using go mod.
RUN go mod download
RUN go mod verify

# Build the binary.
# RUN go build -o /go/bin/sitebot
# Optimize the binary, by removing debug information and
# compiling only for linux target and
# disabling cross compilation.
ENV CGO_ENABLED=0
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/sitebot

############################
# STEP 2 build a small image
############################
FROM scratch

# Copy zoneinfo for timezones
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy SSL certs
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Import the user and group files from the builder.
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Copy our static executable.
COPY --from=builder /go/bin/sitebot /go/bin/sitebot

# Use an unprivileged user.
USER appuser:appuser

# Expose a port  > 1024
EXPOSE 8080

# Run the binary.
ENTRYPOINT ["/go/bin/sitebot"]