package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type config struct {
	RootURL         string
	ListenAddress   string
	DefaultRedirect string
	DatabasePath    string
	AdminKey        string
	Public          bool
}

func writeDefaultConfig(path string) (*config, error) {
	dc := config{
		RootURL:         "",
		ListenAddress:   ":8080",
		DefaultRedirect: "https://google.com",
		DatabasePath:    "shorty.db3",
		AdminKey:        "",
		Public:          true,
	}

	conf, err := json.MarshalIndent(dc, "", "  ")
	if err != nil {
		return nil, err
	}

	return &dc, ioutil.WriteFile(path, conf, 0600)
}

func getConfig(path string) (*config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var conf config
	err = json.Unmarshal(data, &conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}

func loadConfig(path string) (*config, error) {
	conf, err := getConfig(path)
	if err != nil {
		err = os.MkdirAll(filepath.Dir(path), 0766)
		if err != nil {
			return nil, err
		}

		_, err = os.Create(path)
		if err != nil {
			return nil, err
		}

		conf, err = writeDefaultConfig(path)
		if err != nil {
			return nil, err
		}
	}

	return conf, err
}
