package dbhpr

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"

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

	id, _ := row["id"].(int)
	t.Log("id = ", id)

	state, _ := row["state"].(int)
	t.Log("state = ", state)

	utime, _ := row["utime"].(time.Time)
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
	name, _ := row["name"].(string)
	t.Log("name=", name)

	userid, _ := row["id"].(int64)
	t.Log("userid = ", userid)
	bs, _ := json.Marshal(row)
	t.Log(string(bs))
	email, _ := row["email"].(string)
	t.Log("email = ", email)
}

func Test_Query(t *testing.T) {
	rows, err := Query("select * from user where id>?", 0)
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
