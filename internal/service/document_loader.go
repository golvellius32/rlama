package service

import (
	// Suppression des imports non utilisés
	// "bytes"
	// "encoding/json"
	"fmt"
	// "io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/dontizi/rlama/internal/domain"
)

// DocumentLoader se charge de charger les documents depuis le système de fichiers
type DocumentLoader struct {
	supportedExtensions map[string]bool
	extractorPath       string // Chemin vers l'extracteur externe
}

// NewDocumentLoader crée une nouvelle instance de DocumentLoader
func NewDocumentLoader() *DocumentLoader {
	return &DocumentLoader{
		supportedExtensions: map[string]bool{
			// Texte simple
			".txt":   true,
			".md":    true,
			".html":  true,
			".htm":   true,
			".json":  true,
			".csv":   true,
			".log":   true,
			".xml":   true,
			".yaml":  true,
			".yml":   true,
			// Code source
			".go":    true,
			".py":    true,
			".js":    true,
			".java":  true,
			".c":     true,
			".cpp":   true,
			".h":     true,
			".rb":    true,
			".php":   true,
			".rs":    true,
			".swift": true,
			".kt":    true,
			// Documents
			".pdf":   true,
			".docx":  true,
			".doc":   true,
			".rtf":   true,
			".odt":   true,
			".pptx":  true,
			".ppt":   true,
			".xlsx":  true,
			".xls":   true,
			".epub":  true,
		},
		// On utilisera pdftotext si disponible
		extractorPath: findExternalExtractor(),
	}
}

// findExternalExtractor cherche des outils d'extraction externes
func findExternalExtractor() string {
	// Priorité des extracteurs de texte
	extractors := []string{
		"pdftotext", // Pour les PDF (Poppler-utils)
		"textutil",  // macOS
		"catdoc",    // Pour .doc
		"unrtf",     // Pour .rtf
	}

	for _, extractor := range extractors {
		path, err := exec.LookPath(extractor)
		if err == nil {
			fmt.Printf("Extracteur externe trouvé: %s\n", path)
			return path
		}
	}

	fmt.Println("Aucun extracteur externe trouvé. L'extraction de texte sera limitée.")
	return ""
}

// LoadDocumentsFromFolder charge tous les documents supportés du dossier spécifié
func (dl *DocumentLoader) LoadDocumentsFromFolder(folderPath string) ([]*domain.Document, error) {
	var documents []*domain.Document
	var supportedFiles []string
	var unsupportedFiles []string

	// Vérifie si le dossier existe
	info, err := os.Stat(folderPath)
	if os.IsNotExist(err) {
		// Essayer de créer le dossier
		if err := os.MkdirAll(folderPath, 0755); err != nil {
			return nil, fmt.Errorf("le dossier '%s' n'existe pas et ne peut pas être créé: %w", folderPath, err)
		}
		fmt.Printf("Le dossier '%s' a été créé.\n", folderPath)
		// Récupérer les informations du dossier nouvellement créé
		info, err = os.Stat(folderPath)
		if err != nil {
			return nil, fmt.Errorf("impossible d'accéder au dossier '%s': %w", folderPath, err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("impossible d'accéder au dossier '%s': %w", folderPath, err)
	}
	
	if !info.IsDir() {
		return nil, fmt.Errorf("le chemin spécifié n'est pas un dossier: %s", folderPath)
	}

	// Vérification préliminaire des fichiers
	err = filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Ignore les dossiers et les fichiers cachés (commençant par .)
		if info.IsDir() || strings.HasPrefix(filepath.Base(path), ".") {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if dl.supportedExtensions[ext] {
			supportedFiles = append(supportedFiles, path)
		} else {
			unsupportedFiles = append(unsupportedFiles, path)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("erreur lors de l'analyse du dossier: %w", err)
	}

	// Affiche des infos sur les fichiers trouvés
	if len(supportedFiles) == 0 {
		if len(unsupportedFiles) == 0 {
			return nil, fmt.Errorf("le dossier '%s' est vide. Veuillez y ajouter des documents avant de créer un RAG", folderPath)
		} else {
			extensionsMsg := "Extensions supportées: "
			for ext := range dl.supportedExtensions {
				extensionsMsg += ext + " "
			}
			return nil, fmt.Errorf("aucun fichier supporté trouvé dans '%s'. %d fichiers non supportés détectés.\n%s", 
				folderPath, len(unsupportedFiles), extensionsMsg)
		}
	}

	fmt.Printf("Trouvé %d fichiers supportés et %d fichiers non supportés.\n", len(supportedFiles), len(unsupportedFiles))
	
	// Tentative d'installation des dépendances si possible
	dl.tryInstallDependencies()
	
	// Traiter les fichiers supportés
	for _, path := range supportedFiles {
		ext := strings.ToLower(filepath.Ext(path))
		
		// Extraction du texte à l'aide de plusieurs méthodes
		textContent, err := dl.extractText(path, ext)
		if err != nil {
			fmt.Printf("Avertissement: impossible d'extraire le texte de %s: %v\n", path, err)
			fmt.Println("Tentative d'extraction en tant que texte brut...")
			
			// Tentative de lecture en tant que fichier texte
			rawContent, err := ioutil.ReadFile(path)
			if err != nil {
				fmt.Printf("Échec de la lecture brute de %s: %v\n", path, err)
				continue
			}
			
			textContent = string(rawContent)
		}

		// Vérifier que le contenu n'est pas vide
		if strings.TrimSpace(textContent) == "" {
			fmt.Printf("Avertissement: aucun texte extrait de %s\n", path)
			
			// Pour les PDF, tenter une dernière méthode
			if ext == ".pdf" {
				fmt.Println("Tentative d'extraction avec OCR (si installé)...")
				ocrText, err := dl.extractWithOCR(path)
				if err != nil || strings.TrimSpace(ocrText) == "" {
					fmt.Println("Échec de l'OCR ou non disponible.")
					continue
				}
				textContent = ocrText
			} else {
				continue
			}
		}

		// Créer un document
		doc := domain.NewDocument(path, textContent)
		documents = append(documents, doc)
		fmt.Printf("Document ajouté: %s (%d caractères)\n", filepath.Base(path), len(textContent))
	}

	if len(documents) == 0 {
		return nil, fmt.Errorf("aucun document avec contenu valide trouvé dans le dossier '%s'", folderPath)
	}

	return documents, nil
}

// extractText extrait le texte d'un fichier en utilisant la méthode appropriée selon le type
func (dl *DocumentLoader) extractText(path string, ext string) (string, error) {
	switch ext {
	case ".pdf":
		return dl.extractFromPDF(path)
	case ".docx", ".doc", ".rtf", ".odt":
		return dl.extractFromDocument(path, ext)
	case ".pptx", ".ppt":
		return dl.extractFromPresentation(path, ext)
	case ".xlsx", ".xls":
		return dl.extractFromSpreadsheet(path, ext)
	default:
		// Traiter comme un fichier texte
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return "", err
		}
		return string(data), nil
	}
}

// extractFromPDF extrait le texte d'un PDF en utilisant différentes méthodes
func (dl *DocumentLoader) extractFromPDF(path string) (string, error) {
	// Méthode 1: Utiliser pdftotext si disponible
	if strings.Contains(dl.extractorPath, "pdftotext") {
		fmt.Printf("Extraction du PDF avec pdftotext: %s\n", filepath.Base(path))
		out, err := exec.Command(dl.extractorPath, "-layout", path, "-").Output()
		if err == nil && len(out) > 0 {
			return string(out), nil
		}
		fmt.Printf("pdftotext a échoué: %v\n", err)
	}
	
	// Méthode 2: Tentative avec d'autres outils (pdfinfo, pdftk)
	for _, tool := range []string{"pdfinfo", "pdftk"} {
		toolPath, err := exec.LookPath(tool)
		if err == nil {
			fmt.Printf("Tentative d'extraction avec %s\n", tool)
			var out []byte
			if tool == "pdfinfo" {
				out, err = exec.Command(toolPath, path).Output()
			} else {
				out, err = exec.Command(toolPath, path, "dump_data").Output()
			}
			if err == nil && len(out) > 0 {
				return string(out), nil
			}
		}
	}
	
	// Méthode 3: Dernière tentative, ouvrir comme fichier binaire et extraire les chaînes
	fmt.Println("Extraction des chaînes de caractères du PDF...")
	return dl.extractStringsFromBinary(path)
}

// extractFromDocument extrait le texte d'un document Word ou similaire
func (dl *DocumentLoader) extractFromDocument(path string, ext string) (string, error) {
	// Méthode 1: Utiliser textutil sur macOS
	if strings.Contains(dl.extractorPath, "textutil") && (ext == ".docx" || ext == ".doc" || ext == ".rtf") {
		fmt.Printf("Extraction du document avec textutil: %s\n", filepath.Base(path))
		out, err := exec.Command(dl.extractorPath, "-convert", "txt", "-stdout", path).Output()
		if err == nil && len(out) > 0 {
			return string(out), nil
		}
	}
	
	// Méthode 2: Utiliser catdoc pour .doc
	if ext == ".doc" {
		catdocPath, err := exec.LookPath("catdoc")
		if err == nil {
			out, err := exec.Command(catdocPath, path).Output()
			if err == nil && len(out) > 0 {
				return string(out), nil
			}
		}
	}
	
	// Méthode 3: Utiliser unrtf pour .rtf
	if ext == ".rtf" {
		unrtfPath, err := exec.LookPath("unrtf")
		if err == nil {
			out, err := exec.Command(unrtfPath, "--text", path).Output()
			if err == nil && len(out) > 0 {
				return string(out), nil
			}
		}
	}
	
	// Méthode 4: Extraire les chaînes
	return dl.extractStringsFromBinary(path)
}

// extractFromPresentation extrait le texte d'une présentation PowerPoint
func (dl *DocumentLoader) extractFromPresentation(path string, ext string) (string, error) {
	// Les outils externes pour PowerPoint sont limités
	return dl.extractStringsFromBinary(path)
}

// extractFromSpreadsheet extrait le texte d'un tableur Excel
func (dl *DocumentLoader) extractFromSpreadsheet(path string, ext string) (string, error) {
	// Tentative d'utiliser xlsx2csv pour .xlsx
	if ext == ".xlsx" {
		xlsx2csvPath, err := exec.LookPath("xlsx2csv")
		if err == nil {
			out, err := exec.Command(xlsx2csvPath, path).Output()
			if err == nil && len(out) > 0 {
				return string(out), nil
			}
		}
	}
	
	// Tentative d'utiliser xls2csv pour .xls
	if ext == ".xls" {
		xls2csvPath, err := exec.LookPath("xls2csv")
		if err == nil {
			out, err := exec.Command(xls2csvPath, path).Output()
			if err == nil && len(out) > 0 {
				return string(out), nil
			}
		}
	}
	
	// Extraction des chaînes
	return dl.extractStringsFromBinary(path)
}

// extractStringsFromBinary extrait les chaînes de caractères d'un fichier binaire
func (dl *DocumentLoader) extractStringsFromBinary(path string) (string, error) {
	// Utiliser l'outil 'strings' si disponible (Unix/Linux/macOS)
	stringsPath, err := exec.LookPath("strings")
	if err == nil {
		out, err := exec.Command(stringsPath, path).Output()
		if err == nil && len(out) > 0 {
			return string(out), nil
		}
	}
	
	// Implémentation basique de 'strings' en Go
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	
	var result strings.Builder
	var currentWord strings.Builder
	
	for _, b := range data {
		if (b >= 32 && b <= 126) || b == '\n' || b == '\t' || b == '\r' {
			currentWord.WriteByte(b)
		} else {
			if currentWord.Len() >= 4 { // Mot d'au moins 4 caractères
				result.WriteString(currentWord.String())
				result.WriteString(" ")
			}
			currentWord.Reset()
		}
	}
	
	if currentWord.Len() >= 4 {
		result.WriteString(currentWord.String())
	}
	
	return result.String(), nil
}

// extractWithOCR tente d'extraire le texte en utilisant un OCR
func (dl *DocumentLoader) extractWithOCR(path string) (string, error) {
	// Vérifier si tesseract est disponible
	tesseractPath, err := exec.LookPath("tesseract")
	if err != nil {
		return "", fmt.Errorf("OCR non disponible: tesseract non trouvé")
	}
	
	// Créer un fichier de sortie temporaire
	tempDir, err := ioutil.TempDir("", "rlama-ocr")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tempDir)
	
	outBasePath := filepath.Join(tempDir, "out")
	
	// Pour les PDF, convertir d'abord en images si possible
	ext := strings.ToLower(filepath.Ext(path))
	if ext == ".pdf" {
		// Vérifier si pdftoppm est disponible
		pdftoppmPath, err := exec.LookPath("pdftoppm")
		if err == nil {
			// Convertir le PDF en images
			fmt.Println("Conversion du PDF en images pour OCR...")
			cmd := exec.Command(pdftoppmPath, "-png", path, filepath.Join(tempDir, "page"))
			if err := cmd.Run(); err != nil {
				return "", fmt.Errorf("échec de la conversion PDF en images: %w", err)
			}
			
			// OCR sur chaque image
			var allText strings.Builder
			imgFiles, _ := filepath.Glob(filepath.Join(tempDir, "page-*.png"))
			for _, imgFile := range imgFiles {
				fmt.Printf("OCR sur %s...\n", filepath.Base(imgFile))
				cmd := exec.Command(tesseractPath, imgFile, outBasePath, "-l", "fra+eng")
				if err := cmd.Run(); err != nil {
					fmt.Printf("Avertissement: OCR échoué pour %s: %v\n", imgFile, err)
					continue
				}
				
				// Lire le texte extrait
				textBytes, err := ioutil.ReadFile(outBasePath + ".txt")
				if err != nil {
					continue
				}
				
				allText.WriteString(string(textBytes))
				allText.WriteString("\n\n")
			}
			
			return allText.String(), nil
		}
	}
	
	// OCR direct sur le fichier (pour les images)
	cmd := exec.Command(tesseractPath, path, outBasePath, "-l", "fra+eng")
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("échec de l'OCR: %w", err)
	}
	
	// Lire le texte extrait
	textBytes, err := ioutil.ReadFile(outBasePath + ".txt")
	if err != nil {
		return "", err
	}
	
	return string(textBytes), nil
}

// tryInstallDependencies tente d'installer des dépendances si nécessaire
func (dl *DocumentLoader) tryInstallDependencies() {
	// Vérifier si pip est disponible (pour les outils Python)
	pipPath, err := exec.LookPath("pip3")
	if err != nil {
		pipPath, err = exec.LookPath("pip")
	}
	
	if err == nil {
		fmt.Println("Vérification des outils Python d'extraction de texte...")
		// Tenter d'installer des packages utiles
		for _, pkg := range []string{"pdfminer.six", "docx2txt", "xlsx2csv"} {
			cmd := exec.Command(pipPath, "show", pkg)
			if err := cmd.Run(); err != nil {
				fmt.Printf("Installation de %s...\n", pkg)
				installCmd := exec.Command(pipPath, "install", "--user", pkg)
				installCmd.Run() // Ignorer les erreurs
			}
		}
	}
} 