package client

import (
	"config-client/pkg/cc_common"
	"config-client/pkg/mysql"
	"github.com/paashzj/cc_api"
	"log"
	"reflect"
	"sync"
	"time"
	"unsafe"
)

type DataHolder struct {
	mux          sync.Mutex
	dataMap      map[string]*reflect.Value
	structType   reflect.Type
	tableName    string
	listener     cc_api.Listener
	structFields []string
	queryFields  []string
}

func newDataHolder(structType reflect.Type, tableName string, listener cc_api.Listener) *DataHolder {
	d := &DataHolder{
		dataMap:    make(map[string]*reflect.Value),
		structType: structType,
		tableName:  tableName,
		listener:   listener,
	}
	d.run()
	return d
}

// run 定时扫描同步
func (d *DataHolder) run() {
	for i := 0; i < d.structType.NumField(); i++ {
		structField := d.structType.Field(i)
		if structField.Name == "BaseConfig" {
			continue
		}
		configTag := structField.Tag.Get("config")
		if configTag == "" {
			continue
		}
		d.structFields = append(d.structFields, structField.Name)
		d.queryFields = append(d.queryFields, configTag)
	}
	ticker := time.NewTicker(1 * time.Second)
	for range ticker.C {
		d.timerSync()
	}
}

func (d *DataHolder) timerSync() {
	idVersionMap, err := mysql.GetIdVersionMap(d.tableName)
	if err != nil {
		log.Println("query id version failed ", err)
		return
	}
	var needSyncSlice []string
	// idVersionMap 和 dataMap对比
	// 此时只是对比出需要和远端同步的数据，并不需要严格同步，也可以做同步处理
	d.mux.Lock()
	for index, idVersionElem := range idVersionMap {
		value, ok := d.dataMap[index]
		if !ok {
			needSyncSlice = append(needSyncSlice, index)
		} else {
			if cc_common.GetVersion(value) != idVersionElem.Version {
				needSyncSlice = append(needSyncSlice, index)
			}
		}
	}
	d.mux.Unlock()
	d.syncList(needSyncSlice)
}

// 接收notify通知
func (d *DataHolder) notify(configItemId string) {
	log.Printf("notify %s\n", configItemId)
	d.sync(configItemId)
}

func (d *DataHolder) syncList(itemIds []string) {
	if len(itemIds) == 0 {
		log.Println("sync list is empty")
		return
	}
	log.Printf("sync list is %s\n", itemIds)
	queryMap, err := mysql.FindData(d.tableName, d.structFields, d.queryFields, itemIds, d.structType)
	if err != nil {
		log.Printf("find data error %s\n", err)
		return
	}
	defer d.mux.Unlock()
	d.mux.Lock()
	for index, value := range queryMap {
		oldVal := d.dataMap[index]
		if oldVal == nil {
			d.notifyAdd(value)
		} else {
			d.notifyMod(oldVal, value)
		}
		d.dataMap[index] = value
	}
	for _, itemId := range itemIds {
		oldVal := d.dataMap[itemId]
		queryVal := queryMap[itemId]
		if queryVal == nil {
			delete(d.dataMap, itemId)
			d.notifyDel(oldVal)
		}
	}
}

func pointerConvert(value *reflect.Value) unsafe.Pointer {
	return unsafe.Pointer(value.Pointer())
}

func (d *DataHolder) notifyAdd(value *reflect.Value) {
	d.listener.Add(pointerConvert(value))
}

func (d *DataHolder) notifyMod(oldValue *reflect.Value, value *reflect.Value) {
	d.listener.Mod(pointerConvert(oldValue), pointerConvert(value))
}

func (d *DataHolder) notifyDel(value *reflect.Value) {
	d.listener.Del(pointerConvert(value))
}

func (d *DataHolder) sync(itemId string) {
	itemIds := make([]string, 1)
	itemIds[0] = itemId
	d.syncList(itemIds)
}
