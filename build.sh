#! /usr/bin/env bash
set -eo pipefail

trap 'echo -e "\033[33;5mBuild failed on build.sh:$LINENO\033[0m"' ERR

RACE=-race

# Process args
for arg in "$@"
do
    case "$arg" in
        --current-os-only|-c) SKIP_OTHER_OSES=1 ;;
        --skip-linters|-l)    SKIP_LINTERS=1 ;;
        --skip-test|-t)       SKIP_TESTS=1 ;;
        --omit-race|-r)       RACE= ;;
        --fast|-f)            SKIP_OTHER_OSES=1; SKIP_LINTERS=1; SKIP_TESTS=1 ;;
        --help|-h)
            echo "$0 [options]"
            echo "  -f, --fast             Same as -c -l -t -r"
            echo "  -c, --current-os-only  Skip builds for other OSes"
            echo "  -l, --skip-linters     Skip linters"
            echo "  -t, --skip-tests       Skip tests"
            echo "  -r, --omit-race        Omit the -race option in tests"
            echo "  -h, --help             This help text"
            exit 0
            ;;
        *) echo "Invalid argument: $arg"; BAIL=1 ;;
    esac
done
if [ ! -z $BAIL ]; then
    exit 1
fi

# Setup the tools we'll need
TOOLS_DIR=$PWD/tools
mkdir -p $TOOLS_DIR
if [ ! -e $TOOLS_DIR/mkembeddedfs ] || [ "$($TOOLS_DIR/mkembeddedfs -v 2>&1 || true)x" != "1.0.2x" ]; then
    echo -e "\033[33mInstalling version 1.0.2 of mkembeddedfs into $TOOLS_DIR...\033[0m"
    cd $TOOLS_DIR
    git clone --quiet https://github.com/richardwilkes/toolbox
    cd toolbox
    git checkout --quiet v1.20.0
    cd xio/fs/mkembeddedfs
    go build -o ../../../../mkembeddedfs .
    cd ../../../..
    rm -rf toolbox
    cd ..
fi
if [ -z $SKIP_LINTERS ]; then
    if [ ! -e $TOOLS_DIR/golangci-lint ] || [ "$($TOOLS_DIR/golangci-lint version 2>&1 | awk '{ print $4 }' || true)x" != "1.21.0x" ]; then
        echo -e "\033[33mInstalling version 1.21.0 of golangci-lint into $TOOLS_DIR...\033[0m"
        curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $TOOLS_DIR v1.21.0
    fi
fi
export PATH=$TOOLS_DIR:$PATH

# Remove any generated code
find . -iname "*_gen.go" -exec /bin/rm {} \;

# Generate the embedded file system
mkembeddedfs --pkg icons --output icons/efs_gen.go --no-modtime --ignore ".*\\.go" --strip icons icons

# Build the Go code
if [ -z $SKIP_OTHER_OSES ]; then
    for p in darwin linux windows; do
        echo -e "\033[33mBuilding Go code for $p...\033[0m"
        GOOS=$p time go build -v ./...
    done
else
    echo -e "\033[33mBuilding Go code...\033[0m"
    time go build -v ./...
fi

# Run the linters
if [ -z $SKIP_LINTERS ]; then
    echo -e "\033[33mRunning Go linters...\033[0m"
    golangci-lint run
else
    echo -e "\033[33mSkipping Go linters\033[0m"
fi

# Run the tests
if [ -z $SKIP_TESTS ]; then
    echo -e "\033[33mRunning tests...\033[0m"
    go test $RACE ./...
else
    echo -e "\033[33mSkipping tests\033[0m"
fi
