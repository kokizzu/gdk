package echo

import (
	"net/http"
	"time"

	"github.com/rizalgowandy/gdk/pkg/converter"

	"github.com/labstack/echo/v4"
)

// GetDefaultResponse is the default response for http get request
func GetDefaultResponse(c echo.Context) Response {
	return Response{
		Code:       converter.String(http.StatusInternalServerError),
		DisplayMsg: http.StatusText(http.StatusInternalServerError),
		RawMsg:     "default",
		RequestID:  c.Response().Header().Get(echo.HeaderXRequestID),
		ResultData: GetResultData{
			Param:         c.QueryParams(),
			GeneratedDate: time.Now().Format(time.RFC3339),
			TotalData:     "",
			Data:          nil,
		},
	}
}

// GetSuccessResponse is the success response for http get request
func GetSuccessResponse(c echo.Context, totalData int, data any) Response {
	return Response{
		Code:       converter.String(http.StatusOK),
		DisplayMsg: http.StatusText(http.StatusOK),
		RawMsg:     "",
		RequestID:  c.Response().Header().Get(echo.HeaderXRequestID),
		ResultData: GetResultData{
			Param:         c.QueryParams(),
			GeneratedDate: time.Now().Format(time.RFC3339),
			TotalData:     converter.String(totalData),
			Data:          data,
		},
	}
}

// PostDefaultResponse is the default response for http post request
func PostDefaultResponse(c echo.Context) Response {
	return Response{
		Code:       converter.String(http.StatusInternalServerError),
		DisplayMsg: http.StatusText(http.StatusInternalServerError),
		RawMsg:     "default",
		RequestID:  c.Response().Header().Get(echo.HeaderXRequestID),
		ResultData: PostResultData{
			Param:        c.QueryParams(),
			ExecutedDate: time.Now().Format(time.RFC3339),
			RowsAffected: 0,
		},
	}
}

// PostSuccessResponse is the success response for http post request
func PostSuccessResponse(c echo.Context, param any, rowsAffected int) Response {
	return Response{
		Code:       converter.String(http.StatusOK),
		DisplayMsg: http.StatusText(http.StatusOK),
		RawMsg:     "",
		RequestID:  c.Response().Header().Get(echo.HeaderXRequestID),
		ResultData: PostResultData{
			Param:        param,
			ExecutedDate: time.Now().Format(time.RFC3339),
			RowsAffected: rowsAffected,
		},
	}
}
