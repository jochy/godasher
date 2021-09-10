DIR=$PWD

echo "Building plugins"
export GO111MODULE=auto

echo "---- Building Iframe plugin"
cd $DIR/src/plugins/iframe || exit 1
go build -ldflags="-s -w" -buildmode=plugin -o $DIR/plugins/iframe.so

echo "---- Building Jenkins plugin"
cd $DIR/src/plugins/jenkins || exit 1
go build -ldflags="-s -w" -buildmode=plugin -o $DIR/plugins/jenkins.so

echo "---- Building Health check plugin"
cd $DIR/src/plugins/health || exit 1
go build -ldflags="-s -w" -buildmode=plugin -o $DIR/plugins/health.so

echo "Plugins built"
