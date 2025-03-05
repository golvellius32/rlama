#!/bin/bash

set -e

# Couleurs pour les messages
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

echo "
█▀█ █   ▄▀█ █▀▄▀█ ▄▀█
█▀▄ █▄▄ █▀█ █░▀░█ █▀█

Retrieval-Augmented Language Model Adapter
"

# Fonction pour vérifier si une commande existe
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Vérifier si Go est installé
if ! command_exists go; then
    echo -e "${RED}Go n'est pas installé.${NC}"
    echo "RLAMA nécessite Go pour être compilé."
    echo "Installez Go depuis https://golang.org/dl/"
    exit 1
fi

# Vérifier si Ollama est installé
if ! command_exists ollama; then
    echo -e "${YELLOW}⚠️ Ollama n'est pas installé.${NC}"
    echo "RLAMA nécessite Ollama pour fonctionner."
    echo "Vous pouvez installer Ollama avec:"
    echo "curl -fsSL https://ollama.com/install.sh | sh"
    
    read -p "Voulez-vous installer Ollama maintenant? (o/n): " install_ollama
    if [[ "$install_ollama" =~ ^[Oo]$ ]]; then
        echo "Installation d'Ollama..."
        curl -fsSL https://ollama.com/install.sh | sh
    else
        echo "Veuillez installer Ollama avant d'utiliser RLAMA."
    fi
fi

# Vérifier si Ollama est en cours d'exécution
if ! curl -s http://localhost:11434/api/version &>/dev/null; then
    echo -e "${YELLOW}⚠️ Le service Ollama ne semble pas fonctionner.${NC}"
    echo "Veuillez démarrer Ollama avant d'utiliser RLAMA."
fi

# Vérifier si le modèle llama3 est disponible
if command_exists ollama; then
    if ! ollama list 2>/dev/null | grep -q "llama3"; then
        echo -e "${YELLOW}⚠️ Le modèle llama3 n'est pas disponible dans Ollama.${NC}"
        echo "Pour une meilleure expérience, vous devriez l'installer avec:"
        echo "ollama pull llama3"
    fi
fi

# Créer le répertoire d'installation
INSTALL_DIR="/usr/local/bin"
DATA_DIR="$HOME/.rlama"

echo "Installation de RLAMA..."
echo "Clonage du dépôt..."

# Utiliser un répertoire temporaire pour le clonage et la compilation
TEMP_DIR=$(mktemp -d)
cd "$TEMP_DIR"

# Cloner le dépôt RLAMA (à remplacer par votre URL de dépôt)
git clone https://github.com/dontizi/rlama.git .

# Compilation
echo "Compilation de RLAMA..."
go build -o rlama

# Installation
echo "Installation de l'exécutable..."
mkdir -p "$DATA_DIR"

if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS - utiliser sudo si nécessaire pour /usr/local/bin
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

# Nettoyage
cd "$HOME"
rm -rf "$TEMP_DIR"

echo -e "${GREEN}✅ RLAMA a été installé avec succès!${NC}"
echo ""
echo "Vous pouvez maintenant utiliser RLAMA avec les commandes suivantes:"
echo "- rlama rag [modèle] [nom-rag] [chemin-dossier] : Créer un nouveau système RAG"
echo "- rlama run [nom-rag] : Exécuter un système RAG"
echo "- rlama list : Lister tous les systèmes RAG disponibles"
echo "- rlama delete [nom-rag] : Supprimer un système RAG"
echo ""
echo "Exemple: rlama rag llama3 monrag ./documents" 