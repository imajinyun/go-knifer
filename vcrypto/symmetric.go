package vcrypto

import cryptoimpl "github.com/imajinyun/go-knifer/internal/crypto"

// VigenereEncrypt encrypts printable ASCII text using the Vigenere helper algorithm.
func VigenereEncrypt(data, cipherKey string) (string, error) {
	return cryptoimpl.VigenereEncrypt(data, cipherKey)
}

// VigenereDecrypt decrypts text encrypted by VigenereEncrypt.
func VigenereDecrypt(data, cipherKey string) (string, error) {
	return cryptoimpl.VigenereDecrypt(data, cipherKey)
}

// XXTEAEncrypt encrypts data using XXTEA.
func XXTEAEncrypt(data, key []byte) []byte { return cryptoimpl.XXTEAEncrypt(data, key) }

// XXTEADecrypt decrypts data using XXTEA.
func XXTEADecrypt(data, key []byte) ([]byte, error) { return cryptoimpl.XXTEADecrypt(data, key) }
