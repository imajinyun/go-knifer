package conf

import (
	"encoding/base64"
	"strings"
)

// Base64Decrypt decrypts ENC(base64:...) values using standard base64 decoding.
func Base64Decrypt(cipherText string) (string, error) {
	cipherText = strings.TrimPrefix(cipherText, "base64:")
	b, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// DecryptValues returns a copy with ENC(...) values decrypted.
func (s *Conf) DecryptValues(decrypt DecryptFunc) (*Conf, error) {
	out := New()
	if s == nil || s.data == nil {
		return out, nil
	}
	for group, m := range s.data {
		for key, value := range m {
			plain, err := decryptValue(value, decrypt)
			if err != nil {
				return nil, wrapConfigParse("decrypt config value "+group+"."+key, err)
			}
			out.SetByGroup(group, key, plain)
		}
	}
	return out, nil
}

func decryptValue(value string, decrypt DecryptFunc) (string, error) {
	if !strings.HasPrefix(value, "ENC(") || !strings.HasSuffix(value, ")") {
		return value, nil
	}
	if decrypt == nil {
		decrypt = Base64Decrypt
	}
	return decrypt(strings.TrimSpace(value[4 : len(value)-1]))
}
