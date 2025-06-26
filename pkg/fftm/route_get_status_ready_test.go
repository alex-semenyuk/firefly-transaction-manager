// Copyright © 2022 Kaleido, Inc.
//
// SPDX-License-Identifier: Apache-2.0
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

package fftm

import (
	"fmt"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/hyperledger/firefly-transaction-manager/mocks/ffcapimocks"
	"github.com/hyperledger/firefly-transaction-manager/pkg/apitypes"
	"github.com/hyperledger/firefly-transaction-manager/pkg/ffcapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetReadyStatus(t *testing.T) {
	_, m, done := newTestManagerWithMetrics(t, false)
	defer done()

	require.True(t, m.monitoringEnabled)
	require.NotNil(t, getReadiness(m).JSONOutputValue())
	url := fmt.Sprintf("http://%s", m.monitoringServer.Addr())

	mfc := m.connector.(*ffcapimocks.API)
	mfc.On("IsReady", mock.Anything).Return(&ffcapi.ReadyResponse{Ready: true}, ffcapi.ErrorReason(""), nil)

	err := m.Start()
	assert.NoError(t, err)

	var liv apitypes.LiveStatus
	res, err := resty.New().R().
		SetResult(&liv).
		Get(url + "/readyz")
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode())
}

func TestGetReadyStatusDeprecated(t *testing.T) {
	url, m, done := newTestManager(t)
	defer done()

	mfc := m.connector.(*ffcapimocks.API)
	mfc.On("IsReady", mock.Anything).Return(&ffcapi.ReadyResponse{Ready: true}, ffcapi.ErrorReason(""), nil)

	err := m.Start()
	assert.NoError(t, err)

	var liv apitypes.LiveStatus
	res, err := resty.New().R().
		SetResult(&liv).
		Get(url + "/status/ready")
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode())
}
