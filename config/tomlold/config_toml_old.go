package tomlold

import (
	"fmt"

	"github.com/pelletier/go-toml"
)

type Reader struct {
	conf    *toml.Tree
	loadErr error
}

func NewReader() *Reader {
	conf, err := toml.LoadFile("config.toml")
	return &Reader{conf: conf, loadErr: err}
}

func NewReaderWithPath(path string) *Reader {
	conf, err := toml.LoadFile(path)
	return &Reader{conf: conf, loadErr: err}
}

func (r *Reader) Reload() error {
	var err error
	r.conf, err = toml.LoadFile("config.toml")
	return err
}

func (r *Reader) Get(key string) (interface{}, error) {
	if r.loadErr != nil {
		return nil, r.loadErr
	}

	v := r.conf.Get(key)
	if v == nil {
		return nil, fmt.Errorf("no config value for %s", key)
	}
	return v, nil
}

func (r *Reader) GetDef(key string, def interface{}) (interface{}, error) {
	if r.loadErr != nil {
		return def, r.loadErr
	}

	val := r.conf.Get(key)
	if val == nil {
		return def, nil
	}
	return val, nil
}
