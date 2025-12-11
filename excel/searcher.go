package excel

type Searcher interface {
	Headers() []Header
	Search(skip, limit int64) ([]byte, error)
}
