#!/bin/bash

# Script d'installation des dépendances pour RLAMA
# Ce script tente d'installer les outils nécessaires pour l'extraction de texte

echo "Installation des dépendances pour RLAMA..."

# Détection du système d'exploitation
OS=$(uname -s)
echo "Système d'exploitation détecté: $OS"

# Fonction pour vérifier si un programme est installé
is_installed() {
  command -v "$1" >/dev/null 2>&1
}

# macOS
if [ "$OS" = "Darwin" ]; then
  echo "Installation des dépendances pour macOS..."
  
  # Vérifier si Homebrew est installé
  if ! is_installed brew; then
    echo "Homebrew non trouvé. Installation de Homebrew..."
    /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
  fi
  
  # Installer les outils
  echo "Installation des outils d'extraction de texte..."
  brew install poppler  # Pour pdftotext
  brew install tesseract  # Pour OCR
  brew install tesseract-lang  # Langues supplémentaires pour Tesseract
  
  # Python et outils
  if ! is_installed pip3; then
    brew install python
  fi
  
  pip3 install pdfminer.six docx2txt xlsx2csv
  
# Linux
elif [ "$OS" = "Linux" ]; then
  echo "Installation des dépendances pour Linux..."
  
  # Essayer apt-get (Debian/Ubuntu)
  if is_installed apt-get; then
    echo "Gestionnaire de paquets apt-get détecté"
    sudo apt-get update
    sudo apt-get install -y poppler-utils tesseract-ocr tesseract-ocr-fra python3-pip
    sudo apt-get install -y catdoc unrtf
  
  # Essayer yum (CentOS/RHEL)
  elif is_installed yum; then
    echo "Gestionnaire de paquets yum détecté"
    sudo yum update
    sudo yum install -y poppler-utils tesseract tesseract-langpack-fra python3-pip
    sudo yum install -y catdoc
  
  # Essayer pacman (Arch Linux)
  elif is_installed pacman; then
    echo "Gestionnaire de paquets pacman détecté"
    sudo pacman -Syu
    sudo pacman -S poppler tesseract tesseract-data-fra python-pip
  
  # Essayer zypper (openSUSE)
  elif is_installed zypper; then
    echo "Gestionnaire de paquets zypper détecté"
    sudo zypper refresh
    sudo zypper install poppler-tools tesseract-ocr python3-pip
  
  else
    echo "Aucun gestionnaire de paquets connu détecté. Veuillez installer manuellement les dépendances."
  fi
  
  # Installer les packages Python
  pip3 install --user pdfminer.six docx2txt xlsx2csv

# Windows (via WSL)
elif [[ "$OS" == MINGW* ]] || [[ "$OS" == MSYS* ]] || [[ "$OS" == CYGWIN* ]]; then
  echo "Système Windows détecté."
  echo "Il est recommandé d'utiliser WSL (Windows Subsystem for Linux) pour de meilleures performances."
  echo "Vous pouvez installer les dépendances manuellement:"
  echo "1. Installez Python: https://www.python.org/downloads/windows/"
  echo "2. Installez les packages Python: pip install pdfminer.six docx2txt xlsx2csv"
  echo "3. Pour l'OCR, installez Tesseract: https://github.com/UB-Mannheim/tesseract/wiki"
fi

echo "Installation des dépendances Python..."
pip3 install --user pdfminer.six docx2txt xlsx2csv

echo "Installation terminée!" 