package SpiderPIC

import ( 
  "math/rand"
  "strings"
  "strconv"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RoleDice(r int) bool {

  if r >= rand.Intn(r * 10) % 10{
    return true
  }

  return false

}

func RandomValues (inst string) string {

     randomRegister := RandomRegister()
     randomByte := rand.Intn(255)
     randomString := randStringRunes(10)

     inst = strings.ReplaceAll(inst, "[REG]", randomRegister)
     inst = strings.ReplaceAll(inst, "[REG1]", randomRegister)
     inst = strings.ReplaceAll(inst, "[NUM]", strconv.Itoa(randomByte))
     inst = strings.ReplaceAll(inst, "[STR]", randomString)

     return inst

}

func randomLabel() string {
  return randStringRunes(10)
}

func randStringRunes(n int) string {
    b := make([]rune, n)
    for i := range b {
        b[i] = letterRunes[rand.Intn(len(letterRunes))]
    }
    return string(b)
}

