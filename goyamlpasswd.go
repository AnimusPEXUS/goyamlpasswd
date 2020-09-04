package goyamlpasswd

import (
	"errors"
	"io/ioutil"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

var (
	ERR_NOT_FOUND        = errors.New("not found")
	ERR_PASSWORD_NOT_SET = errors.New("password not sed")
)

type YAMLAuthFileSRecord struct {
	User       string
	Password   *string `yaml:",omitempty"`
	UnifiedKey *string `yaml:",omitempty"`
}

type YAMLAuthFileS struct {
	Records []*YAMLAuthFileSRecord
}

type YAMLAuthFile struct {
	filename string
	lastload *time.Time
	data     *YAMLAuthFileS
}

func NewYAMLAuthFile(filename string) *YAMLAuthFile {
	ret := &YAMLAuthFile{
		filename: filename,
		data:     &YAMLAuthFileS{},
	}
	return ret
}

func (self *YAMLAuthFile) LoadIfChanged() (bool, error) {

	var err error

	if self.lastload == nil {
		goto load
	}

	{
		stat, err := os.Stat(self.filename)
		if err != nil {
			return false, err
		}

		if !stat.ModTime().After(*self.lastload) {
			return false, nil
		}
	}

load:

	err = self.Load()
	if err != nil {
		return false, err
	}

	return true, nil
}

func (self *YAMLAuthFile) Load() error {

	data, err := ioutil.ReadFile(self.filename)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, self.data)
	if err != nil {
		return err
	}

	stat, err := os.Stat(self.filename)
	if err != nil {
		return err
	}

	self.lastload = &([]time.Time{stat.ModTime()}[0])

	return nil
}

func (self *YAMLAuthFile) Save() error {

	data, err := yaml.Marshal(self.data)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(self.filename, data, 0o700)
	if err != nil {
		return err
	}

	return nil
}

func (self *YAMLAuthFile) PutRecord(r *YAMLAuthFileSRecord) {
	for i := len(self.data.Records) - 1; i != -1; i += -1 {
		if self.data.Records[i].User == r.User {
			self.data.Records[i].Password = r.Password
			self.data.Records[i].UnifiedKey = r.UnifiedKey
		}
	}
	return
}

func (self *YAMLAuthFile) RemoveRecord(name string) {
	for i := len(self.data.Records) - 1; i != -1; i += -1 {
		if self.data.Records[i].User == name {
			self.data.Records = append(self.data.Records[:i], self.data.Records[i+1:]...)
		}
	}
	return
}

func (self *YAMLAuthFile) UserByKey(key string) (user string, err error) {
	for i := len(self.data.Records) - 1; i != -1; i += -1 {
		if *(self.data.Records[i].UnifiedKey) == key {
			return self.data.Records[i].User, nil
		}
	}
	return "", ERR_NOT_FOUND
}

func (self *YAMLAuthFile) CheckUserPasswordMatch(user, password string) (ok bool, err error) {
	for i := len(self.data.Records) - 1; i != -1; i += -1 {
		if self.data.Records[i].User == user {
			if self.data.Records[i].Password == nil {
				return false, ERR_PASSWORD_NOT_SET
			}
			return *(self.data.Records[i].Password) == password, nil
		}
	}
	return false, ERR_NOT_FOUND
}
