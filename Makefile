clean:
	rm -rf build/

build:
	go build -o build/timeow-mac

build-mac-app:
	mkdir -p build/app
	cd Timeow/ && xcodebuild
	cp -r Timeow/build/Release/Timeow.app build/app/

copy-binary-to-mac:
	cp build/timeow-mac build/app/Timeow.app/Contents/MacOS/Timeow

create-dmg:
	create-dmg \
	  --volname "Timeow Installer" \
	  --window-pos 200 120 \
	  --window-size 800 400 \
	  --icon-size 100 \
	  --icon "Timeow.app" 200 190 \
	  --hide-extension "Timeow.app" \
	  --app-drop-link 600 185 \
	  "build/Timeow.dmg" \
	  "build/app/"

dist-mac-app: clean build build-mac-app copy-binary-to-mac

dist-mac-dmg: dist-mac-app create-dmg

.PHONY: build
