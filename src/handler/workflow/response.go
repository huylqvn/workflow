package workflow

import (
	"net/http"

	"github.com/huylqvn/httpserver"
	"go.elastic.co/apm/v2"
)

type Mess string

const (
	CreateSuccess       Mess = "create success"
	Success             Mess = "success"
	RequestError        Mess = "request error"
	Exception           Mess = "exception error"
	InternalServerError Mess = "internal server error"
	NotFound            Mess = "not found"
)

func customResponse(w http.ResponseWriter, message Mess, data interface{}) {
	code := toCode(message)
	if code >= 400 {
		httpserver.NewResponse(code, string(message), data.(string), "").ToJson(w)
	} else {
		httpserver.NewResponse(code, string(message), "", data).ToJson(w)
	}
}

func toCode(message Mess) int {
	switch message {
	case CreateSuccess:
		return 201
	case Success:
		return 200
	case RequestError:
		return 400
	case InternalServerError, Exception:
		return 500
	case NotFound:
		return 404
	default:
		return 204
	}
}

func apmSendErr(span *apm.Span, err error) {
	e := apm.DefaultTracer().
		NewError(err)
	e.SetSpan(span)
	e.Send()
}
