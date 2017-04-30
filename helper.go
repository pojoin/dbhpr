package dbhpr

import (
	"errors"
	"fmt"
	"time"
)

type Row map[string]interface{}

type Helper interface {
	Insert(sql string, args ...interface{}) (lastInsterId int64, err error)
	Update(sql string, args ...interface{}) (rowsAffected int64, err error)
	Delete(sql string, args ...interface{}) (rowsAffected int64, err error)
	Get(sql string, args ...interface{}) (Row, error)
	Query(sql string, args ...interface{}) ([]Row, error)
}

func (r Row) GetInt64(col string) (int64, error) {
	if cv, ok := r[col]; ok {
		if cvptr, ok := cv.(*interface{}); ok {
			if i, ok := (*cvptr).(int64); ok {
				return i, nil
			} else {
				return 0, errors.New(fmt.Sprintf("column [%s] type is not int64\n", col))
			}
		} else {
			err_msg := fmt.Sprintf("column [%s] type is not *interface{}\n", col)
			return 0, errors.New(err_msg)
		}
	}
	err_msg := fmt.Sprintf("can not found column [%s]\n", col)
	return 0, errors.New(err_msg)
}

func (r Row) GetBytes(col string) ([]byte, error) {
	if cv, ok := r[col]; ok {
		if cvprt, ok := cv.(*interface{}); ok {
			if bs, ok := (*cvprt).([]byte); ok {
				return bs, nil
			} else {
				err_msg := fmt.Sprintf("column [%s] type is not []byte\n", col)
				return nil, errors.New(err_msg)
			}
		} else {
			err_msg := fmt.Sprintf("column [%s] type is not *interface{}\n", col)
			return nil, errors.New(err_msg)
		}
	}
	err_msg := fmt.Sprintf("can not found column [%s]\n", col)
	return nil, errors.New(err_msg)
}

func (r Row) GetString(col string) (string, error) {
	bs, err := r.GetBytes(col)
	if err != nil {
		return "", err
	}
	return string(bs), err
}

func (r Row) GetDate(col string) (time.Time, error) {
	if cv, ok := r[col]; ok {
		if cvptr, ok := cv.(*interface{}); ok {
			if dt, ok := (*cvptr).(time.Time); ok {
				return dt, nil
			} else {
				err_msg := fmt.Sprintf("column [%s] type is not time.Time\n", col)
				return time.Now(), errors.New(err_msg)
			}
		} else {
			err_msg := fmt.Sprintf("column [%s] type is not *interface{}\n", col)
			return time.Now(), errors.New(err_msg)
		}
	}
	err_msg := fmt.Sprintf("can not found column [%s]\n", col)
	return time.Now(), errors.New(err_msg)
}

//func (r Row) MarshalJSON() ([]byte, error) { }
