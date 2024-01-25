package tools

import (
	"bytes"
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
		if err != nil {
			t.Fatalf("%v: decode: %v", in, err)
		}
		if !bytes.Equal(in, out) {
			t.Fatalf("%v: not equal after round trip: %v", in, out)
		}
	})
}

func TestBase64Encode(t *testing.T) {
	bytess := []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit")
	expected := "TG9yZW0gaXBzdW0gZG9sb3Igc2l0IGFtZXQsIGNvbnNlY3RldHVyIGFkaXBpc2NpbmcgZWxpdA=="
	result := Base64Encode(bytess)
	if result != expected {
		t.Fatalf("'%s' was expected bput the result was %s", expected, result)
	}
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
		_, err := Base64Decode(uCase.Data)
		if uCase.Error {
			if err != nil {
				if err.Error() != uCase.Message {
					t.Fatalf("%d - The error message '%s' was expected but was'%s'", i, uCase.Message, err.Error())
				}
			} else {
				t.Fatalf("%d - An error was expected", i+1)
			}
		} else {
			if err != nil {
				t.Fatalf("%d - Unexpected error: '%s'", i, err.Error())
			}
		}
	}
}

func TestEncrypt(t *testing.T) {
	word := string(randstr.Bytes(16))
	_, err := Encrypt(word, secret)
	if err != nil {
		t.Fatalf("Unexpected error encrypting '%s': '%s'", word, err.Error())
	}
}

func TestDecrypt_ok(t *testing.T) {
	//Decrypting correct word
	word := string(randstr.Bytes(16))
	encWord, _ := Encrypt(word, secret)
	_, err := Decrypt(encWord, secret)
	if err != nil {
		t.Fatalf("Unexpected error decrypting '%s': '%s'", encWord, err.Error())
	}
}

func TestDecrypt_ko_errorB64(t *testing.T) {
	//Decrypting incorrect word
	word := string(randstr.Bytes(16))
	encWord, _ := Encrypt(word, secret)
	expectedErr := "Error base64 decoding\nCaused by illegal base64 data at input byte 45"
	_, err := Decrypt(encWord[14:], secret)
	if err != nil {
		if err.Error() != expectedErr {
			t.Fatalf("Error message '%s' was expected instead of '%s", expectedErr, err.Error())
		}
	} else {
		t.Fatalf("An error was expected")
	}
}

func TestDecrypt_ko_errorDecrypting(t *testing.T) {
	//Decrypting incorrect word
	word := string(randstr.Bytes(16))
	wordToDecrypt := Base64Encode([]byte(word))

	expectedErr := "cipher: message authentication failed"
	_, err := Decrypt(wordToDecrypt, secret)
	if err != nil {
		if err.Error() != expectedErr {
			t.Fatalf("Error message '%s' was expected instead of '%s", expectedErr, err.Error())
		}
	} else {
		t.Fatalf("An error was expected")
	}
}

func TestMarshalCrypt(t *testing.T) {
	type myStruct struct {
		Prop1 string `crypt:"true"`
		Prop2 string
	}
	prop1 := string(randstr.Bytes(16))
	prop2 := string(randstr.Bytes(16))
	str := myStruct{
		Prop1: prop1,
		Prop2: prop2,
	}
	encProp1, _ := Encrypt(prop1, secret)

	var result myStruct
	err := MarshalCrypt(str, &result, secret)

	if err != nil {
		t.Fatalf("Unexpected Error Marshaling: '%s'", err.Error())
	}
	if result.Prop1 != encProp1 {
		t.Fatalf("It was expected '%s' for Prop1 but was '%s'", encProp1, result.Prop1)
	}
	if result.Prop2 != prop2 {
		t.Fatalf("It was expected '%s' for Prop2 but was '%s'", prop2, result.Prop2)
	}
}

func TestUnMarshalCrypt_ok(t *testing.T) {
	type myStruct struct {
		Prop1 string `crypt:"true"`
		Prop2 string
	}

	prop := string(randstr.Bytes(16))
	prop1, _ := Encrypt(prop, secret)
	prop2 := string(randstr.Bytes(16))
	str := myStruct{
		Prop1: prop1,
		Prop2: prop2,
	}

	var result myStruct

	err := UnMarshalCrypt(str, &result, secret)
	if err != nil {
		t.Fatalf("Unexpected Error unmarshaling: '%s'", err.Error())
	}
	if result.Prop1 != prop {
		t.Fatalf("It was expected '%s' for Prop1 but was '%s'", prop, result.Prop1)
	}
	if result.Prop2 != prop2 {
		t.Fatalf("It was expected '%s' for Prop2 but was '%s'", prop2, result.Prop2)
	}
}
func TestUnMarshalCrypt_ko(t *testing.T) {
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
	expectedError := "Error decrypting field [Prop3]"
	var result myStruct

	err := UnMarshalCrypt(str, &result, secret)
	if err == nil {
		t.Fatalf("Un error was expected")
	}
	if expectedError != err.Error() {
		t.Fatalf("The error '%s' was expected but was '%s'", expectedError, err.Error())
	}
}

func TestMarshalHash(t *testing.T) {
	type myStruct struct {
		Prop1 string `hash:"true"`
		Prop2 string `hash:"FalSe"`
	}
	prop1 := string(randstr.Bytes(16))
	prop2 := string(randstr.Bytes(16))
	str := myStruct{
		Prop1: prop1,
		Prop2: prop2,
	}
	hash := mdHashing(prop1)

	var result myStruct
	MarshalHash(str, &result)

	if result.Prop1 != hash {
		t.Fatalf("It was expected '%s' for Prop1 but was '%s'", hash, result.Prop1)
	}
	if result.Prop2 != prop2 {
		t.Fatalf("It was expected '%s' for Prop2 but was '%s'", prop2, result.Prop2)
	}
}
