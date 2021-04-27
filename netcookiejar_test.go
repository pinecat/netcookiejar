package netcookiejar

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

// Create some values for testing cookies
const (
	domain     string = ".github.com"
	inclSub    bool   = false
	inclSubStr string = "FALSE"
	path       string = "/"
	secure     bool   = true
	secureStr  string = "TRUE"
	expStr     string = "1462299218"
	name       string = "wom"
	value      string = "bat"
)

// Go doesn't support const structs
var (
	exp time.Time = time.Unix(int64(1462299218), int64(0))
)

// TestNew Ensures netcookiejar.New() is returning the correct value.
func TestNew(t *testing.T) {
	// Get the type of &NetCookieJar to compare with the type we get from New()
	ncjType := reflect.TypeOf(&NetCookieJar{})
	newType := reflect.TypeOf(New())

	// Compare types, error if the types are different
	if ncjType != newType {
		t.Errorf("Call to 'netcookiejar.New()' failed.  Expected it to return '%s', but instead returned '%s'.\n", "*netcookiejar.NetCookieJar", newType.Name())
	}
}

// TestReadString Tests netcookiejar.Read() with MODE_STRING, with a valid cookie string.
func TestReadValidString(t *testing.T) {
	// Create test cookie string
	cookieStr := fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\t%s\n", domain, inclSubStr, path, secureStr, expStr, name, value)

	// Create a new instance of NetCookieJar, and read from a cookie string
	ncj := New()
	cookies, err := ncj.Read(MODE_STRING, cookieStr)

	// If there was error, create a testing error, as there shouldn't be
	// an error in the cookie string used for testing in this test case.
	if err != nil {
		t.Errorf("Invalid cookie string used.\n")
	}

	// Now test and make sure the values match
	checkValues(t, cookies)
}

// TestReadString Tests netcookiejar.Read() with MODE_FILE, with a valid cookie string.
func TestReadValidFile(t *testing.T) {

}

// checkValues Tests the values of the first cookie in an array of NetCookies
//			   against a set of predetermined values.
func checkValues(t *testing.T, cookies []*NetCookie) {
	// Make sure we have at least 1 cookie
	if len(cookies) < 1 {
		t.Errorf("Can't find any cookies to test on.")
		return
	}

	// Easier to read
	c := cookies[0]

	// Check each value
	if c.InclSub != inclSub {
		t.Errorf("Bad 'Includ Subdomain' parsing.")
	} else if c.Cookie.Domain != domain {
		t.Errorf("Bad 'Domain' parsing.")
	} else if c.Cookie.Path != path {
		t.Errorf("Bad 'Path' parsing.")
	} else if c.Cookie.Secure != secure {
		t.Errorf("Bad 'Secure' parsing.")
	} else if c.Cookie.Expires != exp {
		t.Errorf("Bad 'Expires' parsing.")
	} else if c.Cookie.Name != name {
		t.Errorf("Bad 'Name' parsing.")
	} else if c.Cookie.Value != value {
		t.Errorf("Bad 'Value' parsing.")
	}
}
