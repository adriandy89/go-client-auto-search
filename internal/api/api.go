package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
)

const apiURL = "http://localhost:8080/diagnostics" // 🌐 URL de la API
const apiKey = "tu_clave_secreta"                  // 🔐 Clave API

// SendToAPI envía el archivo encriptado a la API
func SendToAPI(encryptedData string, filePath string) error {
	jsonData := map[string]string{
		"file_path":          filePath,
		"original_file_name": filepath.Base(filePath),
		"encrypted_data":     encryptedData,
		"api_key":            apiKey,
	}
	body, _ := json.Marshal(jsonData)

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 📌 Leer la respuesta completa para ver el error
	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		fmt.Println("❌ Error en la API:", resp.Status, "➡", string(respBody))
		return fmt.Errorf("error en la API: %s", resp.Status)
	}

	fmt.Println("✅ Archivo enviado correctamente:", filePath)
	return nil
}

func CheckBackendConnection() bool {
	resp, err := http.Get(apiURL)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}
