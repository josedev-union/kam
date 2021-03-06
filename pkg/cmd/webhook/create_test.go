package webhook

import (
	"bytes"
	"fmt"
	"regexp"
	"testing"

	"github.com/spf13/cobra"
)

type keyValuePair struct {
	key   string
	value string
}

func TestValidateForCreate(t *testing.T) {
	testcases := []struct {
		options *createOptions
		errMsg  string
	}{
		{
			&createOptions{
				options{isCICD: true, serviceName: "foo"},
			},
			"Only one of 'cicd' or 'env-name/service-name' can be specified",
		},
		{
			&createOptions{
				options{isCICD: true, envName: "foo"},
			},
			"Only one of 'cicd' or 'env-name/service-name' can be specified",
		},
		{
			&createOptions{
				options{isCICD: true, envName: "foo", serviceName: "bar"},
			},
			"Only one of 'cicd' or 'env-name/service-name' can be specified",
		},
		{
			&createOptions{
				options{isCICD: false},
			},
			"One of 'cicd' or 'env-name/service-name' must be specified",
		},
		{
			&createOptions{
				options{isCICD: false, serviceName: "foo"},
			},
			"One of 'cicd' or 'env-name/service-name' must be specified",
		},
		{
			&createOptions{
				options{isCICD: false, serviceName: "foo", envName: "gau"},
			},
			"",
		},
		{
			&createOptions{
				options{isCICD: true, serviceName: ""},
			},
			"",
		},
	}

	for i, tt := range testcases {
		t.Run(fmt.Sprintf("Test %d", i), func(t *testing.T) {
			err := tt.options.Validate()
			if err != nil && tt.errMsg == "" {
				t.Errorf("Validate() got an unexpected error: %s", err)
			} else {
				if !matchError(t, tt.errMsg, err) {
					t.Errorf("Validate() failed to match error: got %s, want %s", err, tt.errMsg)
				}
			}
		})
	}
}

func executeCommand(cmd *cobra.Command, flags ...keyValuePair) (output string, err error) {
	buf := new(bytes.Buffer)
	cmd.SetOutput(buf)
	for _, flag := range flags {
		err := cmd.Flags().Set(flag.key, flag.value)
		if err != nil {
			return "", err
		}
	}
	_, err = cmd.ExecuteC()
	return buf.String(), err
}

func flag(k, v string) keyValuePair {
	return keyValuePair{
		key:   k,
		value: v,
	}
}

func matchError(t *testing.T, s string, e error) bool {
	t.Helper()
	if s == "" && e == nil {
		return true
	}
	if s != "" && e == nil {
		return false
	}
	match, err := regexp.MatchString(s, e.Error())
	if err != nil {
		t.Fatal(err)
	}
	return match
}
