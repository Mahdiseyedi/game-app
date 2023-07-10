package uservalidator

import (
	"game-app/dto"
	"game-app/pkg/errmsg"
	"game-app/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation"
	"regexp"
)

type Validator struct {
	repo Repository
}

type Repository interface {
	IsPhoneNumberUnique(phoneNumber string) (bool, error)
}

func New(repository Repository) Validator {
	return Validator{repo: repository}
}

func (v Validator) ValidateRegisterRequest(req dto.RegisterRequest) error {
	const op = "userValidator.ValidateRegisterRequest"

	//if !phoneNumber.IsValid(req.PhoneNumber) {
	//	return richerror.New(op).WithMessage(errmsg.ErrorMsgPhoneNumberIsNotValid).
	//		WithKind(richerror.KindInvalid).
	//		WithMeta(map[string]interface{}{"phone_number": req.PhoneNumber})
	//}

	if isUnique, err := v.repo.IsPhoneNumberUnique(req.PhoneNumber); !isUnique || err != nil {
		if err != nil {
			return richerror.New(op).WithErr(err)
		}

		if !isUnique {
			return richerror.New(op).
				WithMessage(errmsg.ErrorMsgPhoneNumberIsNotUnique).
				WithKind(richerror.KindInvalid).
				WithMeta(map[string]interface{}{"phone_number": req.PhoneNumber})
		}
	}

	return validation.ValidateStruct(&req,
		validation.Field(&req.Name,
			validation.Required.Error(errmsg.ErrorMsgUserNameRequired),
			validation.Length(5, 50).Error(errmsg.ErrorMsgNameLength)),
		validation.Field(&req.PhoneNumber,
			validation.Required.Error(errmsg.ErrorMsgPhoneNumberRequired),
			validation.Match(regexp.MustCompile(`"^09[0-9]{9}$"`)),
			validation.Length(5, 50)),
		validation.Field(&req.Password,
			validation.Required.Error(errmsg.ErrorMsgPasswordRequired),
			validation.Match(regexp.MustCompile(`^[A-Za-z0-9!@#%^&*]{8,}$`)).
				Error(errmsg.ErrorMsgPasswordRegexValidate)),
	)
}
