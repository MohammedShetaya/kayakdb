package fixtures

type When struct {
	*Common
}

func (w *When) Then() *Then {
	return &Then{Common: w.Common}
}
