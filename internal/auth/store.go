package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// File é o conteúdo de ~/.breno/config.json — credenciais, overrides opcionais
// de endpoint e o último token obtido. Texto plano, permissão 0600.
type File struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"senha,omitempty"`
	BaseURL  string `json:"base_url,omitempty"`
	APIKey   string `json:"api_key,omitempty"`
	Token    *Token `json:"token,omitempty"`
}

// Store persiste o File no disco.
type Store interface {
	Load() (File, error)
	Save(File) error
	Clear() error
	Path() string
}

// DefaultPath devolve ~/.breno/config.json.
func DefaultPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolvendo home do usuário: %w", err)
	}
	return filepath.Join(home, ".breno", "config.json"), nil
}

// FileStore implementa Store lendo/gravando um único arquivo JSON.
type FileStore struct {
	path string
}

func NewFileStore(path string) *FileStore {
	return &FileStore{path: path}
}

func (s *FileStore) Path() string { return s.path }

// Load devolve um File vazio (sem erro) quando o arquivo ainda não existe.
func (s *FileStore) Load() (File, error) {
	data, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return File{}, nil
	}
	if err != nil {
		return File{}, fmt.Errorf("lendo %s: %w", s.path, err)
	}
	var f File
	if err := json.Unmarshal(data, &f); err != nil {
		return File{}, fmt.Errorf("parseando %s: %w", s.path, err)
	}
	return f, nil
}

func (s *FileStore) Save(f File) error {
	dir := filepath.Dir(s.path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("criando %s: %w", dir, err)
	}
	data, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return fmt.Errorf("serializando config: %w", err)
	}
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o600); err != nil {
		return fmt.Errorf("gravando %s: %w", tmp, err)
	}
	if err := os.Rename(tmp, s.path); err != nil {
		return fmt.Errorf("movendo %s -> %s: %w", tmp, s.path, err)
	}
	return nil
}

func (s *FileStore) Clear() error {
	err := os.Remove(s.path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("removendo %s: %w", s.path, err)
	}
	return nil
}
