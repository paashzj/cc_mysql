package client

import (
	"config-client/pkg/mysql"
	"github.com/paashzj/cc_api"
	"log"
	"reflect"
	"sync"
	"time"
)

func GetMysqlCcImpl() *MysqlCcImpl {
	return &MysqlCcImpl{}
}

type MysqlCcImpl struct {
	listenerMap   sync.Map
	dataHolderMap sync.Map
}

func (m *MysqlCcImpl) Run() {
	m.listenNotify()
}

func (m *MysqlCcImpl) RegisterConfig(value interface{}, listener cc_api.Listener) {
	typeOf := reflect.TypeOf(value)
	typeName := typeOf.Elem().Name()
	log.Printf("type name is %s\n", typeName)
	m.listenerMap.Store(typeName, listener)
	dataHolder := newDataHolder(typeOf.Elem(), "config_"+typeName, listener)
	dataHolder.run()
	m.dataHolderMap.Store(typeName, dataHolder)
}

func (m *MysqlCcImpl) listenNotify() {
	ticker := time.NewTicker(1 * time.Second)
	var maxId int64 = -1
	go func() {
		for range ticker.C {
			// 先查询最大id
			if maxId == -1 {
				aux, err := mysql.MaxConfigNotifyId()
				if err != nil {
					log.Printf("select max id error %s \n", err)
				} else {
					maxId = aux
				}
				continue
			}
			notifyList := mysql.NotifyList(maxId)
			maxId += (int64)(len(notifyList))
			for _, notifyConfig := range notifyList {
				value, ok := m.dataHolderMap.Load(notifyConfig.ConfigName)
				if !ok {
					continue
				}
				dataHolder := value.(*DataHolder)
				dataHolder.notify(notifyConfig.ConfigItemId)
			}
		}
	}()
}
