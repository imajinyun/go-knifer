package vcrypto

import cryptoimpl "github.com/imajinyun/go-knifer/internal/crypto"

// VigenereEncrypt encrypts printable ASCII text using the Vigenere helper algorithm.
//
// Security: Vigenere is a historical compatibility algorithm, not modern
// cryptography. Do not use it for confidential or attacker-controlled data;
// prefer AESSealGCM for new encrypted data.
func VigenereEncrypt(data, cipherKey string) (string, error) {
	return cryptoimpl.VigenereEncrypt(data, cipherKey)
}

// VigenereDecrypt decrypts text encrypted by VigenereEncrypt.
//
// Security: Vigenere is a historical compatibility algorithm, not modern
// cryptography. Do not use it for confidential or attacker-controlled data;
// prefer AESOpenGCM for new encrypted data.
func VigenereDecrypt(data, cipherKey string) (string, error) {
	return cryptoimpl.VigenereDecrypt(data, cipherKey)
}

// XXTEAEncrypt encrypts data using XXTEA.
//
// Security: XXTEA is provided for legacy compatibility. It is not an
// authenticated encryption scheme and should not be used for new encrypted data;
// prefer AESSealGCM.
func XXTEAEncrypt(data, key []byte) []byte { return cryptoimpl.XXTEAEncrypt(data, key) }

// XXTEADecrypt decrypts data using XXTEA.
//
// Security: XXTEA is provided for legacy compatibility. It is not an
// authenticated encryption scheme and should not be used for new encrypted data;
// prefer AESOpenGCM.
func XXTEADecrypt(data, key []byte) ([]byte, error) { return cryptoimpl.XXTEADecrypt(data, key) }
