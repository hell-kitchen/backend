package config

import "fmt"

type Controller struct {
	BindAddress string `config:"bind-address,short=a"`
	BindPort    int    `config:"bind-port,short=p"`
}

func (c Controller) GetBindAddress() string {
	return fmt.Sprintf("%s:%d", c.BindAddress, c.BindPort)
}
