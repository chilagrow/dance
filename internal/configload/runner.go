// Copyright 2021 FerretDB Inc.
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

package configload

import "github.com/FerretDB/dance/internal/config"

// runnerParams is common interface for runner parameters.
//
//sumtype:decl
type runnerParams interface {
	convert() config.RunnerParams // seal for sumtype
}

// runnerParamsCommand represents `command` runner parameters in the project configuration YAML file.
type runnerParamsCommand struct {
	Dir   string `yaml:"dir"`
	Setup string `yaml:"setup"`
	Tests []struct {
		Name string `yaml:"name"`
		Cmd  string `yaml:"cmd"`
	} `yaml:"tests"`
}

// convert implements [runnerParams] interface.
func (rp *runnerParamsCommand) convert() config.RunnerParams {
	res := &config.RunnerParamsCommand{
		Dir:   rp.Dir,
		Setup: rp.Setup,
	}

	for _, t := range rp.Tests {
		res.Tests = append(res.Tests, config.RunnerParamsCommandTest{
			Name: t.Name,
			Cmd:  t.Cmd,
		})
	}

	return res
}

// check interfaces
var (
	_ runnerParams = (*runnerParamsCommand)(nil)
)
