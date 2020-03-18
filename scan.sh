#!/usr/bin/env bash
set -eo pipefail

function usage(){
echo -e "usage: ${0} [path] [release]\n"
  echo -e "[path]    Optional path to scan. Defaults to current directory. Path can be relative or absolute.\n"
  echo -e "[release] Release version that will be set included in the projcet version (defaults to 'current')\n"
}

if [ $# -gt 2 ] || [ "$1" = "-h" ] || [ "$1" = "--help" ]; then
  usage
  exit
fi

command -v realpath &>/dev/null || { echo >&2 "realpath is required to run. Make sure coreutils are installed. Aborting."; exit 1; }
if [ ! -z $1 ]; then
  SCAN_PATH=$( realpath $1 )
else
  SCAN_PATH=$( realpath ./ )
fi

if [ ! -z $2 ]; then
  RELEASE=$2
fi

BLACKDUCK_URL=${BLACKDUCK_URL:-https://blackduck.eng.netapp.com}
OFFLINE_MODE=${OFFLINE_MODE:-false}
CLEANUP=${CLEANUP:-true}
LOGLEVEL=${LOGLEVEL:-INFO} # TRACE,DEBUG,INFO,WARN,ERROR

RELEASE=${RELEASE:-"current"}
PROJECT_NAME=${PROJECT_NAME:-"NKS"}
PROJECT_VER=${PROJECT_VER:-$(basename ${SCAN_PATH})_${RELEASE}}

BD_TOKEN=${BD_TOKEN}
if [ -z "$BD_TOKEN" ]; then
  read -s -p "Enter token: " BD_TOKEN
  echo
fi

# The GO_DEP detoctor seems very broken so we exclude it by default
BD_EXCLUDE=${BD_EXCLUDE:-GO_DEP}

bash <( curl -s https://detect.synopsys.com/detect.sh ) \
  --blackduck.url=${BLACKDUCK_URL} \
  --blackduck.trust.cert=true \
  --blackduck.api.token=${BD_TOKEN} \
  --blackduck.offline.mode=${OFFLINE_MODE} \
  --detect.project.name=\'${PROJECT_NAME}\' \
  --detect.project.version.name=\'${PROJECT_VER}\' \
  --detect.code.location.name=\'${PROJECT_NAME}_${PROJECT_VER}_code\' \
  --detect.bom.aggregate.name=\'${PROJECT_NAME}_${PROJECT_VER}_bom\' \
  --detect.cleanup=${CLEANUP} \
  --detect.source.path=${SCAN_PATH} \
  --detect.parallel.processors=4 \
  --detect.detector.search.depth=50 \
  --detect.detector.search.continue=true \
  --detect.excluded.detector.types=${BD_EXCLUDE} \
  --detect.blackduck.signature.scanner.exclusion.name.patterns='directoryDoesNotExist' \
  --detect.blackduck.signature.scanner.exclusion.pattern.search.depth=50 \
  --detect.blackduck.signature.scanner.memory=8192 \
  --detect.blackduck.signature.scanner.paths=${SCAN_PATH} \
  --detect.output.path=\'/tmp/blackduck\' \
  --logging.level.com.synopsys.integration=${LOGLEVEL}
