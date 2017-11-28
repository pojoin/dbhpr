package dbhpr

import (
	Sql "database/sql"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"
)

type DBHelper struct {
	dbname string
}

func (h *DBHelper) Exec(sql string, args ...interface{}) (rowsAffected int64, err error) {
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
	return h.Update(sql, args...)
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
		return nil, errors.New("not found row")
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
			// fmt.Println("name=", t.Name(), ",type=", t.ScanType(), ",databaseTypeName=", t.DatabaseTypeName(), t.ScanType())
			switch t.DatabaseTypeName() {
			case "CHAR", "VARCHAR", "TINYTEXT", "TEXT", "MEDIUMTEXT", "LONGTEXT":
				values = append(values, new(Sql.NullString))
			case "INT", "TINYINT", "SMALLINT", "MEDIUMINT", "BIGINT":
				values = append(values, new(Sql.NullInt64))
			case "FLOAT", "DOUBLE", "DECIMAL":
				values = append(values, new(Sql.NullFloat64))
			case "DATETIME", "TIMESTAMP":
				values = append(values, new(time.Time))
			default:
				values = append(values, reflect.New(t.ScanType()).Interface())
			}
		}
		err := rows.Scan(values...)
		if err != nil {
			return nil, err
		}
		for i, t := range columnTypes {
			ptr := values[i]
			log.Println(t.Name(), "==", reflect.TypeOf(ptr))
			switch v := (ptr).(type) {
			case *Sql.NullString:
				row[t.Name()] = v.String
			case *Sql.NullInt64:
				row[t.Name()] = v.Int64
			case *Sql.NullFloat64:
				row[t.Name()] = v.Float64
			case *int64:
				row[t.Name()] = v
			case Sql.RawBytes:
				row[t.Name()] = string(v)
			case []byte:
				row[t.Name()] = string(v)
			case time.Time:
				row[t.Name()] = Time(v)
			case *float32:
				row[t.Name()] = float32(*v)
			case *float64:
				row[t.Name()] = float64(*v)
			case nil:
				row[t.Name()] = ""
			default:
				//fmt.Println("数据库类型非 数字，字符串，时间,使用请自行转换")
				row[t.Name()] = values[i]
			}
			// if ptr, ok := values[i].(*interface{}); ok {
			// 	log.Println(t.Name(), "==", reflect.TypeOf(*ptr))

			// 	switch v := (*ptr).(type) {
			// 	case Sql.NullString:
			// 		row[t.Name()] = v.String
			// 	case Sql.NullInt64:
			// 		row[t.Name()] = v.Int64
			// 	case Sql.NullFloat64:
			// 		row[t.Name()] = v.Float64
			// 	case int64:
			// 		row[t.Name()] = v
			// 	case Sql.RawBytes:
			// 		row[t.Name()] = string(v)
			// 	case []byte:
			// 		row[t.Name()] = string(v)
			// 	case time.Time:
			// 		row[t.Name()] = Time(v)
			// 	case float32:
			// 		row[t.Name()] = float32(v)
			// 	case float64:
			// 		row[t.Name()] = float64(v)
			// 	case nil:
			// 		row[t.Name()] = ""
			// 	default:
			// 		//fmt.Println("数据库类型非 数字，字符串，时间,使用请自行转换")
			// 		row[t.Name()] = values[i]
			// 	}
			// } else {
			// 	row[t.Name()] = values[i]
			// }
		}
		results = append(results, row)
	}
	log.Println(results)
	return results, nil
}

func (h *DBHelper) IsExists(sql string, args ...interface{}) (ok bool, err error) {
	c, err := h.Count(sql, args...)
	if err != nil {
		return false, err
	}
	if c > 0 {
		return true, err
	}
	return false, err
}

func (h *DBHelper) Count(sql string, args ...interface{}) (c int64, err error) {
	// if tmpsql := strings.ToUpper(sql); !strings.Contains(tmpsql, "COUNT(") {
	// 	if fromIndex := strings.Index(tmpsql, "FROM"); fromIndex > 0 {
	// 		sql = fmt.Sprintf("select count(*) %s", []byte(sql)[fromIndex:])
	// 	}
	// }
	tmpsql := strings.ToUpper(sql)
	if fromIndex := strings.Index(tmpsql, "FROM"); fromIndex > 0 {
		sql = fmt.Sprintf("select count(*) %s", []byte(sql)[fromIndex:])
	}
	r := dbHive[h.dbname].QueryRow(sql, args...)
	err = r.Scan(&c)
	return c, err
}

func (h *DBHelper) QueryPage(page *Page, sql string, args ...interface{}) error {
	//get count
	count, err := h.Count(sql, args...)
	if err != nil {
		return err
	}
	page.Count = count
	if count == 0 {
		page.List = make([]Row, 0)
		return nil
	}

	//sql limit
	// if !strings.Contains(strings.ToUpper(sql), "LIMIT") {
	// 	sql = fmt.Sprintf("%s limit %d,%d", sql, page.StartRow(), page.PageSize)
	// } else {
	// 	return errors.New("QueryPage [" + sql + "] contains limit")
	// }
	sql = fmt.Sprintf("%s limit %d,%d", sql, page.StartRow(), page.PageSize)

	//stmt
	stmt, err := dbHive[h.dbname].Prepare(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	//query rows
	rows, err := stmt.Query(args...)
	if err != nil {
		return err
	}
	defer rows.Close()
	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return err
	}
	results := make([]Row, 0, page.Count)
	for rows.Next() {
		row := make(map[string]interface{})
		values := make([]interface{}, 0, len(columnTypes))
		for _, t := range columnTypes {
			//fmt.Println("name=", t.Name(), ",type=", t.ScanType(), ",databaseTypeName=", t.DatabaseTypeName())
			values = append(values, reflect.New(t.ScanType()).Interface())
		}
		err := rows.Scan(values...)
		if err != nil {
			return err
		}
		for i, t := range columnTypes {
			if ptr, ok := values[i].(*interface{}); ok {
				switch v := (*ptr).(type) {
				case int64:
					row[t.Name()] = v
				case []byte:
					row[t.Name()] = string(v)
				case time.Time:
					row[t.Name()] = Time(v)
				case float32:
					row[t.Name()] = float32(v)
				case float64:
					row[t.Name()] = float64(v)
				case nil:
					row[t.Name()] = ""
				default:
					//fmt.Println("数据库类型非 数字，字符串，时间,使用请自行转换")
					row[t.Name()] = values[i]
				}
			} else {
				row[t.Name()] = values[i]
			}
		}
		results = append(results, row)
	}
	page.List = results
	return nil
}
