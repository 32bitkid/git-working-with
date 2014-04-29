package main

import "fmt"
import "os"
import "flag"
import "strings"
import "os/exec"
import "github.com/joshlf13/term"

type gitConfiguration struct {
	env string
}

func (config *gitConfiguration) Exec(args ...string) []byte {
	envFlag := fmt.Sprintf("--%s", config.env)
	args = append([]string{"config", envFlag}, args...)
	cmd := exec.Command("git", args...)
	out, _ := cmd.Output()
	return out
}

func (config *gitConfiguration) get(key string) string {
	return strings.Trim(string(config.Exec(key)), "\r\n ")
}

func (config *gitConfiguration) set(key, val string) {
	config.Exec(key, val)
}

func (config *gitConfiguration) unset(key string) {
	config.Exec("--unset", key)
}

var gitLocal = gitConfiguration{"local"}
var gitGlobal = gitConfiguration{"global"}

func printWorking(who, what string) {
	fmt.Fprint(os.Stdout, "\n  ")
	term.LightGreen(os.Stdout, who)
	fmt.Fprintf(os.Stdout, " %s\n", what)

}

func workingAlone(who string) {
	printWorking(who, "is working alone")
}

func workingTogether(who string) {
	printWorking(who, "are working as a pair")
}

var whomFlag = flag.Bool("who", false, "Display who is presently working")
var noFlag = flag.Bool("no", false, "Display who is presently working")

const usernameKey = "user.name"

func main() {
	flag.Parse()
	if *whomFlag == true || flag.NArg() == 0 {
		if pair := gitLocal.get(usernameKey); pair != "" {
			workingTogether(pair)
		} else {
			me := gitGlobal.get(usernameKey)
			workingAlone(me)
		}
	} else if *noFlag == true || flag.Arg(0) == "-" {
		gitLocal.unset(usernameKey)
		workingAlone(gitGlobal.get(usernameKey))
	} else {
		pair := fmt.Sprintf("%s & %s", gitGlobal.get(usernameKey), flag.Arg(0))
		gitLocal.set(usernameKey, pair)
		workingTogether(pair)
	}
}
