package mysql

import (
	"config-client/pkg/cc_common"
	"errors"
	"log"
	"reflect"
	"strings"
)

func FindData(tableName string, structFields, queryFields, itemIds []string, structType reflect.Type) (map[string]*reflect.Value, error) {
	inSql := "SELECT id, version," + strings.Join(queryFields, ",") + " FROM " + tableName + " WHERE id in (" + getPlaceHolders(len(itemIds)) + ")"
	stmt, err := db.Prepare(inSql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(cc_common.ToInterfaceSlice(itemIds)...)
	if err != nil {
		return nil, err
	}
	queryMap := make(map[string]*reflect.Value)
	for rows.Next() {
		newStruct := reflect.New(structType)
		// 和mysql查询内容对接
		baseConfigField := newStruct.Elem().FieldByName("BaseConfig")
		sqlDestSlice := make([]interface{}, 2)
		sqlDestSlice[0] = baseConfigField.FieldByName("Id").Addr().Interface()
		sqlDestSlice[1] = baseConfigField.FieldByName("Version").Addr().Interface()
		for i := 0; i < len(structFields); i++ {
			sqlDestSlice = append(sqlDestSlice, newStruct.Elem().FieldByName(structFields[i]).Addr().Interface())
		}
		err := rows.Scan(sqlDestSlice...)
		if err != nil {
			return nil, err
		}
		aux := baseConfigField.FieldByName("Id").Addr().Interface().(*string)
		// 存map
		queryMap[*aux] = &newStruct
	}
	return queryMap, nil
}

func GetIdVersionMap(tableName string) (map[string]cc_common.IdVersion, error) {
	querySql := "SELECT id, version FROM " + tableName
	stmt, err := db.Prepare(querySql)
	if err != nil {
		log.Printf("select version error don't sync %s", err)
		return nil, errors.New("query data errored")
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		log.Printf("error is %s", err)
		return nil, errors.New("query data errored")
	}
	var idVersion cc_common.IdVersion
	idVersionMap := make(map[string]cc_common.IdVersion)
	for rows.Next() {
		err := rows.Scan(&idVersion.Id, &idVersion.Version)
		if err != nil {
			log.Printf("error is %s end sync", err)
			return nil, errors.New("load data errored")
		}
		idVersionMap[idVersion.Id] = idVersion
	}
	return idVersionMap, nil
}
