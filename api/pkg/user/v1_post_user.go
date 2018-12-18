package user

import (
	"context"
	"errors"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"gitlab.com/pinterkode/pinterkode/api/pkg/api"
)

type v1PostUserRequest struct {
	Username       string `json:"username"`
	Password       string `json:"password"`
	PasswordRepeat string `json:"password_repeat"`
	Email          string `json:"email"`
}

func (r *v1PostUserRequest) validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Username,
			validation.Required.Error("username tidak boleh kosong"),
			validation.Length(4, 0).Error("username harus lebih dari 4 karakter"),
		),
		validation.Field(&r.Password,
			validation.Required.Error("password tidak boleh kosong"),
			validation.Length(8, 0).Error("password harus lebih dari 8 karakter"),
		),
		validation.Field(&r.PasswordRepeat,
			validation.Required.Error("password_repeat tidak boleh kosong"),
			validation.By(validatePasswordRepeat(r.Password)),
		),
		validation.Field(&r.Email,
			validation.Required.Error("email tidak boleh kosong"),
			is.Email.Error("bukan valid email"),
		),
	)
}

func validatePasswordRepeat(password string) func(interface{}) error {
	return func(value interface{}) error {
		val, ok := value.(string)
		if !ok || val != password {
			return errors.New("password_repeat tidak sama dengan password")
		}
		return nil
	}
}

func V1PostUser(ctx context.Context, req api.Request) api.Response {
	var payload v1PostUserRequest
	if err := req.Bind(&payload); err != nil {
		return api.InternalServerErrResp(err)
	}
	if valErrors := payload.validate(); valErrors != nil {
		return api.OKResp(valErrors)
	}

	user, err := createUser(payload.Username, payload.Password, payload.Email)
	if err != nil {
		return api.InternalServerErrResp(err)
	}

	return api.OKResp(user)
}
