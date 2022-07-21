package kommando

// Kommando defines the potential kommandos of a Ladung
type Kommando int

// Kommando supported
const (
	Unbekannt Kommando = iota
	Data
	Stop
)

// mapping of kommandos
var kommandos = map[Kommando]string{
	Data: "data",
	Stop: "stop",
}

// Name returns the string representation of a kommando
func Name(k Kommando) string {
	return kommandos[k]
}

// Code returns the Kommando based on
func Code(i int) Kommando {
	candidate := Kommando(i)
	if _, ok := kommandos[candidate]; ok {
		return candidate
	}
	return Unbekannt
}
