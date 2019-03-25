package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Logger  Logger  `yaml:"logger"`
	Network Network `yaml:"network"`
}

type Logger struct {
	Type string `yaml:"type"`
}

type Network struct {
	Name            string `yaml:"name"`
	Type            string `yaml:"type"`
	Address         string `yaml:"address"`
	Cert            string `yaml:"cert"`
	Macaroon        string `yaml:"macaroon"`
	MacaroonTimeOut int64  `yaml:"macaroon_timeout"`
	MacaroonIP      string `yaml:"macaroon_ip"`
	MaxMsgRecvSize  int    `yaml:"max_msg_recv_size"`
	ConnTimeout     int    `yaml:"conn_timeout"`
	PoolCapacity    int    `yaml:"pool_capacity"`
}

func Load(path string) (*Config, error) {
	c := &Config{}

	if path == "" {
		dir, err := getAppDir()
		if err != nil {
			return nil, err
		}
		path = fmt.Sprintf("%s/config.yml", dir)
	}

	err := loadFromPath(path, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// loadFromPath loads the configuration from configuration file path.
func loadFromPath(path string, out interface{}) error {
	var err error

	f, err := os.Open(path)
	if f != nil {
		defer func() {
			ferr := f.Close()
			if ferr != nil {
				log.Println(ferr)
			}
		}()
	}

	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, out)
}

// getappDir creates if not exists the app directory where the config file
// as well as the log file will be stored. In case of failure the current dir
// will be used.
func getAppDir() (string, error) {
	usr, _ := user.Current()
	dir := path.Join(usr.HomeDir, ".lntop")
	_, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			oserr := os.Mkdir(dir, 0700)
			if oserr != nil {
				return "", oserr
			}
		} else {
			return "", err
		}
	}
	return dir, nil
}
