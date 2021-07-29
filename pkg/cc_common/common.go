package cc_common

import "reflect"

type IdVersion struct {
	Id      string
	Version int
}

type NotifyHelper struct {
	ConfigName   string
	ConfigItemId string
}

func GetVersion(pointer *reflect.Value) int {
	value := (*pointer).Elem().FieldByName("BaseConfig").FieldByName("Version")
	return *(value.Addr().Interface().(*int))
}

func ToInterfaceSlice(stringSlice []string) []interface{} {
	destSlice := make([]interface{}, len(stringSlice))
	for i, s := range stringSlice {
		destSlice[i] = s
	}
	return destSlice
}