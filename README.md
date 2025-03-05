# RLAMA - Guide d'utilisation

## Table des matières
- [Installation](#installation)
- [Commandes disponibles](#commandes-disponibles)
  - [rag - Créer un système RAG](#rag---créer-un-système-rag)
  - [run - Utiliser un système RAG](#run---utiliser-un-système-rag)
  - [list - Lister les systèmes RAG](#list---lister-les-systèmes-rag)
  - [delete - Supprimer un système RAG](#delete---supprimer-un-système-rag)
  - [update - Mettre à jour RLAMA](#update---mettre-à-jour-rlama)
  - [version - Afficher la version](#version---afficher-la-version)
- [Désinstallation](#désinstallation)
- [Formats de documents supportés](#formats-de-documents-supportés)
- [Dépannage](#dépannage)

## Installation

### Prérequis
- [Go](https://golang.org/doc/install) 1.21 ou supérieur
- [Ollama](https://ollama.ai/) installé et en cours d'exécution

### Installation depuis GitHub

```bash
# Cloner le dépôt
git clone https://github.com/dontizi/rlama.git
cd rlama

# Compiler et installer
go install

# Installer les dépendances (optionnel mais recommandé)
chmod +x scripts/install_deps.sh
./scripts/install_deps.sh
```

### Installation depuis les releases

1. Visitez la [page des releases](https://github.com/dontizi/rlama/releases) sur GitHub.
2. Téléchargez le binaire correspondant à votre système d'exploitation et architecture.
3. Extrayez le binaire et placez-le dans un dossier de votre `PATH`.

```bash
# Exemple pour Linux/macOS
chmod +x rlama
sudo mv rlama /usr/local/bin/
```

## Commandes disponibles

Vous pouvez obtenir de l'aide sur toutes les commandes en utilisant:

```bash
rlama --help
```

### rag - Créer un système RAG

Crée un nouveau système RAG en indexant tous les documents du dossier spécifié.

```bash
rlama rag [modèle] [nom-rag] [chemin-dossier]
```

**Paramètres:**
- `modèle`: Nom du modèle Ollama à utiliser (ex: llama3, mistral, gemma).
- `nom-rag`: Nom unique pour identifier votre système RAG.
- `chemin-dossier`: Chemin vers le dossier contenant vos documents.

**Exemple:**

```bash
rlama rag llama3 documentation ./docs
```

### run - Utiliser un système RAG

Démarre une session interactive pour interagir avec un système RAG existant.

```bash
rlama run [nom-rag]
```

**Paramètres:**
- `nom-rag`: Nom du système RAG à utiliser.

**Exemple:**

```bash
rlama run documentation
> Comment installer le projet?
> Quelles sont les fonctionnalités principales?
> exit
```

### list - Lister les systèmes RAG

Affiche la liste de tous les systèmes RAG disponibles.

```bash
rlama list
```

### delete - Supprimer un système RAG

Supprime définitivement un système RAG et tous ses documents indexés.

```bash
rlama delete [nom-rag] [--force/-f]
```

**Paramètres:**
- `nom-rag`: Nom du système RAG à supprimer.
- `--force` ou `-f`: (Optionnel) Supprimer sans demander de confirmation.

**Exemple:**

```bash
rlama delete ancien-projet
```

Ou pour supprimer sans confirmation:

```bash
rlama delete ancien-projet --force
```

### update - Mettre à jour RLAMA

Vérifie si une nouvelle version de RLAMA est disponible et l'installe.

```bash
rlama update [--force/-f]
```

**Options:**
- `--force` ou `-f`: (Optionnel) Mettre à jour sans demander de confirmation.

### version - Afficher la version

Affiche la version actuelle de RLAMA.

```bash
rlama --version
```

ou

```bash
rlama -v
```

## Désinstallation

Pour désinstaller RLAMA:

### Suppression du binaire

Si vous avez installé via `go install`:

```bash
rlama uninstall
```

### Suppression des données

RLAMA stocke ses données dans `~/.rlama`. Pour les supprimer:

```bash
rm -rf ~/.rlama
```

## Formats de documents supportés

RLAMA prend en charge de nombreux formats de fichiers:

- **Texte**: `.txt`, `.md`, `.html`, `.json`, `.csv`, `.yaml`, `.yml`, `.xml`
- **Code**: `.go`, `.py`, `.js`, `.java`, `.c`, `.cpp`, `.h`, `.rb`, `.php`, `.rs`, `.swift`, `.kt`
- **Documents**: `.pdf`, `.docx`, `.doc`, `.rtf`, `.odt`, `.pptx`, `.ppt`, `.xlsx`, `.xls`, `.epub`

L'installation des dépendances via `install_deps.sh` est recommandée pour améliorer le support de certains formats.

## Dépannage

### Ollama n'est pas accessible

Si vous rencontrez des erreurs de connexion à Ollama:
1. Vérifiez qu'Ollama est en cours d'exécution.
2. Ollama doit être accessible à l'adresse `http://localhost:11434`.
3. Vérifiez les logs d'Ollama pour d'éventuelles erreurs.

### Problèmes d'extraction de texte

Si vous rencontrez des problèmes avec certains formats:
1. Installez les dépendances via `./scripts/install_deps.sh`.
2. Vérifiez que votre système possède les outils requis (`pdftotext`, `tesseract`, etc.).

### Le RAG ne trouve pas d'informations pertinentes

Si les réponses ne sont pas pertinentes:
1. Vérifiez que les documents sont bien indexés avec `rlama list`.
2. Assurez-vous que le contenu des documents est bien extrait.
3. Essayez de reformuler votre question de manière plus précise.

### Autres problèmes

Pour tout autre problème, veuillez ouvrir une issue sur le [dépôt GitHub](https://github.com/dontizi/rlama/issues) en fournissant:
1. La commande exacte utilisée.
2. La sortie complète de la commande.
3. Votre système d'exploitation et architecture.
4. La version de RLAMA (`rlama --version`).

