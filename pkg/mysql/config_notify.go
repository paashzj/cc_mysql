package mysql

import (
	"config-client/pkg/cc_common"
	"log"
)

func MaxConfigNotifyId() (int64, error) {
	sql := "SELECT id from config_notify ORDER BY id DESC LIMIT 0,1"
	stmt, err := db.Prepare(sql)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	query, err := stmt.Query()
	if err != nil {
		log.Printf("error is %s", err)
	}
	if !query.Next() {
		return 0, nil
	}
	var result int64
	err = query.Scan(result)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func NotifyList(maxId int64) []cc_common.NotifyHelper {
	notifyHelpers := make([]cc_common.NotifyHelper, 0)
	// 扫描notify数据表
	querySql := "SELECT config_name, config_item_id FROM config_notify WHERE id > ? ORDER BY id ASC LIMIT 0,500"
	stmt, err := db.Prepare(querySql)
	if err != nil {
		log.Println("prepare failed ", err)
		return notifyHelpers
	}
	defer stmt.Close()
	query, err := stmt.Query(maxId)
	if err != nil {
		log.Println("prepare failed ", err)
		return notifyHelpers
	}
	for query.Next() {
		maxId++
		var configName string
		var configItemId string
		err := query.Scan(&configName, &configItemId)
		if err != nil {
			log.Println("scan failed", err)
			continue
		}
		notifyHelpers = append(notifyHelpers, cc_common.NotifyHelper{ConfigName: configName, ConfigItemId: configItemId})
	}
	return notifyHelpers
}
