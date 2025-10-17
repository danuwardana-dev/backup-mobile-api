package config

import "time"

type App struct {
	ServiceName         string        `envconfig:"APP_SERVICE_NAME" default:"TRY"`
	Mode                string        `envconfig:"APP_MODE" default:"development"`
	Env                 string        `envconfig:"APP_ENV" default:"local"`
	ContextTimeout      time.Duration `envconfig:"APP_CONTEXT_TIMEOUT" default:"2s"`
	OtpExpire           time.Duration `envconfig:"APP_OTP_EXPIRE" default:"60"`
	AccessKeyExpire     time.Duration `envconfig:"APP_ACCESS_KEY_EXPIRE" default:"900s"`
	XsessionExpire      time.Duration `envconfig:"APP_XSESSION_EXPIRE" default:"60s"`
	BiometricPrivateKey string        `envconfig:"APP_BIOMETRIC_PRIVATE_KEY" default:""`
	TimeZone            string        `envconfig:"APP_TIMEZONE" default:"Asia/Jakarta"`
}
