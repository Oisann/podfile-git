package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
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
		if i == 0 {
			continue
		}
		captures[name] = match[i]

	}
	return captures
}

func cmd(cmd string, path string) {
	command := exec.Command("bash", "-c", cmd)
	command.Dir = path
	out, err := command.Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(out))
}

func main() {
	argsWithoutProg := os.Args[1:]

	filter := ""
	skipSelf := false
	for index, arg := range argsWithoutProg {
		if arg == "--filter" {
			filter = argsWithoutProg[index + 1]
			argsWithoutProg = append(argsWithoutProg[:index], argsWithoutProg[index+2:]...)
		}

		if arg == "--skipself" {
			skipSelf = true
			argsWithoutProg = append(argsWithoutProg[:index], argsWithoutProg[index+1:]...)
		}
	}

	args := fmt.Sprintf("%s", argsWithoutProg)
	args = strings.Replace(args, "[", "", -1)
	args = strings.Replace(args, "]", "", -1)

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	if !skipSelf {
		command_self := fmt.Sprintf("git -c color.status=always -c color.ui=always %s", args)
		fmt.Printf("Running the git command 'git %s' in %s.\n", args, "this project")
		cmd(command_self, dir)
	}

	file, err := os.Open(dir + "/Podfile")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data := myExp.FindStringSubmatchMap(scanner.Text())
		if len(data) > 0 {

			if filter != "" {
				var filterRegex = myRegexp{regexp.MustCompile(filter)}
				if len(filterRegex.FindStringSubmatch(strings.ToLower(data["name"]))) == 0 {
					if len(filterRegex.FindStringSubmatch(strings.ToLower(data["path"]))) == 0 {
						fmt.Printf("%s does not match the filter \"%s\", skipping.\n", data["name"], filter)
						continue
					}
				}
			}

			localDir := dir
			localPath := data["path"]

			dirFolders := strings.Split(localDir, "/")
			pathFolders := strings.Split(localPath, "/")
			for _, folder := range pathFolders {
				if folder == ".." {
					dirFolders = dirFolders[:len(dirFolders)-1]
				} else {
					dirFolders = append(dirFolders, folder)
				}
			}

			command := fmt.Sprintf("git -c color.status=always -c color.ui=always %s", args)
			fmt.Printf("Running the git command 'git %s' in %s.\n", args, data["name"])
			cmd(command, fmt.Sprintf("%s/", strings.Join(dirFolders[:], "/")))
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
