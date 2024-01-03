package persistence

import (
	"encoding/gob"
	"encoding/json"
	"io"
	"os"
)

type Persistence interface {
	Encode(v any)
	Decode(v any)
}

type JSON string

func (receiver JSON) Encode(v any) {
	encode(string(receiver), func(f *os.File) {
		bytes, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			panic(err)
		}
		f.Truncate(0)
		f.Write(bytes)
	})
}

func (receiver JSON) Decode(v any) {
	decode(string(receiver), func(f *os.File) {
		bytes, err := io.ReadAll(f)
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(bytes, v)
		if err != nil {
			panic(err)
		}
	})
}

type GOB string

func (receiver GOB) Encode(v any) {
	encode(string(receiver), func(f *os.File) {
		encoder := gob.NewEncoder(f)
		err := encoder.Encode(v)
		if err != nil {
			panic(err)
		}
	})
}

func (receiver GOB) Decode(v any) {
	decode(string(receiver), func(f *os.File) {
		decoder := gob.NewDecoder(f)
		err := decoder.Decode(v)
		if err != nil {
			panic(err)
		}
	})
}

func encode(filename string, fn func(f *os.File)) {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	fn(f)
}

func decode(filename string, fn func(f *os.File)) {
	if exist(filename) {
		f, err := os.OpenFile(filename, os.O_RDONLY, 0777)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		fn(f)
	}
}

func exist(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}
