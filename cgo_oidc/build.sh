go build -buildmode=c-archive tokenverify.go
gcc example.c  ./tokenverify.a -o example -l pthread

