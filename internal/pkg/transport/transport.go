//go:build !js || !wasm

package transport

import (
	"net/http"
)

type Transport http.Transport
