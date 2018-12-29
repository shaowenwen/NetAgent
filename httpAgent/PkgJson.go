package httpAgent

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Smap []*SortMapNode

type SortMapNode struct {
	Key string
	Val interface{}
}

func (c *Smap) Put(key string, val interface{}) {
	index, _, ok := c.get(key)
	if ok {
		(*c)[index].Val = val
	} else {
		node := &SortMapNode{Key: key, Val: val}
		*c = append(*c, node)
	}
}

func (c *Smap) Get(key string) (interface{}, bool) {
	_, val, ok := c.get(key)
	return val, ok
}

func (c *Smap) get(key string) (int, interface{}, bool) {
	for index, node := range *c {
		if node.Key == key {
			return index, node.Val, true
		}
	}
	return -1, nil, false
}

func PkgKafkaJsonData(content interface{},msgtype string) string {
	var arr []interface{}
	arr = append(arr, content)
	smap := &Smap{}
	smap.Put("timestamp", time.Now().Format("2006-01-02 03:04:05"))
	smap.Put("frim", "rm")
	smap.Put("prov", "LiaoNingSheng")
	smap.Put("city", "ShenYangShi")
	smap.Put("area_code", "210100")
	smap.Put("car_type", "cz")
	smap.Put("command", msgtype)
	smap.Put("content", arr)

	s := ToSortedMapJson(smap)

	return s
}


func ToSortedMapJson(smap *Smap) string {
	s := "{"
	for _, node := range *smap {
		v := node.Val
		isSamp := false
		str := ""
		switch v.(type) {
		case *Smap:
			isSamp = true
			str = ToSortedMapJson(v.(*Smap))
		}

		if !isSamp {
			b, _ := json.Marshal(node.Val)
			str = string(b)
		}
		s = fmt.Sprintf("%s\"%s\":%s,", s, node.Key, str)
	}
	s = strings.TrimRight(s, ",")

	s = fmt.Sprintf("%s}", s)
	return s
}