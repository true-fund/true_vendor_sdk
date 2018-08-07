package true_vendor_sdk

import "github.com/alexmay23/httputils"

func reserveVerificationValidatorMap()httputils.VMap{
	return httputils.VMap{
		"id": httputils.RequiredStringValidators("id"),
		"user_id": httputils.RequiredStringValidators("user_id"),
	}
}

func idValidatorMap()httputils.VMap{
	return httputils.VMap{
		"id": httputils.RequiredStringValidators("id"),
	}
}

func resolveVerificationValidatorMap()httputils.VMap{
	return httputils.VMap{
		"id": httputils.RequiredStringValidators("id"),
		"token": httputils.RequiredStringValidators("token"),
	}
}
