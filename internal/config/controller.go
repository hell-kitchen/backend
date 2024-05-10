package config

import "fmt"

type Controller struct {
	BindAddress          string `config:"bind-address,short=a"`
	BindPort             int    `config:"bind-port,short=p"`
	AccessCookieName     string `config:"access-cookie-name"`
	AccessCookieLifetime int    `config:"access-cookie-lifetime"`
}

func (c Controller) GetBindAddress() string {
	return fmt.Sprintf("%s:%d", c.BindAddress, c.BindPort)
}
