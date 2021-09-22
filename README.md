# validator
Golang 参数验证器，目前只支持POST请求，JSON格式参数验证

# 亮点
1、验证时只要有一个错误，错误信息立即返回

2、可自定义参数别名显示错误信息；详情见_example文件

# 使用
```
go mod -u github.com/one-gold-coin/validator
```

# Gin 对接示例
```

1、Gin框架参数接收验证（只支持POST+参数是json格式验证）
package apicontroller

func Validator(ctx *gin.Context, obj interface{}) error {
	body, _ := ioutil.ReadAll(ctx.Request.Body)
	v := validator.Validator()
	v.SetConfig(&validator.Config{
		ValidationTag:    "binding",
		FieldDescribeTag: "desc",
	})
	return v.Binding(string(body), obj).Error()
}

2、参数验证定义
package api

type UserRegisterForm struct {
	Username string `json:"username"`                            // 登录用户名
	Email    string `json:"email"`                               // 登录Email
	Phone    string `json:"phone" binding:"required" desc:"手机号"` // 登录手机号
}

3、参数验证
package api

req := api.UserRegisterForm{}
if err := controller.Validator(ctx, &req); err != nil {
    // code ...
    return
}
```

# 支持的验证方式,目前只支持常用验证类型
```
required
len
eq
ne
lt
lte
gt
gte
email
min
max
oneof
```


# 其他
参考复用github.com/go-playground/validator/v10部分代码逻辑
