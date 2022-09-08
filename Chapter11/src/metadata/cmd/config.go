package main

type config struct {
	API        apiConfig        `yaml:"api"`
	Jaeger     jaegerConfig     `yaml:"jaeger"`
	Prometheus prometheusConfig `yaml:"prometheus"`
}

type apiConfig struct {
	Port int `yaml:"port"`
}

type jaegerConfig struct {
	URL string `yaml:"url"`
}

type prometheusConfig struct {
	MetricsPort int `yaml:"metricsPort"`
}
