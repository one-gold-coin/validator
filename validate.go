package validator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func New() *Validator {
	return &Validator{
		config: &Config{
			FieldDescribeTag: defaultFieldDescribeTag,
			ValidationTag:    defaultValidationTag,
			OmitemptyTag:     defaultOmitemptyTag,
		},
		translate: NewZhTranslate(),
	}
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
	isValidationFuncErr := v.extractStruct(value)
	// 解析参数校验错误信息
	if isValidationFuncErr {
		v.SetError(v.translate.Translate(v).GetErr())
	}
	return v
}

// 提取 Struct 字段信息
func (v *Validator) extractStruct(current reflect.Value) bool {
	//获取字段数量
	if current.Kind() == reflect.Ptr || current.Kind() == reflect.Interface {
		current = current.Elem()
	}
	// 结构体信息
	currentType := current.Type()
	// 获取结构体字段数量
	numFields := currentType.NumField()
	for i := 0; i < numFields; i++ {
		currentField := current.Field(i)
		// 获取每个字段信息
		currentStructField := currentType.Field(i)
		//获取字段名
		fieldName := currentStructField.Name
		//是否空字段"-", struct{-}
		if !currentStructField.Anonymous && currentStructField.PkgPath != blank {
			continue
		}
		//获取验证标签
		validateTag := currentStructField.Tag.Get(v.GetConfig().ValidationTag)
		//验证标签是否忽略或者为空
		if validateTag == skipValidationTag || validateTag == blank {
			continue
		}
		//如果有验证Tag,则进行数据验证
		if len(validateTag) > 0 {
			tags := v.parseFieldTags(current.Field(i), validateTag, currentStructField.Name)
			if tags != nil && tags.isHaveErr == true {
				//如果设置字段别名
				descTag := currentStructField.Tag.Get(v.GetConfig().FieldDescribeTag)
				v.field = &Field{Idx: i, AliasName: fieldName, Sf: &currentStructField, Tags: tags}
				if descTag != blank {
					v.field.AliasName = descTag
				}
				return true
			}
		}
		//递归处理,深层级逻辑
		if v.handleCurrentField(currentField) {
			return true
		}
	}
	return false
}

// 递归处理,深层级逻辑
func (v *Validator) handleCurrentField(current reflect.Value) bool {
	switch current.Type().Kind() {
	case reflect.Ptr, reflect.Interface:
		if v.handleCurrentField(current.Elem()) {
			return true
		}
	case reflect.Struct, reflect.Map:
		if v.extractStruct(current) {
			return true
		}
	case reflect.Slice, reflect.Array:
		for j := 0; j < current.Len(); j++ {
			//如果Slice是值类型时不校验,eg:[1],["a"]
			switch current.Index(j).Kind() {
			case reflect.String, reflect.Int, reflect.Int64:
				return false
			case reflect.Ptr, reflect.Interface:
				if v.handleCurrentField(current.Index(j).Elem()) {
					return true
				}
				return false
			}
			if v.extractStruct(current.Index(j)) {
				return true
			}
		}
	}
	return false
}

//验证数据
func (v *Validator) parseFieldTags(current reflect.Value, tagStr string, fieldName string) *Tag {
	var validaTag string
	var kind reflect.Kind
	var tags []string
	var tag Tag
	var vals []string
	//获取真实数据类型
	current, kind = v.extractTypeInternal(current)
	// 获取验证Tag列表
	tags = strings.Split(tagStr, tagSeparator)
	for i := 0; i < len(tags); i++ {
		validaTag = tags[i]
		// 当Tag == OmitemptyTag 时，再验证
		if v.GetConfig().OmitemptyTag != "" && validaTag == v.GetConfig().OmitemptyTag {
			switch kind {
			case reflect.Slice, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Chan, reflect.Func:
				if !current.IsNil() {
					continue
				}
			default:
				if current.IsValid() && current.Interface() != reflect.Zero(current.Type()).Interface() {
					continue
				}
			}
			return nil
		}
		tag.rv = &current
		orVials := strings.Split(validaTag, orSeparator)
		for j := 0; j < len(orVials); j++ {
			//获取验证值
			vals = strings.SplitN(orVials[j], tagKeySeparator, 2)
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
				validationFuncResult := validationFunc(&tag)
				//验证
				if !validationFuncResult {
					tag.isHaveErr = true
					return &tag
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
