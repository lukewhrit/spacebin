/*
* Copyright 2020-2023 Luke Whritenour

* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at

*     http://www.apache.org/licenses/LICENSE-2.0

* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package util_test

import (
	"html/template"
	"testing"

	"github.com/orca-group/spirit/internal/util"
	"github.com/stretchr/testify/require"
)

func TestValidateBody(t *testing.T) {
	require.NoError(t, util.ValidateBody(100, util.CreateRequest{
		Content: "Test",
	}))

	require.Error(t, util.ValidateBody(2, util.CreateRequest{
		Content: "Test",
	}))

	require.Error(t, util.ValidateBody(2, util.CreateRequest{
		Content: "",
	}))
}

func TestCountLines(t *testing.T) {
	content := "Line 1\nLine 2"

	lines := util.CountLines(content)

	require.Equal(t, lines, template.HTML("<div>1</div><div>2</div>"))
}
