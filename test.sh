echo "Step 2: Testing..."
export GOPATH=$(cd .. && pwd)
echo "Setting GOPATH to " $GOPATH

echo "Testing navm" && \
go test github.com/drellem2/navm
    
