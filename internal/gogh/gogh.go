package gogh

import (
	"fmt"
	"gogh/internal/database"
	"gogh/internal/models"
)

type Gogh struct {
	storage *database.Storage
	Data    *models.Data
}

func New() (*Gogh, error) {
	_storage, err := database.New()
	if err != nil {
		return nil, fmt.Errorf("init gogh database: %w", err)
	}
	_data, err := _storage.Load()
	if err != nil {
		return nil, fmt.Errorf("init gogh: %w", err)
	}
	return &Gogh{
		storage: _storage,
		Data:    _data,
	}, nil
}

func (g *Gogh) SaveData() error {
	return g.storage.Save(g.Data)
}

func (g *Gogh) LoadData() error {
	_data, err := g.storage.Load()
	if err != nil {
		return err
	}
	g.Data = _data
	return nil
}
