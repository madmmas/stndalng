package utils

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

// Function GetJSON returns query result in json object
// based on params: (sqlString)
func GetJSON(rows *sql.Rows) (string, error) {
	tableData, err := GetMapData(rows)
	if err != nil {
		return "", err
	}
	jsonData, err := json.Marshal(tableData)
	if err != nil {
		return "", err
	}
	fmt.Println(string(jsonData))
	return string(jsonData), nil
}

// Function GetJSONTable returns query result in rows and total formatted json object
// based on params: (sqlString)
func GetJSONTable(rows *sql.Rows) ([]byte, error) {

	tableData, err := GetMapData(rows)
	if err != nil {
		return nil, err
	}

	return json.Marshal(map[string]interface{}{
		"rows":  tableData,
		"total": len(tableData),
	})
}

// Function GetMapData returns query result in map formatted data
// based on params: (sqlString)
func GetMapData(rows *sql.Rows) ([]map[string]interface{}, error) {

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	count := len(columns)
	tableData := make([]map[string]interface{}, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				t := string(b)
				//
				if isJSONString(t) {
					fmt.Println("JSON String::")
					var js interface{}
					json.Unmarshal(b, &js)
					v = js
				} else {
					v = string(b)
				}
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}
	return tableData, nil
}

func isJSONString(s string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(s), &js) == nil
}

func isJSON(s string) bool {
	var js interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}
