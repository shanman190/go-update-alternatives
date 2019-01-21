package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v2"
)

type Alternative struct {
	Path string
	Priority int64
}

type configuration struct {
	Link string
	Alternatives []Alternative
}

func SaveAlternative(
	link string,
	groupName string,
	path string,
	priority int64,
) error {
	groupConfigFile := filepath.Join(GetAdminDir(), groupName)

	var alternatives *configuration
	if _, err := os.Stat(groupConfigFile); err == nil {
		alternatives, err = LoadAlternatives(groupName)
		if err != nil {
			return err
		}
	} else if os.IsNotExist(err) {
		alternatives = &configuration{
			Link: link,
			Alternatives: []Alternative{},
		}
	} else {
		return err
	}

	found := false
	for _, alternative := range alternatives.Alternatives {
		if alternative.Path == path {
			found = true
		}
	}

	if !found {
		newIndex := 0
		for index, alternative := range alternatives.Alternatives {
			if priority < alternative.Priority {
				break
			}
			newIndex = index + 1
		}
		if newIndex < len(alternatives.Alternatives) {
			alternatives.Alternatives = append(alternatives.Alternatives, Alternative{})
			copy(alternatives.Alternatives[(newIndex + 1):], alternatives.Alternatives[newIndex:])
			alternatives.Alternatives[newIndex] = Alternative{path, priority}	
		} else {
			alternatives.Alternatives = append(alternatives.Alternatives, Alternative{path, priority})
		}
		
		return writeAlternatives(groupConfigFile, alternatives)
	}

	return nil
}

func LoadAlternatives(groupName string) (*configuration, error) {
	var alternatives *configuration

	groupConfigFile := filepath.Join(GetAdminDir(), groupName)
	if _, err := os.Stat(groupConfigFile); err == nil {
		bytes, err := ioutil.ReadFile(groupConfigFile)
		if err != nil {
			return nil, err
		}
		if err = yaml.Unmarshal(bytes, &alternatives); err != nil {
			return nil, err
		}
	}
	
	if alternatives == nil {
		return &configuration{
			Alternatives: []Alternative{},
		}, nil
	}

	return alternatives, nil
}

func DeleteAlternative(
	groupName string,
	path string,
) error {
	groupConfigFile := filepath.Join(GetAdminDir(), groupName)

	alternatives, err := LoadAlternatives(groupName)
	if err != nil {
		return err
	}

	results := []Alternative{}
	for _, alternative := range alternatives.Alternatives {
		if alternative.Path != path {
			results = append(results, alternative)
		}
	}

	alternatives.Alternatives = results

	if len(alternatives.Alternatives) > 0 {
		return writeAlternatives(groupConfigFile, alternatives)
	}

	if err := os.Remove(groupConfigFile); err != nil {
		return fmt.Errorf("could not delete empty group file: '%s'", err)
	}

	return nil
}

func GetAlternativesDir() string {
	return filepath.Join(userHomeDir(), ".local/etc/alternatives")
}

func GetAdminDir() string {
	return filepath.Join(userHomeDir(), ".local/var/lib/alternatives")
}

func userHomeDir() string {
	home := os.Getenv("HOME")
	if home != "" {
		return home
	}

	if runtime.GOOS == "windows" {
		home = os.Getenv("USERPROFILE")
		if home != "" {
			return home
		}

		home = os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home != "" {
			return home
		}
	}

	panic("could not detect home directory for .alternativesrc")
}

func writeAlternatives(configFileLocation string, alternativesToWrite *configuration) error {
	dir, _ := filepath.Split(configFileLocation)
	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(dir, os.ModePerm); err != nil {
				return fmt.Errorf("could not create admin directory: %s", err)
			}
		} else {
			return fmt.Errorf("unknown error: %s", err)
		}
	}
	yamlBytes, err := yaml.Marshal(alternativesToWrite)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(configFileLocation, yamlBytes, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}