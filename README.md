# NetCookieJar

Read and write cookies of/to the Netscape format.  You can read more about the Netscape format at the links below:
- [HTTP Cookies (curl.se)](https://curl.se/docs/http-cookies.html)
- [Format of cookies when using wget (StackExchange)](https://unix.stackexchange.com/questions/36531/format-of-cookies-when-using-wget/210282)

## Reading Cookies
You can read in cookies from a file...:
```go
func main() {
    ncj := netcookiejar.New()
    cookies, err := ncj.Read(netcookiejar.MODE_FILE, <filepath>)
    if err != nil {
        panic(err)
    }
    ...
}
```

...or you can read in cookies straight from a string:
```go
func main() {
    cookieStr := ".github.com\tFALSE\t/\tTRUE\t1462299218\twom\tbat"
    ncj := netcookiejar.New()
    cookies, err := ncj.Read(netcookiejar.MODE_STRING, cookieStr)
    if err != nil {
        panic(err)
    }
    ...
}
```

The function `netcookiejar.Read()` will return the type `[]*netcookiejar.NetCookie`. A `NetCookie` is descirbed below:
```go
type NetCookie struct {
    InclSub bool
    Cookie  *http.Cookie
}
```

The reason this is done, is because Go's `http.Cookie` does not have a member for `include-subdomains`, one of the values in a Netscape formatted cookie.

## Writing Cookies

You can write cookies to a file, and return a string...:
```go
func main() {
    var path string = "/home/<user>/cookies.txt"
    var cookies []*netcookiejar.NetCookie
    cookies = append(cookies, &netcookiejar.NetCookie{
	InclSub: false,
	Cookie: &http.Cookie{
		Domain:  ".github.com",
		Path:    "/",
		Secure:  true,
		Expires: time.Now(),
		Name:    "foo",
		Value:   "bar",
	},
    })
    cookies = append(cookies, &netcookiejar.NetCookie{
	InclSub: true,
	Cookie: &http.Cookie{
		Domain:  ".github.com",
		Path:    "/",
		Secure:  false,
		Expires: time.Now(),
		Name:    "wom",
		Value:   "bat",
	},
    })
    ncj := netcookiejar.New()
    cookieString, err := ncj.Write(netcookiejar.MODE_FILE, path, cookies)
    if err != nil {
	panic(err)
    }
    fmt.Println(cookieString)
    ...
}
```

...or you can just return a string:
```go
func main() {
    var cookies []*netcookiejar.NetCookie
    cookies = append(cookies, &netcookiejar.NetCookie{
        InclSub: false,
        Cookie: &http.Cookie{
            Domain:  ".github.com",
            Path:    "/",
            Secure:  true,
            Expires: time.Now(),
            Name:    "foo",
            Value:   "bar",
        },
    })
    cookies = append(cookies, &netcookiejar.NetCookie{
        InclSub: true,
        Cookie: &http.Cookie{
            Domain:  ".github.com",
            Path:    "/",
            Secure:  false,
            Expires: time.Now(),
            Name:    "wom",
            Value:   "bat",
        },
    })
    ncj := netcookiejar.New()
    cookieString, err := ncj.Write(netcookiejar.MODE_STRING, "", cookies)
    if err != nil {
        panic(err)
    }
    fmt.Println(cookieString)
    ...
}
```

## Quirks

Describes default behaviors of the program.

##### String Mode vs File Mode
If using `netcookiejar.MODE_FILE`, pass in the filepath of your Netscape formatted cookie file (typically `cookie.txt`) as the `data` parameter for the `netcookiejar.Read()` function.  

If using `netcookie.MODE_STRING`, pass in a Netscape formatted string of cookies as the `data` parameter for the `netcookiejar.Read()` function.

##### Invalid Unix Time Format
If the Unix timestamp from the cookie cannot be parsed for whatever reason, netcookiejar will simply set the `Expires` member of the `http.Cookie` struct to `0`.