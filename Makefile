build:
	go build -o build/timeow-mac

build-mac-app:
	echo "TODO: build with xcodebuild"

copy-binary-to-mac:
	cp build/timeow-mac build/Timeow.app/Contents/MacOS/Timeow

dist-mac: build build-mac-app copy-binary-to-mac
