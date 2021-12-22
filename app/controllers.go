// Code generated by goagen v1.3.1, DO NOT EDIT.
//
// API "jwt-signin": Application Controllers
//
// Command:
// $ goagen
// --design=github.com/Microkubes/jwt-issuer/design
// --out=$(GOPATH)/src/github.com/Microkubes/jwt-issuer
// --version=v1.3.1

package app

import (
	"context"
	"github.com/keitaroinc/goa"
	"github.com/keitaroinc/goa/cors"
	"github.com/keitaroinc/goa/encoding/form"
	"net/http"
)

// initService sets up the service encoders, decoders and mux.
func initService(service *goa.Service) {
	// Setup encoders and decoders
	service.Encoder.Register(goa.NewJSONEncoder, "application/json")
	service.Encoder.Register(goa.NewGobEncoder, "application/gob", "application/x-gob")
	service.Encoder.Register(goa.NewXMLEncoder, "application/xml")
	service.Decoder.Register(form.NewDecoder, "application/x-www-form-urlencoded")

	// Setup default encoder and decoder
	service.Encoder.Register(goa.NewJSONEncoder, "*/*")
	service.Decoder.Register(form.NewDecoder, "*/*")
}

// JWTController is the controller interface for the JWT actions.
type JWTController interface {
	goa.Muxer
	Signin(*SigninJWTContext) error
}

// MountJWTController "mounts" a JWT resource controller on the given service.
func MountJWTController(service *goa.Service, ctrl JWTController) {
	initService(service)
	var h goa.Handler
	service.Mux.Handle("OPTIONS", "/signin", ctrl.MuxHandler("preflight", handleJWTOrigin(cors.HandlePreflight()), nil))

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewSigninJWTContext(ctx, req, service)
		if err != nil {
			return err
		}
		// Build the payload
		if rawPayload := goa.ContextRequest(ctx).Payload; rawPayload != nil {
			rctx.Payload = rawPayload.(*Credentials)
		} else {
			return goa.MissingPayloadError()
		}
		return ctrl.Signin(rctx)
	}
	h = handleJWTOrigin(h)
	service.Mux.Handle("POST", "/signin", ctrl.MuxHandler("signin", h, unmarshalSigninJWTPayload))
	service.LogInfo("mount", "ctrl", "JWT", "action", "Signin", "route", "POST /signin")
}

// handleJWTOrigin applies the CORS response headers corresponding to the origin.
func handleJWTOrigin(h goa.Handler) goa.Handler {

	return func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		origin := req.Header.Get("Origin")
		if origin == "" {
			// Not a CORS request
			return h(ctx, rw, req)
		}
		if cors.MatchOrigin(origin, "*") {
			ctx = goa.WithLogContext(ctx, "origin", origin)
			rw.Header().Set("Access-Control-Allow-Origin", origin)
			rw.Header().Set("Access-Control-Allow-Credentials", "false")
			if acrm := req.Header.Get("Access-Control-Request-Method"); acrm != "" {
				// We are handling a preflight request
				rw.Header().Set("Access-Control-Allow-Methods", "OPTIONS")
			}
			return h(ctx, rw, req)
		}

		return h(ctx, rw, req)
	}
}

// unmarshalSigninJWTPayload unmarshals the request body into the context request data Payload field.
func unmarshalSigninJWTPayload(ctx context.Context, service *goa.Service, req *http.Request) error {
	payload := &credentials{}
	if err := service.DecodeRequest(req, payload); err != nil {
		return err
	}
	goa.ContextRequest(ctx).Payload = payload.Publicize()
	return nil
}
