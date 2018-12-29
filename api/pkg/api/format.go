package api

import (
	"net/http"
)

// OKResp creates http ok response
func OKResp(p interface{}) Response {
	return JSONResponse(http.StatusOK, p)
}

// UnauthorizedResp creates http unauthorized response
func UnauthorizedResp() Response {
	return JSONResponse(http.StatusUnauthorized, map[string]string{
		"error": "unauthorized",
	})
}

// InternalServerErrResp create http internal server error response
func InternalServerErrResp(err error) Response {
	return JSONResponse(http.StatusInternalServerError, map[string]string{
		"error": "something wrong with our servers :(",
	})
}

// ValidationErrResp creates validation error response. The key of errors should be
// the field of error and the value should be the error message.
func ValidationErrResp(errors map[string]string) Response {
	return JSONResponse(http.StatusBadRequest, errors)
}

func CodeResp(code int) {

}
