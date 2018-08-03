set +x;

function makedist {
    echo "Building $1-$2 ...\c"
    GOOS=$1 GOARCH=$2 go build -o koala/koala -ldflags "-X main.version=$3" cmd/koala/main.go
    tar czf artifacts/koala-$3-$1-$2.tar.gz koala/
    rm koala/koala
    echo " done."
}

if [ "$1" == "" ]
then
    echo "Missing version number ... abort!"
    exit 1
fi


echo "Generating artifacts for koala $1 ..."
mkdir -p koala artifacts
cp -r web/dist koala/web

platforms="linux:amd64,arm darwin:amd64"
for pfs in $(echo $platforms | tr " " "\n")
do
    params=$(echo $pfs | tr ":" "\n")
    goos=$(echo $params | awk '{print $1}')
    archs=$(echo $params | awk '{print $2}' | tr "," "\n")
    for goarch in $archs
    do
        makedist $goos $goarch $1
    done
done

rm -rf koala
echo "Done."
