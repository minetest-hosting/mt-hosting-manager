package server_setup

import "fmt"

type ComposeModel struct {
	MTUIVersion    string
	MTUIKey        string
	Hostname       string
	HTTPRouterName string
}

func Setup() {

	m := &ComposeModel{
		MTUIVersion:    "1.34",
		MTUIKey:        "my-secret-key",
		Hostname:       "node-prod-xxx.minetest.ch",
		HTTPRouterName: "router_xxx",
	}

	//TODO
	fmt.Println(m)
}
