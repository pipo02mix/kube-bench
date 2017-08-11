// Copyright © 2017 Aqua Security Software Ltd. <info@aquasec.com>
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

package cmd

import (
	"regexp"
	"testing"
)

func TestCheckVersion(t *testing.T) {
	kubeoutput := `Client Version: version.Info{Major:"1", Minor:"7", GitVersion:"v1.7.0", GitCommit:"d3ada0119e776222f11ec7945e6d860061339aad", GitTreeState:"clean", BuildDate:"2017-06-30T09:51:01Z", GoVersion:"go1.8.3", Compiler:"gc", Platform:"darwin/amd64"}
	Server Version: version.Info{Major:"1", Minor:"7", GitVersion:"v1.7.0", GitCommit:"d3ada0119e776222f11ec7945e6d860061339aad", GitTreeState:"clean", BuildDate:"2017-07-26T00:12:31Z", GoVersion:"go1.8.3", Compiler:"gc", Platform:"linux/amd64"}`
	cases := []struct {
		t     string
		s     string
		major string
		minor string
		exp   string
	}{
		{t: "Client", s: kubeoutput, major: "1", minor: "7"},
		{t: "Server", s: kubeoutput, major: "1", minor: "7"},
		{t: "Client", s: kubeoutput, major: "1", minor: "6", exp: "Unexpected Client version 1.7"},
		{t: "Client", s: kubeoutput, major: "2", minor: "0", exp: "Unexpected Client version 1.7"},
		{t: "Server", s: "something unexpected", major: "2", minor: "0", exp: "Couldn't find Server version from kubectl output 'something unexpected'"},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			m := checkVersion(c.t, c.s, c.major, c.minor)
			if m != c.exp {
				t.Fatalf("Got: %s, expected: %s", m, c.exp)
			}
		})
	}

}

func TestVersionMatch(t *testing.T) {
	minor := regexVersionMinor
	major := regexVersionMajor
	client := `Client Version: version.Info{Major:"1", Minor:"7", GitVersion:"v1.7.0", GitCommit:"d3ada0119e776222f11ec7945e6d860061339aad", GitTreeState:"clean", BuildDate:"2017-06-30T09:51:01Z", GoVersion:"go1.8.3", Compiler:"gc", Platform:"darwin/amd64"}`
	server := `Server Version: version.Info{Major:"1", Minor:"7", GitVersion:"v1.7.0", GitCommit:"d3ada0119e776222f11ec7945e6d860061339aad", GitTreeState:"clean", BuildDate:"2017-07-26T00:12:31Z", GoVersion:"go1.8.3", Compiler:"gc", Platform:"linux/amd64"}`

	cases := []struct {
		r   *regexp.Regexp
		s   string
		exp string
	}{
		{r: major, s: server, exp: "1"},
		{r: minor, s: server, exp: "7"},
		{r: major, s: client, exp: "1"},
		{r: minor, s: client, exp: "7"},
		{r: major, s: "Some unexpected string"},
		{r: minor}, // Checking that we don't fall over if the string is empty
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			m := versionMatch(c.r, c.s)
			if m != c.exp {
				t.Fatalf("Got %s expected %s", m, c.exp)
			}
		})
	}
}
