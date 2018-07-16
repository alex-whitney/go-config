package config

import (
	"io/ioutil"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

type Options struct {
	Environment string
	Directory   string
}

func readFile(file string, out interface{}) error {
	fileContents, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(fileContents, out)
	if err != nil {
		return err
	}

	return nil
}

func generateFilePrefixes(deployment string) []string {
	loadOrder := []string{
		"default.",
		"{deployment}.",
		"local.",
		"local-{deployment}.",
	}

	for idx, pattern := range loadOrder {
		loadOrder[idx] = strings.Replace(pattern, "{deployment}", deployment, -1)
	}

	return loadOrder
}

func Load(opts *Options, out interface{}) error {
	loadOrder := generateFilePrefixes(opts.Environment)

	directory := "./config"
	if opts.Directory != "" {
		directory = opts.Directory
	}

	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return err
	}

	fileHash := map[string]string{}
	for _, file := range files {
		fileName := file.Name()

		// strip extension off the key
		rexp := regexp.MustCompile("^(.*\\.).*$")
		result := rexp.FindStringSubmatch(fileName)
		if len(result) > 0 {
			fileHash[result[1]] = fileName
		}
	}

	for _, prefix := range loadOrder {
		fileName, exists := fileHash[prefix]
		if exists {
			err := readFile(directory+"/"+fileName, out)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
