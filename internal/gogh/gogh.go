package gogh

import (
	"fmt"
	"gogh/internal/database"
	"gogh/internal/models"
	"gogh/internal/service/file"
	"log"
	"sync"

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

func (g *Gogh) UploadParalel(path string) error {
	_file, err := file.New(path)
	if err != nil {
		return fmt.Errorf("init file: %w", err)
	}
	defer _file.Clear()

	type Result struct {
		res string
		err error
	}

	var wg sync.WaitGroup
	reschan := make(chan Result, len(_file.Pieces))

	for _, f := range _file.Pieces {
		wg.Add(1)
		go func(f string) {
			defer wg.Done()
			log.Println(f)
			res, err := g.github.UploadFromPath(f)
			log.Println(res.GithubLink)
			reschan <- Result{
				res: res.GithubLink,
				err: err,
			}
		}(f)
	}
	wg.Wait()
	close(reschan)

	var results []string
	for r := range reschan {
		if r.err != nil {
			return r.err
		}
		results = append(results, r.res)
	}

	log.Println(results)
	return nil
}

func (g *Gogh) Upload(path string) error {
	_file, err := file.New(path)
	if err != nil {
		return fmt.Errorf("init file: %w", err)
	}
	defer _file.Clear()

	for _, f := range _file.Pieces {
		res, err := g.github.UploadFromPath(f)
		if err != nil {
			return err
		}
		g.Data.Filestore.Add(
			_file.Id,
			_file.Filname,
			res.GithubLink,
			_file.Size,
		)
	}

	if err := g.SaveData(); err != nil {
		return err
	}

	return nil
}
