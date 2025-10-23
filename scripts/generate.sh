#!/bin/bash

# FMP Code Generation Script
# Генерирует Go код из OpenAPI спецификации

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
API_SPEC_FILE="api-spec.yaml"
GENERATED_DIR="generated"
FMP_CORE_DIR="fmp-core"
MINAPP_BACKEND_DIR="minapp/backend"

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if API spec exists
if [ ! -f "$API_SPEC_FILE" ]; then
    log_error "API specification file $API_SPEC_FILE not found!"
    exit 1
fi

log_info "Starting code generation from $API_SPEC_FILE..."

# Create generated directory
mkdir -p "$GENERATED_DIR"

# Generate models
log_info "Generating models..."
oapi-codegen -generate types -package models -o "$GENERATED_DIR/models.go" "$API_SPEC_FILE"

# Generate server interface
log_info "Generating server interface..."
oapi-codegen -generate gin -package api -o "$GENERATED_DIR/server.go" "$API_SPEC_FILE"

# Generate client
log_info "Generating client..."
oapi-codegen -generate client -package client -o "$GENERATED_DIR/client.go" "$API_SPEC_FILE"

# Copy generated files to fmp-core
log_info "Copying generated files to fmp-core..."
cp "$GENERATED_DIR/models.go" "$FMP_CORE_DIR/internal/generated/"
cp "$GENERATED_DIR/server.go" "$FMP_CORE_DIR/internal/generated/"

# Copy generated files to minapp backend
log_info "Copying generated files to minapp backend..."
cp "$GENERATED_DIR/models.go" "$MINAPP_BACKEND_DIR/internal/generated/"
cp "$GENERATED_DIR/client.go" "$MINAPP_BACKEND_DIR/internal/generated/"

# Generate Swagger documentation
log_info "Generating Swagger documentation..."
swag init -g "$FMP_CORE_DIR/main.go" -o "$FMP_CORE_DIR/docs"

log_info "Code generation completed successfully!"
log_info "Generated files:"
log_info "  - Models: $GENERATED_DIR/models.go"
log_info "  - Server: $GENERATED_DIR/server.go"
log_info "  - Client: $GENERATED_DIR/client.go"
log_info "  - Swagger docs: $FMP_CORE_DIR/docs/"

log_warn "Don't forget to:"
log_warn "  1. Implement the generated server interface"
log_warn "  2. Update your handlers to use generated models"
log_warn "  3. Run 'go mod tidy' in both fmp-core and minapp/backend"
