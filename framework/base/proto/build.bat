go run .\build.go  -srcDir=.\ -outDir=.\cmd -protoc=.\protoc.exe
move .\cmd\proto_msg.pb.go .\..\
exit
