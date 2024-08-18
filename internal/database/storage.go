package database

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"gogh/internal/models"
	"os"
)

type Storage struct {
	filename string
}

func New() (*Storage, error) {
	_storage := &Storage{
		// buff:     *bytes.NewBuffer([]byte{}),
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
	buff := *bytes.NewBuffer([]byte{})
	enc := gob.NewEncoder(&buff)
	if err := enc.Encode(data); err != nil {
		return fmt.Errorf("encode data: %w", err)
	}
	if err := os.WriteFile(s.filename, buff.Bytes(), 0644); err != nil {
		return fmt.Errorf("save to file %s: %w", s.filename, err)
	}
	return nil
}

func (s *Storage) Load() (*models.Data, error) {
	readbuff, err := os.ReadFile(s.filename)
	if err != nil {
		return nil, fmt.Errorf("read data from file %s: %w", s.filename, err)
	}
	buff := *bytes.NewBuffer(readbuff)
	dec := gob.NewDecoder(&buff)
	data := models.Data{}
	if err := dec.Decode(&data); err != nil {
		return nil, fmt.Errorf("decode data: %w", err)
	}
	return &data, nil
}
