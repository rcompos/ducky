# nks-devops-blackduck

This repository contains tools for scanning code with Synopsis BlackDuck.  Includes a CLI app and a REST API service.

Specify git repo and optional version as account/repo:version as input argument.  If no tag is supplied, the latest  

# Black Duck Server URL

Specify the Black Duck server URL as an environment variable (BLACKDUCK_URL) or command-line arg (-u).

```
⇒ export BLACKDUCK_URL="<blackduck_url>"
```

# Black Duck User Access Token

A Black Duck user access token is required for all usages.  You can generate a new user access token in your Black Duck user profile settings.

Supply the Black Duck server user access token as an environment variable (BD_TOKEN) or command-line arg (-t).

```
⇒ export BD_TOKEN="<blackduck_user_auth_token>"
```

# NKS Ducky CLI

CLI Usage:

```
⇒ go run cmd/ducky/main.go -h
Supply single argument with Git account and repository name separated by a slash. i.e. netapp/myrepo
  -d	Debugging output
  -f	Force without confirmation prompt
  -h	Help
  -l string
    	Input file
  -r string
    	Git repo (account/repo:version) to clone, prep and scan (required) (envvar DUCKY_GIT_REPO)
  -s	Perform Black Duck scan (default true)
  -t string
    	Black Duck auth token (required) (envvar BD_TOKEN)
exit status 1
```

Run on command-line with out building binary:
Scan single repo:

```
⇒ go run cmd/ducky/main.go -s -r netapp/<repo_name>:<optional_tag>
```

Specify repo list from file:

```
⇒ go run cmd/ducky/main.go -s -l repos.txt
```


Alternatley, build a binary and run scans:

```
⇒ cd cmd/ducky
⇒ go build -o ducky

⇒ ./ducky -s -r netapp/chandler:v1.4.0
```

# NKS Ducky API

Build API service binary

```
⇒ cd cmd/duckyapi
⇒ go build -o duckyapi

⇒ ./duckyapi
```

Alternately, run API service binary without building binary

```
⇒ go run cmd/duckyapi/main.go
```

Browse API endpoints:

```
⇒ curl http://localhost:8080/api/v1
```

Scan repo:

```
⇒ curl http://localhost:8080/api/v1/scan/netapp/<repo-name>:<optional_tag>
```

# NKS Ducky Docker

Build API service binary for Linux operating systems on AMD64 architecture.

Build Docker image from repo top-level directory
```
⇒ docker build -t duckyapi .
```

Run Docker image
```
⇒ docker run -it --rm -p 8080:8080 duckyapi
```
