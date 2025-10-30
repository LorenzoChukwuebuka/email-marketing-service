#!/usr/bin/env bash
set -e

# ==========================================================
# Usage: ./deploy.sh [environment] [--push]
# Example: ./deploy.sh staging --push
# ==========================================================

if [ -f ".deploy.env" ]; then
  echo "📦 Loading .deploy.env..."
  set -a
  source .deploy.env
  set +a
fi

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
docker build -t $REGISTRY/$OWNER/$MG_NAME:$TAG_SUFFIX -f $MG_DOCKERFILE ./backend &

wait
echo "✅ All builds completed successfully."

# ================================
# Push images (optional)
# ================================
if [ "$PUSH" == "--push" ]; then
  echo "🔐 Authenticating with GitHub Container Registry..."
  
  # Check if GH_TOKEN is set
  if [ -z "$GH_TOKEN" ]; then
    echo "❌ GH_TOKEN not found. Please set it in .deploy.env"
    echo "   Create a token at: https://github.com/settings/tokens"
    echo "   Required scopes: write:packages, read:packages"
    exit 1
  fi
  
  # Login to GHCR
  echo "$GH_TOKEN" | docker login ghcr.io -u "$OWNER" --password-stdin
  
  if [ $? -eq 0 ]; then
    echo "✅ Authentication successful"
  else
    echo "❌ Authentication failed. Please check your GH_TOKEN"
    exit 1
  fi
  
  echo "📦 Pushing images to GHCR..."
  docker push $REGISTRY/$OWNER/$FE_NAME:$TAG_SUFFIX
  docker push $REGISTRY/$OWNER/$BE_NAME:$TAG_SUFFIX
  docker push $REGISTRY/$OWNER/$MG_NAME:$TAG_SUFFIX
  echo "✅ Images pushed to $REGISTRY/$OWNER"
  
  # Logout after push
  docker logout ghcr.io
  echo "🔓 Logged out from GHCR"
fi

echo ""
echo "🎉 Deployment process completed for $ENVIRONMENT environment!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"