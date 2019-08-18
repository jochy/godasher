package main

import (
	"../../godasher"
	"bytes"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
)

var passphrase, _ = godasher.GenerateRandomString(16)

func proxy(w http.ResponseWriter, r *http.Request) {
	cryptedToken := strings.TrimSuffix(r.URL.Path[len("/iframeproxy/"):], "/")
	tokenString, err := godasher.DecodeBase64AndDecrypt(cryptedToken, passphrase)

	if err != nil {
		_, _ = fmt.Fprintf(w, "%v", err)
		return
	}

	token, err := jwt.Parse(string(tokenString), func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(passphrase), nil
	})

	if err != nil {
		_, _ = fmt.Fprintf(w, "%v", err)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		remoteUrl, _ := url.Parse(fmt.Sprint(claims["url"]))
		user := claims["basicAuthUser"]
		password := claims["basicAuthPassword"]

		if user != "" {
			r.SetBasicAuth(fmt.Sprint(user), fmt.Sprint(password))
		}

		p := httputil.NewSingleHostReverseProxy(remoteUrl)
		r.URL = remoteUrl
		r.Host = ""
		r.Header.Del("Accept-Encoding")
		p.ModifyResponse = func(response *http.Response) error {
			b, err := ioutil.ReadAll(response.Body) //Read html
			if err != nil {
				return err
			}
			err = response.Body.Close()
			if err != nil {
				return err
			}
			if bytes.Contains(b, []byte("<head>")) {
				b = bytes.Replace(b, []byte("<head>"), []byte("<head><base href=\""+fmt.Sprint(claims["url"])+"\" />"), -1)
			} else if bytes.Contains(b, []byte("</head>")) {
				b = bytes.Replace(b, []byte("</head>"), []byte("<base href=\""+fmt.Sprint(claims["url"])+"\" /></head>"), -1)
			}

			body := ioutil.NopCloser(bytes.NewReader(b))
			response.Body = body
			response.ContentLength = int64(len(b))
			response.Header.Set("Content-Length", strconv.Itoa(len(b)))
			return nil
		}
		p.ServeHTTP(w, r)
	}
}

func init() {
	http.HandleFunc("/iframeproxy/", proxy)
}

func Setup(config godasher.Config) {
	// Nothing
}

func Render(component godasher.Component) template.HTML {
	attribute := "src"
	value := component.Data["url"]

	if val, ok := component.Data["src"]; ok {
		attribute = "srcdoc"
		value = val
	} else if shouldProxy := component.Data["proxy"]; strings.ToLower(shouldProxy) == "true" {
		token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
			"url":               component.Data["url"],
			"basicAuthUser":     component.Data["basic-auth-user"],
			"basicAuthPassword": component.Data["basic-auth-password"],
		})
		str, err := token.SignedString([]byte(passphrase))
		if err != nil {
			return template.HTML(fmt.Sprint(err))
		}

		tokenCrypted, er := godasher.EncryptAndEncodeBase64([]byte(str), passphrase)
		if er != nil {
			return template.HTML(fmt.Sprint(er))
		}
		value = "/iframeproxy/" + tokenCrypted + "/"
	}

	return template.HTML(fmt.Sprintf("<iframe %v=\"%v\" style=\"width: 100%%; height: 100%%; border:none;\"></iframe>",
		attribute,
		strings.ReplaceAll(value, "\"", "'")))
}
