#!/usr/bin/env bash

SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ]; do # resolve $SOURCE until the file is no longer a symlink
  DIR="$( cd -P "$( dirname "$SOURCE" )" && pwd )"
  SOURCE="$(readlink "$SOURCE")"
  [[ $SOURCE != /* ]] && SOURCE="$DIR/$SOURCE" # if $SOURCE was a relative symlink, we need to resolve it relative to the path where the symlink file was located
done
DIR="$( cd -P "$( dirname "$SOURCE" )" && pwd )"

echo "Cross building, outputs are in $DIR/build"
GOOS="linux windows darwin"
GOARCH="386 amd64 arm"
rm -rf build/ tmp/
mkdir -p build

for OS in $GOOS
do
    for ARCH in $GOARCH
    do
        cd $DIR
        if GOOS=$OS GOARCH=$ARCH go build >& /dev/null
            then
                echo "Building for $OS on $ARCH..."
                # Create temp dir
                mkdir -p tmp/build/$OS/$ARCH
                # Move statics and config in this temp dir
                cp -R static tmp/build/$OS/$ARCH
                cp configs/base.ini tmp/build/$OS/$ARCH/config.ini
                # Build and move to temp dir

                mv mewpipe* tmp/build/$OS/$ARCH/ # Can't use go build -o, it's don't append .exe if windows.

                #Compress !
                cd tmp/build/$OS/$ARCH
                tar -cjf $DIR/build/mewpipe_${OS}_${ARCH}.tar.bz2 *
                echo "Compress build/mewpipe_${OS}_${ARCH}.tar.bz2"
                echo "---------------------------"
        fi
    done
done
rm -rf tmp/