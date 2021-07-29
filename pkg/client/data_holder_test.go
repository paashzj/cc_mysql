package client

import (
	"fmt"
	"github.com/paashzj/cc_api"
	"log"
	"reflect"
	"testing"
	"unsafe"
)

type DataHolderTestConfig struct {
	Name *string `config:"name"`
	Age  string `config:"age"`
	cc_api.BaseConfig
}

type DataHolderTestConfigListener struct {

}

func (t *DataHolderTestConfigListener) Add(value unsafe.Pointer) {
	val := (*DataHolderTestConfig)(value)
	log.Printf("add value %s %d", val.Id, val.Version)
}

func (t *DataHolderTestConfigListener) Mod(oldValue, value unsafe.Pointer) {
	oldVal := (*DataHolderTestConfig)(oldValue)
	val := (*DataHolderTestConfig)(value)
	log.Printf("old value version %d new value version %d", oldVal.Version, val.Version)
}

func (t *DataHolderTestConfigListener) Del(value unsafe.Pointer) {
	val := (*DataHolderTestConfig)(value)
	log.Printf("delete value id %s", val.Id)
}


func TestDataHolder(t *testing.T) {
	elem := reflect.TypeOf(&DataHolderTestConfig{}).Elem()
	d := newDataHolder(elem, "config_TestConfig", &DataHolderTestConfigListener{})
}
