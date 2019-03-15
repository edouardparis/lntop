package config

import (
	"io/ioutil"
	"log"
	"os"

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
	ID              string `yaml:"id"`
	Type            string `yaml:"type"`
	Status          string `yaml:"status"`
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
