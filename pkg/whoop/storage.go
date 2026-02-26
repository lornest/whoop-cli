package whoop

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/zalando/go-keyring"
	"golang.org/x/crypto/argon2"
)

const (
	keychainService = "whoop-cli"
	keychainUser    = "default"
	configDir       = ".whoop-cli"
	tokenFile       = "tokens.enc"
	saltFile        = "salt"
	saltSize        = 16
	keySize         = 32 // AES-256
)

// Storage defines the interface for persisting tokens.
type Storage interface {
	Save(data *TokenData) error
	Load() (*TokenData, error)
	Delete() error
}

// KeychainStorage stores tokens in the OS keychain.
type KeychainStorage struct{}

func NewKeychainStorage() *KeychainStorage {
	return &KeychainStorage{}
}

func (k *KeychainStorage) Save(data *TokenData) error {
	b, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal token: %w", err)
	}
	return keyring.Set(keychainService, keychainUser, string(b))
}

func (k *KeychainStorage) Load() (*TokenData, error) {
	s, err := keyring.Get(keychainService, keychainUser)
	if err != nil {
		return nil, fmt.Errorf("keychain load: %w", err)
	}
	var data TokenData
	if err := json.Unmarshal([]byte(s), &data); err != nil {
		return nil, fmt.Errorf("unmarshal token: %w", err)
	}
	return &data, nil
}

func (k *KeychainStorage) Delete() error {
	return keyring.Delete(keychainService, keychainUser)
}

// FileStorage stores tokens in an AES-256-GCM encrypted file with Argon2id KDF.
type FileStorage struct {
	dir        string
	passphrase string
}

func NewFileStorage(passphrase string) *FileStorage {
	home, _ := os.UserHomeDir()
	return &FileStorage{
		dir:        filepath.Join(home, configDir),
		passphrase: passphrase,
	}
}

// NewFileStorageWithDir creates a FileStorage with a custom directory (for testing).
func NewFileStorageWithDir(dir, passphrase string) *FileStorage {
	return &FileStorage{dir: dir, passphrase: passphrase}
}

func (f *FileStorage) Save(data *TokenData) error {
	if err := os.MkdirAll(f.dir, 0o700); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}

	salt, err := f.loadOrCreateSalt()
	if err != nil {
		return err
	}

	key := deriveKey(f.passphrase, salt)

	plaintext, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal token: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("create GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return fmt.Errorf("generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	path := filepath.Join(f.dir, tokenFile)
	return os.WriteFile(path, ciphertext, 0o600)
}

func (f *FileStorage) Load() (*TokenData, error) {
	path := filepath.Join(f.dir, tokenFile)
	ciphertext, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read token file: %w", err)
	}

	salt, err := f.loadSalt()
	if err != nil {
		return nil, err
	}

	key := deriveKey(f.passphrase, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("decrypt: %w", err)
	}

	var data TokenData
	if err := json.Unmarshal(plaintext, &data); err != nil {
		return nil, fmt.Errorf("unmarshal token: %w", err)
	}
	return &data, nil
}

func (f *FileStorage) Delete() error {
	path := filepath.Join(f.dir, tokenFile)
	return os.Remove(path)
}

func (f *FileStorage) loadOrCreateSalt() ([]byte, error) {
	salt, err := f.loadSalt()
	if err == nil {
		return salt, nil
	}

	salt = make([]byte, saltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, fmt.Errorf("generate salt: %w", err)
	}

	path := filepath.Join(f.dir, saltFile)
	if err := os.WriteFile(path, salt, 0o600); err != nil {
		return nil, fmt.Errorf("write salt: %w", err)
	}
	return salt, nil
}

func (f *FileStorage) loadSalt() ([]byte, error) {
	path := filepath.Join(f.dir, saltFile)
	return os.ReadFile(path)
}

func deriveKey(passphrase string, salt []byte) []byte {
	return argon2.IDKey([]byte(passphrase), salt, 1, 64*1024, 4, keySize)
}

// NewStorage returns the best available storage backend.
// It tries keychain first, falling back to encrypted file.
func NewStorage() Storage {
	// Test keychain availability
	testKey := keychainService + "-test"
	err := keyring.Set(testKey, keychainUser, "test")
	if err == nil {
		_ = keyring.Delete(testKey, keychainUser)
		return NewKeychainStorage()
	}
	return NewFileStorage(keychainService)
}
