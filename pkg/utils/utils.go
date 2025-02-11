package utils

import (
	"fmt"
	"os"
)

// CheckFileSize verifica si el tamaño del archivo es menor que el tamaño máximo permitido
func CheckFileSize(filePath string, maxFileSize int64) (bool, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return false, err
	}
	return fileInfo.Size() <= maxFileSize, nil
}

// ReadFile lee el contenido de un archivo y lo devuelve como un slice de bytes
func ReadFile(filePath string) ([]byte, error) {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return fileData, nil
}

// LogError imprime un mensaje de error en la consola
func LogError(err error) {
	if err != nil {
		fmt.Println("❌ Error:", err)
	}
}
