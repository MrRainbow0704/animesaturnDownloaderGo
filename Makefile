NAME = animesaturn-downloader
PACKAGE = github.com/MrRainbow0704/animesaturnDownloaderGo
VERSION = $(shell cat ./version.txt)
SRC_DIR = ./cmd/animesaturn-downloader
SRC_DIR_GUI = ./cmd/animesaturn-downloader-gui
END_DIR = ./bin
END_DIR_GUI = ../../bin
LDFLAGS = -ldflags="-X '$(PACKAGE)/internal/version.version=$(VERSION)'"

.PHONY: clean build 

all:
	make cli
	make gui

cli:
	make linux-cli
	make win-cli
	make mac-cli

gui:
	make linux-gui
	make win-gui
	make mac-gui

linux:
	make linux-cli
	make linux-gui

win:
	make win-cli
	make win-gui

mac:
	make mac-cli
	make mac-gui

linux-cli: export GOOS=linux
linux-cli: export GOARCH=amd64
linux-cli: 
	go build -o $(END_DIR)/$(NAME)-$(VERSION)-linux $(LDFLAGS) $(SRC_DIR)

linux-gui: export GOOS=linux
linux-gui: export GOARCH=amd64
linux-gui:
	cd $(SRC_DIR_GUI) && wails build -tags webkit2_41 -o $(END_DIR_GUI)/$(NAME)-$(VERSION)-linux-gui $(LDFLAGS)

win-cli: export GOOS=windows
win-cli: export GOARCH=amd64
win-cli:
	go build -o $(END_DIR)/$(NAME)-$(VERSION)-windows.exe $(LDFLAGS) $(SRC_DIR)

win-gui: export GOOS=windows
win-gui: export GOARCH=amd64
win-gui:
	cd $(SRC_DIR_GUI) && wails build -o $(END_DIR_GUI)/$(NAME)-$(VERSION)-windows-gui.exe $(LDFLAGS)

mac-cli: export GOOS=darwin
mac-cli: export GOARCH=amd64
mac-cli:
	go build -o $(END_DIR)/$(NAME)-$(VERSION)-darwin $(LDFLAGS) $(SRC_DIR)

mac-gui: export GOOS=darwin
mac-gui: export GOARCH=amd64
mac-gui:
	cd $(SRC_DIR_GUI) && wails build -tags webkit2_41 -o $(END_DIR_GUI)/$(NAME)-$(VERSION)-darwin-gui $(LDFLAGS)

