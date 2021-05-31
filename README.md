# httperr [![Go Reference](https://pkg.go.dev/badge/github.com/ClavinJune/httperr.svg)](https://pkg.go.dev/github.com/ClavinJune/httperr)
Helper error for wrapping golang error with HTTP Status Code and stacktrace.

## Usage

```shell
go get -u github.com/ClavinJune/httperr@latest
```

## Example

```go
package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/ClavinJune/httperr"
)

func main() {
	// simple HTTP Error
	_ = httperr.From(http.StatusNotFound)
	base := httperr.New(sql.ErrNoRows, http.StatusNotFound, "user not found")
	wrapBase := httperr.Wrap(base, "wrap base with message")
	wrapAgain := httperr.Wrap(wrapBase, "once again, wrapped")

	fmt.Println(wrapAgain)
	/* stdout:
	   {
	     "cause": {
	       "cause": {
	         "error": "sql: no rows in result set",
	         "message": "user not found",
	         "caller": "main.go:13"
	       },
	       "message": "wrap base with message",
	       "caller": "main.go:14"
	     },
	     "message": "once again, wrapped",
	     "caller": "main.go:15"
	   }
	*/
	fmt.Println(errors.Is(wrapAgain, sql.ErrNoRows)) // stdout: true

	var ee *httperr.Error
	if errors.As(wrapAgain, &ee) {
		fmt.Println(ee.StatusCode()) // stdout: 404
	}
}
```