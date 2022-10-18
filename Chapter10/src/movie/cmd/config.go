package main

type config struct {
	API apiConfig `yaml:"api"`
}

type apiConfig struct {
	Port int `yaml:"port"`
}
