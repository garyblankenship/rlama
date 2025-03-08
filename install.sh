#!/bin/bash

set -e

# Colors for messages
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

echo "
█▀█ █   ▄▀█ █▀▄▀█ ▄▀█
█▀▄ █▄▄ █▀█ █░▀░█ █▀█

Retrieval-Augmented Language Model Adapter
"

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Determine OS and architecture
get_os_arch() {
    local os=$(uname -s | tr '[:upper:]' '[:lower:]')
    local arch=$(uname -m)
    
    # Convert architecture to standard format
    case "$arch" in
        x86_64)
            arch="amd64"
            ;;
        aarch64|arm64)
            arch="arm64"
            ;;
        *)
            echo -e "${RED}Unsupported architecture: $arch${NC}"
            exit 1
            ;;
    esac
    
    # Handle macOS naming
    if [ "$os" = "darwin" ]; then
        os="darwin"
    elif [ "$os" = "linux" ]; then
        os="linux"
    else
        echo -e "${RED}Unsupported operating system: $os${NC}"
        exit 1
    fi
    
    echo "${os}_${arch}"
}

# Check if Ollama is installed
if ! command_exists ollama; then
    echo -e "${YELLOW}⚠️ Ollama is not installed.${NC}"
    echo "RLAMA requires Ollama to function."
    echo "You can install Ollama with:"
    echo "curl -fsSL https://ollama.com/install.sh | sh"
    
    read -p "Do you want to install Ollama now? (y/n): " install_ollama
    if [[ "$install_ollama" =~ ^[Yy]$ ]]; then
        echo "Installing Ollama..."
        curl -fsSL https://ollama.com/install.sh | sh
    else
        echo "Please install Ollama before using RLAMA."
    fi
fi

# Check if Ollama is running
if ! curl -s http://localhost:11434/api/version &>/dev/null; then
    echo -e "${YELLOW}⚠️ The Ollama service doesn't seem to be running.${NC}"
    echo "Please start Ollama before using RLAMA."
fi

# Check if the llama3 model is available
if command_exists ollama; then
    if ! ollama list 2>/dev/null | grep -q "llama3"; then
        echo -e "${YELLOW}⚠️ The llama3 model is not available in Ollama.${NC}"
        echo "For a better experience, you should install it with:"
        echo "ollama pull llama3"
    fi
fi

# Create installation directory
INSTALL_DIR="/usr/local/bin"
DATA_DIR="$HOME/.rlama"

echo "Installing RLAMA..."

# Determine OS and architecture for downloading the correct binary
OS_ARCH=$(get_os_arch)
BINARY_NAME="rlama_${OS_ARCH}"
DOWNLOAD_URL="https://github.com/yibie/rlama/releases/latest/download/${BINARY_NAME}"

echo "Downloading RLAMA for $OS_ARCH..."

# Create a temporary directory for downloading
TEMP_DIR=$(mktemp -d)
cd "$TEMP_DIR"

# Download the binary
if command_exists curl; then
    curl -L -o rlama $DOWNLOAD_URL
elif command_exists wget; then
    wget -O rlama $DOWNLOAD_URL
else
    echo -e "${RED}Error: Neither curl nor wget is installed.${NC}"
    exit 1
fi

# Make it executable
chmod +x rlama

# Install
echo "Installing executable..."
mkdir -p "$DATA_DIR"

# Try to install to /usr/local/bin, fall back to ~/.local/bin if permission denied
if [ -w "$INSTALL_DIR" ]; then
    mv rlama "$INSTALL_DIR/"
else
    echo "Cannot write to $INSTALL_DIR, trying alternative location..."
    LOCAL_BIN="$HOME/.local/bin"
    mkdir -p "$LOCAL_BIN"
    mv rlama "$LOCAL_BIN/"
    
    # Add to PATH if not already there
    if [[ ":$PATH:" != *":$LOCAL_BIN:"* ]]; then
        echo "Adding $LOCAL_BIN to your PATH..."
        echo 'export PATH="$HOME/.local/bin:$PATH"' >> "$HOME/.bashrc"
        echo 'export PATH="$HOME/.local/bin:$PATH"' >> "$HOME/.zshrc" 2>/dev/null || true
        export PATH="$HOME/.local/bin:$PATH"
    fi
    
    INSTALL_DIR="$LOCAL_BIN"
fi

# Clean up
cd - > /dev/null
rm -rf "$TEMP_DIR"

echo -e "${GREEN}RLAMA has been successfully installed to $INSTALL_DIR/rlama!${NC}"
echo "You can now use RLAMA by running the 'rlama' command."
echo "Run 'rlama --help' to see available commands."
echo ""
echo "You can also use RLAMA with the following commands:"
echo "- rlama rag [model] [rag-name] [folder-path] : Create a new RAG system"
echo "- rlama run [rag-name] : Run a RAG system"
echo "- rlama list : List all available RAG systems"
echo "- rlama delete [rag-name] : Delete a RAG system"
echo ""
echo "Example: rlama rag llama3 myrag ./documents" 
