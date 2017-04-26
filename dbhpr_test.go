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

	id := row.GetInt64("id")
	t.Log("id = ", id)

	state := row.GetInt64("state")
	t.Log("state = ", state)

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
	name := row.GetString("name")
	t.Log("name=", name)
}
