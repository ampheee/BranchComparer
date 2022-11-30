builder =go build

all: mainBuild

mainBuild: branchComparer.go
	$(builder)	branchComparer.go