package storage

import (
	"errors"
	"os"
	"path"
	"track/layout"
	"track/types"
)

const defaultDataDir = ".tracker"
const defaultDataFile = "data.json"

type FileStorage struct {
	filepath string
	layout   types.Layout
	hasData  bool
}

func NewDefaultFileStorage() *FileStorage {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll(path.Join(home, defaultDataDir), os.ModePerm)
	if err != nil {
		panic(err)
	}
	dataPath := path.Join(home, defaultDataDir, defaultDataFile)
	hasData := false
	stat, err := os.Stat(dataPath)
	if errors.Is(err, os.ErrNotExist) {
		file, err := os.Create(dataPath)
		if err != nil {
			panic(err)
		}
		defer file.Close()
	} else {
		if stat.IsDir() {
			panic("data file is a directory")
		}
		hasData = true
	}
	return &FileStorage{
		filepath: dataPath,
		layout:   layout.JSONLayout,
		hasData:  hasData,
	}
}

func (f FileStorage) Save(intervals types.Intervals) error {
	data, err := f.layout.FromIntervals(intervals)
	if err != nil {
		return err
	}
	bytes, err := data.MarshalBinary()
	if err != nil {
		return err
	}
	return os.WriteFile(f.filepath, bytes, 600)
}

func (f FileStorage) Load() (types.Intervals, error) {
	if !f.hasData {
		return types.Intervals{}, nil
	}
	data, err := os.ReadFile(f.filepath)
	if err != nil {
		return types.Intervals{}, err
	}
	return f.layout.ToIntervals(types.Bytes(data))
}
