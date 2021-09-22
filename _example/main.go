package main

import (
	"fmt"
	"github.com/one-gold-coin/validator"
)

// User contains user information
type User struct {
	FirstName string     `json:"fname" validate:"required,min=1" desc:"姓氏"`
	LastName  string     `json:"lname" validate:"required" desc:"名称"`
	Age       int        `json:"age" validate:"gte=0,lte=100" desc:"年龄"`
	Sex       int        `json:"sex" validate:"required,oneof=1 2" desc:"性别"`
	Email     string     `json:"email" validate:"required,email" desc:"邮件"`
	Job       *Job       `json:"job" validate:"required" desc:"工作"`
	Addresses []*Address `json:"addresses" validate:"required,min=1" desc:"地址"`
}

type Job struct {
	Id   int    `json:"id" validate:"required,min=1" desc:"工作ID"`
	Name string `json:"name" desc:"工作名称"`
}

// Address houses a users address information
type Address struct {
	Street string `json:"street" validate:"required,max=10" desc:"街道"`
	City   string `json:"city" validate:"required" desc:"城市"`
	Planet string `json:"planet" validate:"required" desc:"星球"`
	Phone  string `json:"phone" validate:"required,max=11" desc:"联系手机号"`
}

func main() {

	req := `
{
	"fname":"我F",
	"lname":"我L",
	"age":1,
	"sex":2,
	"email":"a@163.com",
	"favourite_color":"rgb",
	"job":{
		"id":1,
		"name":"职位"
	},
	"addresses":[
		{
			"street":"杨浦",
			"city":"上海",
			"planet":"planetStr1",
			"phone":"11111111111"
		},
		{
			"street":"杨浦2",
			"city":"上海2",
			"planet":"planetStr2",
			"phone":"11111111112"
		}
	]
}
`
	user := User{}
	err := validator.Validator().Binding(req, &user).Error()
	if err != nil {
		fmt.Printf("%s: %s", "err", err.Error())
		return
	}
	fmt.Printf("%s: %+v", "user", user)
}