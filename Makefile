.PHONY: build build-lin build-win-x86 build-win-x64 build-mac deploy

build: build-lin build-win-x86 build-win-x64 #build-mac

build-lin:
	mkdir -p dist
	go build -o dist/libktanemod-remote-math-interface.so -buildmode=c-shared .

build-win-x86:
	mkdir -p dist
	$(eval A=$(abspath $(lastword $(MAKEFILE_LIST))))
	$(eval B=$(dir $(A)))
	docker run --rm -it -v $(B):/go/work -w /go/work -e GOARCH=386 x1unix/go-mingw:1.18 go build -o dist/ktanemod-remote-math-interface-x86.dll -buildmode=c-shared .

build-win-x64:
	mkdir -p dist
	$(eval A=$(abspath $(lastword $(MAKEFILE_LIST))))
	$(eval B=$(dir $(A)))
	docker run --rm -it -v $(B):/go/work -w /go/work -e GOARCH=amd64 x1unix/go-mingw:1.18 go build -o dist/ktanemod-remote-math-interface-x64.dll -buildmode=c-shared .

build-mac:
	mkdir -p dist
	GOOS=darwin go build -o dist/ktanemod-remote-math-interface.dylib -buildmode=c-shared .

deploy: build
	$(eval A=$(abspath "../ktanemod-remote-math/Assets/Plugins/dlls"))
	mkdir -p $(A)
	mkdir -p $(A)/x86
	mkdir -p $(A)/x86_64
	cp dist/libktanemod-remote-math-interface.so $(A)/libktanemod-remote-math-interface.so
	cp dist/ktanemod-remote-math-interface-x86.dll $(A)/x86/ktanemod-remote-math-interface.dll
	cp dist/ktanemod-remote-math-interface-x64.dll $(A)/x86_64/ktanemod-remote-math-interface.dll
	#cp dist/ktanemod-remote-math-interface.dylib $(A)/ktanemod-remote-math-interface.dylib
