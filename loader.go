package config

import (
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

type Options struct {
	Environment string
	Directory   string
}

func applyFile(file string, out interface{}) error {
	fileContents, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	return apply(fileContents, out)
}

func applyEnvVars(file string, out interface{}) error {
	fileContents, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	// Read the yaml file. Assume every value is a string that indicates the
	// env var that should be used for the corresponding key.

	// This is hacky, and likely only supports a subset of yaml syntax
	// I couldn't find a parsing package after a pretty quick search,
	// and this is good enough for my needs for now.

	withReplacements := ""

	lines := strings.Split(string(fileContents), "\n")
	for _, line := range lines {
		envVar := ""
		prefix := ""
		if strings.Contains(line, ":") {
			rexp := regexp.MustCompile("^(.*):\\s+((?:\\\".*\\\"\\s+)|(?:.*))$")
			matches := rexp.FindStringSubmatch(line)
			if len(matches) > 0 {
				prefix = matches[1]
				envVar = strings.Trim(matches[2], "\" \t")
			}
		}

		if envVar != "" {
			val := os.Getenv(envVar)
			if val != "" {
				line = prefix + ": " + val
			} else {
				line = ""
			}
		}

		withReplacements = withReplacements + "\n" + line
	}

	return apply([]byte(withReplacements), out)
}

func apply(yamlBuf []byte, out interface{}) error {
	return yaml.Unmarshal(yamlBuf, out)
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
			err := applyFile(directory+"/"+fileName, out)
			if err != nil {
				return err
			}
		}
	}

	fileName, exists := fileHash["custom-environment-variables."]
	if exists {
		err := applyEnvVars(directory+"/"+fileName, out)
		if err != nil {
			return err
		}
	}

	return nil
}
