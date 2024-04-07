package requestmodels

import "errors"

type ICreateSingleUserRequestModel interface {
	Validate() error
}
type CreateSingleUserRequestModel struct {
	Username string
	Email    string
}

func BuildCreateSingleUserRequestModel() ICreateSingleUserRequestModel {
	return &CreateSingleUserRequestModel{}
}
func (createSingleUserRequestModel *CreateSingleUserRequestModel) Validate() error {
	if createSingleUserRequestModel.Username == "" {
		return errors.New("username field is required")
	} else if createSingleUserRequestModel.Email == "" {
		return errors.New("email field is required")
	}
	return nil
}
