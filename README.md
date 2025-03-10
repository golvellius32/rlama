# RLAMA - User Guide
RLAMA is a powerful AI-driven question-answering tool for your documents, seamlessly integrating with your local Ollama models. It enables you to create, manage, and interact with Retrieval-Augmented Generation (RAG) systems tailored to your documentation needs.


[![RLAMA Demonstration](https://img.youtube.com/vi/EIsQnBqeQxQ/0.jpg)](https://www.youtube.com/watch?v=EIsQnBqeQxQ)

## Table of Contents
- [Installation](#installation)
- [Available Commands](#available-commands)
  - [rag - Create a RAG system](#rag---create-a-rag-system)
  - [run - Use a RAG system](#run---use-a-rag-system)
  - [list - List RAG systems](#list---list-rag-systems)
  - [delete - Delete a RAG system](#delete---delete-a-rag-system)
  - [update - Update RLAMA](#update---update-rlama)
  - [version - Display version](#version---display-version)
- [Uninstallation](#uninstallation)
- [Supported Document Formats](#supported-document-formats)
- [Troubleshooting](#troubleshooting)

## Installation

### Prerequisites
- [Ollama](https://ollama.ai/) installed and running

### Installation from terminal

```bash
curl -fsSL https://raw.githubusercontent.com/dontizi/rlama/main/install.sh | sh
```


## Available Commands

You can get help on all commands by using:

```bash
rlama --help
```

### rag - Create a RAG system

Creates a new RAG system by indexing all documents in the specified folder.

```bash
rlama rag [model] [rag-name] [folder-path]
```

**Parameters:**
- `model`: Name of the Ollama model to use (e.g., llama3, mistral, gemma).
- `rag-name`: Unique name to identify your RAG system.
- `folder-path`: Path to the folder containing your documents.

**Example:**

```bash
rlama rag llama3 documentation ./docs
```

### run - Use a RAG system

Starts an interactive session to interact with an existing RAG system.

```bash
rlama run [rag-name]
```

**Parameters:**
- `rag-name`: Name of the RAG system to use.

**Example:**

```bash
rlama run documentation
> How do I install the project?
> What are the main features?
> exit
```

### list - List RAG systems

Displays a list of all available RAG systems.

```bash
rlama list
```

### delete - Delete a RAG system

Permanently deletes a RAG system and all its indexed documents.

```bash
rlama delete [rag-name] [--force/-f]
```

**Parameters:**
- `rag-name`: Name of the RAG system to delete.
- `--force` or `-f`: (Optional) Delete without asking for confirmation.

**Example:**

```bash
rlama delete old-project
```

Or to delete without confirmation:

```bash
rlama delete old-project --force
```

### update - Update RLAMA

Checks if a new version of RLAMA is available and installs it.

```bash
rlama update [--force/-f]
```

**Options:**
- `--force` or `-f`: (Optional) Update without asking for confirmation.

### version - Display version

Displays the current version of RLAMA.

```bash
rlama --version
```

or

```bash
rlama -v
```

## Uninstallation

To uninstall RLAMA:

### Removing the binary

If you installed via `go install`:

```bash
rlama uninstall
```

### Removing data

RLAMA stores its data in `~/.rlama`. To remove it:

```bash
rm -rf ~/.rlama
```

## Supported Document Formats

RLAMA supports many file formats:

- **Text**: `.txt`, `.md`, `.html`, `.json`, `.csv`, `.yaml`, `.yml`, `.xml`
- **Code**: `.go`, `.py`, `.js`, `.java`, `.c`, `.cpp`, `.h`, `.rb`, `.php`, `.rs`, `.swift`, `.kt`
- **Documents**: `.pdf`, `.docx`, `.doc`, `.rtf`, `.odt`, `.pptx`, `.ppt`, `.xlsx`, `.xls`, `.epub`

Installing dependencies via `install_deps.sh` is recommended to improve support for certain formats.

## Troubleshooting

### Ollama is not accessible

If you encounter connection errors to Ollama:
1. Check that Ollama is running.
2. Ollama must be accessible at `http://localhost:11434`.
3. Check Ollama logs for potential errors.

### Text extraction issues

If you encounter problems with certain formats:
1. Install dependencies via `./scripts/install_deps.sh`.
2. Verify that your system has the required tools (`pdftotext`, `tesseract`, etc.).

### The RAG doesn't find relevant information

If the answers are not relevant:
1. Check that the documents are properly indexed with `rlama list`.
2. Make sure the content of the documents is properly extracted.
3. Try rephrasing your question more precisely.

### Other issues

For any other issues, please open an issue on the [GitHub repository](https://github.com/golvellius32/rlama/issues) providing:
1. The exact command used.
2. The complete output of the command.
3. Your operating system and architecture.
4. The RLAMA version (`rlama --version`).
