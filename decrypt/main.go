package main

import (
	"crypto/aes"
	"crypto/cipher"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func main() {
	key := [32]byte{174, 76, 83, 109, 164, 175, 162, 120, 19, 152, 172, 207, 133, 183, 184, 30, 250, 5, 116, 121, 196, 111, 53, 21, 177, 194, 20, 37, 246, 118, 44, 235}

	usr, _ := user.Current()
	downloads := filepath.Join(usr.HomeDir, "Downloads/target")
	foldersPath := []string{downloads}

	for _, folderPath := range foldersPath {
		files, _ := os.ReadDir(folderPath)
		for _, file := range files {
			filePath := filepath.Join(folderPath, file.Name())
			if strings.HasSuffix(file.Name(), ".NT_sdh") {
				ciphertext, _ := os.ReadFile(filePath)
				block, err := aes.NewCipher(key[:])
				if err != nil {
					return
				}
				if len(ciphertext) < aes.BlockSize {
					return
				}
				iv := ciphertext[:aes.BlockSize]
				ciphertext = ciphertext[aes.BlockSize:]
				stream := cipher.NewCFBDecrypter(block, iv)
				stream.XORKeyStream(ciphertext, ciphertext)

				destinationPath := filepath.Join(folderPath, file.Name()[:len(file.Name())-7])
				err = os.WriteFile(destinationPath, ciphertext, 0644)
				if err != nil {
					return
				}

				os.Remove(filePath)
			}
		}
	}
}
