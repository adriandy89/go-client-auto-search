package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

var aesKey = []byte{
	0x4F, 0x2A, 0xD3, 0x9C, 0xE5, 0x71, 0x8B, 0xA4,
	0xFF, 0x17, 0xB2, 0x93, 0xC6, 0x5D, 0x8E, 0x1F,
	0x4A, 0xE3, 0x71, 0xD2, 0xC9, 0xA8, 0x7F, 0xB6,
	0x88, 0x0C, 0x3D, 0xE5, 0x91, 0xFA, 0x12, 0x7E,
}

// EncryptData encripta los datos con AES-GCM
func EncryptData(data []byte) (string, error) {
	block, err := aes.NewCipher(aesKey) // 🔐 Usa la clave fija
	if err != nil {
		return "", err
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	encrypted := aesGCM.Seal(nonce, nonce, data, nil)
	return base64.StdEncoding.EncodeToString(encrypted), nil
}
