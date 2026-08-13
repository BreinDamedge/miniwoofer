package impl

type Doc struct {
	Id  string
	Tok []string
}

type Search struct {
	Id    string
	Score float64
}

type Index interface {
	Append(id string, name []string) error
	Search(query []string) ([]Search, error)
}
