package main

import (
	"fmt"
	"testing"

	"github.com/bancodobrasil/featws-resolver-adapter-go/types"
)

type InputResolverTest struct {
	ipStr string
	input types.ResolveInput
}

func TestResovler(t *testing.T) {

	var params []InputResolverTest

	params = appendResult(params, "189.40.76.241")
	params = appendResult(params, "170.66.1.232")

	for _, param := range params {
		testName := fmt.Sprintf("%s", param.ipStr)
		t.Run(testName, func(t *testing.T) {
			output := types.ResolveOutput{}
			resolver(param.input, &output)
			if output.Errors != nil {
				t.Errorf("got error %v", output.Errors)
			}
		})
	}

}

func appendResult(sliceIn []InputResolverTest, ipStr string) (sliceOut []InputResolverTest) {
	sliceOut = append(sliceIn, InputResolverTest{
		ipStr: ipStr,
		input: types.ResolveInput{
			Context: map[string]interface{}{
				"remote_ip": ipStr,
			},
			Load: []string{"geoip"},
		},
	})
	return
}
