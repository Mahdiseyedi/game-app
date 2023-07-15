package category

type Category string

const (
	FootballCategory = "football"
	HistoryCategory  = "history"
	SinceCategory    = "since"
	MathCategory     = "math"
	EnglishCategory  = "english"
	TechCategory     = "Tech"
	FoodCategory     = "Food"
	CinemaCategory   = "cinema"
	MusicCategory    = "music"
)

func (c Category) IsValid() bool {
	switch c {
	case FootballCategory:
		return true
	case HistoryCategory:
		return true
	case SinceCategory:
		return true
	case MathCategory:
		return true
	case EnglishCategory:
		return true
	case TechCategory:
		return true
	case FoodCategory:
		return true
	case CinemaCategory:
		return true
	case MusicCategory:
		return true
	}

	return false
}

func CategoryList() []Category {
	return []Category{FootballCategory,
		HistoryCategory,
		SinceCategory,
		MathCategory,
		EnglishCategory,
		TechCategory,
		FoodCategory,
		CinemaCategory,
		MusicCategory,
	}
}
