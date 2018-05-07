package main

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
    "io/ioutil"
    "regexp"
)

type myRegexp struct {
	*regexp.Regexp
}

var myExp = myRegexp{regexp.MustCompile(`pod '(?P<name>.*)', :path => '(?P<path>.*)'`)}

func (r *myRegexp) FindStringSubmatchMap(s string) map[string]string {
	captures := make(map[string]string)

	match := r.FindStringSubmatch(s)
	if match == nil {
		return captures
	}

	for i, name := range r.SubexpNames() {
		// 
		if i == 0 {
			continue
		}
		captures[name] = match[i]

	}
	return captures
}

func main() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
    if err != nil {
            log.Fatal(err)
    }

    dat, err := ioutil.ReadFile(dir + "/Podfile")
    if err != nil {
    	fmt.Println("No Podfile found.")
    	return
    }
    //fmt.Print(string(dat))
	fmt.Printf("%+v\n", myExp.FindStringSubmatchMap(string(dat)))
}

