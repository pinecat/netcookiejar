package netcookiejar

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	MODE_STRING = 0
	MODE_FILE   = 1
)

// CookieJar Struct holds an array of http cookies, and can read in (parse), as well as write.
type NetCookieJar struct {
}

// NetCookie Special type which implements http.Cookie (since http.Cookie doesn't have InclSub member).
type NetCookie struct {
	InclSub bool
	Cookie  *http.Cookie
}

type HeaderOptions struct {
	Secure   bool
	HttpOnly bool
}

// New The default constructor for NetCookieJar.
func New() *NetCookieJar {
	return &NetCookieJar{}
}

// Read attempts to read in cookies from a string or file of the Netscape format (https://curl.se/docs/http-cookies.html).
// Returns the slice of cookies and nil on success, nil and error on failure.
func (c NetCookieJar) Read(mode int, data string) ([]*NetCookie, error) {
	switch mode {
	case MODE_STRING:
		return c.readString(data)
	case MODE_FILE:
		return c.readFile(data)
	}
	return nil, errors.New("invalid mode, please use 'netcookiejar.MODE_STRING' or 'netcookiejar.MODE_FILE'")
}

// Write attempts to write Cookies []*http.Cookie to the specified file in Netscape format (https://curl.se/docs/http-cookies.html).
// Returns the written string and nil on success, returns empty string and error on failure.
func (c NetCookieJar) Write(mode int, path string, cookies []*NetCookie) (string, error) {
	switch mode {
	case MODE_STRING:
		return c.writeString(cookies)
	case MODE_FILE:
		return c.writeFile(path, cookies)
	}
	return "", nil
}

// Header builds a string from the cookie that can be used in an http request.
func (c NetCookie) Header(options *HeaderOptions) string {
	var resp string

	if options == nil {
		options = &HeaderOptions{
			Secure:   false,
			HttpOnly: false,
		}
	}

	resp = fmt.Sprintf(
		"%s=%s; Expires=%s; Domain=%s; Path=%s;",
		c.Cookie.Name,
		c.Cookie.Value,
		c.Cookie.Expires.Format("Mon, 02 Jan 2006 15:04:05 GMT"),
		c.Cookie.Domain,
		c.Cookie.Path,
	)

	if options.Secure {
		resp += " Secure;"
	}

	if options.HttpOnly {
		resp += " HttpOnly;"
	}

	if resp[len(resp)-1] == ';' {
		resp = resp[:len(resp)-1]
	}

	return resp
}

// Creates a new scanner on the string, then passes it to the parsing function.
func (c NetCookieJar) readString(data string) ([]*NetCookie, error) {
	// Create a new scanner on the data string
	r := bufio.NewScanner(strings.NewReader(data))
	return c.parse(r)
}

// readFile Attempts to open the specified file, then create a new scanner on it, before passing it to the parsing function.
func (c NetCookieJar) readFile(data string) ([]*NetCookie, error) {
	// Try opening the file
	f, err := os.OpenFile(data, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Create a scanner from the file
	r := bufio.NewScanner(f)

	// Try to parse the data
	return c.parse(r)
}

// parse Parses a file/string with Netscape formatted cookies to an array of NetCookies ([]*NetCookie)
func (c NetCookieJar) parse(s *bufio.Scanner) ([]*NetCookie, error) {
	var cookies []*NetCookie
	for s.Scan() {
		fields := strings.Split(s.Text(), "\t")
		if len(fields) < 7 {
			return nil, errors.New("invalid format for cookie file, cookie missing a field")
		}

		// Get domain
		domain := fields[0]

		// Get include-subdomains
		inclSub, err := strconv.ParseBool(fields[1])
		if err != nil {
			return nil, err
		}

		// Get path
		path := fields[2]

		// Get secure
		secure, err := strconv.ParseBool(fields[3])
		if err != nil {
			return nil, err
		}

		// Get expiration date
		expFloat, err := strconv.ParseFloat(fields[4], 64)
		if err != nil {
			expFloat = 0
		}
		sec, dec := math.Modf(expFloat)
		exp := time.Unix(int64(sec), int64(dec*(1e9)))

		// Get name
		name := fields[5]

		// Get value
		value := fields[6]

		// Create cookie
		cookie := &NetCookie{
			InclSub: inclSub,
			Cookie: &http.Cookie{
				Domain:  domain,
				Path:    path,
				Secure:  secure,
				Expires: exp,
				Name:    name,
				Value:   value,
			},
		}

		// Add cookie to slice
		cookies = append(cookies, cookie)
	}
	return cookies, nil
}

// writeString Generates a Netscape formatted cookie string from a slice of NetCookies.
func (c NetCookieJar) writeString(cookies []*NetCookie) (string, error) {
	// Make string to concatenate onto
	var netscape string = ""

	// Loop through the slice of cookies
	for i := 0; i < len(cookies); i++ {
		// Shorter variable name, so easier to read
		c := cookies[i]

		// Have to do some extra work for the boolean values
		inclSub := "FALSE"
		if c.InclSub {
			inclSub = "TRUE"
		}
		secure := "FALSE"
		if c.Cookie.Secure {
			secure = "TRUE"
		}

		// Determine whether or not to tack on newline
		var nl string = "\n"
		if i == len(cookies)-1 {
			nl = ""
		}

		// Concatenate properly formatted Netscape cookie onto our string
		netscape += fmt.Sprintf("%s\t%s\t%s\t%s\t%d\t%s\t%s%s", c.Cookie.Domain, inclSub, c.Cookie.Path, secure, c.Cookie.Expires.Unix(), c.Cookie.Name, c.Cookie.Value, nl)
	}

	// Return the string
	return netscape, nil
}

// writeFile Gets a Netscape formatted cookie string which it returns, and also writes it to the specified file.
func (c NetCookieJar) writeFile(path string, cookies []*NetCookie) (string, error) {
	// Get netscape formatted string of cookies for writing (writeString() should never return an error)
	netscape, _ := c.writeString(cookies)

	// Try to open/create the file at the path specified for writing
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return "", err
	}
	defer f.Close()

	// Try writing the string to the file
	_, err = f.WriteString(netscape)
	if err != nil {
		return "", err
	}

	// Return the Netscape formatted string
	return netscape, nil
}
