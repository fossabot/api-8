package auth

import (
	"context"
	"errors"
	"net/http"

	"github.com/devlover-id/api/pkg/api"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type v1PostRegisterPayload struct {
	Name            string `json:"name"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
	Email           string `json:"email"`
}

func (p *v1PostRegisterPayload) validate() error {
	return validation.ValidateStruct(p,
		validation.Field(&p.Name,
			validation.Required.Error("nama tidak boleh kosong"),
			validation.Length(4, 0).Error("nama harus lebih dari 3 karakter"),
		),
		validation.Field(&p.Username,
			validation.Required.Error("username tidak boleh kosong"),
			validation.Length(4, 0).Error("username harus lebih dari 4 karakter"),
		),
		validation.Field(&p.Password,
			validation.Required.Error("password tidak boleh kosong"),
			validation.Length(8, 0).Error("password harus lebih dari 8 karakter"),
		),
		validation.Field(&p.PasswordConfirm,
			validation.Required.Error("konfirmasi password tidak boleh kosong"),
			validation.By(validatePasswordConfirm(p.Password)),
		),
		validation.Field(&p.Email,
			validation.Required.Error("email tidak boleh kosong"),
			is.Email.Error("bukan valid email"),
		),
	)
}

func validatePasswordConfirm(password string) func(interface{}) error {
	return func(value interface{}) error {
		val, ok := value.(string)
		if !ok || val != password {
			return errors.New("konfirmasi password tidak sama dengan password")
		}
		return nil
	}
}

func V1PostRegister(ctx context.Context, req api.Request) api.Response {
	var payload v1PostRegisterPayload
	if err := req.Bind(&payload); err != nil {
		return api.InternalServerErrResp(err)
	}
	if validationErrors := payload.validate(); validationErrors != nil {
		return api.JSONResponse(http.StatusBadRequest, validationErrors)
	}

	user, err := createUser(payload.Username, payload.Password, payload.Name, payload.Email)
	if err != nil {
		return api.InternalServerErrResp(err)
	}

	return api.JSONResponse(http.StatusCreated, user)
}
