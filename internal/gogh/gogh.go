package gogh

import (
	"fmt"
	"gogh/internal/database"
	"gogh/internal/models"
	"gogh/internal/service/file"

	gh "github.com/j178/github-s3"
)

type Gogh struct {
	Data    *models.Data
	storage *database.Storage
	github  *gh.GitHub
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
		Data:    _data,
		storage: _storage,
		github:  gh.New(_data.Settings.SessionToken, ""),
	}, nil
}

func (g *Gogh) SaveData() error {
	return g.storage.Save(g.Data)
}

func (g *Gogh) SetToken(token string) {
	g.Data.Settings.SessionToken = token
	g.github = gh.New(g.Data.Settings.SessionToken, "")
}

func (g *Gogh) LoadData() error {
	_data, err := g.storage.Load()
	if err != nil {
		return err
	}
	g.Data = _data
	return nil
}

func (g *Gogh) Upload(path string) error {
	_file, err := file.New(path)
	if err != nil {
		return fmt.Errorf("init file: %w", err)
	}
	defer _file.Clear()
	for _, p := range _file.Pieces {
		res, err := g.github.UploadFromPath(p)
		if err != nil {
			return fmt.Errorf("failed to upload on github %s: %w", path, err)
		}
		fmt.Println(res.GithubLink)
	}
	// g.Data.Filestore.AddFile(path)
	// if err := g.SaveData(); err != nil {
	// 	return fmt.Errorf("failed to save data: %w", err)
	// }
	return nil
}
