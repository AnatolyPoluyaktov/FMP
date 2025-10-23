#!/bin/bash

# FMP Deployment Script
# This script deploys the Financial Manager Platform using Ansible

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
INVENTORY_FILE="inventory.yml"
PLAYBOOK_DIR="."

# Functions
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

check_requirements() {
    log_info "Checking requirements..."
    
    if ! command -v ansible &> /dev/null; then
        log_error "Ansible is not installed. Please install it first."
        exit 1
    fi
    
    if ! command -v ansible-playbook &> /dev/null; then
        log_error "ansible-playbook is not installed. Please install it first."
        exit 1
    fi
    
    if [ ! -f "$INVENTORY_FILE" ]; then
        log_error "Inventory file $INVENTORY_FILE not found."
        exit 1
    fi
    
    log_info "Requirements check passed."
}

deploy_core() {
    log_info "Deploying FMP Core API..."
    ansible-playbook -i "$INVENTORY_FILE" "$PLAYBOOK_DIR/fmp-core.yml" --ask-become-pass
    log_info "FMP Core API deployment completed."
}

deploy_minapp() {
    log_info "Deploying FMP Mini App..."
    ansible-playbook -i "$INVENTORY_FILE" "$PLAYBOOK_DIR/fmp-minapp.yml" --ask-become-pass
    log_info "FMP Mini App deployment completed."
}

deploy_analytics() {
    log_info "Deploying FMP Analytics..."
    ansible-playbook -i "$INVENTORY_FILE" "$PLAYBOOK_DIR/fmp-analytics.yml" --ask-become-pass
    log_info "FMP Analytics deployment completed."
}

deploy_all() {
    log_info "Deploying all FMP components..."
    deploy_core
    deploy_minapp
    deploy_analytics
    log_info "All deployments completed."
}

# Main script
main() {
    log_info "Starting FMP deployment..."
    
    check_requirements
    
    case "${1:-all}" in
        "core")
            deploy_core
            ;;
        "minapp")
            deploy_minapp
            ;;
        "analytics")
            deploy_analytics
            ;;
        "all")
            deploy_all
            ;;
        *)
            log_error "Unknown deployment target: $1"
            log_info "Usage: $0 [core|minapp|analytics|all]"
            exit 1
            ;;
    esac
    
    log_info "Deployment completed successfully!"
}

# Run main function with all arguments
main "$@"
