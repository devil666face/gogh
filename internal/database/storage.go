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

func New(_filename string) (*Storage, error) {
	var data = models.Data{
		Settings: models.Settings{
			Compress: true,
		},
	}
	_storage := &Storage{
		filename: _filename,
	}
	if _, err := os.Stat(_storage.filename); os.IsNotExist(err) {
		if err := _storage.Save(&data); err != nil {
			return nil, fmt.Errorf("create database: %w", err)
		}

	}
	return _storage, nil
}

func (s *Storage) saveToFile(filename string, data interface{}) error {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	if err := enc.Encode(data); err != nil {
		return fmt.Errorf("encode data: %w", err)
	}
	if err := os.WriteFile(filename, buff.Bytes(), 0644); err != nil {
		return fmt.Errorf("save to file %s: %w", filename, err)
	}
	return nil
}

func (s *Storage) Save(data *models.Data) error {
	return s.saveToFile(s.filename, data)
}

func (s *Storage) SaveFile(file *models.File) error {
	filename := file.Filename + ".gogh"
	return s.saveToFile(filename, file)
}

// func (s *Storage) Save(data *models.Data) error {
// 	buff := *bytes.NewBuffer([]byte{})
// 	enc := gob.NewEncoder(&buff)
// 	if err := enc.Encode(data); err != nil {
// 		return fmt.Errorf("encode data: %w", err)
// 	}
// 	if err := os.WriteFile(s.filename, buff.Bytes(), 0644); err != nil {
// 		return fmt.Errorf("save to file %s: %w", s.filename, err)
// 	}
// 	return nil
// }

// func (s *Storage) SaveFile(file *models.File) error {
// 	var filename = file.Filename + ".gogh"
// 	buff := *bytes.NewBuffer([]byte{})
// 	enc := gob.NewEncoder(&buff)
// 	if err := enc.Encode(file); err != nil {
// 		return fmt.Errorf("encode data: %w", err)
// 	}
// 	if err := os.WriteFile(filename, buff.Bytes(), 0644); err != nil {
// 		return fmt.Errorf("save to file %s: %w", filename, err)
// 	}
// 	return nil
// }

// func (s *Storage) Load() (*models.Data, error) {
// 	var data models.Data

// 	readbuff, err := os.ReadFile(s.filename)
// 	if err != nil {
// 		return nil, fmt.Errorf("read data from file %s: %w", s.filename, err)
// 	}
// 	buff := *bytes.NewBuffer(readbuff)
// 	dec := gob.NewDecoder(&buff)
// 	if err := dec.Decode(&data); err != nil {
// 		return nil, fmt.Errorf("decode data: %w", err)
// 	}
// 	return &data, nil
// }

// func (s *Storage) LoadFile(filename string) (*models.File, error) {
// 	var file models.File

// 	readbuff, err := os.ReadFile(filename)
// 	if err != nil {
// 		return nil, fmt.Errorf("read data from file %s: %w", filename, err)
// 	}
// 	buff := *bytes.NewBuffer(readbuff)
// 	dec := gob.NewDecoder(&buff)
// 	if err := dec.Decode(&file); err != nil {
// 		return nil, fmt.Errorf("decode data: %w", err)
// 	}
// 	return &file, nil
// }

func (s *Storage) loadFromFile(filename string, data interface{}) error {
	readbuff, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("read data from file %s: %w", filename, err)
	}
	buff := bytes.NewBuffer(readbuff)
	dec := gob.NewDecoder(buff)
	if err := dec.Decode(data); err != nil {
		return fmt.Errorf("decode data: %w", err)
	}
	return nil
}

func (s *Storage) Load() (*models.Data, error) {
	var data models.Data
	if err := s.loadFromFile(s.filename, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

func (s *Storage) LoadFile(filename string) (*models.File, error) {
	var file models.File
	if err := s.loadFromFile(filename, &file); err != nil {
		return nil, err
	}
	return &file, nil
}
