package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"my/bot/storage"
	"os"
	"path/filepath"
	"time"
)

const (
	defaultPerm = 0774
)

type Storage struct {
	basePath string
}

func New(basePath string) Storage {
	return Storage{basePath: basePath}
}

func (s Storage) Save(page *storage.Page) error {
	fPath := filepath.Join(s.basePath, page.UserName)

	err := os.MkdirAll(fPath, defaultPerm)
	if err != nil {
		return fmt.Errorf("can not save page: %w", err)
	}

	fName, err := fileName(page)
	if err != nil {
		return fmt.Errorf("can not save page: %w", err)
	}

	fPath = filepath.Join(fPath, fName)

	file, err := os.Create(fPath)
	if err != nil {
		return fmt.Errorf("can not save page: %w", err)
	}
	defer func() { _ = file.Close() }()

	err = gob.NewEncoder(file).Encode(page)
	if err != nil {
		return fmt.Errorf("can not save page: %w", err)
	}
	return nil
}

func (s Storage) PickRandom(userName string) (*storage.Page, error) {

	path := filepath.Join(s.basePath, userName)

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("can not pic random page: %w", err)
	}

	if len(files) == 0 {
		return nil, storage.ErrNoSavedPages
	}

	rand.Seed(time.Now().UnixNano())

	n := rand.Intn(len(files))

	file := files[n]

	return s.decodePage(filepath.Join(path, file.Name()))

}

func (s Storage) Remove(p *storage.Page) error {

	fName, err := fileName(p)
	if err != nil {
		return fmt.Errorf("can not remove file: %w", err)
	}
	path := filepath.Join(s.basePath, p.UserName, fName)

	err = os.Remove(path)
	if err != nil {
		return fmt.Errorf("can not remove file %s: %w", path, err)
	}

	return nil
}

func (s Storage) IsExists(p *storage.Page) (bool, error) {

	fName, err := fileName(p)
	if err != nil {
		return false, fmt.Errorf("can not check if file: %w", err)
	}
	path := filepath.Join(s.basePath, p.UserName, fName)

	switch _, err = os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("can not check if file %s: %w", path, err)

	}
	return true, nil
}

func (s Storage) decodePage(filePath string) (*storage.Page, error) {

	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("can not decode page: %w", err)
	}

	defer func() { _ = f.Close() }()

	var p storage.Page

	err = gob.NewDecoder(f).Decode(&p)
	if err != nil {
		return nil, fmt.Errorf("can not decode page: %w", err)
	}

	return &p, nil
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}
