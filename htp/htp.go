package htp

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/btagrass/gobiz/app"
	"github.com/btagrass/gobiz/utl"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

var (
	Ip      string
	Port    uint
	Timeout time.Duration
)

func init() {
	Ip = os.Getenv("HTTP_IP")
	if Ip == "" {
		Ip = viper.GetString("http.ip")
	}
	if Ip == "" {
		Ip, _ = utl.GetIp()
	}
	Port = cast.ToUint(os.Getenv("HTTP_PORT"))
	if Port == 0 {
		Port = viper.GetUint("http.port")
	}
	Timeout = cast.ToDuration(os.Getenv("HTTP_TIMEOUT"))
	if Timeout == 0 {
		Timeout = viper.GetDuration("http.timeout")
	}
}

func Delete(url string, headers map[string]string, r ...any) (*resty.Response, error) {
	req := resty.New().
		SetTimeout(Timeout).
		SetTLSClientConfig(&tls.Config{
			InsecureSkipVerify: true,
		}).
		OnAfterResponse(respond).
		R().
		SetHeaders(headers).
		ForceContentType("application/json")
	if len(r) > 0 {
		req.SetResult(r[0])
	}
	res, err := req.Delete(GetFullUrl(url))
	slog.Debug("", "method", req.Method, "url", req.URL, "result", res)
	return res, err
}

func Get(url string, headers map[string]string, r ...any) (*resty.Response, error) {
	req := resty.New().
		SetTimeout(Timeout).
		SetTLSClientConfig(&tls.Config{
			InsecureSkipVerify: true,
		}).
		OnAfterResponse(respond).
		R().
		SetHeaders(headers).
		ForceContentType("application/json")
	if len(r) > 0 {
		req.SetResult(r[0])
	}
	res, err := req.Get(GetFullUrl(url))
	slog.Debug("", "method", req.Method, "url", req.URL, "result", res)
	return res, err
}

func GetFileUrl(filePath string) string {
	return fmt.Sprintf("http://%s:%d/%s", Ip, Port, utl.Replace(filePath, app.DataDir, "data"))
}

func GetFullUrl(url string) string {
	if strings.HasPrefix(url, "http") {
		return url
	} else if strings.HasPrefix(url, ":") {
		return fmt.Sprintf("http://%s%s", Ip, url)
	} else if strings.HasPrefix(url, "/") {
		return fmt.Sprintf("http://%s:%d%s", Ip, Port, url)
	}
	return fmt.Sprintf("http://%s:%d/%s", Ip, Port, url)
}

func Post(url string, headers map[string]string, data any, r ...any) (*resty.Response, error) {
	req := resty.New().
		SetTimeout(Timeout).
		SetTLSClientConfig(&tls.Config{
			InsecureSkipVerify: true,
		}).
		OnAfterResponse(respond).
		R().
		SetHeader("Accept", "*/*").
		SetHeader("Content-Type", "application/json").
		SetHeaders(headers).
		SetBody(data).
		ForceContentType("application/json")
	if len(r) > 0 {
		req.SetResult(r[0])
	}
	res, err := req.Post(GetFullUrl(url))
	slog.Debug("", "method", req.Method, "url", req.URL, "data", data, "result", res)
	return res, err
}

func PostFile(url string, headers map[string]string, files map[string]string, r ...any) (*resty.Response, error) {
	req := resty.New().
		SetTLSClientConfig(&tls.Config{
			InsecureSkipVerify: true,
		}).
		OnAfterResponse(respond).
		R().
		SetHeaders(headers).
		SetFiles(files).
		ForceContentType("application/json")
	if len(r) > 0 {
		req.SetResult(r[0])
	}
	res, err := req.Post(GetFullUrl(url))
	slog.Debug("", "method", req.Method, "url", req.URL, "data", files, "result", res)
	return res, err
}

func PostForm(url string, headers map[string]string, data map[string]string, r ...any) (*resty.Response, error) {
	req := resty.New().
		SetTimeout(Timeout).
		SetTLSClientConfig(&tls.Config{
			InsecureSkipVerify: true,
		}).
		OnAfterResponse(respond).
		R().
		SetHeader("Accept", "*/*").
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetHeaders(headers).
		SetFormData(data).
		ForceContentType("application/json")
	if len(r) > 0 {
		req.SetResult(r[0])
	}
	res, err := req.Post(GetFullUrl(url))
	slog.Debug("", "method", req.Method, "url", req.URL, "data", data, "result", res)
	return res, err
}

func Put(url string, headers map[string]string, data any, r ...any) (*resty.Response, error) {
	req := resty.New().
		SetTimeout(Timeout).
		SetTLSClientConfig(&tls.Config{
			InsecureSkipVerify: true,
		}).
		OnAfterResponse(respond).
		R().
		SetHeader("Accept", "*/*").
		SetHeader("Content-Type", "application/json").
		SetHeaders(headers).
		SetBody(data).
		ForceContentType("application/json")
	if len(r) > 0 {
		req.SetResult(r[0])
	}
	res, err := req.Put(GetFullUrl(url))
	slog.Debug("", "method", req.Method, "url", req.URL, "data", data, "result", res)
	return res, err
}

func SaveFile(url string, headers map[string]string, file ...string) (*resty.Response, error) {
	filePath := filepath.Base(url)
	if len(file) > 0 {
		filePath = file[0]
	}
	req := resty.New().
		SetTLSClientConfig(&tls.Config{
			InsecureSkipVerify: true,
		}).
		R().
		SetHeaders(headers).
		SetOutput(filePath)
	res, err := req.Get(GetFullUrl(url))
	slog.Debug("", "method", req.Method, "url", req.URL, "result", res)
	return res, err
}

func respond(c *resty.Client, res *resty.Response) error {
	if res.StatusCode() != http.StatusOK {
		return errors.New(res.Status())
	}
	r := cast.ToStringMap(res.String())
	code, ok := r["code"]
	if !ok {
		code = r["error_code"]
	}
	if !ok {
		code = r["status"]
	}
	code = cast.ToInt(code)
	if code != http.StatusOK && code != 0 {
		msg, ok := r["msg"]
		if !ok {
			msg = r["desp"]
		}
		if !ok {
			msg = r["message"]
		}
		return fmt.Errorf("api error: %s -> %d", msg, code)
	}
	return nil
}
