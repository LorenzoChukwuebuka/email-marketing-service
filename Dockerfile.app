# Stage 1: Build the Go backend
FROM golang:1.20 as backend-builder

WORKDIR /app

# Copy the backend source code
COPY ./backend/ ./

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final Stage: Production image
FROM alpine:3.17

RUN apk update && apk add --no-cache \
    bash \
    curl 

WORKDIR /root/

# Copy the Go binary from the build stage
COPY --from=backend-builder /app/main .

# Copy .env file from the current directory in the build context
COPY ./backend/.env .env

# Copy healthcheck script
COPY ./healthcheck.sh /healthcheck.sh
RUN chmod +x /healthcheck.sh

# Expose port
EXPOSE 9000 1025

# Add healthcheck
HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=3 CMD /healthcheck.sh

# Start the Go application
CMD ["sh", "-c", ". .env && ./main"]