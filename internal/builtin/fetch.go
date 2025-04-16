package builtin

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/flarelang/flare/internal/errs"
	"github.com/flarelang/flare/internal/models"
	"github.com/flarelang/flare/lang"
)

type fetchConfig struct {
	url     string
	method  string
	body    io.Reader
	headers http.Header
	debug   *models.Debug
}

func fnFetch(args []lang.Object) (lang.Object, error) {
	url, ok := args[0].Value().(string)
	if !ok {
		return nil, errs.WithDebug(fmt.Errorf("invalid argument type for url, want: string, got: %s", args[0].Type()), args[0].Debug())
	}

	var conf = &fetchConfig{
		url:   url,
		debug: args[0].Debug(),
	}

	variadicArgs := args[1].Value().([]lang.Object)

	if len(variadicArgs) != 1 {
		return fnRealFetch(conf)
	}

	argConfig, ok := variadicArgs[0].(*lang.Array)
	if !ok || argConfig.Type() != lang.TArray {
		return nil, errs.WithDebug(fmt.Errorf("invalid argument type for config, want: array, got: %s", args[1].Type()), args[1].Debug())
	}

	if method, ok := argConfig.Access("method"); ok {
		if method.Type() != lang.TString {
			return nil, errs.WithDebug(fmt.Errorf("invalid argument type for config.method, want: string, got: %s", method.Type()), method.Debug())
		}
		conf.method = strings.ToUpper(method.Value().(string))
	}

	if body, ok := argConfig.Access("body"); ok {
		conf.body = strings.NewReader(body.String())
	}

	if rawHeaders, ok := argConfig.Access("headers"); ok {
		headers, ok := rawHeaders.(*lang.Array)
		if !ok || rawHeaders.Type() != lang.TArray {
			return nil, errs.WithDebug(fmt.Errorf("invalid argument type for config.headers, want: array, got: %s", headers.Type()), headers.Debug())
		}
		conf.headers = make(http.Header)

		for _, headerKey := range headers.Keys {
			header, ok := headers.Access(headerKey.Value())
			if !ok {
				continue
			}

			if header.Type() == lang.TList {
				for _, item := range header.Value().([]lang.Object) {
					conf.headers.Add(headerKey.String(), item.String())
				}
			} else {
				conf.headers.Set(headerKey.String(), header.String())
			}
		}
	}

	return fnRealFetch(conf)
}

func fnRealFetch(conf *fetchConfig) (lang.Object, error) {
	req, err := http.NewRequest(conf.method, conf.url, conf.body)
	req.Header.Set("User-Agent", "Flare-Http-Fetch/1.0")
	if err != nil {
		return nil, err
	}

	for key, values := range conf.headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	res := make(map[string]lang.Object)

	res["body"] = fnFetchGetBody(body)
	headerMap := make(map[string]lang.Object)

	for key, values := range resp.Header {
		if len(values) == 1 {
			headerMap[key] = lang.NewString("headers", values[0], conf.debug)
		} else {
			items := make([]lang.Object, len(values))
			for i, value := range values {
				items[i] = lang.NewString(key, value, nil)
			}
			headerMap[key] = lang.NewList("headers", items, conf.debug)
		}
	}

	res["headers"] = lang.NewArrayMap("headers", conf.debug, headerMap)
	res["status"] = lang.NewString("status", resp.Status, conf.debug)
	res["statusCode"] = lang.NewInteger("statusCode", resp.StatusCode, conf.debug)
	res["redirectURL"] = lang.NewString("redirectURL", resp.Request.URL.String(), conf.debug)

	return lang.NewArrayMap("fetch", conf.debug, res), nil
}

func fnFetchGetBody(body []byte) lang.Object {
	var data any

	err := json.Unmarshal(body, &data)
	if err == nil {
		obj, err := lang.FromValue(data)
		if err == nil {
			return obj
		}
	}

	return lang.NewString("body", string(body), nil)
}
