// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Copied from: https://github.com/golang/go/blob/go1.18.3/src/net/http/roundtrip_js.go

//go:build js && wasm

package transport

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"syscall/js"
)

var uint8Array = js.Global().Get("Uint8Array")

// jsFetchMode is a Request.Header map key that, if present,
// signals that the map entry is actually an option to the Fetch API mode setting.
// Valid values are: "cors", "no-cors", "same-origin", "navigate"
// The default is "same-origin".
//
// Reference: https://developer.mozilla.org/en-US/docs/Web/API/WindowOrWorkerGlobalScope/fetch#Parameters
const jsFetchMode = "js.fetch:mode"

// jsFetchRedirect is a Request.Header map key that, if present,
// signals that the map entry is actually an option to the Fetch API redirect setting.
// Valid values are: "follow", "error", "manual"
// The default is "follow".
//
// Reference: https://developer.mozilla.org/en-US/docs/Web/API/WindowOrWorkerGlobalScope/fetch#Parameters
const jsFetchRedirect = "js.fetch:redirect"

// jsFetchMissing will be true if the Fetch API is not present in
// the browser globals.
var jsFetchMissing = js.Global().Get("fetch").IsUndefined()

// RoundTrip implements the RoundTripper interface using the WHATWG Fetch API.
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	ac := js.Global().Get("AbortController")
	if !ac.IsUndefined() {
		// Some browsers that support WASM don't necessarily support
		// the AbortController. See
		// https://developer.mozilla.org/en-US/docs/Web/API/AbortController#Browser_compatibility.
		ac = ac.New()
	}

	opt := js.Global().Get("Object").New()
	// See https://developer.mozilla.org/en-US/docs/Web/API/WindowOrWorkerGlobalScope/fetch
	// for options available.
	opt.Set("method", req.Method)
	if h := req.Header.Get(jsFetchMode); h != "" {
		opt.Set("mode", h)
		req.Header.Del(jsFetchMode)
	}
	if h := req.Header.Get(jsFetchRedirect); h != "" {
		opt.Set("redirect", h)
		req.Header.Del(jsFetchRedirect)
	}
	if !ac.IsUndefined() {
		opt.Set("signal", ac.Get("signal"))
	}
	headers := js.Global().Get("Headers").New()
	for key, values := range req.Header {
		for _, value := range values {
			headers.Call("append", key, value)
		}
	}
	opt.Set("headers", headers)

	if req.Body != nil {
		// TODO(johanbrandhorst): Stream request body when possible.
		// See https://bugs.chromium.org/p/chromium/issues/detail?id=688906 for Blink issue.
		// See https://bugzilla.mozilla.org/show_bug.cgi?id=1387483 for Firefox issue.
		// See https://github.com/web-platform-tests/wpt/issues/7693 for WHATWG tests issue.
		// See https://developer.mozilla.org/en-US/docs/Web/API/Streams_API for more details on the Streams API
		// and browser support.
		body, err := io.ReadAll(req.Body)
		if err != nil {
			req.Body.Close() // RoundTrip must always close the body, including on errors.
			return nil, err
		}
		req.Body.Close()
		if len(body) != 0 {
			buf := uint8Array.New(len(body))
			js.CopyBytesToJS(buf, body)
			opt.Set("body", buf)
		}
	}

	fetchPromise := js.Global().Call("fetch", req.URL.String(), opt)
	var (
		respCh           = make(chan *http.Response, 1)
		errCh            = make(chan error, 1)
		success, failure js.Func
	)
	success = js.FuncOf(func(this js.Value, args []js.Value) any {
		success.Release()
		failure.Release()

		result := args[0]
		header := http.Header{}
		// https://developer.mozilla.org/en-US/docs/Web/API/Headers/entries
		headersIt := result.Get("headers").Call("entries")
		for {
			n := headersIt.Call("next")
			if n.Get("done").Bool() {
				break
			}
			pair := n.Get("value")
			key, value := pair.Index(0).String(), pair.Index(1).String()
			ck := http.CanonicalHeaderKey(key)
			header[ck] = append(header[ck], value)
		}

		contentLength := int64(0)
		clHeader := header.Get("Content-Length")
		switch {
		case clHeader != "":
			cl, err := strconv.ParseInt(clHeader, 10, 64)
			if err != nil {
				errCh <- fmt.Errorf("net/http: ill-formed Content-Length header: %v", err)
				return nil
			}
			if cl < 0 {
				// Content-Length values less than 0 are invalid.
				// See: https://datatracker.ietf.org/doc/html/rfc2616/#section-14.13
				errCh <- fmt.Errorf("net/http: invalid Content-Length header: %q", clHeader)
				return nil
			}
			contentLength = cl
		default:
			// If the response length is not declared, set it to -1.
			contentLength = -1
		}

		b := result.Get("body")
		var body io.ReadCloser
		// The body is undefined when the browser does not support streaming response bodies (Firefox),
		// and null in certain error cases, i.e. when the request is blocked because of CORS settings.
		if !b.IsUndefined() && !b.IsNull() {
			body = &streamReader{stream: b.Call("getReader")}
		} else {
			// Fall back to using ArrayBuffer
			// https://developer.mozilla.org/en-US/docs/Web/API/Body/arrayBuffer
			body = &arrayReader{arrayPromise: result.Call("arrayBuffer")}
		}

		code := result.Get("status").Int()
		respCh <- &http.Response{
			Status:        fmt.Sprintf("%d %s", code, http.StatusText(code)),
			StatusCode:    code,
			Header:        header,
			ContentLength: contentLength,
			Body:          body,
			Request:       req,
		}

		return nil
	})
	failure = js.FuncOf(func(this js.Value, args []js.Value) any {
		success.Release()
		failure.Release()
		errCh <- fmt.Errorf("net/http: fetch() failed: %s", args[0].Get("message").String())
		return nil
	})

	fetchPromise.Call("then", success, failure)
	select {
	case <-req.Context().Done():
		if !ac.IsUndefined() {
			// Abort the Fetch request.
			ac.Call("abort")
		}
		return nil, req.Context().Err()
	case resp := <-respCh:
		return resp, nil
	case err := <-errCh:
		return nil, err
	}
}
