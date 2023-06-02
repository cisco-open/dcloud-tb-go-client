// Copyright 2023 Cisco Systems, Inc. and its affiliates
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

package tbclient

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"strings"
)

type ClientError struct {
	VndErrors []vndError
	Status    string
}

func (e *ClientError) Error() string {

	messages := make([]string, len(e.VndErrors))
	var logref string

	for i, e := range e.VndErrors {
		messages[i] = e.Message
		logref = e.Logref // Will be the same for every vndError in a given response
	}

	var formattedError = ""
	if len(messages) != 0 {
		formattedError = fmt.Sprintf(" (logref: %s):\n%s", logref, strings.Join(messages, "\n"))
	}

	return fmt.Sprintf("Error processing requst [Server Response: %s]%s", e.Status, formattedError)
}

func generateError(response resty.Response) *ClientError {
	return &ClientError{
		VndErrors: *response.Error().(*[]vndError),
		Status:    response.Status(),
	}
}
