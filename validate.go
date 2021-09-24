package validator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func New() *Validator {
	v := &Validator{
		config: &Config{
			FieldDescribeTag: defaultFieldDescribeTag,
			ValidationTag:    defaultValidationTag,
			OmitemptyTag:     defaultOmitemptyTag,
		},
		translate: NewZhTranslate(),
	}
	return v
}

type Validator struct {
	config    *Config
	field     *Field
	translate *ZhTranslate
	err       error
}

func (v *Validator) GetField() *Field {
	return v.field
}

func (v *Validator) SetConfig(conf *Config) *Validator {
	v.config = conf
	return v
}

func (v *Validator) GetConfig() *Config {
	return v.config
}

func (v *Validator) Error() error {
	return v.err
}

func (v *Validator) SetError(err error) *Validator {
	v.err = err
	return v
}

/**
验证完数据后，返回接收对象
1、给reqValidate赋值
2、解析reqValidate每一个字段信息
*/
func (v *Validator) Binding(obj interface{}) *Validator {
	value := reflect.ValueOf(obj)
	//确保 obj 是struct
	if value.Kind() == reflect.Ptr && !value.IsNil() {
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct && value.Kind() != reflect.Interface {
		v.SetError(errors.New(mustStruct))
		return v
	}
	// 遍历 Struct 字段结构 & 校验数据
	isValidationFuncErr := v.extractStruct(value, false)
	// 解析参数校验错误信息
	if isValidationFuncErr {
		v.SetError(v.translate.Translate(v).GetErr())
	}
	return v
}

// 提取 Struct 字段信息
func (v *Validator) extractStruct(current reflect.Value, isValidationFuncErr bool) bool {
	if isValidationFuncErr {
		return isValidationFuncErr
	}
	//获取字段数量
	if current.Kind() == reflect.Ptr || current.Kind() == reflect.Interface {
		current = current.Elem()
	}
	numFields := current.NumField()
	typ := current.Type()

	var fld reflect.StructField
	var validateTag string
	var customName string

	for i := 0; i < numFields; i++ {

		fld = typ.Field(i)

		//是否空字段"-", struct{-}
		if !fld.Anonymous && fld.PkgPath != blank {
			continue
		}
		//获取验证标签
		validateTag = fld.Tag.Get(v.GetConfig().ValidationTag)
		//验证标签是否忽略 或者为空
		if validateTag == skipValidationTag {
			continue
		}
		//验证数据
		if len(validateTag) > 0 {
			tags := v.parseFieldTags(current.Field(i), validateTag, fld.Name)
			if tags != nil && tags.isHaveErr == true {
				//获取字段名
				customName = fld.Name
				//如果设置字段别名
				descTag := fld.Tag.Get(v.GetConfig().FieldDescribeTag)
				if descTag != blank {
					customName = descTag
				}
				v.field = &Field{Idx: i, AliasName: customName, Sf: &fld, Tags: tags}
				return true
			}
		}
		//在此完善其他结构逻辑
		switch current.Field(i).Type().Kind() {
		case reflect.Ptr:
			if current.Field(i).Type().Elem().Kind() == reflect.Struct {
				if v.extractStruct(current.Field(i).Elem(), false) {
					return true
				}
			}
		case reflect.Struct, reflect.Map:
			if v.extractStruct(current.Field(i), false) {
				return true
			}
		case reflect.Slice, reflect.Array:
			for j := 0; j < current.Field(i).Len(); j++ {
				if v.extractStruct(current.Field(i).Index(j), false) {
					return true
				}
			}
		}
	}
	return false
}

//验证数据
func (v *Validator) parseFieldTags(current reflect.Value, tagStr string, fieldName string) *Tag {
	var t string
	var kind reflect.Kind
	//获取真实数据类型
	current, kind = v.extractTypeInternal(current)
	// 获取验证Tag列表
	tags := strings.Split(tagStr, tagSeparator)
	for i := 0; i < len(tags); i++ {
		t = tags[i]
		if t == v.GetConfig().OmitemptyTag {
			switch kind {
			case reflect.Slice, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Chan, reflect.Func:
				if !current.IsNil() {
					continue
				}
			default:
				if current.IsValid() {
					continue
				}
			}
			return nil
		}
		tag := &Tag{}
		tag.rv = &current
		orVials := strings.Split(t, orSeparator)
		for j := 0; j < len(orVials); j++ {
			//获取验证值
			vals := strings.SplitN(orVials[j], tagKeySeparator, 2)
			tag.tag = vals[0]
			if len(tag.tag) == 0 {
				v.SetError(errors.New(strings.TrimSpace(fmt.Sprintf(invalidValidation, fieldName))))
				return nil
			}
			if len(vals) > 1 {
				tag.param = strings.Replace(strings.Replace(vals[1], utf8HexComma, ",", -1), utf8Pipe, "|", -1)
			}
			// 验证
			if validationFunc, ok := validationFuncS[tag.tag]; !ok {
				v.SetError(errors.New(strings.TrimSpace(fmt.Sprintf(undefinedValidation, fieldName))))
				return nil
			} else {
				validationFuncResult := validationFunc(tag)
				//验证
				if !validationFuncResult {
					tag.isHaveErr = true
					return tag
				}
			}
		}
	}
	return nil
}

//获取真实数据类型
func (v *Validator) extractTypeInternal(current reflect.Value) (reflect.Value, reflect.Kind) {

BEGIN:
	switch current.Kind() {
	case reflect.Ptr:

		if current.IsNil() {
			return current, reflect.Ptr
		}

		current = current.Elem()
		goto BEGIN

	case reflect.Interface:

		if current.IsNil() {
			return current, reflect.Interface
		}

		current = current.Elem()
		goto BEGIN
	case reflect.Invalid:
		return current, reflect.Invalid
	default:
		return current, current.Kind()
	}
}
