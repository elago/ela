package ela

import (
	"errors"
	"github.com/gogather/com"
	// "github.com/gogather/com/log"
	"path/filepath"
)

type Session struct {
	path string
}

func NewSession(path string) Session {
	return Session{
		path: path,
	}
}

func (this *Session) getSessionObject(sid string) (map[string]interface{}, error) {
	fullpath := this.getPath(sid)

	data, err := com.ReadFileByte(fullpath)
	if err != nil {
		return nil, err
	}

	to := map[string]interface{}{}
	err = com.Decode(data, &to)
	return to, err
}

func (this *Session) saveSession(sid string, object interface{}) error {
	fullpath := this.getPath(sid)
	str, err := com.Encode(object)
	err = com.WriteFileWithCreatePath(fullpath, string(str))
	return err
}

func (this *Session) Get(sid string, key string) (interface{}, error) {
	mapObject, err := this.getSessionObject(sid)

	if err != nil {
		return nil, err
	}

	value, ok := mapObject[key]

	if ok {
		return value, nil
	} else {
		return nil, errors.New("key value does not exist")
	}
}

func (this *Session) Set(sid string, key string, value interface{}) error {
	mapObject, _ := this.getSessionObject(sid)
	if mapObject == nil {
		mapObject = map[string]interface{}{}
	}
	mapObject[key] = value
	return this.saveSession(sid, mapObject)
}

func (this *Session) getPath(key string) string {
	length := len(key) - 3
	dir1 := com.SubString(key, 0, 1)
	dir2 := com.SubString(key, 1, 2)
	file := com.SubString(key, 3, length)
	return filepath.Join(this.path, dir1, dir2, file)
}
