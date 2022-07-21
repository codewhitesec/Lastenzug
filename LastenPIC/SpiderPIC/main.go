package main

import (

	"bufio"
	b64 "encoding/base64"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

  spiderPIC "SpiderPIC/src"

)

func main() {

  fName := flag.String("asm", "", "Assembly file")
  fNameOut := flag.String("o", "", "Outputfile")
  disphelp := flag.Bool("help", false, "display this message")
  silent := flag.Bool("silent", false, "do not diplay logo and stuff")
  polyfactor := flag.Int("pf", 1, "Polymorphism factor [1-10]")
  flag.Parse()

  if len(os.Args) < 2 || *disphelp {
		help()
	}

  if *silent == false {
    logo()
  }

  file, err := os.Open(*fName)
  if err != nil {
    spiderPIC.LogFatal(err.Error())
  }

  defer file.Close()

  rand.Seed(time.Now().UTC().UnixNano())

  spiderPIC.LogInfo("Parsing file ... ")
  scanner := bufio.NewScanner(file)
  output := ""
  for scanner.Scan() {
    
    line := strings.TrimSpace(scanner.Text())
    if strings.HasPrefix(line, ".") || strings.HasSuffix(line, ":") {
      if *silent == false {
        spiderPIC.LogInfo("Ignoring: %s", line)
      }
      
      output += scanner.Text() + "\n"
      continue
    }

    components := strings.Split(line, "\t") 
    if len(components) == 0 {
      output += scanner.Text() + "\n"
      continue
    }

    if spiderPIC.RoleDice(*polyfactor){ /* Randomly decide if this line should be mutated */

      substition := spiderPIC.SubstituteInstruction(components) /* Prefering substition */
      if substition != "" {
        output += substition
      } else {

          switch rand.Intn(2) {
            case 0:
              output += spiderPIC.RandomInstructions()
            case 1:
              output += spiderPIC.Trash()
          }

          output += scanner.Text() + "\n"

      }
    } else {
      output += scanner.Text() + "\n"
    }
  }

  outFile, err := os.Create(*fNameOut)
  if err != nil {
    spiderPIC.LogFatal(err.Error())
    os.Exit(1)
  }

  defer outFile.Close()

  outFile.WriteString(output)
  spiderPIC.LogSuccess("Done")
  spiderPIC.LogSuccess("Check: " + *fNameOut)

}

func help() {

  flag.PrintDefaults()
  os.Exit(1)

}

func logo() {

  logoB64 := "CgogICBfX19fXyAgICAgICBfICAgICBfICAgICAgICAgIF9fX19fXyBfX19fXyBfX19fXyAKICAvICBfX198ICAgICAoXykgICB8IHwgICAgICAgICB8IF9fXyBcXyAgIF8vICBfXyBcCiAgXCBgLS0uIF8gX18gIF8gIF9ffCB8IF9fXyBfIF9ffCB8Xy8gLyB8IHwgfCAvICBcLwogICBgLS0uIFwgJ18gXHwgfC8gX2AgfC8gXyBcICdfX3wgIF9fLyAgfCB8IHwgfCAgICAKICAvXF9fLyAvIHxfKSB8IHwgKF98IHwgIF9fLyB8ICB8IHwgICAgX3wgfF98IFxfXy9cCiAgXF9fX18vfCAuX18vfF98XF9fLF98XF9fX3xffCAgXF98ICAgIFxfX18vIFxfX19fLwogICAgICAgIHwgfCAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICB8X3wgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCgoK"
  logoDec, _ := b64.StdEncoding.DecodeString(logoB64)

  fmt.Print(string(logoDec))

}


