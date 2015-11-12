package misc

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Path string
	Dict map[string]interface{}
}

func parseConfig(path string) (map[string]interface{}, error) {

	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := map[string]interface{}{}

	err = json.Unmarshal(contents, &config)

	if nil != err {

		return nil, err
	}
	return config, nil
}

func NewConfig(path string) (*Config, error) {
	conf := &Config{}
	conf.Path = path
	if err := conf.LoadConfig(); err != nil {
		return nil, err
	}
	return conf, nil
}

func (p *Config) LoadConfig() error {
	dict, err := parseConfig(p.Path)
	if err != nil {
		return err
	}
	p.Dict = dict
	return nil
}

func (p *Config) Reload() error {
	return p.LoadConfig()
}

func (p *Config) GetInt(key string, defalut ...int) int {
	val, ok := p.Dict[key].(float64)
	if !ok {
		if len(defalut) > 0 {
			return defalut[0]
		} else {
			return 0
		}
	}
	return int(val)
}

func (p *Config) GetString(key string, defalut ...string) string {
	val, ok := p.Dict[key].(string)
	if !ok {
		if len(defalut) > 0 {
			return defalut[0]
		} else {
			return ""
		}
	}
	return val
}

func (p *Config) GetBool(key string, defalut ...bool) bool {
	val, ok := p.Dict[key].(bool)
	if !ok {
		if len(defalut) > 0 {
			return defalut[0]
		} else {
			return false
		}
	}
	return val
}

func (p *Config) GetStrings(key string, defalut ...[]string) ([]string, int) {
	val, ok := p.Dict[key].([]interface{})
	if !ok {
		if len(defalut) > 0 {
			return defalut[0], len(defalut[0])
		} else {
			return nil, 0
		}
	}
	length := len(val)
	arr := make([]string, length)
	for i, v := range val {
		f, _ := v.(string)
		arr[i] = f
	}
	return arr, length
}

func (p *Config) GetInts(key string, defalut ...[]int) ([]int, int) {

	val, ok := p.Dict[key].([]interface{})
	if !ok {
		if len(defalut) > 0 {
			return defalut[0], len(defalut[0])
		} else {
			return nil, 0
		}
	}
	length := len(val)
	arr := make([]int, length)
	for i, v := range val {
		f, _ := v.(float64)
		arr[i] = int(f)
	}
	return arr, length
}
