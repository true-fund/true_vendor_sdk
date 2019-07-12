package true_vendor_sdk

import (
	"github.com/wolvesdev/gohttplib/validator"
)

func reserveVerificationValidatorMap() validator.VMap {
	return validator.VMap{
		"id":      validator.RequiredStringValidators("id"),
		"user_id": validator.RequiredStringValidators("user_id"),
	}
}

func idValidatorMap() validator.VMap {
	return validator.VMap{
		"id": validator.RequiredStringValidators("id"),
	}
}

func resolveVerificationValidatorMap() validator.VMap {
	return validator.VMap{
		"id":      validator.RequiredStringValidators("id"),
		"token":   validator.RequiredStringValidators("token"),
		"verdict": validator.RequiredStringValidators("verdict"),
		"note":    validator.RequiredStringValidators("note"),
	}
}
