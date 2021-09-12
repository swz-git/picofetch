package main

import (
	"bytes"
	"embed"
	"fmt"
	"math"
	"regexp"
	"strings"
	"text/template"

	"github.com/zcalusic/sysinfo"
)

const VERSION = "0.1.1"

//go:embed os-ansi
var distroLogos embed.FS

//go:embed infotemplate.txt
var infotemplate string

func horiJoin(str1, str2 string) string {
	removecolor, _ := regexp.Compile(`.\[(\d{1,3};?)+m`)

	str1Lines := strings.Split(str1, "\n")
	var str1Biggest float64 = 0
	for _, v := range str1Lines {
		str1Biggest = math.Max(float64(len(removecolor.ReplaceAllString(v, ""))), str1Biggest)
	}
	str2Lines := strings.Split(str2, "\n")
	result := ""
	for i := 0; i < int(math.Max(float64(len(str1Lines)), float64(len(str2Lines)))); i++ {
		line1 := ""
		if i < len(str1Lines) {
			line1 = str1Lines[i]
		}
		line2 := ""
		if i < len(str2Lines) {
			line2 = str2Lines[i]
		}
		line1pad := strings.Repeat(" ", int(str1Biggest)-len(removecolor.ReplaceAllString(line1, ""))+2)
		result += line1 + line1pad + line2 + "\n"
	}
	return result
}

func main() {
	colorReset := "\033[0m"

	// colorRed := "\033[31m"
	colorGreen := "\033[32m"
	// colorYellow := "\033[33m"
	// colorBlue := "\033[34m"
	// colorPurple := "\033[35m"
	colorCyan := "\033[36m"
	// colorWhite := "\033[37m"

	bold := "\033[1m"

	var si sysinfo.SysInfo
	si.GetSysInfo()

	infopart := ""

	logoFile, err := distroLogos.ReadFile(fmt.Sprintf("os-ansi/%s.ansi", si.OS.Vendor))

	if err != nil {
		logoFile, _ = distroLogos.ReadFile("os-ansi/linux.ansi")
	}

	logo := string(logoFile)

	infopart += colorCyan + bold + "picofetch" + colorGreen + " " + VERSION + colorReset + "\n"

	infotemplate = strings.ReplaceAll(infotemplate, "\n", "\n"+colorGreen+bold)
	infotemplate = strings.ReplaceAll(infotemplate, ": ", ":"+colorReset+" ")

	ut, err := template.New("users").Parse(infotemplate)
	if err != nil {
		panic(err)
	}
	var tpl bytes.Buffer
	err = ut.Execute(&tpl, si)
	if err != nil {
		panic(err)
	}

	infopart += tpl.String()[1:]

	fmt.Print(horiJoin(logo, infopart))
}
