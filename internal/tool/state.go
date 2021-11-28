package tool

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

var (
	SaveFolder string = ".fdlr/"
	stateFile  string = "state.yaml"
)

type State struct {
	URL            string
	DownloadRanges []DownloadRange
}

type DownloadRange struct {
	URL       string
	Path      string
	RangeFrom int64
	RangeTo   int64
}

func (state *State) Save() error {
	// make file Folder, to save all downloaded parts and state.yaml
	folder, err := GetFolderFrom(state.URL)
	if err != nil {
		return errors.WithStack(err)
	}
	fmt.Printf("Saving states data in %s\n", folder)
	err = Mkdir(folder)
	if err != nil {
		return errors.WithStack(err)
	}

	// move current downloading parts to file folder
	for _, part := range state.DownloadRanges {
		err = os.Rename(part.Path, filepath.Join(folder, filepath.Base(part.Path)))
		if err != nil {
			return errors.WithStack(err)
		}
	}

	// save state to state.yaml
	y, err := yaml.Marshal(state)
	if err != nil {
		return errors.WithStack(err)
	}

	return ioutil.WriteFile(filepath.Join(folder, stateFile), y, 0644)
}

func Read(task string) (*State, error) {
	file := filepath.Join(os.Getenv("HOME"), SaveFolder, task, stateFile)
	fmt.Printf("Reading state from %s\n", file)

	var err error

	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	state := &State{}
	if err = yaml.Unmarshal(bytes, state); err != nil {
		return nil, errors.WithStack(err)
	}

	return state, nil
}
