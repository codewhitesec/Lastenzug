package SpiderPIC

import (
	"math/rand"
	"strings"
)

var jumpOvers = [...]string {
  "\tjmp [LABEL]\n[BLOB]\t[LABEL]:\n",
} 

var trashInstructions = [...]string {
  "\tpush [REG]\n",
  "\tpop [REG]\n",
  "\tmov [REG], [REG1]\n",
  "\tinc [REG]\n",
  "\tdec [REG]\n",
  "\txor [REG], [REG1]\n",
  "\tmov [REG], [NUM]\n",
  "\tadd [REG], [NUM]\n",
  "\tsub [REG], [NUM]\n",
  "\tneg [REG]\n",
  "\tcdq\n",
}

func Trash() string {

  blob := ""
  num := rand.Intn(6) + 4

  LogInfo("Adding %d trashinstructions", num)

  for i:= 0; i < num; i++ {

    inst := trashInstructions[rand.Intn(len(trashInstructions))]
    inst = RandomValues(inst)

    blob += inst

  }

  trash := strings.ReplaceAll(jumpOvers[rand.Intn(len(jumpOvers))], "[LABEL]", randomLabel()) 
  trash = strings.ReplaceAll(trash, "[BLOB]", blob) 

  _, randomSimpleInstructions := RandomSimpleInstructions()
  trash = strings.ReplaceAll(trash, "[MORE]", randomSimpleInstructions) 

  return trash

}
