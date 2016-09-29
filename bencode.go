// Package bencode implements encoding and decoding of Bencode as
// described in https://zh.wikipedia.org/wiki/Bencode

package bencode

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// Encode returns the Becode encoding text of data.
func Encode(data interface{}) (string, error) {
	v := reflect.Indirect(reflect.ValueOf(data))
	t := v.Type()
	k := t.Kind()

	switch {
	case reflect.Invalid < k && k <= reflect.Int64:
		return encodeInt(int(v.Int()))
	case reflect.Uint <= k && k <= reflect.Uint64:
		return encodeInt(int(v.Uint()))
	case k == reflect.String:
		return encodeString(v.String())
	case k == reflect.Slice || k == reflect.Array:
		return encodeSlice(v)
	case k == reflect.Map:
		return encodeMap(v)
	case k == reflect.Struct:
		return encodeStruct(v)
	default:
		return "", errors.New("data type no support")
	}
	return "", errors.New("encode fail")
}

func encodeInt(data int) (string, error) {
	return fmt.Sprintf("i%de", data), nil
}
func encodeString(data string) (string, error) {
	return strings.Join([]string{strconv.Itoa(len(data)), data}, ":"), nil
}
func encodeSlice(v reflect.Value) (string, error) {
	ret := "l"

	for idx := 0; idx != v.Len(); idx++ {
		itemRet, err := Encode(v.Index(idx).Interface())
		if err != nil {
			return itemRet, err
		}
		ret += itemRet
	}
	ret += "e"

	return ret, nil
}

type keySli []reflect.Value

func (sli keySli) Len() int {
	return len(sli)
}
func (sli keySli) Less(i, j int) bool {
	return sli[i].String() < sli[j].String()
}
func (sli keySli) Swap(i, j int) {
	sli[i], sli[j] = sli[j], sli[i]
}

func encodeMap(v reflect.Value) (string, error) {
	ret := "d"
	vkey := v.MapKeys()
	sort.Sort(keySli(vkey))

	for _, val := range vkey {
		itemRet, err := Encode(v.MapIndex(val).Interface())
		if err != nil {
			return itemRet, err
		}
		keyRet, err := encodeString(val.String())
		if err != nil {
			return keyRet, err
		}
		ret += keyRet + itemRet
	}
	return ret + "e", nil
}
func encodeMapIt(data interface{}) (string, error) {
	v := reflect.Indirect(reflect.ValueOf(data))
	return encodeMap(v)
}
func encodeStruct(v reflect.Value) (string, error) {
	m := make(map[string]interface{})
	t := v.Type()

	for idx := 0; idx != t.NumField(); idx++ {
		if v.Field(idx).CanInterface() {
			name := t.Field(idx).Name

			if tagName := t.Field(idx).Tag.Get("json"); tagName != "" {
				name = tagName
			}

			m[name] = v.Field(idx).Interface()
		}
	}
	return encodeMapIt(m)
}
