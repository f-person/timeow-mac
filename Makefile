build:
	go build -o build/timeow-mac

build-mac-app:
	cd Timeow/ && xcodebuild
	cp -r Timeow/build/Release/Timeow.app build/

copy-binary-to-mac:
	cp build/timeow-mac build/Timeow.app/Contents/MacOS/Timeow

dist-mac: build build-mac-app copy-binary-to-mac

.PHONY: build
