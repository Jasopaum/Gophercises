package secrets

import (
	"encoding/json"
	"errors"
	"fmt"
	"gophercises/secrets/encryption"
	"io"
	"os"
	"sync"
)

type Vault struct {
	key      string
	filepath string
	table    map[string]string
	mutex    sync.Mutex
}

func File(key, path string) *Vault {
	return &Vault{
		key:      key,
		filepath: path,
	}
}

func (v *Vault) load() error {
	f, err := os.Open(v.filepath)
	if err != nil {
		v.table = make(map[string]string)
		return nil
	}
	defer f.Close()
	r, err := encryption.DecryptReader(v.key, f)
	if err != nil {
		return err
	}
	return v.readTable(r)
}

func (v *Vault) readTable(r io.Reader) error {
	dec := json.NewDecoder(r)
	return dec.Decode(&v.table)
}

func (v *Vault) save() error {
	f, err := os.OpenFile(v.filepath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer f.Close()
	w, err := encryption.EncryptWriter(v.key, f)
	if err != nil {
		return err
	}
	return v.writeTable(w)
}

func (v *Vault) writeTable(w io.Writer) error {
	end := json.NewEncoder(w)
	return end.Encode(v.table)
}

func (v *Vault) Set(entry, value string) error {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	err := v.load()
	if err != nil {
		return err
	}
	v.table[entry] = value
	err = v.save()
	return err
}

func (v *Vault) Get(entry string) (string, error) {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	err := v.load()
	if err != nil {
		return "", err
	}
	if value, ok := v.table[entry]; ok {
		return value, nil
	}
	return "", errors.New("Entry not found.")
}

func (v *Vault) List() ([]string, error) {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	err := v.load()
	if err != nil {
		return nil, err
	}
	var res []string
	for entry, value := range v.table {
		res = append(res, fmt.Sprintf("%s: %s", entry, value))
	}
	return res, nil
}
