package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	type simpleConfig struct {
		A string
		B bool
		C float64
		D int
		E string
	}

	cases := []struct {
		directory string
		env       string
		envVar    map[string]string
		host      string
		input     interface{}
		expect    interface{}
	}{{
		directory: "simple",
		input:     &simpleConfig{},
		expect: &simpleConfig{
			A: "foo",
			B: true,
			C: 123.34,
			D: 111,
			E: "default",
		},
	}, {
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
	}}

	for _, test := range cases {
		opts := &Options{
			Environment: test.env,
			Directory:   "./testConfigs/" + test.directory,
		}
		err := Load(opts, test.input)
		require.Nil(t, err)
		require.Equal(t, test.expect, test.input)
	}
}
