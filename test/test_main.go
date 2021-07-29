package main

import (
	"config-client/pkg/client"
	"github.com/paashzj/cc_api"
	"log"
	"unsafe"
)

type TestConfig struct {
	Name *string `config:"name"`
	Age  string `config:"age"`
	cc_api.BaseConfig
}

type TestConfigListener struct {
}

func (t *TestConfigListener) Add(value unsafe.Pointer) {
	val := (*TestConfig)(value)
	log.Printf("add value %s %d", val.Id, val.Version)
}

func (t *TestConfigListener) Mod(oldValue, value unsafe.Pointer) {
	oldVal := (*TestConfig)(oldValue)
	val := (*TestConfig)(value)
	log.Printf("old value version %d new value version %d", oldVal.Version, val.Version)
}

func (t *TestConfigListener) Del(value unsafe.Pointer) {
	val := (*TestConfig)(value)
	log.Printf("delete value id %s", val.Id)
}

func main() {
	closeChan := make(chan struct{}, 1)
	mysqlClient := client.GetMysqlCcImpl()
	mysqlClient.Run()
	listener := &TestConfigListener{}
	mysqlClient.RegisterConfig(&TestConfig{}, listener)
	<-closeChan
}
