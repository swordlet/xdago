package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

func AesEncrypt(data, key, iv []byte) ([]byte, error) {
	aesBlockEncryptor, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	content := PKCS5Padding(data, aesBlockEncryptor.BlockSize())
	encrypted := make([]byte, len(content))

	aesEncryptor := cipher.NewCBCEncrypter(aesBlockEncryptor, iv)
	aesEncryptor.CryptBlocks(encrypted, content)

	return encrypted, nil
}

func AesDecrypt(encrypted, key, iv []byte) ([]byte, error) {

	decrypted := make([]byte, len(encrypted))
	var aesBlockDecrypt cipher.Block
	aesBlockDecrypt, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesDecrypt := cipher.NewCBCDecrypter(aesBlockDecrypt, iv)
	aesDecrypt.CryptBlocks(decrypted, encrypted)
	content := PKCS5Trimming(decrypted)
	return content, nil

}

func PKCS5Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

func PKCS5Trimming(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1]
	return encrypt[:len(encrypt)-int(padding)]
}
