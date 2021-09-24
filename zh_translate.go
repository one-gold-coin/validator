package validator

import (
	"errors"
	"reflect"
	"strings"
)

var (
	// {0} == field.AliasName
	// {1} == field.Tags.param
	defaultTranslateMap = map[string]string{
		"required": "{0}为必填字段",
		"eq":       "{0}不等于{1}",
		"ne":       "{0}不能等于{1}",
		"email":    "{0}必须是一个有效的邮箱",
		"oneof":    "{0}必须是[{1}]中的一个",
		//len相关
		"len-string": "{0}长度必须是{1}个字符",
		"len-int":    "{0}必须等于{1}",
		"len-int64":  "{0}必须等于{1}",
		"len-slice":  "{0}必须包含{1}项",
		//min相关
		"min-string": "{0}长度必须至少为{1}个字符",
		"min-int":    "{0}最小只能为{1}",
		"min-int64":  "{0}最小只能为{1}",
		"min-slice":  "{0}至少包含{1}项",
		//max相关
		"max-string": "{0}长度不超过{1}个字符",
		"max-int":    "{0}必须小于或等于{1}",
		"max-int64":  "{0}必须小于或等于{1}",
		"max-slice":  "{0}最多包含{1}项",
		//lt相关
		"lt-string": "{0}长度必须小于{1}个字符",
		"lt-int":    "{0}必须小于{1}",
		"lt-int64":  "{0}必须小于{1}",
		"lt-slice":  "{0}必须少于{1}项",
		//lte相关
		"lte-string": "{0}长度不能超过{1}个字符",
		"lte-int":    "{0}必须小于或等于{1}",
		"lte-int64":  "{0}必须小于或等于{1}",
		"lte-slice":  "{0}只能包含{1}项",
		//gt相关
		"gt-string": "{0}长度必须大于{1}个字符",
		"gt-int":    "{0}必须大于{1}",
		"gt-int64":  "{0}必须大于{1}",
		"gt-slice":  "{0}必须大于{1}项",
		//gte相关
		"gte-string": "{0}长度必须至少为{1}个字符",
		"gte-int":    "{0}必须大于或等于{1}",
		"gte-int64":  "{0}必须大于或等于{1}",
		"gte-slice":  "{0}必须至少包含{1}项",
	}
)

func NewZhTranslate() *ZhTranslate {
	t := &ZhTranslate{
		translateMap: defaultTranslateMap,
	}
	return t
}

type ZhTranslate struct {
	translateMap map[string]string
	err          error
}

func (m *ZhTranslate) GetErr() error {
	return m.err
}

func (m *ZhTranslate) SetErr(err error) *ZhTranslate {
	m.err = err
	return m
}

func (m *ZhTranslate) GetTranslateMap() map[string]string {
	return m.translateMap
}

func (m *ZhTranslate) SetTranslateMap(translateMap map[string]string) *ZhTranslate {
	m.translateMap = translateMap
	return m
}

func (m *ZhTranslate) Translate(v *Validator) *ZhTranslate {
	field := v.GetField()
	tMap := m.GetTranslateMap()
	//如果是指针类型则取指针对应真实类型
	tKind := field.Sf.Type.Kind()
	if tKind == reflect.Ptr {
		tKind = field.Sf.Type.Elem().Kind()
	}
	// 判断Tags.tag是否有定义
	// 再判断Tags.tag + tKind 是否有定义
	if val, isOk := tMap[field.Tags.tag]; isOk {
		m.GetStr(val, field.AliasName, field.Tags.param)
		return m
	}
	groupKey := field.Tags.tag + "-" + tKind.String()
	if val, isOk := tMap[groupKey]; isOk {
		m.GetStr(val, field.AliasName, field.Tags.param)
		return m
	}
	m.SetErr(errors.New("参数异常"))
	return m
}

func (m *ZhTranslate) GetStr(translate, altName, tagParam string) {
	if strings.ContainsAny(translate, "{0}&{1}") {
		translate = strings.Replace(translate, "{0}", altName, 1)
		translate = strings.Replace(translate, "{1}", tagParam, 1)
		m.SetErr(errors.New(translate))
		return
	}
	if strings.Contains(translate, "{0}") {
		translate := strings.Replace(translate, "{0}", altName, 1)
		m.SetErr(errors.New(translate))
		return
	}
	m.SetErr(errors.New("解析异常"))
	return
}
