# Stage 1: Build the Go backend
FROM golang:1.21 as backend-builder
WORKDIR /app

# Copy the backend source code
COPY . ./

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final Stage: Production image
FROM alpine:3.17
RUN apk --no-cache add ca-certificates wget bash curl

WORKDIR /root/


# Copy the Go binary from the build stage
COPY --from=backend-builder /app/main .

# Copy .env file from the current directory in the build context
COPY .env .env


RUN cat .env

# Copy healthcheck script
COPY ../healthcheck.sh /healthcheck.sh
RUN chmod +x /healthcheck.sh

RUN ls -l /healthcheck.sh

# Expose the Go application port
EXPOSE 9000 1025

# Add healthcheck
HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=3 CMD /healthcheck.sh

# Start the application and load the environment variables
CMD ["sh", "-c", ". .env; ./main"]
