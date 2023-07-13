package errmsg

const (
	//TODO - before launch production modifiy password error
	ErrorMsgInvalidInput               = "invalid input"
	ErrorMsgNotFound                   = "record not found"
	ErrorMsgCantScanQueryResult        = "can't scan query result"
	ErrorMsgSomethingWentWrong         = "something went wrong"
	ErrorMsgPhoneNumberIsNotUnique     = "phone number is not unique"
	ErrorMsgCantOpenDatabase           = "cant open database"
	ErrorMsgCantInsertUserIntoDatabase = "cant inseret into DB"
	ErrorMsgWrongPassword              = "wrong password"
	ErrorMsgPhoneNumberIsNotValid      = "phone number is not valid"
	ErrorMsgPhoneNumberRequired        = "phone number Required"
	ErrorMsgNameLength                 = "name length should grater than 3"
	ErrorMsgUserNameRequired           = "Username Required"
	ErrorMsgPasswordRequired           = "password Required"
	ErrorMsgPasswordRegexValidate      = "password length should be 10 at least"
	ErrorMsgUserNotAllowed             = "user not allowed"
)
