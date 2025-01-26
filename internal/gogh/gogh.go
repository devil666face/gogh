package gogh

import (
	"fmt"

	"gogh/internal/config"
	"gogh/internal/database"
	"gogh/internal/models"
	"gogh/internal/service/localfile"
	"gogh/internal/service/remotefile"

	gh "github.com/j178/github-s3"
)

type Gogh struct {
	Data    *models.Data
	config  *config.Config
	storage *database.Storage
	github  *gh.GitHub
}

func New() (*Gogh, error) {
	_config, err := config.New()
	if err != nil {
		return nil, fmt.Errorf("init config: %w", err)
	}
	_storage, err := database.New(_config.Database)
	if err != nil {
		return nil, fmt.Errorf("init gogh database: %w", err)
	}
	_data, err := _storage.Load()
	if err != nil {
		return nil, fmt.Errorf("init gogh storage: %w", err)
	}
	return &Gogh{
		Data:    _data,
		config:  _config,
		storage: _storage,
		github: gh.New(
			gh.Credential{
				UserSession: _data.Settings.SessionToken,
				DeviceID:    _data.Settings.DeviceID,
			},
			""),
	}, nil
}

func (g *Gogh) SetToken(token string) error {
	g.Data.Settings.SessionToken = token
	g.github = gh.New(gh.Credential{
		UserSession: g.Data.Settings.SessionToken,
		DeviceID:    g.Data.Settings.DeviceID,
	}, "")
	return g.SaveData()
}

func (g *Gogh) SetDeviceID(id string) error {
	g.Data.Settings.DeviceID = id
	g.github = gh.New(gh.Credential{
		UserSession: g.Data.Settings.SessionToken,
		DeviceID:    g.Data.Settings.DeviceID,
	}, "")
	return g.SaveData()
}

func (g *Gogh) SaveData() error {
	return g.storage.Save(g.Data)
}

func (g *Gogh) SaveFile(file *models.File) error {
	return g.storage.SaveFile(file)
}

func (g *Gogh) LoadData() error {
	_data, err := g.storage.Load()
	if err != nil {
		return err
	}
	g.Data = _data
	return nil
}

func (g *Gogh) LoadFile(filename string) (*models.File, error) {
	return g.storage.LoadFile(filename)
}

func (g *Gogh) Upload(path string, compress bool) error {
	file, err := localfile.New(
		g.Data.Filestore.ID(),
		path,
		compress,
	)
	if err != nil {
		return fmt.Errorf("init local file: %w", err)
	}
	defer file.Clear()

	for _, piece := range file.Pieces {
		res, err := g.github.UploadFromPath(piece.Filename)
		if err != nil {
			return err
		}
		g.Data.Filestore.Add(
			file.Id,
			file.Filname,
			file.Size,
			file.Compress,
			piece.Key,
			res.GithubLink,
		)
	}

	if err := g.SaveData(); err != nil {
		return err
	}
	return nil
}

func (g *Gogh) Delete(id int) error {
	if err := g.Data.Filestore.Delete(id); err != nil {
		return err
	}
	if err := g.SaveData(); err != nil {
		return err
	}
	return nil
}

func (g *Gogh) Download(id int) error {
	_file := g.Data.Filestore.Files[id]
	file, err := remotefile.New(
		id,
		_file.Filename,
		_file.Compress,
		_file.Pieces)
	if err != nil {
		return fmt.Errorf("init remote file: %w", err)
	}
	defer file.Clear()
	if err := file.Download(); err != nil {
		return err
	}
	return nil
}

// func (g *Gogh) UploadParalel(path string) error {
// 	_file, err := localfile.New(g.Data.Filestore.ID(), path, true)
// 	if err != nil {
// 		return fmt.Errorf("init file: %w", err)
// 	}
// 	defer _file.Clear()

// 	type Result struct {
// 		res string
// 		err error
// 	}

// 	var wg sync.WaitGroup
// 	reschan := make(chan Result, len(_file.Pieces))

// 	for _, f := range _file.Pieces {
// 		wg.Add(1)
// 		go func(f string) {
// 			defer wg.Done()
// 			log.Println(f)
// 			res, err := g.github.UploadFromPath(f)
// 			log.Println(res.GithubLink)
// 			reschan <- Result{
// 				res: res.GithubLink,
// 				err: err,
// 			}
// 		}(f)
// 	}
// 	wg.Wait()
// 	close(reschan)

// 	var results []string
// 	for r := range reschan {
// 		if r.err != nil {
// 			return r.err
// 		}
// 		results = append(results, r.res)
// 	}

// 	log.Println(results)
// 	return nil
// }
