// Copyright © 2022 UCloud. All rights reserved.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"reflect"
	"strings"
	"time"

	"git.ucloudadmin.com/ubase/ubase"
	"github.com/google/uuid"
)

func httpDemo() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := ubase.SimpleContext()

		rr, err := httputil.DumpRequest(r, true)
		if err != nil {
			fmt.Println("DumpRequest error:", err)
		} else {
			ctx.Info("DumpRequest", "request", string(rr))
			fmt.Println("DumpRequest result:", string(rr))
		}

		fmt.Printf("[server]incoming URL is: %+v\n", r.URL)
		typ := r.Header.Get("Content-Type")
		fmt.Println(w, "[server]current Content-Type is", typ)
		fmt.Fprintln(w, "[server]parsed Content-Type", typ)
		if strings.Contains(typ, "application/json") {
			if r.Form == nil {
				fmt.Println("not carry form")
			}
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				fmt.Println("[server]read body error:", err)
				fmt.Fprintln(w, "[server]read body error:", err)
			} else {
				fmt.Println("[server]###overload body", body)
				fmt.Println("[server]overload body", string(body))
				ctx.Info("", "###body", body)
				ctx.Info("", "body", string(body))
				fmt.Fprintln(w, "[server]overload body", string(body))
			}

		} else if strings.Contains(typ, "application/x-www-form-urlencoded") {
			if r.Body == nil {
				fmt.Println("not carry body")
			}
			// body, err := ioutil.ReadAll(r.Body)
			// fmt.Println("[server]overload body", string(body), "error", err)

			if err := r.ParseForm(); err != nil {
				fmt.Println("[server]parse form error:", err)
				fmt.Fprintln(w, "[server]parse form error:", err)
			} else {
				if r.FormValue("proxy") == "true" {
					fmt.Println("[server]forward to proxy")
					fmt.Fprintln(w, "[server]forward to proxy")
					NewProxy(w, r, "localhost:8081/proxy")
				}

				fmt.Println(w, "[server]overload form", r.Form)
				fmt.Fprintln(w, "[server]overload form", r.Form)
				if data, err := TransformToJson(r.Form); err != nil {
					fmt.Fprintln(w, "[server]transform to json error", err)
				} else {
					fmt.Println(w, "[server]transform to json", string(data))
					fmt.Fprintln(w, "[server]transform to json", string(data))
				}
			}
		} else {
			fmt.Fprintln(w, "[server]not support content-type", typ)
		}
		// ww, err := httputil.DumpResponse(w.Header(), true)
		// if err != nil {
		// 	fmt.Println("DumpRequest error:", err)
		// } else {
		// 	fmt.Println("DumpRequest result:", rr)
		// }
	})
	log.Fatalln(http.ListenAndServe(":8080", nil))
}

func httpProxy() {
	http.HandleFunc("/proxy", func(w http.ResponseWriter, r *http.Request) {
		// fmt.Printf("[server]incoming request is: %+v\n", r.URL)
		contentType := r.Header.Get("Content-Type")
		fmt.Println(w, "[server]current Content-Type is", contentType)
		fmt.Fprintln(w, "[server]parsed Content-Type", contentType)
		if strings.Contains(contentType, "application/json") {
			if r.Form == nil {
				fmt.Println("not carry form")
			}
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				fmt.Println("[server]read body error:", err)
				fmt.Fprintln(w, "[server]read body error:", err)
			} else {
				fmt.Println("[server]overload body", string(body))
				fmt.Fprintln(w, "[server]overload body", string(body))
			}

		} else if strings.Contains(contentType, "application/x-www-form-urlencoded") {
			if r.Body == nil {
				fmt.Println("not carry body")
			}
			// body, err := ioutil.ReadAll(r.Body)
			// fmt.Println("[server]overload body", string(body), "error", err)

			if err := r.ParseForm(); err != nil {
				fmt.Println("[server]parse form error:", err)
				fmt.Fprintln(w, "[server]parse form error:", err)
			} else {
				fmt.Println(w, "[server]overload form", r.Form)
				fmt.Fprintln(w, "[server]overload form", r.Form)
				if data, err := TransformToJson(r.Form); err != nil {
					fmt.Fprintln(w, "[server]transform to json error", err)
				} else {
					fmt.Println(w, "[server]transform to json", string(data))
					fmt.Fprintln(w, "[server]transform to json", string(data))
				}
			}
		} else {
			fmt.Fprintln(w, "[server]not support content-type", contentType)
		}
	})
	log.Fatalln(http.ListenAndServe(":8081", nil))
}

func NewProxy(w http.ResponseWriter, r *http.Request, upstream string) (code int) {
	ctx := ubase.NamedContext("proxy")
	log := ctx.Logger()
	log.Infof("forward request from[RemoteAddr:%s, Host:%s] to upstream[%s]", r.RemoteAddr, r.Host, upstream)
	proxy := &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			r.URL.Scheme = "http"
			r.URL.Host = upstream
			r.Header.Add("X-Forwarded-For", r.RemoteAddr)
			r.Header.Add("X-Origin-For", r.URL.Host)
		},
		ModifyResponse: func(resp *http.Response) error {
			code = resp.StatusCode
			return nil
		},
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			u := uuid.New().String()
			w.Header().Set("X-Session-Id", u)
			// code = http.StatusBadGateway
			// w.WriteHeader(code)
		},
	}
	var transport http.RoundTripper = &http.Transport{
		Dial: func(network, addr string) (net.Conn, error) {
			c, err := net.DialTimeout(network, addr, time.Second*time.Duration(10000))
			if err != nil {
				return nil, err
			}
			return c, nil
		},
		ResponseHeaderTimeout: time.Second * time.Duration(10000),
	}

	proxy.Transport = transport
	proxy.ServeHTTP(w, r)
	return
}

func TransformToJson(formValue map[string][]string) ([]byte, error) {
	for k, v := range formValue {
		fmt.Println("###key:", k, ", value:", v, "type:", reflect.ValueOf(v).Type().String())
	}
	// return json.Marshal(formValue)
	var kvMap = make(map[string]interface{})
	for k, v := range formValue {
		if len(v) == 0 {
			continue
		}
		// if e, ok := kvMap[k]; ok {
		// 	val := reflect.ValueOf(e)
		// 	if val.Type().Kind() == reflect.String {
		// 		tmp := []string{val.String()}
		// 		tmp = append(tmp, v[0])
		// 		kvMap[k] = tmp
		// 	} else if val.Type().Kind() == reflect.Slice {
		// 		// val.Elem().SetLen(val.Elem().Len() + 1)
		// 		// val.Elem().Index(val.Elem().Len()).Set(reflect.ValueOf(v[0]))
		// 	}
		// } else
		{
			val := reflect.ValueOf(v)
			if val.Len() == 1 {
				kvMap[k] = v[0]
			} else {
				kvMap[k] = v
			}
		}
	}
	fmt.Println(strings.Repeat("-", 60))
	for k, v := range kvMap {
		fmt.Println("debug: key -> ", k, ", value -> ", v, ", type -> ", reflect.TypeOf(v).String())
	}
	fmt.Println(strings.Repeat("-", 60))
	return json.Marshal(kvMap)
}

type textHandler struct {
	response string
}

func (textHdl *textHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("textHandler request: ", r)
	fmt.Println()
	fmt.Printf("request: %v\n", r)

	fmt.Fprintf(w, textHdl.response)
}

type indexHandler struct{}

func (indexHdl *indexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("indexHandler request: ", r)
	fmt.Println()
	fmt.Printf("request: %v\n", r)
	w.Header().Set("Content-Type", "text/html")

	html := `<doctype html>
<html>
<head>
	<title>Hello world</title>
</head>
<body>
<p>
	<a href="/welcome">Welcome</a> | <a href="/message">Message</a>
</p>
</body>
</html>`

	fmt.Fprintln(w, html)
}

func muxDemo() {
	mux := http.NewServeMux()
	mux.Handle("/", &indexHandler{})
	textHdl := &textHandler{"Text Handler!"}
	mux.Handle("/text", textHdl)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalln("encounter error: ", err)
	}
}
func main() {
	httpDemo()
	httpProxy()
}
