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

# Check if Go is installed
if ! command_exists go; then
    echo -e "${RED}Go is not installed.${NC}"
    echo "RLAMA requires Go to be compiled."
    echo "Install Go from https://golang.org/dl/"
    exit 1
fi

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
echo "Cloning repository..."

# Use a temporary directory for cloning and building
TEMP_DIR=$(mktemp -d)
cd "$TEMP_DIR"

# Clone the RLAMA repository (replace with your repository URL)
git clone https://github.com/dontizi/rlama.git .

# Build
echo "Building RLAMA..."
go build -o rlama

# Install
echo "Installing executable..."
mkdir -p "$DATA_DIR"

if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS - use sudo if necessary for /usr/local/bin
    if [ -w "$INSTALL_DIR" ]; then
        cp rlama "$INSTALL_DIR/"
    else
        sudo cp rlama "$INSTALL_DIR/"
    fi
    chmod +x "$INSTALL_DIR/rlama"
else
    # Linux
    sudo cp rlama "$INSTALL_DIR/"
    sudo chmod +x "$INSTALL_DIR/rlama"
fi

# Cleanup
cd "$HOME"
rm -rf "$TEMP_DIR"

echo -e "${GREEN}✅ RLAMA has been successfully installed!${NC}"
echo ""
echo "You can now use RLAMA with the following commands:"
echo "- rlama rag [model] [rag-name] [folder-path] : Create a new RAG system"
echo "- rlama run [rag-name] : Run a RAG system"
echo "- rlama list : List all available RAG systems"
echo "- rlama delete [rag-name] : Delete a RAG system"
echo ""
echo "Example: rlama rag llama3 myrag ./documents" 