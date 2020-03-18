package main

// Clone Github repo and pre-process for open-source scanning

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/fatih/color"
	"github.com/gorilla/mux"
	"github.com/kyokomi/emoji"
	"github.com/rcompos/ducky/ducky"
	"github.com/rcompos/ducky/duckyapi"
)

const sleepytime time.Duration = 1

var baseDir string
var debug bool
var httpDucky string

// BaseRepoDir Default base directory for repos
var BaseRepoDir string = "repos"
var HttpURL string = "localhost:8080"

//var thirdpartyDir string = "third_party"
//var opensourceDir string = "third_party/open_source"

func main() {

	var tokenBD string
	var help bool

	//flag.StringVar(&input, "r", os.Getenv("GIT_REPO"), "Git repository (account/repo) to prep")
	//flag.StringVar(&tag, "t", os.Getenv("GIT_TAG"), "Git tag to clone (envvar GIT_TAG)")
	flag.StringVar(&baseDir, "b", os.Getenv("DUCKY_BASEDIR"), "Base directory where files stored (envvar DUCKY_BASEDIR)")
	flag.StringVar(&tokenBD, "t", os.Getenv("BD_TOKEN"), "Black Duck auth token (envvar BD_TOKEN)")
	flag.BoolVar(&debug, "d", false, "Debugging output")
	flag.BoolVar(&help, "h", false, "Help")
	flag.Parse()

	defaultBaseDir := "repos"
	if baseDir == "" {
		baseDir = defaultBaseDir
		duckyapi.BaseRepoDir = defaultBaseDir
	} else { // baseDir specified
		duckyapi.BaseRepoDir = baseDir
	}

	// If token supplied on command-line, then set the envvar
	if tokenBD != "" && os.Getenv("BD_TOKEN") == "" {
		os.Setenv("BD_TOKEN", tokenBD)
	}

	tokenClean := strings.Repeat("*", utf8.RuneCountInString(tokenBD))
	//fmt.Printf("Black Duck token: %s\n", tokenBD) // NOT FOR PROD
	fmt.Printf("Black Duck token: %s\n", tokenClean)

	color.Set(color.FgMagenta)
	emoji.Println("\n          NKS  :duck: Ducky ")
	color.Unset()

	if help == true {
		fmt.Println()
		ducky.HelpMe("")
	}

	ducky.DuckArt()
	fmt.Println("\n   NKS Ducky Black Duck Scanner API starting")

	if debug == true {
		fmt.Println("Debugging mode")
		fmt.Printf("baseDir: %s\n", baseDir)
	}

	path, errP := os.Getwd()
	if errP != nil {
		log.Println(errP)
	}
	fmt.Printf("Path: %s\n", path)

	elon := mux.NewRouter().StrictSlash(true)

	//
	//		Add git repo cache?
	//

	tester := func(w http.ResponseWriter, r *http.Request) {
		msg := "NKS Ducky"
		w.Write([]byte(msg))
		w.Write([]byte("\n"))
	}

	apiV1 := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("/api/v1\n"))
		w.Write([]byte("/api/v1/tester\n"))
		w.Write([]byte("/api/v1/info\n"))
		w.Write([]byte("/api/v1/ping\n"))
		w.Write([]byte("/api/v1/err\n"))
		w.Write([]byte("/api/v1/count\n"))
		w.Write([]byte("/api/v1/counttest\n"))
		w.Write([]byte("/api/v1/scan\n"))
		w.Write([]byte("/api/v1/view\n"))
	}

	elon.HandleFunc("/", tester)
	elon.HandleFunc("/api", apiV1)
	elon.HandleFunc("/api/", apiV1)
	elon.HandleFunc("/api/v1", apiV1)
	elon.HandleFunc("/api/v1/", apiV1)
	elon.HandleFunc("/api/v1/tester", tester)
	elon.HandleFunc("/api/v1/tester/", tester)
	elon.HandleFunc("/api/v1/info", duckyapi.Info)
	elon.HandleFunc("/api/v1/info/", duckyapi.Info)
	elon.HandleFunc("/api/v1/ping", duckyapi.Ping)
	elon.HandleFunc("/api/v1/ping/", duckyapi.Ping)
	elon.HandleFunc("/api/v1/count", duckyapi.Counter)
	elon.HandleFunc("/api/v1/count/", duckyapi.Counter)
	elon.HandleFunc("/api/v1/counttest", duckyapi.TestCounter)
	elon.HandleFunc("/api/v1/counttest/", duckyapi.TestCounter)
	elon.HandleFunc("/api/v1/scan/{account}/{repo}:{tag}", duckyapi.ScanBD)
	elon.HandleFunc("/api/v1/scan/{account}/{repo}:{tag}/", duckyapi.ScanBD)
	elon.HandleFunc("/api/v1/scan/{account}/{repo}", duckyapi.ScanBD)
	elon.HandleFunc("/api/v1/scan/{account}/{repo}/", duckyapi.ScanBD)
	elon.HandleFunc("/api/v1/view/{account}/{repo}", duckyapi.ViewBD)
	elon.HandleFunc("/api/v1/view/{account}/{repo}/", duckyapi.ViewBD)
	//elon.HandleFunc("/api/v1/err", viewErr)
	//elon.HandleFunc("/api/v1/err/", viewErr)
	//elon.HandleFunc("/api/v1/upload", uploader)
	//elon.HandleFunc("/api/v1/upload/", uploader)

	//
	//  FileServer for log files?
	//
	// File server for upload and download
	//maxUploadSize := 2 * 1024 // 2 MB
	fileServerPath := BaseRepoDir
	//httpfs := http.FileServer(http.Dir(fileServerPath))
	//http.Handle("/repos/", http.StripPrefix("/repos", httpfs))
	elon.PathPrefix("/repos/").Handler(http.StripPrefix("/repos/", http.FileServer(http.Dir(fileServerPath))))
	fmt.Println("API endpoints /repos downloading.")

	//http.Handle("/swagger/", http.StripPrefix("/swagger/", fs))

	//httpDucky := ":80"
	//httpDucky = "localhost:8080"
	httpDucky = HttpURL
	duckyapi.HttpURL = HttpURL
	listenMsg := "Listening on http://" + httpDucky + " ..."
	fmt.Println(listenMsg)
	//log.Println(listenMsg)
	log.Fatal(http.ListenAndServe(httpDucky, elon))

	//
	//
	//

	color.Set(color.FgMagenta)
	emoji.Println("\n   NKS  :duck: Ducky  started")
	color.Unset()

} // End main
