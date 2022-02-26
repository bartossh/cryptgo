package ciphers

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
)

const (
	bitSize = 4096
)

type (
	// Encrypt offers encryption pipe
	Encrypt struct {
		k *rsa.PrivateKey
	}
	// Decrypt offers decryption pipe
	Decrypt struct {
		k *rsa.PrivateKey
	}
)

// NewEncrypt returns instance of Encrypt or error otherwise
func NewEncrypt(priv, passwd []byte) (*Encrypt, error) {
	k, err := rsaConfigSetup(priv, passwd)
	if err != nil {
		return nil, fmt.Errorf("cannot create instance of encrypt, %s", err)
	}
	return &Encrypt{k}, nil
}

// NewDecrypt returns instance of Encrypt or error otherwise
func NewDecrypt(priv, passwd []byte) (*Decrypt, error) {
	k, err := rsaConfigSetup(priv, passwd)
	if err != nil {
		return nil, fmt.Errorf("cannot create instance of decrypt, %s", err)
	}
	return &Decrypt{k}, nil
}

// Pipe pipes bytes from reader to writer performing cryptographic encryption
func (e *Encrypt) Pipe(rd io.Reader, wr io.Writer) error {
	buf, err := io.ReadAll(rd)
	if err != nil {
		return fmt.Errorf("cannot pipe and encrypt, reading failed, %s", err)
	}
	cipherbuf, err := encryptWithPublicKey(buf, &e.k.PublicKey)
	if err != nil {
		return fmt.Errorf("cannot pipe end encrypt, encryption failed, %s", err)
	}
	if _, err := wr.Write(cipherbuf); err != nil {
		return fmt.Errorf("cannot pipe and encrypt, writing to failed, %s", err)
	}
	return nil
}

// Pipe pipes bytes from reader to writer performing cryptographic decryption
func (d *Decrypt) Pipe(rd io.Reader, wr io.Writer) error {
	buf, err := io.ReadAll(rd)
	if err != nil {
		return fmt.Errorf("cannot pipe and decrypt, reading failed, %s", err)
	}
	plainbuf, err := decryptWithPrivateKey(buf, d.k)
	if err != nil {
		return fmt.Errorf("cannot pipe end decrypt, decryption failed, %s", err)
	}
	if _, err := wr.Write(plainbuf); err != nil {
		return fmt.Errorf("cannot pipe and decrypt, writing failed, %s", err)
	}
	return nil
}

// GeneratePrivateKey generates new rsa private key
func GeneratePrivateKey() (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, err
	}

	err = privateKey.Validate()
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

// EncodePrivateKeyToPEM encodes private key to PEM format
func EncodePrivateKeyToPEM(privateKey *rsa.PrivateKey) []byte {
	privDER := x509.MarshalPKCS1PrivateKey(privateKey)
	privBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privDER,
	}
	privatePEM := pem.EncodeToMemory(&privBlock)
	return privatePEM
}

func rsaConfigSetup(priv, passwd []byte) (*rsa.PrivateKey, error) {
	privPem, _ := pem.Decode(priv)
	var privPemBytes []byte
	var err error
	if privPem.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("provided key is of wrong type")
	}

	if len(passwd) > 0 {
		privPemBytes, err = x509.DecryptPEMBlock(privPem, passwd)
		if err != nil {
			return nil, fmt.Errorf("cannot decrypt private key, %s", err)
		}
	} else {
		privPemBytes = privPem.Bytes
	}

	var privateKey *rsa.PrivateKey
	if privateKey, err = x509.ParsePKCS1PrivateKey(privPemBytes); err != nil {
		return nil, fmt.Errorf("cannot parse private key, %s", err)
	}

	return privateKey, nil
}

func encryptWithPublicKey(msg []byte, pub *rsa.PublicKey) ([]byte, error) {
	hash := sha512.New()
	buf := bytes.NewBuffer([]byte{})
	chunkSize := pub.Size() - 2*hash.Size() - 2

	for i := 0; i < len(msg); i = i + chunkSize {
		end := i + chunkSize
		if end > len(msg) {
			end = len(msg)
		}
		ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, pub, msg[i:end], nil)
		if err != nil {
			return nil, err
		}
		buf.Write(ciphertext)
	}

	return buf.Bytes(), nil
}

func decryptWithPrivateKey(ciphertext []byte, priv *rsa.PrivateKey) ([]byte, error) {
	hash := sha512.New()
	buf := bytes.NewBuffer([]byte{})
	chunkSize := priv.Size()

	for i := 0; i < len(ciphertext); i = i + chunkSize {
		end := i + chunkSize
		if end > len(ciphertext) {
			end = len(ciphertext)
		}
		plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, priv, ciphertext[i:end], nil)
		if err != nil {
			return nil, err
		}
		buf.Write(plaintext)
	}
	return buf.Bytes(), nil
}
