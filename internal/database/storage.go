package database

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"gogh/internal/models"
	"os"
)

type Storage struct {
	buff     bytes.Buffer
	filename string
}

func New() (*Storage, error) {
	_storage := &Storage{
		buff:     *bytes.NewBuffer([]byte{}),
		filename: "data.enc",
	}
	if _, err := os.Stat(_storage.filename); os.IsNotExist(err) {
		if err := _storage.Save(&models.Data{}); err != nil {
			return nil, fmt.Errorf("create database: %w", err)
		}

	}
	return _storage, nil
}

func (s *Storage) Save(data *models.Data) error {
	enc := gob.NewEncoder(&s.buff)
	if err := enc.Encode(data); err != nil {
		return fmt.Errorf("encode %s: %w", data, err)
	}
	if err := os.WriteFile(s.filename, s.buff.Bytes(), 0644); err != nil {
		return fmt.Errorf("save to file %s: %w", s.filename, err)
	}
	return nil
}

func (s *Storage) Load() (*models.Data, error) {
	buff, err := os.ReadFile(s.filename)
	if err != nil {
		return nil, fmt.Errorf("read data from file %s: %w", s.filename, err)
	}
	s.buff = *bytes.NewBuffer(buff)
	dec := gob.NewDecoder(&s.buff)
	data := models.Data{}
	if err := dec.Decode(&data); err != nil {
		return nil, fmt.Errorf("decode data: %w", err)
	}
	return &data, nil
}
