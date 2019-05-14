package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Logger  Logger  `toml:"logger"`
	Network Network `toml:"network"`
	Views   Views   `toml:"views"`
}

type Logger struct {
	Type string `toml:"type"`
	Dest string `toml:"dest"`
}

type Network struct {
	Name            string `toml:"name"`
	Type            string `toml:"type"`
	Address         string `toml:"address"`
	Cert            string `toml:"cert"`
	Macaroon        string `toml:"macaroon"`
	MacaroonTimeOut int64  `toml:"macaroon_timeout"`
	MacaroonIP      string `toml:"macaroon_ip"`
	MaxMsgRecvSize  int    `toml:"max_msg_recv_size"`
	ConnTimeout     int    `toml:"conn_timeout"`
	PoolCapacity    int    `toml:"pool_capacity"`
}

type Views struct {
	Channels     *View `toml:"channels"`
	Transactions *View `toml:"transactions"`
}

type View struct {
	Columns []string `toml:"columns"`
}

func Load(path string) (*Config, error) {
	c := &Config{}

	if path == "" {
		dir, err := getAppDir()
		if err != nil {
			return nil, err
		}
		path = fmt.Sprintf("%s/config.toml", dir)
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

	return toml.Unmarshal(data, out)
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
			err := os.Mkdir(dir, 0700)
			if err != nil {
				return "", err
			}
			err = ioutil.WriteFile(dir+"/config.toml",
				[]byte(DefaultFileContent()), 0644)
			if err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}
	return dir, nil
}
