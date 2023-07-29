package sys_enum_captcha

import "github.com/kysion/base-library/utility/enum"

// 验证码类型：1注册，2登录，4找回用户名/修改用户名，8找回密码/重置密码，16设置手机号码 ，32设置邮箱, 64忘记用户名&密码

type CaptchaTypeEnum enum.IEnumCode[int]

type captchaType struct {
	Register                  CaptchaTypeEnum
	Login                     CaptchaTypeEnum
	SetUserName               CaptchaTypeEnum
	SetPassword               CaptchaTypeEnum
	SetMobile                 CaptchaTypeEnum
	SetMail                   CaptchaTypeEnum
	ForgotUserNameAndPassword CaptchaTypeEnum
}

var CaptchaType = captchaType{
	Register:                  enum.New[CaptchaTypeEnum](1, "register"),
	Login:                     enum.New[CaptchaTypeEnum](2, "login"),
	SetUserName:               enum.New[CaptchaTypeEnum](4, "setUserName"),
	SetPassword:               enum.New[CaptchaTypeEnum](8, "setPassword"),
	SetMobile:                 enum.New[CaptchaTypeEnum](16, "setMobile"),
	SetMail:                   enum.New[CaptchaTypeEnum](32, "setMail"),
	ForgotUserNameAndPassword: enum.New[CaptchaTypeEnum](64, "forgotUserNameAndPassword"),
	// 可拓展.....
}

func (e captchaType) New(code int, description string) CaptchaTypeEnum {
	if (code & e.Register.Code()) == e.Register.Code() {
		return e.Register
	}
	if (code & e.Login.Code()) == e.Login.Code() {
		return e.Login
	}
	if (code & e.SetUserName.Code()) == e.SetUserName.Code() {
		return e.SetUserName
	}
	if (code & e.SetPassword.Code()) == e.SetPassword.Code() {
		return e.SetPassword
	}
	if (code & e.SetMobile.Code()) == e.SetMobile.Code() {
		return e.SetMobile
	}
	if (code & e.SetMail.Code()) == e.SetMail.Code() {
		return e.SetMail
	}
	if (code & e.ForgotUserNameAndPassword.Code()) == e.ForgotUserNameAndPassword.Code() {
		return e.ForgotUserNameAndPassword
	}

	return enum.New[CaptchaTypeEnum](code, description)
}
