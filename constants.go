package validator

// 默认配置
const (
	defaultValidationTag    = "validate"
	defaultFieldDescribeTag = "desc"
	defaultOmitemptyTag     = "omitempty"
)

const (
	blank               = ""
	utf8HexComma        = "0x2C"
	utf8Pipe            = "0x7C"
	tagSeparator        = ","
	orSeparator         = "|"
	tagKeySeparator     = "="
	skipValidationTag   = "-"
	invalidValidation   = "Invalid validation tag on field %s"
	undefinedValidation = "Undefined validation function on field %s"
)
