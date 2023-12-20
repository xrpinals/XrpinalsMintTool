package utils

import "github.com/fatih/color"

var AsciiPicure = `
██╗  ██╗██████╗ ██████╗ ██╗███╗   ██╗ █████╗ ██╗     ███████╗
╚██╗██╔╝██╔══██╗██╔══██╗██║████╗  ██║██╔══██╗██║     ██╔════╝
 ╚███╔╝ ██████╔╝██████╔╝██║██╔██╗ ██║███████║██║     ███████╗
 ██╔██╗ ██╔══██╗██╔═══╝ ██║██║╚██╗██║██╔══██║██║     ╚════██║
██╔╝ ██╗██║  ██║██║     ██║██║ ╚████║██║  ██║███████╗███████║
╚═╝  ╚═╝╚═╝  ╚═╝╚═╝     ╚═╝╚═╝  ╚═══╝╚═╝  ╚═╝╚══════╝╚══════╝

`

var FgWhiteBgGreen = color.New(color.FgWhite, color.BgGreen).SprintFunc()
var FgWhiteBgBlue = color.New(color.FgWhite, color.BgBlue).SprintFunc()
var FgWhiteBgYellow = color.New(color.FgWhite, color.BgYellow).SprintFunc()
var FgWhiteBgRed = color.New(color.FgWhite, color.BgRed).SprintFunc()

var Bold = color.New(color.Bold).SprintFunc()
var BoldGreen = color.New(color.FgGreen, color.Bold).SprintFunc()
var BoldYellow = color.New(color.FgYellow, color.Bold).SprintFunc()
var BoldRed = color.New(color.FgRed, color.Bold).SprintFunc()
