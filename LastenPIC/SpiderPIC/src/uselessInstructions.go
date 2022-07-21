package SpiderPIC

import (
	"math/rand"
)

var uselessInst = [...]string{
	"\tnop\n",
	"\txchg [REG], [REG]\n",
	"\tcmova [REG], [REG]\n",
	"\tcmovb [REG], [REG]\n",
	"\tcmovc [REG], [REG]\n",
	"\tcmove [REG], [REG]\n",
	"\tcmovg [REG], [REG]\n",
	"\tcmovl [REG], [REG]\n",
	"\tpush [REG]\n\tpop [REG]\n",
}

func RandomSimpleInstructions() (int, string) {

	uselessInstructions := ""
	n := rand.Intn(6) + 1

	for i := 0; i <= n; i++ {

		randomSimpleInstruction := uselessInst[rand.Intn(len(uselessInst))]
		randomSimpleInstruction = RandomValues(randomSimpleInstruction)

		uselessInstructions += randomSimpleInstruction

	}

	return n, uselessInstructions

}

func RandomInstructions() string {

	n, ret := RandomSimpleInstructions()
	LogInfo("Adding %d useless instructions", n)

	return ret

}
