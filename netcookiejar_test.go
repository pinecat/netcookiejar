package netcookiejar

import (
	"fmt"
	"testing"
)

func TestHeader(t *testing.T) {
	ncj := New()
	cookies, err := ncj.Read(MODE_STRING, ".humblebundle.com\tFALSE\t/\tFALSE\t1621483493.456876\thbuid\tXUTE6GAA4XN8T")
	if err != nil {
		t.Errorf("Couldn't read in cookie string: %s\n", err.Error())
	}
	fmt.Printf("Cookie String: %s\n", cookies[0].Header(&HeaderOptions{
		Secure:   true,
		HttpOnly: true,
	}))
}
