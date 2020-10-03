package util

import (
	"github.com/deckarep/golang-set"
	"reflect"
	"strings"
)

// SetStructFieldByMap 遍历 结构体指针 ptr 中的每个属性,检测在 fields 中 是否存在.存在且类型相同 则 赋值.
// 如果 allowFieldsName 非 nil , 则属性名在 allowFieldsName 列表中 时才会进行赋值。
func SetStructFieldByMap(ptr interface{}, fields map[string]interface{}, allowFieldsName []string) {
	v := reflect.ValueOf(ptr).Elem() // the struct variable

	var set mapset.Set
	if allowFieldsName != nil {
		var fieldsName []interface{}
		for _, v := range allowFieldsName {
			fieldsName = append(fieldsName, v)
		}
		set = mapset.NewSetFromSlice(fieldsName)
	}
	for i := 0; i < v.NumField(); i++ {
		fieldInfo := v.Type().Field(i) // fieldInfo reflect.StructField
		tag := fieldInfo.Tag           // tag reflect.StructTag
		name := tag.Get("json")

		if name == "" {
			name = strings.ToLower(fieldInfo.Name) // 没有 json tag 时 使用 原始的 值
		}

		name = strings.Split(name, ",")[0] //去掉逗号后面内容 如 `json:"voucher_usage,omitempty"
		if set != nil {
			if !set.Contains(name) {
				continue
			}
		}
		if value, ok := fields[name]; ok {
			if reflect.ValueOf(value).Type() == fieldInfo.Type { // 类型 相同
				v.FieldByName(fieldInfo.Name).Set(reflect.ValueOf(value))
			}
		}
	}
}
