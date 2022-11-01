package config

import (
	"fmt"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type RawEndpointsConfig struct {
	Docker    bool                `fig:"docker"`
	Endpoints []rawEndpointConfig `fig:"services"`
}

type rawEndpointConfig struct {
	Endpoint    string   `fig:"service"`
	EntryPoints []string `fig:"entry_points"`
}

type EndpointsConfig struct {
	Endpoints map[string]string
}

type EndpointsConfiger interface {
	EndpointsConfig() *EndpointsConfig
}

func NewEndpointConfiger(getter kv.Getter) EndpointsConfiger {
	return &endpointsConfig{
		getter: getter,
	}
}

type endpointsConfig struct {
	getter kv.Getter
	once   comfig.Once
}

func (c *endpointsConfig) EndpointsConfig() *EndpointsConfig {
	return c.once.Do(func() interface{} {
		raw := kv.MustGetStringMap(c.getter, "Endpoints")
		config, err := parseRawEndpointsConfig(raw)
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out"))
		}
		return &config
	}).(*EndpointsConfig)
}

func parseRawEndpointsConfig(raw map[string]interface{}) (EndpointsConfig, error) {
	services := raw["services"]
	docker := raw["docker"].(bool)
	resConf := EndpointsConfig{}
	resConf.Endpoints = make(map[string]string)
	for _, service := range services.([]interface{}) {
		mapService := service.(map[interface{}]interface{})
		serviceName := mapService["service"].(string)
		serviceEntries := parseEntryPoints(mapService["entry_points"].([]interface{}))
		if docker {
			resConf.Endpoints[serviceName] = fmt.Sprintf("http://%v", serviceEntries[1])
		} else {
			resConf.Endpoints[serviceName] = fmt.Sprintf("http://%v", serviceEntries[0])
		}
	}

	return resConf, nil
}

func parseEntryPoints(entryPoints []interface{}) []string {
	s := make([]string, len(entryPoints))
	for i, v := range entryPoints {
		s[i] = fmt.Sprint(v)
	}
	return s
}
