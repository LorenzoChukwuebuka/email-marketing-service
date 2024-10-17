# Stage 1: Build the React frontend
FROM node:16 as frontend-builder
WORKDIR /app/frontend
COPY ./frontend/ ./
RUN npm install
RUN npm run build

# Stage 2: Build the Go backend
FROM golang:1.20 as backend-builder
WORKDIR /app/backend
COPY ./backend/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Stage 3: Final image
FROM nginx:alpine
WORKDIR /root/

# Copy the frontend build from frontend-builder
COPY --from=frontend-builder /app/frontend/dist /usr/share/nginx/html

# Copy the Go binary from the backend-builder
COPY --from=backend-builder /app/backend/main /root/

# Copy .env file for the backend
COPY ./backend/.env /root/.env

# Copy healthcheck script
COPY ./healthcheck.sh /healthcheck.sh
RUN chmod +x /healthcheck.sh

# Copy Nginx configuration
COPY ./config/nginx.conf /etc/nginx/nginx.conf

EXPOSE 80 9000 1025

# Add healthcheck
HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=3 CMD /healthcheck.sh

# Start Nginx and the Go application
CMD ["sh", "-c", "nginx && . /root/.env && /root/main"]