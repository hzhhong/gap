package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync/atomic"

	"gopkg.in/yaml.v3"
)

type Config interface {
	Load() error
	Value(key string) (interface{}, error)
}

type config struct {
	path   string
	values map[string]interface{}
}

// keyvalue is config key value.
type keyvalue struct {
	key    string
	value  []byte
	format string
}

// New create Config
func New(filepath string) Config {

	cfg := &config{
		path:   filepath,
		values: make(map[string]interface{}),
	}
	return cfg
}

// Load load config from file
func (f *config) Load() error {
	fi, err := os.Stat(f.path)
	if err != nil {
		return err
	}

	if fi.IsDir() {
		return fmt.Errorf("incorrect config file path: %s", f.path)
	}
	kv, err := f.loadFile(f.path)
	if err != nil {
		return err
	}

	if err := decoder(kv, f.values); err != nil {
		return err
	}

	return nil
}

func decoder(src *keyvalue, target map[string]interface{}) error {
	if src.format == "" {
		// expand key "aaa.bbb" into map[aaa]map[bbb]interface{}
		keys := strings.Split(src.key, ".")
		for i, k := range keys {
			if i == len(keys)-1 {
				target[k] = src.value
			} else {
				sub := make(map[string]interface{})
				target[k] = sub
				target = sub
			}
		}
		return nil
	}

	if strings.ToLower(src.format) != "yaml" {
		return fmt.Errorf("unsupported key: %s format: %s", src.key, src.format)
	}
	return yaml.Unmarshal(src.value, &target)
}

func (f *config) loadFile(path string) (*keyvalue, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	info, err := file.Stat()
	if err != nil {
		return nil, err
	}
	return &keyvalue{
		key:    info.Name(),
		format: format(info.Name()),
		value:  data,
	}, nil
}

func format(name string) string {
	if p := strings.Split(name, "."); len(p) > 1 {
		return p[len(p)-1]
	}
	return ""
}

// Value get value from config by key
func (c *config) Value(key string) (interface{}, error) {
	if v, ok := readValue(c.values, key); ok {
		return v, nil
	}
	return nil, fmt.Errorf("config key not found: %s", key)
}

// readValue read Value in given map[string]interface{}
// by the given path, will return false if not found.
func readValue(values map[string]interface{}, path string) (interface{}, bool) {
	var (
		next = values
		keys = strings.Split(path, ".")
		last = len(keys) - 1
	)
	for idx, key := range keys {
		value, ok := next[key]
		if !ok {
			return nil, false
		}
		if idx == last {
			av := &atomic.Value{}
			av.Store(value)
			return av, true
		}
		switch vm := value.(type) {
		case map[string]interface{}:
			next = vm
		default:
			return nil, false
		}
	}
	return nil, false
}
