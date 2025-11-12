package vault

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Meta struct {
	Env         string `json:"env"`
	Salt        string `json:"salt"`
	FingerPrint string `json:"fingerprint"`
}

type Entry struct {
	Value     string `json:"value"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type Vault struct {
	Meta    Meta             `json:"meta"`
	Path    string           `json:"-"`
	Entries map[string]Entry `json:"entries"`
}

func Create(vaultPath, env, passphrase, salt string) (*Vault, error) {
	var vault Vault
	fingerprint, err := vault.GenerateFingerprint(passphrase)
	if err != nil {
		return nil, err
	}

	vault.Path = vaultPath
	vault.Meta.Env = env
	vault.Meta.FingerPrint = fingerprint
	vault.Meta.Salt = salt
	vault.Entries = make(map[string]Entry)
	err = vault.Save()
	if err != nil {
		return nil, err
	}

	return &vault, nil
}

func Open(env string) (*Vault, error) {
	vaultPath := filepath.Join(".lockify", env+".vault.enc")
	if _, err := os.Stat(vaultPath); err != nil {
		return nil, err
	}

	jsonFile, err := os.Open(vaultPath)
	if err != nil {
		return nil, err
	}

	defer jsonFile.Close()

	var vault Vault
	byteValue, _ := io.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &vault)
	vault.Path = vaultPath

	return &vault, nil
}

func (v *Vault) GenerateFingerprint(passphrase string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(passphrase), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (v *Vault) VerifyFingerPrint(passphrase string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(v.Meta.FingerPrint), []byte(passphrase))
	return err == nil
}

func (v *Vault) DeleteEntry(key string) {
	delete(v.Entries, key)
}

func (v *Vault) GetEntry(key string) (Entry, error) {
	entry, exists := v.Entries[key]
	if !exists {
		return Entry{}, fmt.Errorf("key %s has no value", key)
	}

	return entry, nil
}

func (v *Vault) SetEntry(key, value string) {
	_, exists := v.Entries[key]
	var updatedAt string
	if exists {
		updatedAt = time.Now().Format("2006-1-2T15:04:05")
	}

	v.Entries[key] = Entry{
		Value:     value,
		CreatedAt: time.Now().Format("2006-1-2T15:04:05"),
		UpdatedAt: updatedAt,
	}
}

func (v *Vault) Save() error {
	content, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(v.Path, content, 0600)
	if err != nil {
		return err
	}

	return nil
}
