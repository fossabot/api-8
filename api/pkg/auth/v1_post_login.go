package auth

import (
	"context"
	"net/http"

	"github.com/devlover-id/api/pkg/api"
	validation "github.com/go-ozzo/ozzo-validation"
)

type v1PostLoginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (p *v1PostLoginPayload) validate() error {
	return validation.ValidateStruct(p,
		validation.Field(&p.Username,
			validation.Required.Error("username tidak boleh kosng"),
		),
		validation.Field(&p.Password,
			validation.Required.Error("password tidak boleh kosong"),
		),
	)
}

func V1PostLogin(ctx context.Context, req api.Request) api.Response {
	var payload v1PostLoginPayload
	if err := req.Bind(&payload); err != nil {
		return api.InternalServerErrResp(err)
	}
	if validationErrors := payload.validate(); validationErrors != nil {
		return api.JSONResponse(http.StatusBadRequest, validationErrors)
	}

	uid, token, err := login(payload.Username, payload.Password)
	if err == errNilUserPassword {
		return api.ValidationErrResp(map[string]string{
			"password": "password belum di set",
		})
	}
	if err != nil {
		return api.InternalServerErrResp(err)
	}
	if len(token) == 0 {
		return api.JSONResponse(http.StatusUnauthorized, nil)
	}

	if err := createUserSession(uid, token, req.Header().Get("user-agent"), req.ClientIP()); err != nil {
		return api.InternalServerErrResp(err)
	}

	return api.JSONResponse(http.StatusOK, map[string]string{
		"token": token,
	})
}
