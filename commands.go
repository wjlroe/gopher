package main

import (
	"regexp"
	"time"
	"strings"
	"fmt"
	"http"
	"io/ioutil"
	"strconv"
	"json"
	"exec"
	"os"
)

var man_page = regexp.MustCompile("man (.*)")
var hoogle = regexp.MustCompile("hoogle (.*)")
var what_time = regexp.MustCompile(".*(what).*(time).*")
var who_are_you = regexp.MustCompile(".*(who).*(are).*(you).*[?]*.*")
var gist = regexp.MustCompile(".*gist[^0-9]*([0-9][0-9]*)(.*)")

type GistInfo struct {
	Files []string
	CreatedAt string
	Description string
	Owner string
}

type GistsInfo struct {
	Gists []*GistInfo
}

func GetGistInfo(gist_id uint) []*GistInfo {
	url := ParseURL(fmt.Sprintf("http://gist.github.com/api/v1/json/%d", gist_id))
	response, err := http.Get(url.String())
	if err != nil {
		fmt.Printf("Error: %s", err.String())
		return nil
	}

	defer response.Body.Close()

	var body []byte
	gistsInfo := &GistsInfo{Gists: nil}
	body, err = ioutil.ReadAll(response.Body)
	err = json.Unmarshal(body, &gistsInfo)
	if err != nil {
		fmt.Printf("Error getting Gift info: %s\n", err.String())
	}
	return gistsInfo.Gists
}

func GetGist(gist_id uint, filename string) []string {
	if filename == "" {
		gistsInfo := GetGistInfo(gist_id)
		if gistsInfo != nil && len(gistsInfo) > 0 {
			filename = gistsInfo[0].Files[0]
		}
	}
	fmt.Printf("Gist file again: %s\n", filename)
	url := ParseURL(fmt.Sprintf("http://gist.github.com/raw/%d/%s", gist_id, filename))
	response, err := http.Get(url.String())
	if err != nil {
		return []string{fmt.Sprintf("Error: %s", err.String())}
	}

	defer response.Body.Close()

	var body []byte
	body, err = ioutil.ReadAll(response.Body)
	return []string{"", string(body)}
}

func RunShellCommand(cmd string, cmdargs []string) string {
	cmd_name, err := exec.LookPath(cmd)
	if err != nil {
		return fmt.Sprintf("cmd: %s could not be found in your PATH", cmd)
	}
	//curr_env := os.Environ()
	//cwd, err := os.Getwd()
	cmd_obj := exec.Command(cmd, strings.Join(cmdargs, " "))
	fmt.Printf("Going to run: `%s %s`\n", cmd_name, strings.Join(cmdargs, " "))
	body, err := cmd_obj.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("Error running cmd: %s : %s", cmd, err.String())
	}
	// _,err = command.Wait(0)
	// if err != nil {
	// 	return ""
	// }
	// var body []byte
	// body, err = ioutil.ReadAll(command.Stdout)
	// if err != nil {
	// 	return ""
	// }
	return string(body)
}

func Manpage(manpage []string) []string {
	output := RunShellCommand("man", manpage)
	return []string{"", output}
}

func Hoogle(query []string) []string {
	output := RunShellCommand("hoogle", query)
	return []string{"", output}
}

func Command(cmd string) []string {
	command := []byte(strings.ToLower(cmd))
	if what_time.Match(command) {
		current_time := time.LocalTime()
		formatted_time := current_time.Format(time.Kitchen)
		return []string{formatted_time}
	} else if who_are_you.Match(command) {
		return []string{"I am Gopher, I'm here to help."}
	} else if cmd == "gopher quit" {
		os.Exit(0)
	} else if matches := gist.FindSubmatch(command); matches != nil {
		for _,match := range matches {
			fmt.Printf("Match: %s\n", string(match))
		}

		if len(matches) > 1 {
			gist_id, err := strconv.Atoui(string(matches[1]))
			if err != nil {
				fmt.Printf("Could not parse the gist_id: %s\n", string(matches[1]))
				return []string{}
			}
			gist_file := strings.TrimSpace(string(matches[2]))
			fmt.Printf("Gist file: %s\n", gist_file)
			gist_content := GetGist(gist_id, gist_file)
			fmt.Printf("Gist Content: %s\n", gist_content)
			return gist_content
		}
		return []string{}
	} else if matches := hoogle.FindSubmatch([]byte(cmd)); matches != nil {
		return Hoogle(strings.Fields(string(matches[1])))
	} else if matches := man_page.FindSubmatch([]byte(cmd)); matches != nil {
		return Manpage(strings.Fields(string(matches[1])))
	}
	return []string{"Dunno what you are asking me. WTF dude?"}
}