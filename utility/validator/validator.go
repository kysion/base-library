package validator

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/util/gvalid"
	"strings"
)

// RegisterServicePhone 函数用于注册一个自定义的验证规则，该规则用于校验“服务电话”是否符合指定格式。
// 该函数不接受任何参数，也不返回任何值。
func RegisterServicePhone() {
	// 初始化一个验证器实例。
	validator := gvalid.New()

	// 注册一个新的验证规则“service-phone”，并定义该规则的验证逻辑。
	gvalid.RegisterRule("service-phone", func(ctx context.Context, in gvalid.RuleFuncInput) error {
		// 获取待验证的电话号码。
		phoneNumber := in.Value.String()

		// 使用已有的验证器对电话号码进行“电话”规则的验证，如果通过则返回。
		if err := validator.Data(phoneNumber).Rules("phone").Run(ctx); err == nil {
			return nil
		}
		// 检查电话号码是否以“400”或“800”开头且长度为11，满足条件则返回。
		if (strings.HasPrefix(phoneNumber, "400") || strings.HasPrefix(phoneNumber, "800")) && len(phoneNumber) == 11 {
			return nil
		}
		// 检查电话号码是否以“95”开头且长度为5，满足条件则返回。
		if strings.HasPrefix(phoneNumber, "95") && len(phoneNumber) == 5 {
			return nil
		}
		// 使用已有的验证器对电话号码进行“电话”规则的验证，如果通过则返回。
		if err := validator.Data(phoneNumber).Rules("telephone").Run(ctx); err == nil {
			return nil
		}

		// 如果有自定义的错误信息，则使用该信息生成一个错误。
		if in.Message != "" {
			return gerror.New(in.Message)
		}

		// 如果以上验证均未通过，返回一个错误，说明电话号码不是有效的服务电话号码。
		return gerror.New("The ServicePhone value `" + phoneNumber + "` is not a valid ServicePhone number")
	})
}
