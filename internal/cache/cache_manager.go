package cache

import (
	"bytes"
	"encoding/gob"
	"errors"
	"log"
	"os"
	"time"
)

const CACHE_DURATION = 3600

func GetData[K interface{}](filename string, f func() K) K {
	data, err := getCache[K](filename)
	if err != nil {
		data = f()
		dump(filename, data)
	}
	return data
}

func getCache[K interface{}](filename string) (K, error) {
	var ret K
	infos, err := os.Stat("ext/cache/" + filename + ".bin")
	if err != nil {
		return ret, errors.New("No cache found")
	}
	if time.Now().Sub(infos.ModTime()).Seconds() > CACHE_DURATION {
		return ret, errors.New("Cache too old")
	}
	ret, err = pick[K](filename)
	if err != nil {
		return ret, err
	}
	return ret, nil
}

func dump(filename string, data any) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(data)
	if err != nil {
		log.Fatal("Error encoding cache" + err.Error())
	}
	f, err := os.Create("ext/cache/" + filename + ".bin")
	if err != nil {
		log.Fatal("Couldn't open cache file " + err.Error())
	}
	defer f.Close()
	f.Write(buf.Bytes())
}

func pick[K interface{}](filename string) (K, error) {
	var ret K
	f, err := os.ReadFile("ext/cache/" + filename + ".bin")
	if err != nil {
		return ret, err
	}
	dec := gob.NewDecoder(bytes.NewReader(f))
	err = dec.Decode(&ret)
	if err != nil {
		return ret, err
	}
	return ret, nil
}
