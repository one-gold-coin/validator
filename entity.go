package validator

import (
	"reflect"
)

type Config struct {
	FieldDescribeTag string
	ValidationTag    string
	OmitemptyTag     string
}

// Field 字段信息
type Field struct {
	Idx       int                  //字段下标
	AliasName string               //字段别名
	Sf        *reflect.StructField //字段类型
	Tags      *Tag                 //字段tag信息
}

//Tag 解析信息
type Tag struct {
	tag       string         //tag 名称
	param     string         //验证tag标签值 eg: max=100 ; param=100
	isHaveErr bool           //是否有验证错误
	rv        *reflect.Value //验证struct对应的字段信息
}
