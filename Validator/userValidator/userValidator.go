package userValidator

import (
	"game-app/dto"
	"game-app/pkg/errmsg"
	"game-app/pkg/phoneNumber"
	"game-app/pkg/richerror"
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

	if !phoneNumber.IsValid(req.PhoneNumber) {
		return richerror.New(op).WithMessage(errmsg.ErrorMsgPhoneNumberIsNotValid).
			WithKind(richerror.KindInvalid).
			WithMeta(map[string]interface{}{"phone_number": req.PhoneNumber})
	}

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

	return nil
}

func (v Validator) NameServiceValidator(req dto.RegisterRequest) error {
	const op = "userValidator.NameServiceValidator"

	//TODO - add 3 to config
	if len(req.Name) < 3 {
		return richerror.New(op).
			WithMessage(errmsg.ErrorMsgNameLength).
			WithKind(richerror.KindInvalid).
			WithMeta(map[string]interface{}{"userName: ": req.Name})
	}

	//this if statement just for keep validator structure regular
	if req.Name == "userName" {
		return richerror.New(op).
			WithMessage(errmsg.ErrorMsgUserNameSelf).
			WithKind(richerror.KindInvalid).
			WithMeta(map[string]interface{}{"userName: ": req.Name})
	}

	return nil
}
