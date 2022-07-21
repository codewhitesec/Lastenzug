package SpiderPIC

import (
  "fmt"
	"github.com/fatih/color"
)

func LogFatal(f string, v ...interface{}){
  red := color.New(color.Bold, color.FgRed).PrintfFunc()
	red("[*] ")
	fmt.Printf(f+"\n", v...)
}

func LogInfo(f string, v ...interface{}){
  yellow := color.New(color.Bold, color.FgYellow).PrintfFunc()
	yellow("[*] ")
	fmt.Printf(f+"\n", v...)
}

func LogSuccess(f string, v ...interface{}){
  green := color.New(color.Bold, color.FgGreen).PrintfFunc()
	green("[*] ")
	fmt.Printf(f+"\n", v...)
}
