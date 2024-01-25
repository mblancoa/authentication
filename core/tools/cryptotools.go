package tools

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"github.com/mblancoa/authentication/core/errors"
	"reflect"
	"strings"
)

var bytesToStaticEncryption = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

func Base64Encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}
func Base64Decode(s string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return []byte{}, errors.NewGenericErrorByCause("Error base64 decoding", err)
	}
	return data, nil
}

// Encrypt method is to encrypt or hide any classified text
func Encrypt(text, secret string) (string, error) {
	cipherText, err := encryptIt([]byte(text), secret)
	if err != nil {
		return "", err
	}
	return Base64Encode(cipherText), nil
}

// Decrypt method is to extract back the encrypted text
func Decrypt(text, secret string) (string, error) {
	cipherText, err := Base64Decode(text)
	if err != nil {
		return "", err
	}
	plainText, err := decryptIt(cipherText, secret)
	if err != nil {
		return "", err
	}
	return string(plainText), nil
}

func MarshalCrypt(obj, v any, secret string) error {
	oType := reflect.TypeOf(obj)
	vType := reflect.Indirect(reflect.ValueOf(v))

	for i := 0; i < oType.NumField(); i++ {
		f := oType.Field(i)
		fValue := reflect.Indirect(reflect.ValueOf(obj)).Field(i)
		if cryTag, ok := f.Tag.Lookup("crypt"); ok && strings.ToLower(cryTag) != "false" {
			value := fValue.String()
			encValue, err := Encrypt(value, secret)
			if err != nil {
				return errors.NewGenericErrorf("Error encrypting field %s", f.Name)
			}
			vType.Field(i).Set(reflect.ValueOf(encValue))
		} else {
			vType.Field(i).Set(fValue)
		}
	}
	return nil
}

func UnMarshalCrypt(obj, v any, secret string) error {
	oType := reflect.TypeOf(obj)
	vValue := reflect.Indirect(reflect.ValueOf(v))

	for i := 0; i < oType.NumField(); i++ {
		f := oType.Field(i)
		fValue := reflect.Indirect(reflect.ValueOf(obj)).Field(i)
		if cryTag, ok := f.Tag.Lookup("crypt"); ok && strings.ToLower(cryTag) != "false" {
			value := fValue.String()
			decValue, err := Decrypt(value, secret)
			if err != nil {
				return errors.NewGenericErrorf("Error decrypting field %s", f.Name)
			}
			vValue.Field(i).Set(reflect.ValueOf(decValue))
		} else {
			vValue.Field(i).Set(fValue)
		}
	}
	return nil
}

func MarshalHash(obj any, v any) {
	oType := reflect.TypeOf(obj)
	vValue := reflect.Indirect(reflect.ValueOf(v))

	for i := 0; i < oType.NumField(); i++ {
		f := oType.Field(i)
		fValue := reflect.Indirect(reflect.ValueOf(obj)).Field(i)
		if hashTag, ok := f.Tag.Lookup("hash"); ok && strings.ToLower(hashTag) != "false" {
			value := fValue.String()
			vHashValue := reflect.ValueOf(mdHashing(value))
			vValue.Field(i).Set(vHashValue)
		} else {
			vValue.Field(i).Set(fValue)
		}
	}
}

func mdHashing(input string) string {
	byteInput := []byte(input)
	md5Hash := md5.Sum(byteInput)
	return hex.EncodeToString(md5Hash[:]) // by referring to it as a string
}

func encryptIt(value []byte, keyPhrase string) ([]byte, error) {
	aesBlock, err := aes.NewCipher([]byte(mdHashing(keyPhrase)))
	if err != nil {
		return []byte{}, err
	}

	gcmInstance, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return []byte{}, err
	}

	//	nonce := make([]byte, gcmInstance.NonceSize())
	//_, _ = io.ReadFull(rand.Reader, nonce)
	nonce := bytesToStaticEncryption[:gcmInstance.NonceSize()]
	cipheredText := gcmInstance.Seal(nonce, nonce, value, nil)

	return cipheredText, nil
}

func decryptIt(ciphered []byte, keyPhrase string) ([]byte, error) {
	hashedPhrase := mdHashing(keyPhrase)
	aesBlock, err := aes.NewCipher([]byte(hashedPhrase))
	if err != nil {
		return []byte{}, err
	}
	gcmInstance, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return []byte{}, err
	}
	nonceSize := gcmInstance.NonceSize()
	nonce, cipheredText := ciphered[:nonceSize], ciphered[nonceSize:]
	originalText, err := gcmInstance.Open(nil, nonce, cipheredText, nil)
	if err != nil {
		return []byte{}, err
	}
	return originalText, nil
}
