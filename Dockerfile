# Stage 1: Build the Go backend
FROM golang:1.21 as backend-builder
WORKDIR /app/backend
COPY backend/ ./
RUN make buildfe
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final Stage: Production image
FROM alpine:3.17
RUN apk --no-cache add ca-certificates wget bash

WORKDIR /root/

# Copy the Go binary
COPY --from=backend-builder /app/backend/main .

# Copy .env file
COPY backend/.env .env

# Copy healthcheck script
COPY healthcheck.sh /healthcheck.sh
RUN chmod +x /healthcheck.sh

# Expose the Go application port
EXPOSE 9000

# Add healthcheck
HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=3 CMD /healthcheck.sh

# Start the application and load the environment variables
CMD ["sh", "-c", "set -o allexport; source .env; set +o allexport; ./main"]
