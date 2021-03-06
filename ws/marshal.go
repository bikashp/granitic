// Copyright 2016 Granitic. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be found in the LICENSE file at the root of this project.

package ws

import (
	"context"
	"errors"
	"github.com/graniticio/granitic/httpendpoint"
	"github.com/graniticio/granitic/logging"
	"net/http"
)

// Implemented by components that can convert the supplied data into a form suitable for serialisation and
// write that serialised form to the HTTP output stream.
type MarshalingWriter interface {

	// MarshalAndWrite converts the data to some serialisable form (JSON, XML etc) and writes it to the HTTP output stream.
	MarshalAndWrite(data interface{}, w http.ResponseWriter) error
}

// A response writer that uses automatic marshalling of structs to serialisable forms rather than using templates.
type MarshallingResponseWriter struct {
	// Injected automatically
	FrameworkLogger logging.Logger

	// Component able to calculate the HTTP status code that should be written to the HTTP response.
	StatusDeterminer HttpStatusCodeDeterminer

	// Component able to generate errors if a problem is encountered during marshalling.
	FrameworkErrors *FrameworkErrorGenerator

	// The common and static set of headers that should be written to all responses.
	DefaultHeaders map[string]string

	// Component able to wrap response data in a standardised structure.
	ResponseWrapper ResponseWrapper

	// Component able to dynamically generate additional headers to be written to the response.
	HeaderBuilder WsCommonResponseHeaderBuilder

	// Component able to format services errors in an application specific manner.
	ErrorFormatter ErrorFormatter

	// Component able to serialize the data to the HTTP output stream.
	MarshalingWriter MarshalingWriter
}

// See WsResponseWriter.Write
func (rw *MarshallingResponseWriter) Write(ctx context.Context, state *WsProcessState, outcome WsOutcome) error {

	var ch map[string]string

	if rw.HeaderBuilder != nil {
		ch = rw.HeaderBuilder.BuildHeaders(ctx, state)
	}

	switch outcome {
	case Normal:
		return rw.write(ctx, state.WsResponse, state.HttpResponseWriter, ch)
	case Error:
		return rw.writeErrors(ctx, state.ServiceErrors, state.HttpResponseWriter, ch)
	case Abnormal:
		return rw.writeAbnormalStatus(ctx, state.Status, state.HttpResponseWriter, ch)
	}

	return errors.New("Unsuported WsOutcome value")
}

func (rw *MarshallingResponseWriter) write(ctx context.Context, res *WsResponse, w *httpendpoint.HttpResponseWriter, ch map[string]string) error {

	if w.DataSent {
		//This HTTP response has already been written to by another component - not safe to continue
		if rw.FrameworkLogger.IsLevelEnabled(logging.Debug) {
			rw.FrameworkLogger.LogDebugfCtx(ctx, "Response already written to.")
		}

		return nil
	}

	headers := MergeHeaders(res, ch, rw.DefaultHeaders)
	WriteHeaders(w, headers)

	s := rw.StatusDeterminer.DetermineCode(res)
	w.WriteHeader(s)

	e := res.Errors

	if res.Body == nil && !e.HasErrors() {
		return nil
	}

	ef := rw.ErrorFormatter
	wrap := rw.ResponseWrapper

	fe := ef.FormatErrors(e)
	wrapper := wrap.WrapResponse(res.Body, fe)

	return rw.MarshalingWriter.MarshalAndWrite(wrapper, w)
}

// See AbnormalStatusWriter.WriteAbnormalStatus
func (rw *MarshallingResponseWriter) WriteAbnormalStatus(ctx context.Context, state *WsProcessState) error {
	return rw.Write(ctx, state, Abnormal)
}

func (rw *MarshallingResponseWriter) writeAbnormalStatus(ctx context.Context, status int, w *httpendpoint.HttpResponseWriter, ch map[string]string) error {

	res := new(WsResponse)
	res.HttpStatus = status
	var errors ServiceErrors

	e := rw.FrameworkErrors.HttpError(status)
	errors.AddError(e)

	res.Errors = &errors

	return rw.write(ctx, res, w, ch)

}

func (rw *MarshallingResponseWriter) writeErrors(ctx context.Context, errors *ServiceErrors, w *httpendpoint.HttpResponseWriter, ch map[string]string) error {

	res := new(WsResponse)
	res.Errors = errors

	return rw.write(ctx, res, w, ch)
}
