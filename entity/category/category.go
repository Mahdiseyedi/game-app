package category

type Category string

const (
	FootballCategory = "Football"
	HistoryCategory  = "History"
)

func (c Category) IsValid() bool {
	switch c {
	case FootballCategory:
		return true
	}

	return false
}
