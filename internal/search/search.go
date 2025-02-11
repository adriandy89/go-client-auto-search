package search

import (
	"go-client-auto-search/pkg/utils"
	"io/fs"
	"path/filepath"
)

// SearchFiles busca archivos con una extensión específica en un directorio dado
func SearchFiles(startDir string, fileExt string, maxFileSize int64, processFile func(string, []byte) error) error {
	return filepath.WalkDir(startDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // Ignorar errores de acceso
		}
		if d.IsDir() || filepath.Ext(d.Name()) != fileExt {
			return nil
		}

		// 📌 Verificar el tamaño del archivo
		isValidSize, err := utils.CheckFileSize(path, maxFileSize)
		if err != nil || !isValidSize {
			return nil
		}

		// 📌 Leer el archivo
		fileData, err := utils.ReadFile(path)
		if err != nil {
			utils.LogError(err)
			return nil
		}

		return processFile(path, fileData)
	})
}
