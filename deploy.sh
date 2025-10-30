#!/usr/bin/env bash
set -e

# ==========================================================
# Usage: ./deploy.sh [environment] [--push]
# Example: ./deploy.sh staging --push
# ==========================================================

ENVIRONMENT=${1:-dev}
PUSH=${2:-""}

# GitHub Container Registry settings
REGISTRY="ghcr.io"
OWNER="lorenzochukwuebuka"
FE_NAME="crabmailer-frontend"
BE_NAME="crabmailer-backend"
MG_NAME="crabmailer-migrate"

# Choose Dockerfiles based on environment
case "$ENVIRONMENT" in
  dev)
    FE_DOCKERFILE="./frontend/Dockerfile.frontend.dev"
    BE_DOCKERFILE="./backend/Dockerfile.backend.development"
    MG_DOCKERFILE="./backend/Dockerfile.migrate"
    TAG_SUFFIX="dev"
    ;;
  staging)
    FE_DOCKERFILE="./frontend/Dockerfile.frontend.staging"
    BE_DOCKERFILE="./backend/Dockerfile.backend.staging"
    MG_DOCKERFILE="./backend/Dockerfile.migrate"
    TAG_SUFFIX="staging"
    ;;
  prod)
    FE_DOCKERFILE="./frontend/Dockerfile.frontend.prod"
    BE_DOCKERFILE="./backend/Dockerfile.backend.production"
    MG_DOCKERFILE="./backend/Dockerfile.migrate"
    TAG_SUFFIX="latest"
    ;;
  *)
    echo "❌ Invalid environment. Use: dev | staging | prod"
    exit 1
    ;;
esac

echo "🚀 Building Docker images for $ENVIRONMENT environment..."
echo "🧩 Registry: $REGISTRY"
echo "🧩 Owner: $OWNER"
echo "🧩 Frontend Image: $FE_NAME:$TAG_SUFFIX"
echo "🧩 Backend Image: $BE_NAME:$TAG_SUFFIX"
echo "🧩 Migration Image: $MG_NAME:$TAG_SUFFIX"

# ================================
# ENV VARIABLES for build context
# ================================
export NODE_ENV=$ENVIRONMENT
export API_URL="https://api.${ENVIRONMENT}.example.com"
export DATABASE_URL="postgres://user:password@localhost:5432/${ENVIRONMENT}_db"

# ================================
# Build FE + BE + Migrate simultaneously
# ================================
echo "🔨 Building frontend..."
docker build -t $REGISTRY/$OWNER/$FE_NAME:$TAG_SUFFIX -f $FE_DOCKERFILE ./frontend &

echo "🔨 Building backend..."
docker build -t $REGISTRY/$OWNER/$BE_NAME:$TAG_SUFFIX -f $BE_DOCKERFILE ./backend &

echo "🔨 Building migration image..."
docker build -t $REGISTRY/$OWNER/$MG_NAME:$TAG_SUFFIX -f $MG_DOCKERFILE ./migrations &

wait
echo "✅ All builds completed successfully."

# ================================
# Push images (optional)
# ================================
if [ "$PUSH" == "--push" ]; then
  echo "📦 Pushing images to GHCR..."
  docker push $REGISTRY/$OWNER/$FE_NAME:$TAG_SUFFIX
  docker push $REGISTRY/$OWNER/$BE_NAME:$TAG_SUFFIX
  docker push $REGISTRY/$OWNER/$MG_NAME:$TAG_SUFFIX
  echo "✅ Images pushed to $REGISTRY/$OWNER"
fi
