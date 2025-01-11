echo "Starting build"
./fmt.sh && ./test.sh
export GOPATH=$(cd .. && pwd) && \
GO111MODULE=on go install github.com/drellem2/navm/cmd/...
