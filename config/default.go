package config

import "fmt"

func DefaultFileContent() string {
	cfg := NewDefault()
	return fmt.Sprintf(`
[logger]
type = "%[1]s"
dest = "%[2]s"

[network]
name = "%[3]s"
type = "%[4]s"
address = "%[5]s"
cert = "%[6]s"
macaroon = "%[7]s"
macaroon_timeout = %[8]d
max_msg_recv_size = %[9]d
conn_timeout = %[10]d
pool_capacity = %[11]d
`,
		cfg.Logger.Type,
		cfg.Logger.Dest,
		cfg.Network.Name,
		cfg.Network.Type,
		cfg.Network.Address,
		cfg.Network.Cert,
		cfg.Network.Macaroon,
		cfg.Network.MacaroonTimeOut,
		cfg.Network.MaxMsgRecvSize,
		cfg.Network.ConnTimeout,
		cfg.Network.PoolCapacity,
	)
}

func NewDefault() *Config {
	return &Config{}
}
