package main

// Clone Github repo and pre-process for open-source scanning

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/kyokomi/emoji"
	"github.com/rcompos/ducky/ducky"

)

const sleepytime time.Duration = 1

var debug bool

//var thirdpartyDir string = "third_party"
//var opensourceDir string = "third_party/open_source"

func main() {

	var account, repoIn, repo, gitURL, input, tokenBD, fileBD string
	var tag string = ""
	var force, help, scan bool

	// Input value in format: account/repo:tag  (where tag is optional)
	flag.StringVar(&input, "r", os.Getenv("DUCKY_GIT_REPO"), "Git repo (account/repo:version) to clone, prep and scan (required) (envvar DUCKY_GIT_REPO)")
	flag.StringVar(&tokenBD, "t", os.Getenv("BD_TOKEN"), "Black Duck auth token (required) (envvar BD_TOKEN)")
	flag.StringVar(&fileBD, "l", os.Getenv("DUCKY_FILE"), "Input file")
	// Add file processing capability
	flag.BoolVar(&debug, "d", false, "Debugging output")
	flag.BoolVar(&help, "h", false, "Help")
	flag.BoolVar(&force, "f", false, "Force without confirmation prompt")
	flag.BoolVar(&scan, "s", true, "Perform Black Duck scan")
	flag.Parse()

	if help == true {
		ducky.HelpMe("")
	}

	if scan == true && tokenBD == "" {
		ducky.HelpMe("Must supply Black Duck token")
	}

	if input == "" && fileBD == "" {
		ducky.HelpMe("Must supply repo as argument (-r) or from input file (-l).")
	}

	// If token supplied on command-line, then set the envvar
	if tokenBD != "" && os.Getenv("BD_TOKEN") == "" {
		os.Setenv("BD_TOKEN", tokenBD)
	}

	gitBaseURL := "http://git@www.github.com"
	gitBaseDir := "github.com"

	// Default to Github.com
	if gitURL == "" {
		gitURL = gitBaseURL
	}

	//gitCloneURL := gitURL + "/" + account + "/" + repo
	baseRepoDir := "repos"

	color.Set(color.FgMagenta)
	emoji.Println("\n          NKS  :duck: Ducky ")
	color.Unset()

	// Must provide repo as arg if no input file specified
	//var in, lines []string
	var lines []string
	var inputType string = ""
	if fileBD != "" {
		inputType = "file"
		// Read entire file into memory
		lines = ducky.ReadInFile(fileBD)
	} else {
		inputType = "argument"
		lines = append(lines, input)
	}

	ducky.DuckArt()
	fmt.Println("   NKS Ducky Black Duck Scanner")
	fmt.Printf("Clone, pre-process and scan GitHub repo(s)\n\n")
	fmt.Printf("Repo list from %s\n", inputType)
	for j := range lines {
		fmt.Printf("%s\n", lines[j])
	}

	if force != true {
		ducky.PromptRead()
	}

	if debug == true {
		fmt.Printf("Github account: %s\n", account)
		fmt.Printf("Github repo: %s\n", repo)
		fmt.Printf("Github tag: %s\n", tag)
		fmt.Printf("Repo dir: %s\n\n", baseRepoDir)
		fmt.Printf("Base dir: %s\n\n", gitBaseDir)
	}

	for i := range lines {
		r := lines[i]
		in := strings.Split(r, "/")

		//fmt.Printf("-> %s, %s, %s\n", r, in[0], in[1]) // DEBUG

		if len(in) != 2 {
			fmt.Println("Specify account/repo as input argument")
			continue
		}

		account = in[0]
		repoIn = in[1]
		//fmt.Printf("Git account: %s\n", account)
		//fmt.Printf("Git repository: %s\n", repoIn)

		re := regexp.MustCompile(`:`)
		repoWithTag := re.Split(repoIn, -1)
		if len(repoWithTag) == 2 {
			//for i := range repoWithTag {
			//	fmt.Printf("repoWithTag[%v] %v\n", i, repoWithTag[i])
			//}
			repo = repoWithTag[0]
			tag = repoWithTag[1]
		} else if len(repoWithTag) == 1 {
			repo = repoIn
		} else { // TODO: Change to notice and skip
			ducky.HelpMe("")
		}

		fmt.Printf("Scanning repo: %s\n\n", r)
		ducky.ScanBD(account, repo, tag, gitURL, baseRepoDir, gitBaseDir, scan)
	}

	color.Set(color.FgMagenta)
	emoji.Println("\n   NKS  :duck: Ducky  complete")
	color.Unset()

} // End main
