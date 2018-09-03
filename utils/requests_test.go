package utils

import (
	"fmt"
	"goBoss/config"
	"testing"
)

func TestHttpGet(t *testing.T) {
	url := fmt.Sprintf("%sLATEST_RELEASE", config.Config.DriverUrl)
	ops := RequestData{
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
	r, err := HttpGet(url, ops)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(r.String())
}
