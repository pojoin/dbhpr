package dbhpr

import (
	"reflect"
	"strings"
)

type DBHelper struct {
	dbname string
}

func (h *DBHelper) Insert(sql string, args ...interface{}) (lastInsterId int64, err error) {
	stmt, err := dbHive[h.dbname].Prepare(sql)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	r, err := stmt.Exec(args...)
	if err != nil {
		return 0, err
	}
	lastInsterId, err = r.LastInsertId()
	return
}

func (h *DBHelper) Update(sql string, args ...interface{}) (rowsAffected int64, err error) {
	stmt, err := dbHive[h.dbname].Prepare(sql)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	r, err := stmt.Exec(args...)
	if err != nil {
		return 0, err
	}
	rowsAffected, err = r.RowsAffected()
	return
}

func (h *DBHelper) Delete(sql string, args ...interface{}) (rowsAffected int64, err error) {
	return h.Update(sql, args)
}

func (h *DBHelper) Get(sql string, args ...interface{}) (Row, error) {
	if !strings.Contains(strings.ToLower(sql), "limit") {
		sql += " limit 1 "
	}

	rows, err := h.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, err
	}
	return rows[0], nil
}

func (h *DBHelper) Query(sql string, args ...interface{}) ([]Row, error) {
	stmt, err := dbHive[h.dbname].Prepare(sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}
	results := make([]Row, 0)
	for rows.Next() {
		row := make(map[string]interface{})
		values := make([]interface{}, 0, len(columnTypes))
		for _, t := range columnTypes {
			//fmt.Println("name=", t.Name(), ",type=", t.ScanType(), ",databaseTypeName=", t.DatabaseTypeName())
			values = append(values, reflect.New(t.ScanType()).Interface())
		}
		err := rows.Scan(values...)
		if err != nil {
			return nil, err
		}
		for i, t := range columnTypes {
			row[t.Name()] = values[i]
		}
		results = append(results, row)
	}
	return results, nil
}
