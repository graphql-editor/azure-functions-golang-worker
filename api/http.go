package api

import (
	"net/http"
	"strconv"

	"github.com/graphql-editor/azure-functions-golang-worker/converters"
	"github.com/graphql-editor/azure-functions-golang-worker/rpc"
)

func encodeStatusCode(statusCode int) string {
	if statusCode == 0 {
		return "200"
	}
	return strconv.FormatInt(int64(statusCode), 10)
}

func encodeResponseObject(resp *Response) (td *rpc.TypedData, err error) {
	body, err := converters.Marshal(resp.Body)
	if err == nil {
		var cookies []*rpc.RpcHttpCookie
		cookies, err = encodeCookies(resp.Cookies)
		if err == nil {
			td = &rpc.TypedData{
				Data: &rpc.TypedData_Http{
					Http: &rpc.RpcHttp{
						Headers:    encodeHeaders(resp.Headers),
						Cookies:    cookies,
						StatusCode: encodeStatusCode(resp.StatusCode),
						Body:       body,
					},
				},
			}
		}
	}
	return
}

func encodeHeaders(headers http.Header) map[string]string {
	if headers == nil {
		return nil
	}
	h := make(map[string]string)
	for k := range headers {
		h[k] = headers.Get(k)
	}
	return h
}

func encodeCookie(cookie Cookie) (rpcCookie *rpc.RpcHttpCookie, err error) {
	ts, err := converters.EncodeNullableTimestamp(cookie.Expires)
	if err == nil {
		rpcCookie = &rpc.RpcHttpCookie{
			Name:     cookie.Name,
			Value:    cookie.Value,
			Domain:   converters.EncodeNullableString(cookie.Domain),
			Path:     converters.EncodeNullableString(cookie.Path),
			Expires:  ts,
			Secure:   converters.EncodeNullableBool(cookie.Secure),
			HttpOnly: converters.EncodeNullableBool(cookie.HTTPOnly),
			MaxAge:   converters.EncodeNullableDouble(cookie.MaxAge),
		}
		switch cookie.SameSite {
		case Lax:
			rpcCookie.SameSite = rpc.RpcHttpCookie_Lax
		case Strict:
			rpcCookie.SameSite = rpc.RpcHttpCookie_Strict
		}
	}
	return
}

func encodeCookies(cookies Cookies) (rpcCookies []*rpc.RpcHttpCookie, err error) {
	if len(cookies) != 0 {
		rpcCookies = make([]*rpc.RpcHttpCookie, len(cookies))
		for i, ck := range cookies {
			rpcCookies[i], err = encodeCookie(ck)
			if err != nil {
				break
			}
		}
	}
	return
}
