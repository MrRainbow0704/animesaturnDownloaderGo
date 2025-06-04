NAME = animesaturn-downloader
PACKAGE = github.com/MrRainbow0704/animesaturnDownloaderGo
VERSION = $(shell cat ./version.txt)
SRC_DIR = ./cmd/animesaturn-downloader
SRC_DIR_GUI = ./cmd/animesaturn-downloader-gui
SRC_DIR_FRONTEND = ./frontend
END_DIR = ./bin
END_DIR_GUI = ../../bin
LDFLAGS = -ldflags="-X '$(PACKAGE)/internal/version.version=$(VERSION)'"

.PHONY: all cli gui linux win mac linux-cli linux-gui win-cli win-gui mac-cli mac-gui clean

all:
	$(MAKE) cli
	$(MAKE) gui

cli:
	$(MAKE) linux-cli
	$(MAKE) win-cli
	$(MAKE) mac-cli

gui:
	$(MAKE) linux-gui
	$(MAKE) win-gui
	$(MAKE) mac-gui

linux:
	$(MAKE) linux-cli
	$(MAKE) linux-gui

win:
	$(MAKE) win-cli
	$(MAKE) win-gui

mac:
	$(MAKE) mac-cli
	$(MAKE) mac-gui

linux-cli: export GOOS=linux
linux-cli: export GOARCH=amd64
linux-cli: 
	go build -o $(END_DIR)/$(NAME)-$(VERSION)-linux $(LDFLAGS) $(SRC_DIR)

linux-gui: export GOOS=linux
linux-gui: export GOARCH=amd64
linux-gui:
	cd $(SRC_DIR_FRONTEND) && npm run build
	cd $(SRC_DIR_GUI) && wails build -s -tags webkit2_41 -o $(END_DIR_GUI)/$(NAME)-$(VERSION)-linux-gui $(LDFLAGS)

win-cli: export GOOS=windows
win-cli: export GOARCH=amd64
win-cli:
	go build -o $(END_DIR)/$(NAME)-$(VERSION)-windows.exe $(LDFLAGS) $(SRC_DIR)

win-gui: export GOOS=windows
win-gui: export GOARCH=amd64
win-gui:
	cd $(SRC_DIR_FRONTEND) && npm run build
	cd $(SRC_DIR_GUI) && wails build -s -o $(END_DIR_GUI)/$(NAME)-$(VERSION)-windows-gui.exe $(LDFLAGS)

mac-cli: export GOOS=darwin
mac-cli: export GOARCH=amd64
mac-cli:
	go build -o $(END_DIR)/$(NAME)-$(VERSION)-darwin $(LDFLAGS) $(SRC_DIR)

mac-gui: export GOOS=darwin
mac-gui: export GOARCH=amd64
mac-gui:
	cd $(SRC_DIR_FRONTEND) && npm run build
	cd $(SRC_DIR_GUI) && wails build -s -tags webkit2_41 -o $(END_DIR_GUI)/$(NAME)-$(VERSION)-darwin-gui $(LDFLAGS)

clean:
	bash -c "rm -rf $(END_DIR)"
	bash -c "rm -rf $(SRC_DIR_FRONTEND)/dist"