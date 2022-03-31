// Code generated by go-swagger; DO NOT EDIT.

//
// Copyright 2021 The Sigstore Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package tlog

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

// NewGetInactiveLogInfoParams creates a new GetInactiveLogInfoParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewGetInactiveLogInfoParams() *GetInactiveLogInfoParams {
	return &GetInactiveLogInfoParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewGetInactiveLogInfoParamsWithTimeout creates a new GetInactiveLogInfoParams object
// with the ability to set a timeout on a request.
func NewGetInactiveLogInfoParamsWithTimeout(timeout time.Duration) *GetInactiveLogInfoParams {
	return &GetInactiveLogInfoParams{
		timeout: timeout,
	}
}

// NewGetInactiveLogInfoParamsWithContext creates a new GetInactiveLogInfoParams object
// with the ability to set a context for a request.
func NewGetInactiveLogInfoParamsWithContext(ctx context.Context) *GetInactiveLogInfoParams {
	return &GetInactiveLogInfoParams{
		Context: ctx,
	}
}

// NewGetInactiveLogInfoParamsWithHTTPClient creates a new GetInactiveLogInfoParams object
// with the ability to set a custom HTTPClient for a request.
func NewGetInactiveLogInfoParamsWithHTTPClient(client *http.Client) *GetInactiveLogInfoParams {
	return &GetInactiveLogInfoParams{
		HTTPClient: client,
	}
}

/* GetInactiveLogInfoParams contains all the parameters to send to the API endpoint
   for the get inactive log info operation.

   Typically these are written to a http.Request.
*/
type GetInactiveLogInfoParams struct {
	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the get inactive log info params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetInactiveLogInfoParams) WithDefaults() *GetInactiveLogInfoParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the get inactive log info params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetInactiveLogInfoParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the get inactive log info params
func (o *GetInactiveLogInfoParams) WithTimeout(timeout time.Duration) *GetInactiveLogInfoParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get inactive log info params
func (o *GetInactiveLogInfoParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get inactive log info params
func (o *GetInactiveLogInfoParams) WithContext(ctx context.Context) *GetInactiveLogInfoParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get inactive log info params
func (o *GetInactiveLogInfoParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get inactive log info params
func (o *GetInactiveLogInfoParams) WithHTTPClient(client *http.Client) *GetInactiveLogInfoParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get inactive log info params
func (o *GetInactiveLogInfoParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WriteToRequest writes these params to a swagger request
func (o *GetInactiveLogInfoParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}