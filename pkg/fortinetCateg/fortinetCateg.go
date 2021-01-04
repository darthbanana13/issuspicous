package fortinetCateg

import (
	"fmt"
	"regexp"
	"net/http"
	"io/ioutil"
)

type Categ struct {
	Name	string
	Err		error
}

const webFilterURL = "https://www.fortiguard.com/webfilter?version=8&q="
const categoryRegex = `<h4 class="info_title">Category: (.*)</h4>`
const outputTemplate = `Fortinet Category:	`

func NewCateg(addr string) Categ {
	// TODO: Handle errors
	resp, _ := http.Get(webFilterURL + addr)
	if resp.StatusCode == 200 {
		// TODO: Handle errors
		body, _ := ioutil.ReadAll(resp.Body)
		reg := regexp.MustCompile(categoryRegex)
		match := reg.FindStringSubmatch(string(body))
		if len(match) > 1 {
			return Categ{
				Name:	match[1],
				Err:	nil,
			}
		}
		return Categ{
			Name:	"",
			Err:	fmt.Errorf("Parse error"),
		}
	}
	return Categ{
		Name:	"",
		Err:	fmt.Errorf("Response blocked"),
	}
}

func (c Categ) String() string {
	if c.Err != nil {
		return outputTemplate + c.Err.Error()
	}
	return outputTemplate + c.Name
}
