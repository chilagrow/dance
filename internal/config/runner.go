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

package config

// RunnerType represents the type of test runner used in the project configuration.
type RunnerType string

const (
	// RunnerTypeCommand indicates a command-line test runner.
	RunnerTypeCommand RunnerType = "command"

	// RunnerTypeGoTest indicates a Go test runner.
	RunnerTypeGoTest RunnerType = "gotest"

	// RunnerTypeJSTest indicates a JavaScript test runner.
	RunnerTypeJSTest RunnerType = "jstest"

	// RunnerTypeYCSB indicates a YCSB test runner.
	RunnerTypeYCSB RunnerType = "ycsb"
)

// RunnerParams is common interface for runner parameters.
//
//sumtype:decl
type RunnerParams interface {
	runnerParams() // seal for sumtype
}

// RunnerParamsCommand represents `command` runner parameters.
type RunnerParamsCommand struct {
	Dir   string
	Setup string
	Tests []RunnerParamsCommandTest
}

// RunnerParamsCommandTest represents a single test in `command` runner parameters.
type RunnerParamsCommandTest struct {
	Name string
	Cmd  string
}

// runnerParams implements [RunnerParams] interface.
func (rp *RunnerParamsCommand) runnerParams() {}

// check interfaces
var (
	_ RunnerParams = (*RunnerParamsCommand)(nil)
)
