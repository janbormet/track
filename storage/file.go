package storage

import (
	"errors"
	"os"
	"path"
	"track/layout"
	"track/types"
)

type FileStorage struct {
	filepath string
	layout   types.Layout
	hasData  bool
}

func createFileIfNotExists(dataDirPath string, dataFileName string) bool {
	err := os.MkdirAll(dataDirPath, os.ModePerm)
	if err != nil {
		panic(err)
	}
	dataPath := path.Join(dataDirPath, dataFileName)
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
	return hasData
}

func NewDefaultJSONFileStorage() *FileStorage {
	const defaultDataDir = ".tracker"
	const defaultDataFileName = "data.json"
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	hasData := createFileIfNotExists(path.Join(home, defaultDataDir), defaultDataFileName)
	return &FileStorage{
		filepath: path.Join(home, defaultDataDir, defaultDataFileName),
		layout:   layout.JSONLayout,
		hasData:  hasData,
	}
}

func NewDefaultTimewFileStorage() *FileStorage {
	const defaultDataDir = ".tracker"
	const defaultDataFileName = "data.timew"
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	hasData := createFileIfNotExists(path.Join(home, defaultDataDir), defaultDataFileName)
	return &FileStorage{
		filepath: path.Join(home, defaultDataDir, defaultDataFileName),
		layout:   layout.TimewLayout,
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
