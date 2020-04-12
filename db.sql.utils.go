package m1

import (
	"encoding/json"
	"errors"
	"fmt"
	str "strings"

	log "github.com/sirupsen/logrus"
)

// ExecuteQuery, execute stored procedure with a list of data, not have a single row
func ExecuteQuery(database *DataBase, spName string, args ...interface{}) ([]interface{}, error) {

	const errorMessage string = "Error Exec Sp: %s, line code: %s, error= %s"

	// Execute Query, gives rows and maybe an error
	rows, err := database.DB.QueryContext(*database.Context, spName, args...)

	// Check execute query doesn't have an error
	if err != nil {
		log.Error(fmt.Sprintf(errorMessage, spName, "QueryContext, 21", err.Error()))
		return nil, err
	}

	// Gets columns
	cols, err := rows.Columns()
	if err != nil {
		log.Error(fmt.Sprintf(errorMessage, spName, "Columns, 28", err.Error()))
		return nil, err
	}

	// Describe the list for returns
	var list []interface{}

	// Read to end from the rows
	for rows.Next() {
		data := make(map[string]interface{})
		columns := make([]interface{}, len(cols))
		columnsPointer := make([]interface{}, len(cols))

		for i, _ := range columns {
			columnsPointer[i] = &columns[i]
		}

		// Read current row
		rows.Scan(columnsPointer...)

		for i, colName := range cols {

			// Convert JSON string to model data
			if str.HasPrefix(colName, "json_") {
				// Create new field, eq: json_field => field and convert JSON
				str, ok := columns[i].(string)
				if !ok {
					errorText := fmt.Sprintf("%s is not a string value", colName)
					log.Error(fmt.Sprintf(errorMessage, spName, "columns[i].(string), 59", errorText))
					return nil, errors.New(errorText)
				}

				var jsonFiled interface{}
				if err := json.Unmarshal([]byte(str), &jsonFiled); err != nil {
					log.Error(fmt.Sprintf(errorMessage, spName, "json.Unmarshal, 63 ", err.Error()))
					return nil, err
				}
				data[colName[len("json_"):len(colName)]] = jsonFiled
			} else {
				// Normal field
				data[colName] = columns[i]
			}
		}

		// Append current row to the list
		list = append(list, data)
	}

	defer rows.Close()

	return list, nil

}
