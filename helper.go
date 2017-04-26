package dbhpr

import "log"

type Row map[string]interface{}

type Helper interface {
	Insert(sql string, args ...interface{}) (lastInsterId int64, err error)
	Update(sql string, args ...interface{}) (rowsAffected int64, err error)
	Delete(sql string, args ...interface{}) (rowsAffected int64, err error)
	Get(sql string, args ...interface{}) (Row, error)
	Query(sql string, args ...interface{}) ([]Row, error)
}

func (r Row) GetInt64(col string) int64 {
	if cv, ok := r[col]; ok {
		if cvptr, ok := cv.(*interface{}); ok {
			if i, ok := (*cvptr).(int64); ok {
				return i
			} else {
				log.Printf("column [%s] type is not int64\n", col)
				return 0
			}
		} else {
			log.Printf("column [%s] type is not *interface{}\n", col)
			return 0
		}
	} else {
		log.Printf("can not found column [%s]\n", col)
		return 0
	}
	return 0
}

func (r Row) GetBytes(col string) []byte {
	if cv, ok := r[col]; ok {
		if cvprt, ok := cv.(*interface{}); ok {
			if bs, ok := (*cvprt).([]byte); ok {
				return bs
			} else {
				log.Printf("column [%s] type is not []byte\n", col)
				return nil
			}
		} else {
			log.Printf("column [%s] type is not *interface{}\n", col)
			return nil
		}
	}
	log.Printf("can not found column [%s]\n", col)
	return nil
}

func (r Row) GetString(col string) string {
	if r := r.GetBytes(col); r != nil {
		return string(r)
	}
	return ""
}
