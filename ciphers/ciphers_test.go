package ciphers

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	privatekey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Error(err)
	}

	plainbytes := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
	bufPlain := bytes.NewBuffer(plainbytes)

	bufEncr := bytes.NewBuffer([]byte{})
	bufDecr := bytes.NewBuffer([]byte{})

	e := Encrypt{privatekey}
	d := Decrypt{privatekey}

	e.Pipe(bufPlain, bufEncr)
	d.Pipe(bufEncr, bufDecr)

	if !bytes.Equal(plainbytes, bufDecr.Bytes()) {
		t.Errorf("expected plain buffer %v to equal decrypted buffer %v, got equality false", plainbytes, bufDecr.Bytes())
	}
}
