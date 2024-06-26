echo
echo --- Prepare ---
echo [i] Configuring proxy...
go env -w GOPROXY=https://goproxy.cn,direct
echo [i] Upgrading packages...
go get -u ./...
go mod tidy
echo [i] Formatting files...
go fmt -x ./...
echo
echo --- End ---
