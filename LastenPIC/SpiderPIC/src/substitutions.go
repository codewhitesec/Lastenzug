package SpiderPIC

import (
	"math/rand"
	"strings"
)

var substitutionsPush = [...]string{
	"\tsub rsp, 8\n\tmov [rsp], [REG]\n",
	"\tmov [rsp - 8], [REG]\n\tsub rsp, 8\n",
}

var substitutionsPop = [...]string{
	"\tmov [REG], [rsp]\n\tadd rsp, 8\n",
	"\tadd rsp, 8\n\tmov [REG], [rsp-8]\n",
}

var substitutionsJmp = [...]string{
	"\tpush [WHAT]\n\tret\n",
}

var substitutionsMov = [...]string{
	"\tpush [REG]\n\tpop [REG1]\n",
}

func substitutePush(line []string) string {
	LogInfo("Substituting Push")
	substitution := substitutionsPush[rand.Intn(len(substitutionsPush))]
	return strings.ReplaceAll(substitution, "[REG]", line[1])
}

func substitutePop(line []string) string {
	LogInfo("Substituting Pop")
	substitution := substitutionsPop[rand.Intn(len(substitutionsPop))]
	return strings.ReplaceAll(substitution, "[REG]", line[1])
}

func substituteJmp(line []string) string {

	if len(strings.Split(line[1], " ")) > 2 || !IsRegister(line[1]) {
		return ""
	}

	LogInfo("Substituting Jmp")
	substitution := substitutionsJmp[rand.Intn(len(substitutionsJmp))]

	substitution = strings.ReplaceAll(substitution, "[WHAT]", line[1])

	return strings.ReplaceAll(substitution, "[WHAT]", line[1])

}

func substituteMov(line []string) string {

	if len(strings.Split(line[1], " ")) > 2 {
		return ""
	}

	LogInfo("Substituting Mov")

	regs := strings.Split(line[1], ",")
	regA := strings.TrimSpace(regs[0])
	regB := strings.TrimSpace(regs[1])

	if !IsRegister(regA) || !IsRegister(regB) {
		return ""
	}

	substitution := substitutionsMov[rand.Intn(len(substitutionsMov))]
	substitution = strings.ReplaceAll(substitution, "[REG]", regB)
	substitution = strings.ReplaceAll(substitution, "[REG1]", regA)

	return substitution

}

func SubstituteInstruction(line []string) string {

	switch line[0] {
	case "push":
		return substitutePush(line)
	case "pop":
		return substitutePop(line)
	case "jmp":
		return substituteJmp(line)
	case "mov":
		return substituteMov(line)
	}

	return ""

}
