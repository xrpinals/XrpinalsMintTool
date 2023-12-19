LD_GUI_FLAGS	= -ldflags '-w -s -H windowsgui'
LD_FLAGS	= -ldflags '-w -s'

all: XrpinalsMintTool XrpinalsMintTool-GUI
.PHONY: all

clean: XrpinalsMintTool_clean XrpinalsMintTool-GUI_clean
.PHONY: clean

XrpinalsMintTool:
	go build ${LD_FLAGS} ./cmd/XrpinalsMintTool
.PHONY: XrpinalsMintTool

XrpinalsMintTool-GUI:
	go build ${LD_GUI_FLAGS} ./cmd/XrpinalsMintTool-GUI
.PHONY: XrpinalsMintTool-GUI

XrpinalsMintTool_clean:
	rm -f XrpinalsMintTool
.PHONY: XrpinalsMintTool_clean

XrpinalsMintTool-GUI_clean:
	rm -f XrpinalsMintTool-GUI
.PHONY: XrpinalsMintTool-GUI_clean