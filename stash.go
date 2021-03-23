package froth

import (
	"bytes"
	"encoding/gob"
	"github.com/boggydigital/kvas"
)

type Stash struct {
	dst       string
	asset     string
	keyValues map[string]interface{}
}

func NewStash(dst, asset string) (*Stash, error) {
	kvStash, err := kvas.NewGobLocal(dst)
	if err != nil {
		return nil, err
	}

	stashRC, err := kvStash.Get(asset)
	if err != nil {
		return nil, err
	}

	var keyValues map[string]interface{}

	if stashRC != nil {
		defer stashRC.Close()
		if err := gob.NewDecoder(stashRC).Decode(&keyValues); err != nil {
			return nil, err
		}
	}

	if keyValues == nil {
		keyValues = make(map[string]interface{}, 0)
	}

	return &Stash{
		dst:       dst,
		asset:     asset,
		keyValues: keyValues,
	}, nil
}

func (stash *Stash) All() []string {
	keys := make([]string, 0, len(stash.keyValues))
	for k := range stash.keyValues {
		keys = append(keys, k)
	}
	return keys
}

func (stash *Stash) set(key string, value interface{}) error {
	stash.keyValues[key] = value
	return stash.write()
}

func (stash *Stash) SetString(key string, value string) error {
	return stash.set(key, value)
}

func (stash *Stash) SetStringSlice(key string, values []string) error {
	return stash.set(key, values)
}

func (stash *Stash) SetInt(key string, value int) error {
	return stash.set(key, value)
}

func (stash *Stash) write() error {
	kvStash, err := kvas.NewGobLocal(stash.dst)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(stash.keyValues); err != nil {
		return err
	}

	return kvStash.Set(stash.asset, buf)
}

func (stash *Stash) setMany(keyValues map[string]interface{}) error {
	for k, v := range keyValues {
		stash.keyValues[k] = v
	}
	return stash.write()
}

func (stash *Stash) SetManyStrings(keyValues map[string]string) error {
	for k, v := range keyValues {
		stash.keyValues[k] = v
	}
	return stash.write()
}

func (stash *Stash) SetManyStringSlices(keyValues map[string][]string) error {
	for k, v := range keyValues {
		stash.keyValues[k] = v
	}
	return stash.write()
}

func (stash *Stash) SetManyInts(keyValues map[string]int) error {
	for k, v := range keyValues {
		stash.keyValues[k] = v
	}
	return stash.write()
}

func (stash *Stash) get(key string) (interface{}, bool) {
	if stash == nil || stash.keyValues == nil {
		return "", false
	}
	val, ok := stash.keyValues[key]
	return val, ok
}

func (stash *Stash) GetString(key string) (string, bool) {
	val, ok := stash.get(key)
	return val.(string), ok
}

func (stash *Stash) GetStringSlice(key string) ([]string, bool) {
	val, ok := stash.get(key)
	return val.([]string), ok
}

func (stash *Stash) GetInt(key string) (int, bool) {
	val, ok := stash.get(key)
	return val.(int), ok
}
