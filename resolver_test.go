package main

import (
	"fmt"
	"testing"

	"github.com/bancodobrasil/featws-resolver-adapter-go/types"
	"github.com/stretchr/testify/assert"
)

type InputResolverTest struct {
	ipStr string
	input types.ResolveInput
}

func TestValidIPs(t *testing.T) {

	var params []InputResolverTest

	params = appendResult(params, "81.2.69.142")
	params = appendResult(params, "81.2.69.143")

	for _, param := range params {
		testName := fmt.Sprintf("%s", param.ipStr)
		t.Run(testName, func(t *testing.T) {
			output := types.ResolveOutput{
				Context: make(map[string]interface{}),
				Errors:  make(map[string]interface{}),
			}
			resolver(param.input, &output)
			if len(output.Errors) > 0 {
				t.Errorf("got error %v", output.Errors)
			} else {
				geoRecord := output.Context["geoip"].(*GeoRecord)
				assert.Equal(t, param.ipStr, geoRecord.RemoteIP)
				assert.Equal(t, "United Kingdom", geoRecord.Country)
				assert.Equal(t, "London", geoRecord.City)
				assert.Equal(t, float64(51.5142), geoRecord.Location.Latitude)
				assert.Equal(t, float64(-0.0931), geoRecord.Location.Longitude)
				assert.Equal(t, uint16(10), geoRecord.Location.AccuracyRadius)
			}
		})
	}
}

func TestInvalidIPs(t *testing.T) {

	var params []InputResolverTest

	params = appendResult(params, "2.69.142")
	params = appendResult(params, "381.2.69.143")
	params = appendResult(params, "")

	for _, param := range params {
		testName := fmt.Sprintf("%s", param.ipStr)
		t.Run(testName, func(t *testing.T) {
			output := types.ResolveOutput{
				Context: make(map[string]interface{}),
				Errors:  make(map[string]interface{}),
			}
			resolver(param.input, &output)
			if len(output.Errors) > 0 {
				assert.NotEmpty(t, output.Errors["geoip"])
			} else {
				t.Errorf("Resolved invalid context")
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

func BenchmarkResolver(b *testing.B) {
	input := types.ResolveInput{
		Context: map[string]interface{}{
			"remote_ip": "81.2.69.142",
		},
		Load: []string{"geoip"},
	}
	output := types.ResolveOutput{
		Context: make(map[string]interface{}),
		Errors:  make(map[string]interface{}),
	}
	for i := 0; i < b.N; i++ {
		resolver(input, &output)
	}
}
