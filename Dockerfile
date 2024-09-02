# Stage 1: Build the Go backend
FROM golang:1.21 as backend-builder
WORKDIR /app

# Copy the backend source code
COPY ./backend/ ./

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final Stage: Production image
FROM alpine:3.17

# Install dependencies
RUN apk update && apk add --no-cache \
    postfix \
    opendkim \
    opendmarc \
    spamassassin \
    clamav \
    clamav-libunrar \
    rsyslog \
    bash \
    curl \
    ca-certificates \
    tzdata \
    redis \
    fail2ban \
    supervisor \
    dovecot \
    dovecot-pop3d \
    openssl

WORKDIR /root/

# Copy the Go binary from the build stage
COPY --from=backend-builder /app/main .

# Copy .env file from the current directory in the build context
COPY ./backend/.env .env

# Copy healthcheck script
COPY ./healthcheck.sh /healthcheck.sh
RUN chmod +x /healthcheck.sh

# Generate self-signed SSL certificate for Dovecot
RUN openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
    -keyout /etc/dovecot/dovecot.key -out /etc/dovecot/dovecot.pem \
    -subj "/C=US/ST=State/L=City/O=Organization/CN=localhost"

# Expose ports
EXPOSE 9000 1025 587 465 143 993 110 995

# Add healthcheck
HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=3 CMD /healthcheck.sh

# Debugging steps
RUN which dovecot
RUN find / -name "pop3" 2>/dev/null
RUN ls -l /usr/lib/dovecot
RUN dovecot --version

# Start services
CMD ["sh", "-c", ". .env && postfix start && dovecot && ./main"]