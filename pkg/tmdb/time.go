package tmdb

import "time"

type TMDBTime struct {
	TimeString time.Time
}

func (t *TMDBTime) UnmarshalJSON(data []byte) error {
	str := string(data)

	if data == nil || str == "" || str == "null" || str == `""` {
		t.TimeString = time.Time{}
		return nil
	}

	parsedTime, err := time.Parse(`"2006-01-02"`, str)
	if err != nil {
		t.TimeString = time.Time{}
		return nil
	}
	t.TimeString = parsedTime
	return nil
}

func (t TMDBTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.TimeString.Format("2006-01-02") + `"`), nil
}

func (t TMDBTime) IsZero() bool {
	return t.TimeString.IsZero()
}

func (t TMDBTime) Year() int {
	return t.TimeString.Year()
}

func (t TMDBTime) ToTime() time.Time {
	return t.TimeString
}
