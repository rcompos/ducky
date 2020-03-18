package duckyapi

// Clone Github repo and pre-process for open-source scanning

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/rcompos/ducky/ducky"

	"github.com/fatih/color"
	"github.com/gofrs/flock"
	"github.com/gorilla/mux"
)

//const sleepytime time.Duration = 1

var mu sync.Mutex
var count int
var BaseRepoDir string // Test it out
var HttpURL string

// info handler displays http header
func Info(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
	fmt.Fprintf(w, "%s %s %s\n", r.Method, r.URL, r.Proto)
	for k, v := range r.Header {
		fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
	}
	fmt.Fprintf(w, "Host = %q\n", r.Host)
	fmt.Fprintf(w, "RemoteAddr = %q\n", r.RemoteAddr)
	if err := r.ParseForm(); err != nil {
		log.Print(err)
	}
	for k, v := range r.Form {
		fmt.Fprintf(w, "Form[%q] = %q\n", k, v)
	}
}

// counter displays the page count
func Counter(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Count %d\n", count)
	//log.Printf("Count %d\n", count)
}

func Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
	//log.Println("pong")
}

func TestCounter(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	count++
	mu.Unlock()
	w.Write([]byte("Count incremented"))
}

func ViewBD(w http.ResponseWriter, r *http.Request) {
	color.Set(color.FgMagenta)
	output := "Black Duck viewer"
	w.Write([]byte(output))
	fmt.Println(output)
	color.Unset()
	vars := mux.Vars(r)
	account := vars["account"]
	repo := vars["repo"]
	tag := vars["tag"]
	output = fmt.Sprintf("GitHub Account: %s\nGitHub Repo: %s\nGitHub Tag: %s\n", account, repo, tag)
	w.Write([]byte(output))
	fmt.Println(output)

	//
	// Read and print log file
	//

}

func ScanBD(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("Black Duck scanner\n"))
	vars := mux.Vars(r)
	account := vars["account"]
	repo := vars["repo"]
	tag := vars["tag"]
	accountDir := BaseRepoDir + "/" + account
	logfile := accountDir + "/" + repo + ".log"

	// Delete log file
	err := deleteFile(logfile)
	if err != nil {
		fmt.Printf("ERROR: Couldn't delete log file %s\n", logfile)
	}

	// Create log file
	//logFileDateFormat := "2006-01-02-150405"
	//logStamp := time.Now().Format(logFileDateFormat)
	viewURL := "http://" + HttpURL + "/api/v1/view/" + account + "/" + repo
	fileURL := "http://" + HttpURL + "/" + BaseRepoDir + "/" + account + "/" + repo + ".log"

	logf, err := os.OpenFile(logfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
	mw := io.MultiWriter(os.Stdout, logf, w)
	log.SetOutput(mw)
	if err != nil {
		log.Fatal(err)
	}
	defer logf.Close()
	log.Printf("Starting scan: %s/%s\n", account, repo)
	log.Printf("Logs: %s\n", viewURL)
	logMsg := fmt.Sprintf("Logs: %s\n", viewURL)
	w.Write([]byte(logMsg))
	fileMsg := fmt.Sprintf("Fileserver URL: %s\n", fileURL)
	w.Write([]byte(fileMsg))
	/*
		if gitTag != "" {
			output := fmt.Sprintf("Git tag specified!  %s\n", gitTag)
			w.Write([]byte(output))
		}
	*/
	//output := fmt.Sprintf("GitHub Account: %s\nGitHub Repo: %s\nGitHub Tag: %s\n", account, repo, tag)
	//w.Write([]byte(output))

	//gitCloneCmd := "git clone " + gitTag + url
	gitCloneURL := "http://git@github.com/" + account + "/" + repo
	gitBaseDir := "github.com"
	// Lock file
	//fileName := BaseDir + "/" + gitAccount + "/" + gitRepo + ".lock"
	//baseRepoDir := "repos" // TODO: Use global from main !?

	// Lock file
	fileName := BaseRepoDir + "/" + account + "/" + repo + ".lock"
	fileLock := flock.New(fileName)

	path, errP := os.Getwd()
	if errP != nil {
		log.Println(errP)
	}
	fmt.Printf("Path: %s\n", path)

	ducky.CreateDir(accountDir)
	locked, err := fileLock.TryLock()
	if err != nil {
		// handle locking error
		output := fmt.Sprintf("Locking error! %v\n", err)
		w.Write([]byte(output))
	}

	if locked {
		//defer file unlock
		defer func() {
			if err := fileLock.Unlock(); err != nil {
				// handle unlock error
				output := fmt.Sprintf("Unlock fail: %v\n", err)
				w.Write([]byte(output))
			}
			deleteFile(fileName)
		}()

		// Do work

		ducky.GitClone(gitCloneURL, accountDir, repo, path, tag)
		language, pkgman := ducky.DetermineLanguage(accountDir, repo, path)
		ducky.PullDependencies(accountDir, repo, path, language, pkgman, gitBaseDir)

		fullRepoDir := path + "/" + accountDir + "/" + repo
		if _, err := os.Stat(fullRepoDir); os.IsNotExist(err) {
			//fmt.Printf("Dir not found: %s\n", fullRepoDir)
			w.Write([]byte("Can't scan. Repo dir not found: " + fullRepoDir + "\n"))
			return // repo dir doesn't exist
		}

		//scanCmd := fmt.Sprintf("./scan.sh %s", fullRepoDir)  // TODO: Implement tag
		dir := path + "/" + accountDir + "/" + repo
		scanCmd := fmt.Sprintf("./scan.sh %s", dir) // TODO: Implement tag
		scanExec := exec.Command("bash", "-c", scanCmd)
		//scanExec.Stdout = os.Stdout
		scanExec.Stdout = mw
		scanExec.Stderr = os.Stderr
		err = scanExec.Run()
		if err != nil {
			errMsg := fmt.Sprintf("\nERROR: Black Duck scan execution failed for %s/%s!\n%v\n\n", account, repo, err)
			w.Write([]byte(errMsg))
			color.Set(color.FgRed)
			fmt.Println(errMsg)
			color.Unset()
			//os.Exit(33)
			return
		}

		//output := fmt.Sprintf("Successful lock!\n")
		//w.Write([]byte(output))
		output4 := fmt.Sprintf("path: %s; locked: %v\n", fileLock.Path(), fileLock.Locked())
		w.Write([]byte(output4))
		time.Sleep(1 * time.Second)
		output2 := "File unlocked."
		w.Write([]byte(output2))
	} else {
		output := "File could not be locked."
		w.Write([]byte(output))
	}

	w.Write([]byte(fmt.Sprintf("\nExiting Black Duck scan: %s/%s\n\n", account, repo)))

}

func deleteFile(file string) error {
	// delete file
	var err = os.Remove(file)
	if err != nil {
		log.Println(err)
	}
	log.Printf("File deleted: %s\n", file)
	return err
}
