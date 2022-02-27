#!/usr/bin/env bash

function get_os() {
  unameOut="$(uname -s)"
  case "${unameOut}" in
  Linux*)
    echo -n "linux"
    ;;
  Darwin*)
    echo -n "macos"
    ;;
  CYGWIN*)
    echo -n "cygwin"
    ;;
  MINGW*)
    echo -n "mingw"
    ;;
  *)
    echo "Cannot detect your operating system.  Exiting."
    exit 1
    ;;
  esac
}

os="$(get_os)"

function check_available() {
  which $1 >/dev/null
  if [ $? -ne 0 ]; then
    echo "**** ERROR needed program missing: $1"
    exit 1
  fi
}

check_available 'go'
check_available 'which'

cwd="$(echo "$(pwd)")"
function cleanup() {
  cd "$cwd"
}
# Make sure that we get the user back to where they started
trap cleanup EXIT

function usage() {
  echo "Usage: build.sh [-h|--help] [-c|--clean] [-C|--clean-all]"
  echo "                [-g|--use-go] [-b|--build] [-o|--optimize]"
  echo "                [-w|--path-to-wasm-opt]"
  echo
  echo '    Build wasmdemo.'
  echo
  echo "Arguments:"
  echo "  -h|--help                     This help text"
  echo '  -c|--clean                    Clean generated artifacts.'
  echo "  -C|--clean-all                Clean all the artifacts and the Go module cache."
  echo "  -g|--use-go                   Build using regular go"
  echo "  -o|--optimize                 Run 'wasm-opt' to size optimize the build"
  echo "  -w|--path-to-wasm-opt <path>  Path to the 'wasm-opt' program"
}

clean=0
clean_all=0
tiny_go=true
optimize=0
path_to_wasm_opt=""

while [[ $# -gt 0 ]]; do
  key="$1"

  case $key in
  -h | --help)
    usage
    exit 0
    ;;
  -c | --clean)
    clean=true
    shift
    ;;
  -C | --clean-all)
    clean_all=true
    shift
    ;;
  -o | --optimize)
    optimize=true
    shift
    ;;
  -g | --use-go)
    tiny_go=0
    shift
    ;;
  -w | --path-to-wasm-opt)
    path_to_wasm_opt="$2"
    shift
    shift
    ;;
  *)
    echo "ERROR: unknown argument $1"
    echo
    usage
    exit 1
    ;;
  esac
done

function go_cmd () {
  eval "CGO_ENABLED=0 GOOS=js GOARCH=wasm go build -v -o web/wasm/wasm.wasm ./wasm"
}

function tiny_go_cmd () {
  eval "CGO_ENABLED=0 GOOS=js GOARCH=wasm tinygo build -size full -o web/wasm/wasm.wasm -target wasm ./wasm"
}

go_used="Building with: $(which go)"
build_cmd="go_cmd"
if [ "$tiny_go" = true ]; then
  check_available 'tinygo'
  go_used="Building with: $(which tinygo)"
  build_cmd="tiny_go_cmd"
fi

if [ "$clean_all" = true ]; then
  echo "Deep cleaning..."
  clean=true
  if [ "$tiny_go" = true ]; then
    CGO_ENABLED=0 GOOS=js GOARCH=wasm tinygo clean
  else
    go clean --modcache
    CGO_ENABLED=0 GOOS=js GOARCH=wasm go clean --modcache
  fi
  exit 0
fi

if [ "$clean" = true ]; then
  echo "Regular cleaning..."
	rm -fr ./web/wasm/*.wasm
	rm -f server
  if [ "$tiny_go" = true ]; then
    CGO_ENABLED=0 GOOS=js GOARCH=wasm tinygo clean
  else
    go clean .
    CGO_ENABLED=0 GOOS=js GOARCH=wasm go clean .
  fi
  exit 0
fi

echo "Building server"
go build -v -o server ./cmd/server

echo "Building wasm"
echo "$go_used"
$build_cmd

if [ "$optimize" = true ]; then
  echo "Optimizing..."
  if [ -z "$path_to_wasm_opt" ]; then
    path_to_wasm_opt="$(which wasm-opt)"
  fi
  if [ -z "$path_to_wasm_opt" ]; then
    echo "Could not find the 'wasm-opt' program.  Please supply with the '-w' command"
    exit 2
  fi
  rm -f ./web/wasm/wasm-bak.wasm
	cp ./web/wasm/wasm.wasm ./web/wasm/wasm-bak.wasm
  "$path_to_wasm_opt" web/wasm/wasm-bak.wasm -o web/wasm/wasm.wasm -Oz --strip-dwarf --strip-producers --zero-filled-memory
fi

echo "Done"
