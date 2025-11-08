#!/usr/bin/env bash

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
REGISTRY="${REGISTRY:-ghcr.io/urandomdev}"
IMAGE_NAME="${IMAGE_NAME:-regulation}"
PLATFORMS="${PLATFORMS:-linux/amd64,linux/arm64}"
PUSH="${PUSH:-false}"

# Get version from git
get_version() {
    if git rev-parse --is-inside-work-tree > /dev/null 2>&1; then
        # Check if we're on a tag
        if git describe --exact-match --tags HEAD > /dev/null 2>&1; then
            git describe --tags --exact-match
        else
            # Use branch name and short commit hash
            local branch=$(git rev-parse --abbrev-ref HEAD)
            local commit=$(git rev-parse --short HEAD)
            echo "${branch}-${commit}"
        fi
    else
        echo "unknown"
    fi
}

VERSION=$(get_version)
FULL_IMAGE="${REGISTRY}/${IMAGE_NAME}"

echo -e "${GREEN}Building regulation API${NC}"
echo "Version: ${VERSION}"
echo "Registry: ${REGISTRY}"
echo "Image: ${FULL_IMAGE}"
echo "Platforms: ${PLATFORMS}"
echo ""

# Check if ko is installed
if ! command -v ko &> /dev/null; then
    echo -e "${RED}Error: 'ko' is not installed.${NC}"
    echo "Please install ko: https://ko.build/install/"
    echo "  brew install ko"
    echo "  or"
    echo "  go install github.com/google/ko@latest"
    exit 1
fi

# Change to api directory
cd api

# Build tags
TAGS=("${VERSION}")

# If it's a semver tag, add additional tags
if [[ "${VERSION}" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    # Extract major, minor, patch
    MAJOR=$(echo "${VERSION}" | cut -d. -f1)
    MINOR=$(echo "${VERSION}" | cut -d. -f2 | sed 's/v//')
    PATCH=$(echo "${VERSION}" | cut -d. -f3)

    TAGS+=("${MAJOR}.${MINOR}.${PATCH}")
    TAGS+=("${MAJOR}.${MINOR}")
    TAGS+=("${MAJOR}")
    TAGS+=("latest")
fi

# Build tag arguments
TAG_ARGS=""
for tag in "${TAGS[@]}"; do
    TAG_ARGS="${TAG_ARGS} -t ${tag}"
done

# Set environment variables for ko
export KO_DOCKER_REPO="${FULL_IMAGE}"

# Build command
BUILD_CMD="ko build --platform=${PLATFORMS} --bare ${TAG_ARGS} ./cmd/regulation"

echo -e "${YELLOW}Running: ${BUILD_CMD}${NC}"
echo ""

# Execute build
eval "${BUILD_CMD}"

if [ $? -eq 0 ]; then
    echo ""
    echo -e "${GREEN}✓ Build successful!${NC}"
    echo ""
    echo "Built images:"
    for tag in "${TAGS[@]}"; do
        echo "  - ${FULL_IMAGE}:${tag}"
    done
else
    echo -e "${RED}✗ Build failed!${NC}"
    exit 1
fi
