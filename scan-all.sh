#!/usr/bin/env bash
set -eo pipefail

# This is pretty ugly, but it works on macOS without requiring coreutils for GNU readlink
# Ideally this would be:
# SCRIPT=$( readlink -f ${BASH_SOURCE[0]} )
# SCRIPTPATH=$( dirname $SCRIPT )
SCRIPT=$( cd "$(dirname ${BASH_SOURCE[0]})" &>/dev/null && pwd )/$( basename ${BASH_SOURCE[0]} )
SCRIPTPATH=$( dirname ${SCRIPT} )
WORKDIR=${SCRIPTPATH}/workdir

mkdir -p ${WORKDIR}
cd ${WORKDIR}

declare -a spc_repos=(
  cluster-api-provider-nks
  dispatch
  hub-backplane
  hub-credentials-server
  hub-workers
  quarter-master-rest
  stemson-api
  stemson-ui
)

declare -a netapp_repos=(
  chandler
  cluster-api
  cluster-api-bootstrap-provider-kubeadm
  cluster-api-provider-vsphere
  hci-nks-vsphere
  vsphere-manager
)

repos=("${spc_repos[@]}" "${netapp_repos[@]}")

echo -e "Cloning repositories"
for repo in "${spc_repos[@]}"; do
  echo -e "Cloning stackpointcloud/${repo}"
  git clone https://github.com/stackpointcloud/${repo}.git || true
done
for repo in "${netapp_repos[@]}"; do
  echo -e "Cloning netapp/${repo}"
  git clone https://github.com/netapp/${repo}.git || true
done

echo -e "Begin scans"
for repo in "${repos[@]}"; do
  echo -e "Scanning ${repo}"
  cd ${WORKDIR}/${repo}
  if [ ${repo} == "cluster-api" ] || [ ${repo} == "cluster-api-provider-vsphere" ] || [ ${repo} == "cluster-api-bootstrap-provider-kubeadm" ]; then
    git checkout netapp
    git reset --hard origin/netapp
  else
    git reset --hard origin/master
  fi

  if [ -f ./go.mod ]; then
    go mod vendor
    export BD_EXCLUDE=GO_DEP,GO_MOD,GRADLE,MAVEN,NPM,PIP,RUBYGEMS,YARN,SWIFT
  elif [ -d ./vendor ]; then
    export BD_EXCLUDE=GO_DEP,GO_MOD,GRADLE,MAVEN,NPM,PIP,RUBYGEMS,YARN,SWIFT
  elif [ -f ./package.json ]; then
    npm install
    export BD_EXCLUDE=GO_DEP,GO_MOD,GRADLE,MAVEN,PIP,RUBYGEMS,YARN,SWIFT
  elif [ -f ./requirements.txt ]; then
    pip install -r requirements.txt
    export BD_EXCLUDE=GO_DEP,GO_MOD,GRADLE,MAVEN,NPM,RUBYGEMS,YARN,SWIFT
  fi

  ${SCRIPTPATH}/scan.sh || true
  echo -e "\nDone\n\n"
done
