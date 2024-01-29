package tools

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/thanhpk/randstr"
	"testing"
)

const secret = "ob#qjc|?[^dhg`p9ot"

func FuzzBase64EncDec(f *testing.F) {
	for _, seed := range [][]byte{{}, {0}, {9}, {0xa}, {0xf}, {1, 2, 3, 4}} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, in []byte) {
		enc := Base64Encode(in)
		out, err := Base64Decode(enc)
		assert.NoError(t, err, fmt.Sprintf("%v: decode: %v", in, err))
		assert.Equal(t, in, out, fmt.Sprintf("%v: not equal after round trip: %v", in, out))
	})
}

func TestBase64Encode(t *testing.T) {
	bts := []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit")
	expected := "TG9yZW0gaXBzdW0gZG9sb3Igc2l0IGFtZXQsIGNvbnNlY3RldHVyIGFkaXBpc2NpbmcgZWxpdA=="
	result := Base64Encode(bts)

	assert.Equal(t, expected, result)
}

func TestBase64Decode(t *testing.T) {
	useCases := []struct {
		Data    string
		Error   bool
		Message string
	}{
		{"TG9yZW0gaXBzdW0gZG9sb3Igc2l0IGFtZXQsIGNfassds", true, "Error base64 decoding\nCaused by illegal base64 data at input byte 44"},
		{"TG9yZW0gaXBzdW0gZG9sb3Igc2l0IGFtZXQsIGNvbnNlY3RldHVyIGFkaXBpc2NpbmcgZWxpdA==", false, ""},
	}
	for i, uCase := range useCases {
		result, err := Base64Decode(uCase.Data)
		if uCase.Error {
			assert.Error(t, err, "case:%d", i)
			assert.Equal(t, uCase.Message, err.Error(), "case:%d", i)
			assert.Empty(t, result)
		} else {
			assert.NoError(t, err, "case:%d", i)
			assert.NotEmpty(t, result)
		}
	}
}

func TestEncrypt(t *testing.T) {
	word := string(randstr.Bytes(16))
	result, err := Encrypt(word, secret)

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestDecrypt_ok(t *testing.T) {
	//Decrypting correct word
	word := string(randstr.Bytes(16))
	encWord, _ := Encrypt(word, secret)
	result, err := Decrypt(encWord, secret)

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.Equal(t, word, result)
}

func TestDecrypt_ko_errorB64(t *testing.T) {
	//Decrypting incorrect word
	word := string(randstr.Bytes(16))
	encWord, _ := Encrypt(word, secret)
	expectedErr := "Error base64 decoding\nCaused by illegal base64 data at input byte 45"
	result, err := Decrypt(encWord[14:], secret)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err.Error())
	assert.Empty(t, result)
}

func TestDecrypt_ko_errorDecrypting(t *testing.T) {
	//Decrypting incorrect word
	word := string(randstr.Bytes(16))
	wordToDecrypt := Base64Encode([]byte(word))

	expectedErr := "cipher: message authentication failed"
	result, err := Decrypt(wordToDecrypt, secret)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err.Error())
	assert.Empty(t, result)
}

func TestMarshalCrypt(t *testing.T) {
	type myStruct struct {
		Prop1 string `crypt:"true"`
		Prop2 string
	}
	prop1 := string(randstr.Bytes(16))
	prop2 := string(randstr.Bytes(16))
	str := myStruct{prop1, prop2}
	encProp1, _ := Encrypt(prop1, secret)
	expected := myStruct{encProp1, prop2}

	var result myStruct
	err := MarshalCrypt(str, &result, secret)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestUnmarshalCrypt_ok(t *testing.T) {
	type myStruct struct {
		Prop1 string `crypt:"true"`
		Prop2 string
	}

	prop := string(randstr.Bytes(16))
	prop1, _ := Encrypt(prop, secret)
	prop2 := string(randstr.Bytes(16))
	str := myStruct{prop1, prop2}
	expected := myStruct{prop, prop2}

	var result myStruct
	err := UnmarshalCrypt(str, &result, secret)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}
func TestUnmarshalCrypt_ko(t *testing.T) {
	type myStruct struct {
		Prop1 string `crypt:"true"`
		Prop2 string
		Prop3 string `crypt:"true"`
	}

	prop1, _ := Encrypt(string(randstr.Bytes(16)), secret)
	str := myStruct{
		Prop1: prop1,
		Prop2: string(randstr.Bytes(16)),
		Prop3: string(randstr.Bytes(16)),
	}
	expectedError := "Error decrypting field Prop3"
	var result myStruct

	err := UnmarshalCrypt(str, &result, secret)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err.Error())
}

func TestMarshalHash(t *testing.T) {
	type myStruct struct {
		Prop1 string `hash:"true"`
		Prop2 string `hash:"FalSe"`
	}
	prop1 := string(randstr.Bytes(16))
	prop2 := string(randstr.Bytes(16))
	str := myStruct{prop1, prop2}
	expected := myStruct{mdHashing(prop1), prop2}

	var result myStruct
	MarshalHash(str, &result)

	assert.Equal(t, expected, result)
}
