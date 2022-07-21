package SpiderPIC

import "math/rand"

func RandomRegister() string {
  return registers[rand.Intn(len(registers))]
}

var registers = []string{
  "rax",
  "rcx",
  "rdx",
  "rbx",
  "rsi",
  "rdi",
  "r8",
  "r9",
  "r10",
  "r11",
  "r12",
  "r13",
  "r14",
  "r15",
}

func IsRegister(str string) bool {

   for _, b := range registers {
        if b == str {
            return true
        }
    }

    return false

}
