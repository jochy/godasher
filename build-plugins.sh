DIR=$PWD
export GOOS=linux
export GOARCH=amd64

echo "Building plugins"

echo "---- Building Iframe plugin"
cd $DIR/src/plugins/iframe || exit 1
go build -ldflags="-s -w" -buildmode=plugin -o $DIR/plugins/iframe.so

echo "---- Building Jenkins plugin"
cd $DIR/src/plugins/jenkins || exit 1
go build -ldflags="-s -w" -buildmode=plugin -o $DIR/plugins/jenkins.so

echo "---- Building Metric plugin"
cd $DIR/src/plugins/metric || exit 1
go build -ldflags="-s -w" -buildmode=plugin -o $DIR/plugins/metric.so
cp metric.html $DIR/plugins/metric.html

echo "Plugins built"