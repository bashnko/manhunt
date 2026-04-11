package query

type QueryAssets struct {
	Perfix        string
	DefQueryParam string
	SearchParam   string
}

func (q *QueryAssets) Youtube() string {
	q.Perfix = "https://youtub.com"
	q.DefQueryParam = "/?search"
	q.SearchParam = "how to dance"
	return q.Perfix + q.DefQueryParam + q.SearchParam
}
