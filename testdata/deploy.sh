#!/bin/bash

# set default param
CMDNAME=`basename $0`
CWD=$(cd $(dirname $0) && pwd)
SOURCE_HOME="testdata/ptune"
TARGET_HOME="/tmp/ptune"

# get command option
OPT=
while getopts s:t: OPT
do
    case $OPT in
    s) SOURCE_HOME=$OPTARG
        ;;
    t) TARGET_HOME=$OPTARG
        ;;
    \?)
        echo "Usage" 1>&2
        echo "$CMDNAME [-s testdata/ptune] [-t /tmp/ptune]" 1>&2
        exit 1
        ;;
    esac
done
shift `expr $OPTIND - 1`

# build Go tools
go build -o testdata/stubcmd testdata/stubcmd.go
go build -o testdata/getperfsoap testdata/getperfsoap.go
go build -o testdata/getperf cmd/getperf/main.go

# copy module
mkdir -p $TARGET_HOME
(cd $SOURCE_HOME; cp -r * $TARGET_HOME)
mkdir -p $TARGET_HOME/bin
cp -r testdata/stubcmd $TARGET_HOME/bin/
cp -r testdata/getperfsoap $TARGET_HOME/bin/
cp -r testdata/getperf $TARGET_HOME/bin/

exit 0
