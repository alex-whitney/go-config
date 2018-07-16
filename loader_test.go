package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	type simpleConfig struct {
		A string
		B bool
		C float64
		D int
		E string
	}

	cases := map[string]struct {
		directory string
		env       string
		envVar    map[string]string
		host      string
		input     interface{}
		expect    interface{}
	}{
		"defaults": {
			directory: "simple",
			input:     &simpleConfig{},
			expect: &simpleConfig{
				A: "foo",
				B: true,
				C: 123.34,
				D: 111,
				E: "default",
			},
		},
		"env overrides": {
			directory: "simple",
			env:       "production",
			input:     &simpleConfig{},
			expect: &simpleConfig{
				A: "foo",
				B: true,
				C: 123.34,
				D: 111,
				E: "production",
			},
		},
		"local & env overrides": {
			directory: "local",
			env:       "production",
			input:     &simpleConfig{},
			expect: &simpleConfig{
				A: "newer",
				C: 22,
				D: 33,
				E: "local-prod",
			},
		},
		"env var overrides": {
			directory: "env-var",
			input:     &simpleConfig{},
			envVar: map[string]string{
				"ENV_C": "55.555",
			},
			expect: &simpleConfig{
				A: "foo",
				C: 55.555,
			},
		},
		"missing env var overrides": {
			directory: "env-var",
			input:     &simpleConfig{},
			expect: &simpleConfig{
				A: "foo",
				C: 123.34,
			},
		},
	}

	for name, test := range cases {
		opts := &Options{
			Environment: test.env,
			Directory:   "./testConfigs/" + test.directory,
		}

		originalEnv := map[string]string{}
		for k, v := range test.envVar {
			originalEnv[k] = os.Getenv(k)
			os.Setenv(k, v)
		}

		err := Load(opts, test.input)
		assert.Nil(t, err, name)
		assert.Equal(t, test.expect, test.input, name)

		for k, v := range originalEnv {
			os.Setenv(k, v)
		}
	}
}
