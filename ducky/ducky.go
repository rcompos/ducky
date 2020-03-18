package ducky

// Clone Github repo and pre-process for open-source scanning

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"

	"time"

	"github.com/fatih/color"
	"github.com/gofrs/flock"
)

//const sleepytime time.Duration = 1

var debug bool
var thirdpartyDir string = "third_party"
var opensourceDir string = "third_party/open_source"
var mu sync.Mutex
var count int

var BaseRepoDir string // Test it out

func PullDependencies(dir, repo, path, lang, pman, gitdir string) {
	fmt.Println()
	//fmt.Printf("-> Directory: %s\n", dir)
	//fmt.Printf("-> Repository: %s\n", repo)
	//fmt.Printf("-> Language: %s\n", lang)
	//fmt.Printf("-> Package Manager: %s\n\n", pman)
	if lang == "go" {
		if pman == "mod" {
			pullDepsGo(dir, repo, path, gitdir, pman)
		} else if pman == "dep" {
			pullDepsGo(dir, repo, path, gitdir, pman)
		}
	} else if lang == "python" {
		pullDepsPython(dir, repo, path, gitdir)
	} else if lang == "node" {
		pullDepsNode(dir, repo, path, gitdir)
	} else {
		fmt.Println("ERROR: Can't pull deps!")
		return
	}
}

func pullDepsGo(dir, r, path, gitdir, pman string) {
	repo := dir + "/" + r
	fmt.Printf("Getting dependencies...\n")
	fmt.Printf("Repo: %s\n", repo)

	cmd := "ls -Alf"
	fmt.Printf("# %s\n", cmd)
	cmdExec := exec.Command("bash", "-c", cmd)
	cmdExec.Stdout = os.Stdout
	cmdExec.Stderr = os.Stderr
	cmdExec.Dir = repo
	err := cmdExec.Run()
	if err != nil {
		fmt.Printf("WARNING: List directory failed\n%v\n", err)
	}
	fmt.Println()

	// Pull all Go dependencies to directory vendor
	fmt.Println("Pull all Go dependencies.")
	var cmdGoMod string
	if pman == "mod" {
		cmdGoMod = "go mod vendor"
	} else if pman == "dep" {
		cmdGoMod = "dep ensure"
	}
	fmt.Printf("# %s\n", cmdGoMod)
	cmdGoModExec := exec.Command("bash", "-c", cmdGoMod)
	cmdGoModExec.Stdout = os.Stdout
	cmdGoModExec.Stderr = os.Stderr
	cmdGoModExec.Dir = repo
	errGoMod := cmdGoModExec.Run()
	if errGoMod != nil {
		fmt.Printf("ERROR: Go dependency collection failed\n%v\n", errGoMod)
		return
	}
	fmt.Println()

	//CreateDir(repo + "/" + opensourceDir, path)
	CreateDir(repo + "/" + thirdpartyDir)

	// Move Go dependencies to special dir
	//cmdMv := "mv vendor/* " + opensourceDir
	cmdMv := "mv vendor " + opensourceDir
	//cmdMv := "cp -af vendor/* " + opensourceDir + " && rm -R vendor/*"
	fmt.Printf("# %s\n", cmdMv)
	cmdMvExec := exec.Command("bash", "-c", cmdMv)
	cmdMvExec.Stdout = os.Stdout
	cmdMvExec.Stderr = os.Stderr
	cmdMvExec.Dir = repo
	errMv := cmdMvExec.Run()
	if errMv != nil {
		fmt.Printf("WARNING: Copy failed:\t%v\n", errMv)
	}
	fmt.Println()

	// Create git dir under vender
	cmdMkdir := "mkdir -p vendor/" + gitdir
	fmt.Printf("# %s\n", cmdMkdir)
	cmdMkdirExec := exec.Command("bash", "-c", cmdMkdir)
	cmdMkdirExec.Stdout = os.Stdout
	cmdMkdirExec.Stderr = os.Stderr
	cmdMkdirExec.Dir = repo
	errMkdir := cmdMkdirExec.Run()
	if errMkdir != nil {
		fmt.Printf("WARNING: Mkdir failed:\t%v\n", errMkdir)
	}
	fmt.Println()

	//vendorDir := "vendor/" + gitdir

	//
	// Move NetApp repos back to vendor
	//
	// TODO: Keep dirs in slice and for loop
	naDir := opensourceDir + "/" + gitdir + "/netapp"
	moveVendor(naDir, repo, path, gitdir)

	nDir := opensourceDir + "/" + gitdir + "/NetApp"
	moveVendor(nDir, repo, path, gitdir)

	// Move StackPointCloud repos back to vendor
	spcDir := opensourceDir + "/" + gitdir + "/stackpointcloud"
	moveVendor(spcDir, repo, path, gitdir)

	sDir := opensourceDir + "/" + gitdir + "/StackPointCloud"
	moveVendor(sDir, repo, path, gitdir)

}

func moveVendor(dir, repo, path, gitdir string) {
	//fmt.Printf("dir: %s\n", dir)
	if _, err := os.Stat(path + "/" + repo + "/" + dir); !os.IsNotExist(err) {
		// Move NetApp repos back to vendor
		//cmdN := "cp -af " + dir + " " + vendorDir + " && rm -R " + dir
		cmdN := "mv " + dir + " vendor/" + gitdir
		fmt.Printf("# %s\n", cmdN)
		cmdNExec := exec.Command("bash", "-c", cmdN)
		cmdNExec.Stdout = os.Stdout
		cmdNExec.Stderr = os.Stderr
		cmdNExec.Dir = repo
		errN := cmdNExec.Run()
		if errN != nil {
			fmt.Printf("WARNING: Move NetApp repos failed\n%v\n", errN)
		}
		fmt.Println()
	}
}

func pullDepsPython(dir, r, path, gitdir string) {
	repo := dir + "/" + r
	fmt.Printf("Getting dependencies...\n")
	fmt.Printf("Repo: %s\n", repo)
	cmd := "ls -Alf"
	fmt.Printf("# %s\n", cmd)
	cmdExec := exec.Command("bash", "-c", cmd)
	cmdExec.Stdout = os.Stdout
	cmdExec.Stderr = os.Stderr
	cmdExec.Dir = repo
	err := cmdExec.Run()
	if err != nil {
		fmt.Printf("ERROR: List directory failed\n%v\n", err)
		return
	}
	fmt.Println()

	pwDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pwDir)

	// Pull all Python dependencies
	cmdPip := "pip install -r requirements.txt --target=./third_party/open_source"
	fmt.Printf("# %s\n", cmdPip)
	cmdPipExec := exec.Command("bash", "-c", cmdPip)
	cmdPipExec.Stdout = os.Stdout
	cmdPipExec.Stderr = os.Stderr
	cmdPipExec.Dir = repo
	errPip := cmdPipExec.Run()
	if errPip != nil {
		fmt.Printf("ERROR: Python pip download failed\n%v\n", errPip)
	}
	CreateDir(repo + "/" + opensourceDir)
}

func pullDepsNode(dir, r, path, gitdir string) {
	repo := dir + "/" + r
	fmt.Printf("Getting dependencies...\n")
	fmt.Printf("Repo: %s\n", repo)
	cmd := "ls -Alf"
	fmt.Printf("# %s\n", cmd)
	cmdExec := exec.Command("bash", "-c", cmd)
	cmdExec.Stdout = os.Stdout
	cmdExec.Stderr = os.Stderr
	cmdExec.Dir = repo
	err := cmdExec.Run()
	if err != nil {
		fmt.Printf("ERROR: List directory failed\n%v\n", err)
		return
	}
	fmt.Println()

	CreateDir(repo + "/" + opensourceDir)

	cmdCp := "cp -a package.json " + opensourceDir
	fmt.Printf("# %s\n", cmdCp)
	cmdCpExec := exec.Command("bash", "-c", cmdCp)
	cmdCpExec.Stdout = os.Stdout
	cmdCpExec.Stderr = os.Stderr
	cmdCpExec.Dir = repo
	errCp := cmdCpExec.Run()
	if errCp != nil {
		fmt.Printf("ERROR: Moved package.json failed\n%v\n", errCp)
		return
	}
	fmt.Println()

	// Pull all node dependencies
	cmdNpm := "npm install"
	fmt.Printf("# %s\n", cmdNpm)
	cmdNpmExec := exec.Command("bash", "-c", cmdNpm)
	cmdNpmExec.Stdout = os.Stdout
	cmdNpmExec.Stderr = os.Stderr
	cmdNpmExec.Dir = repo + "/" + opensourceDir
	errNpm := cmdNpmExec.Run()
	if errNpm != nil {
		fmt.Printf("ERROR: Node npm install failed\n%v\n", errNpm)
		return
	}
	fmt.Println()

}

func DetermineLanguage(dir, repo, path string) (string, string) {
	repoDir := dir + "/" + repo
	fmt.Println("Detecting account type from package manager config files.")
	var kind, subkind string

	idGoMod := "go.mod"
	idGoDep := "Gopkg.toml"
	idPython := "requirements.txt"
	idNode := "package.json"

	fileGoMod := repoDir + "/" + idGoMod
	fileGoDep := repoDir + "/" + idGoDep
	filePython := repoDir + "/" + idPython
	fileNode := repoDir + "/" + idNode

	var foundGo, foundGoMod, foundGoDep, foundPython, foundNode bool = false, false, false, false, false
	if _, err := os.Stat(path + "/" + fileGoMod); !os.IsNotExist(err) {
		foundGo = true
		foundGoMod = true
		fmt.Printf("Found: %s (Go mod)\n", fileGoMod)
	}
	if _, err := os.Stat(path + "/" + fileGoDep); !os.IsNotExist(err) {
		foundGo = true
		foundGoDep = true
		fmt.Printf("Found: %s (Go dep)\n", fileGoDep)
	}
	if _, err := os.Stat(path + "/" + filePython); !os.IsNotExist(err) {
		foundPython = true
		fmt.Printf("Found: %s (Python)\n", filePython)
	}
	if _, err := os.Stat(path + "/" + fileNode); !os.IsNotExist(err) {
		foundNode = true
		fmt.Printf("Found: %s (Node)\n", fileNode)
	}

	var detectMulti bool
	if foundGo == false && foundPython == false && foundNode == false {
		fmt.Println("ERROR: Can't determine language type.")
		return "", ""
	}
	if (foundGo == true && foundPython == true) ||
		(foundGo == true && foundNode == true) ||
		(foundNode == true && foundPython == true) {
		detectMulti = true
	} else {
		if foundGoMod == true {
			kind = "go"
			subkind = "mod"
		} else if foundGoDep == true {
			kind = "go"
			subkind = "dep"
		} else if foundPython == true {
			kind = "python"
			subkind = "pip"
		} else if foundNode == true {
			kind = "node"
			subkind = "npm"
		}
	}

	if detectMulti == true {
		color.Set(color.FgRed)
		fmt.Println("ERROR: Can't determine language type. Multiple languages detected.")
		color.Unset()
		return "", ""
	}

	color.Set(color.FgYellow)
	fmt.Printf("\nLanguage: %s\t\tPackage Manager: %s\n", kind, subkind)
	color.Unset()
	return kind, subkind
}

func gitPull(url, dir, repo string) {
	repoDir := dir + "/" + repo
	fmt.Printf("Git repo exists: %s\n\n", repoDir)
	pullCmd := "git pull"
	fmt.Println("Pulling latest changes:")
	fmt.Printf("# %s\n", pullCmd)
	pullExec := exec.Command("bash", "-c", pullCmd)
	pullExec.Stdout = os.Stdout
	pullExec.Stderr = os.Stderr
	pullExec.Dir = repoDir
	err := pullExec.Run()
	if err != nil {
		fmt.Printf("ERROR: Git pull failed\n%v\n", err)
		return
	}
	fmt.Println()
}

func GitClone(url, dir, repo, path, tag string) {
	// Check for existence of repo directory
	repoDir := dir + "/" + repo
	if _, err := os.Stat(path + "/" + repoDir); !os.IsNotExist(err) {
		//gitPull(url, dir, repo)
		//return
		// Repo dir exists - Delete it
		deleteDir(repoDir, path)
	}

	var gitTag string = ""
	if tag != "" {
		gitTag = fmt.Sprintf("--branch %s ", tag)

	}
	gitCloneCmd := "git clone " + gitTag + url
	//gitCloneCmd := "pwd; which git"
	fmt.Println("Cloning git repo:")
	fmt.Printf("# %s\n", gitCloneCmd)
	gitCloneExec := exec.Command("bash", "-c", gitCloneCmd)
	gitCloneExec.Stdout = os.Stdout
	gitCloneExec.Stderr = os.Stderr
	gitCloneExec.Dir = dir
	err := gitCloneExec.Run()
	if err != nil {
		fmt.Printf("ERROR: Git clone failed\n%v\n", err)
		return
	}
	fmt.Println()
}

func HelpMe(msg string) {
	if msg != "" {
		fmt.Printf("%s\n\n", msg)
	}
	fmt.Println("Supply single argument with Git account and repository name separated by a slash. i.e. netapp/myrepo")
	flag.PrintDefaults()
	os.Exit(1)
}

func PromptRead() {
	reader := bufio.NewReader(os.Stdin)
	confirm := "y"
	fmt.Printf("\nTo continue type %s: ", confirm)
	text, _ := reader.ReadString('\n')
	answer := strings.TrimRight(text, "\n")
	//fmt.Printf("answer: %s \n", answer)
	//if answer == a {
	if answer == "y" || answer == "Y" {
		return
	} else {
		log.Fatal("Exiting without action.")
	}
}

func CreateDir(dir string) {
	//	if [ -d "repos" ]; then echo true; else echo false; fi
	//if _, err := os.Stat(path + "/" + dir); !os.IsNotExist(err) {
	//	return // dir exists
	//}
	mkdirCmd := fmt.Sprintf("if [ ! -d %s ]; then mkdir -p -m775 %s; fi", dir, dir)
	fmt.Printf("Creating dir: %s\n", dir)
	if debug == true {
		fmt.Printf("# %s\n", mkdirCmd)
	}
	mkdirExec := exec.Command("bash", "-c", mkdirCmd)
	mkdirExec.Stdout = os.Stdout
	mkdirExec.Stderr = os.Stderr
	//mkdirExec.Dir = path
	err := mkdirExec.Run()
	if err != nil {
		fmt.Printf("ERROR: Create dir failed\n%v\n", err)
		return
	}
	fmt.Println()
}

func deleteDir(dir, path string) {
	if _, err := os.Stat(path + "/" + dir); os.IsNotExist(err) {
		fmt.Printf("Dir not found: %s\n", dir)
		return // dir doesn't exist
	}
	delCmd := fmt.Sprintf("rm -fr %s", dir)
	fmt.Printf("Deleting existing repo dir: %s\n", dir)
	if debug == true {
		fmt.Printf("# %s\n", delCmd)
	}
	delExec := exec.Command("bash", "-c", delCmd)
	delExec.Stdout = os.Stdout
	delExec.Stderr = os.Stderr
	err := delExec.Run()
	if err != nil {
		fmt.Printf("ERROR: Delete dir failed\n%v\n", err)
		return
	}
	fmt.Println()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func ScanBD(account, repo, tag, gitURL, baseRepoDir, gitBaseDir string, scan bool) {
	// Black Duck scans

	pwDir, errDir := os.Getwd()
	if errDir != nil {
		log.Fatal(errDir)
	}
	fmt.Printf("PWD: %s\n", pwDir)

	gitCloneURL := gitURL + "/" + account + "/" + repo
	//baseRepoDir := "repos"  // TODO: Use global from main !?
	accountDir := baseRepoDir + "/" + account

	// Lock file
	fileName := baseRepoDir + "/" + account + "/" + repo + ".lock"
	fileLock := flock.New(fileName)

	path, errP := os.Getwd()
	if errP != nil {
		log.Println(errP)
	}
	fmt.Printf("Path: %s\n", path)

	CreateDir(accountDir)
	locked, err := fileLock.TryLock()
	if err != nil {
		// handle locking error
		fmt.Printf("Locking error! %v\n", err)
		//os.Exit(33)
		return
	}
	if locked {
		//defer file unlock
		defer func() {
			if err := fileLock.Unlock(); err != nil {
				// handle unlock error
				fmt.Printf("Unlock fail: %v\n", err)
			}
		}()
		// Do what needs to be done
		GitClone(gitCloneURL, accountDir, repo, path, tag)
		language, pkgman := DetermineLanguage(accountDir, repo, path)
		// For debugging
		if debug == true {
			fmt.Printf("accountDir: '%s'\n", accountDir)
			fmt.Printf("repo: '%s'\n", repo)
			fmt.Printf("language: '%s'\n", language)
			fmt.Printf("pkgman: '%s'\n", pkgman)
			fmt.Printf("gitBaseDir: '%s'\n", gitBaseDir)
			fmt.Printf("path: '%s'\n", path)
		}
		PullDependencies(accountDir, repo, path, language, pkgman, gitBaseDir)

		// Perform Black Duck scan
		if scan == true {
			fmt.Println("Performing Black Duck scan")
			dir := path + "/" + accountDir + "/" + repo
			//fmt.Printf("dir: %s\n", dir)
			if _, err := os.Stat(dir); os.IsNotExist(err) {
				fmt.Printf("Can't scan. Repo dir not found: %v: %v", dir, err)
				//os.Exit(34) // repo dir doesn't exist
				return
			}

			//pwDir, errDir := os.Getwd()
			//if errDir != nil {
			//	log.Fatal(errDir)
			//}
			//fmt.Printf("PWD: %s\n", pwDir)
			//fmt.Printf("path: %s\n", path)

			//scanCmd := fmt.Sprintf("./scan.sh %s", fullRepoDir)  // TODO: Implement tag
			scanCmd := fmt.Sprintf("./scan.sh %s", dir) // TODO: Implement tag
			fmt.Printf("scanCmd: %s\n", scanCmd)
			scanExec := exec.Command("bash", "-c", scanCmd)
			scanExec.Stdout = os.Stdout
			scanExec.Stderr = os.Stderr
			scanExec.Dir = path
			err := scanExec.Run()
			if err != nil {
				fmt.Printf("ERROR: Black Duck scan execution failed!\n%v\n", err)
				//os.Exit(35)
				return
			}
		} else {
			fmt.Println("Skipping Black Duck scan.")
		}

		//output := fmt.Sprintf("Successful lock!\n")
		//w.Write([]byte(output))
		fmt.Printf("path: %s; locked: %v\n", fileLock.Path(), fileLock.Locked())
		time.Sleep(1 * time.Second)
		fmt.Println("File unlocked.")
	} else {
		fmt.Println("File could not be locked.")
	}
}

func ReadInFile(i string) []string {
	// Read line-by-line
	var lines []string
	file, err := os.Open(i)
	if err != nil {
		log.Println(err)
		return lines
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func DuckArt() {
	color.Set(color.FgYellow)
	fmt.Printf(`
	        __
	       /'{>
	   ____) (____
	 //'--;   ;--'\\
	///////\_/\\\\\\\
	       m m `)
	color.Unset()
	fmt.Println()
}
