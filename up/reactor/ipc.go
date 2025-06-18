package reactor

const (
	Signal_End uint8 = iota
	Signal_Ready
	Signal_Cede
)

type Msg struct {
	Signal uint8
	PId    string // process node Id
}
