package licensei

import (
	"testing"
)

func TestHeaderChecker_Check(t *testing.T) {
	template := `// Copyright © :YEAR: Márk Sági-Kazár
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
// limitations under the License.`

	checker := HeaderChecker{
		IgnorePaths: []string{"path"},
		IgnoreFiles: []string{"*_gen.go", "*_test.go"},
	}

	violations, err := checker.Check("testdata/header", template)
	if err != nil {
		t.Fatal(err)
	}

	for path, violation := range violations {
		t.Errorf("%s: %s", violation, path)
	}
}

func TestHeaderChecker_Author(t *testing.T) {
	template := `// Copyright © :YEAR: :AUTHOR:`

	testdata := []struct {
		valid   bool
		authors [][]string
	}{
		{
			valid: true,
			authors: [][]string{
				{"Márk Sági-Kazár"},
				{"Márk"},
				{"Márk", "Jozsi"},
			},
		},
		{
			valid: false,
			authors: [][]string{
				{"Már"},
				{"Márk Sági-K"},
				{"Jozsi"},
			},
		},
	}

	for _, td := range testdata {
		for _, a := range td.authors {
			checker := HeaderChecker{
				IgnorePaths: []string{"path"},
				IgnoreFiles: []string{"*_gen.go", "*_test.go"},
				Authors:     a,
			}

			violations, err := checker.Check("testdata/header", template)
			if err != nil {
				t.Fatal(err)
			}

			if td.valid && len(violations) > 0 {
				for path, violation := range violations {
					t.Errorf("%s: %s", violation, path)
				}
			} else if !td.valid && len(violations) == 0 {
				t.Errorf("expected error for authors %+v", a)
			}
		}
	}
}
