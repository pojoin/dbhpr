package dbhpr

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func init() {
	err := NewDB("default", "mysql", "shit:shit@tcp(hbhaize.oicp.net:3306)/shit?charset=utf8&parseTime=true&loc=Local")
	if err != nil {
		fmt.Println(err)
	}
}

func Test_Get(t *testing.T) {
	row, err := Get("select id,userid,state,utime,devId from userstate where id=3")
	t.Log(err)
	t.Log(row)

	date := row["utime"]
	t.Log(reflect.TypeOf(date))
	if i, ok := date.(*interface{}); ok {
		t.Log(reflect.TypeOf(*i))
		t.Log(*i)
	}

	id, err := row.GetInt64("id")
	if err != nil {
		t.Error(err)
	}
	t.Log("id = ", id)

	state, err := row.GetInt64("state")
	if err != nil {
		t.Error(err)
	}
	t.Log("state = ", state)

	utime, err := row.GetDate("utime")
	if err != nil {
		t.Error(err)
	}
	t.Log("utime = ", utime)

	bs, err := json.Marshal(&row)
	if err != nil {
		t.Log(err)
		return
	}

	t.Log(string(bs))
}

func Test_GetUser(t *testing.T) {
	row, err := Get("select * from user")
	if err != nil {
		t.Fatal(err)
		return
	}
	if n, ok := row["name"].(*interface{}); ok {
		t.Log(reflect.TypeOf(*n))
		if names, ok := (*n).([]byte); ok {
			t.Log(string(names))
		}
	}
	name, err := row.GetString("name")
	if err != nil {
		t.Error(err)
	}
	t.Log("name=", name)

	userid, err := row.GetInt64("id")
	if err != nil {
		t.Error(err)
	}
	t.Log("userid = ", userid)
	bs, _ := json.Marshal(row)
	t.Log(string(bs))
	email, _ := row.GetString("email")
	t.Log("email = ", email)
}

func Test_Query(t *testing.T) {
	rows, err := Query("select * from user where id=?", 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(rows)
	bs, err := json.Marshal(rows)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(bs))
}
