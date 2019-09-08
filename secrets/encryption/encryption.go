package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"errors"
	"io"
)

func EncryptWriter(keyphrase string, w io.Writer) (*cipher.StreamWriter, error) {
	key := hashKeyphrase(keyphrase)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	s := cipher.NewCFBEncrypter(block, iv)

	n, err := w.Write(iv)
	if n != len(iv) || err != nil {
		return nil, errors.New("encrypt: unable to write full iv to writer")
	}

	return &cipher.StreamWriter{S: s, W: w}, nil
}

func DecryptReader(keyphrase string, r io.Reader) (*cipher.StreamReader, error) {
	key := hashKeyphrase(keyphrase)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	iv := make([]byte, aes.BlockSize)
	n, err := r.Read(iv)
	if err != nil || n != aes.BlockSize {
		return nil, errors.New("Could not read IV.")
	}

	s := cipher.NewCFBDecrypter(block, iv)

	return &cipher.StreamReader{S: s, R: r}, nil
}

func hashKeyphrase(keyphrase string) []byte {
	h := md5.New()
	io.WriteString(h, keyphrase)
	return h.Sum(nil)
}
