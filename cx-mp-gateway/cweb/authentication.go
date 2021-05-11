package cweb

import (
	"errors"
	"github.com/gin-gonic/gin"
)

func getHeaderFromCtx(ctx *gin.Context) (*RequestHeader, error) {

	header := ctx.Request.Header
	reqHeader := &RequestHeader{}
	var hasParamsErr = false
	var paramsErrInfo = ""


	if header[HttpHeaderStringDomainId] == nil {
		hasParamsErr = true
		paramsErrInfo = paramsErrInfo + " " + "not find domain id in header"
	} else {
		reqHeader.DomainID = header[HttpHeaderStringDomainId][0]
	}
	if header[HttpHeaderStringAppId] == nil {
		hasParamsErr = true
		paramsErrInfo = paramsErrInfo + " " + "not find app id in header"
	} else {
		reqHeader.AppID = header[HttpHeaderStringAppId][0]
	}

	if header[HttpHeaderStringCaller] == nil {
		hasParamsErr = true
		paramsErrInfo = paramsErrInfo + " " + "not find caller in header"
	} else {
		reqHeader.Caller = CallerType(header[HttpHeaderStringCaller][0])
	}

	switch reqHeader.Caller {
	case CallerTypeServices:
		if header[HttpHeaderStringSecretKey] == nil {
			hasParamsErr = true
			paramsErrInfo = paramsErrInfo + " " + "not find secret key id in header"
		} else {
			reqHeader.SecretKey = header[HttpHeaderStringSecretKey][0]
		}

		reqHeader.Uid = "0" // grpc 会强转这个字段，导致报错，所以设置为0
	case CallerTypeGeneralUsers:
		if header[HttpHeaderStringUid] == nil {
			hasParamsErr = true
			paramsErrInfo = paramsErrInfo + " " + "not find user id in header"
		} else {
			reqHeader.Uid = header[HttpHeaderStringUid][0]
		}

		if header[HttpHeaderStringAccessToken] == nil {
			hasParamsErr = true
			paramsErrInfo = paramsErrInfo + " " + "not find access token id in header"
		} else {
			reqHeader.AccessToken = header[HttpHeaderStringAccessToken][0]
		}
	default:
		return nil, errors.New("caller was wrong")
	}

	if header[HttpHeaderStringCommandId] == nil {
		hasParamsErr = true
		paramsErrInfo = paramsErrInfo + " " + "not find command id id in header"
	} else {
		reqHeader.MethodName = header[HttpHeaderStringCommandId][0]
	}

	if header[HttpHeaderStringCommandVersion] == nil {
		hasParamsErr = true
		paramsErrInfo = paramsErrInfo + " " + "not find command version id in header"
	} else {
		reqHeader.CommandVersion = header[HttpHeaderStringCommandVersion][0]
	}

	if header[HttpHeaderStringServerName] == nil {
		hasParamsErr = true
		paramsErrInfo = paramsErrInfo + " " + "not find Server name in header"
	} else {
		reqHeader.ServerName = header[HttpHeaderStringServerName][0]
	}

	if hasParamsErr {
		return nil, errors.New(paramsErrInfo)
	}

	return reqHeader, nil
}
