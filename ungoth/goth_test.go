package ungoth

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"testing"
)

type UserRecord struct {
	Username string `html:"td"`
	Count    int    `html:"td"`
}

func TestMarshal(t *testing.T) {
	x := UserRecord{Username: "admin", Count: 20}
	tag := MarshalHTML(x)
	log.Println(tag)
}

func MarshalHTML(t interface{}) string {
	list := []Tag{}
	rt := reflect.TypeOf(t)
	pt := reflect.ValueOf(t)
	for i := 0; i < rt.NumField(); i++ {

		tag := rt.Field(i).Tag.Get("html")
		name := rt.Field(i).Name
		v := Tag{
			Name:  tag,
			Props: make(map[string]interface{}),
		}
		switch pt.FieldByName(name).Kind() {
		case reflect.String:
			v.Props[name] = pt.FieldByName(name).String()
		case reflect.Int:
			v.Props[name] = pt.FieldByName(name).Int()
		}
		list = append(list, v)
	}
	slist := []string{}
	for _, v := range list {
		slist = append(slist, v.MarshalHTML())
	}
	return strings.Join(slist, "")
}

type Tag struct {
	Name     string
	Props    Properties
	Children []Tag
}

func (t *Tag) MarshalHTML() string {
	props := t.Props.MarshalHTML()
	return fmt.Sprintf("<%s %s></%s>", t.Name, props, t.Name)
}

type Properties map[string]interface{}

func (p Properties) MarshalHTML() string {
	l := []string{}
	for k, v := range p {
		key := formatHTMLProperty(k)
		switch x := v.(type) {
		case string:
			l = append(l, fmt.Sprintf("%s=%s", key, x))
		case int:
			l = append(l, fmt.Sprintf("%s=%d", key, x))
		case int64:
			l = append(l, fmt.Sprintf("%s=%d", key, x))
		default:
			log.Println(x, reflect.TypeOf(x).Kind())
			l = append(l, fmt.Sprintf("%s=undefined", key))
		}
	}
	return strings.Join(l, " ")
}

var cUpperCase = "ABCDRFGHIJKLMNOPQRSTUVWXYZ"

func formatHTMLProperty(s string) string {
	var i int
	for {
		i = strings.IndexAny(s, cUpperCase)
		if i < 0 {
			return s
		}
		if i == 0 {
			s = strings.ToLower(string(s[i])) + s[i+1:]
		} else {
			s = s[0:i] + "-" + strings.ToLower(string(s[i])) + s[i+1:]
		}
	}
}
