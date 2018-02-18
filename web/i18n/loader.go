package i18n

// Loader i18n loader
type Loader interface {
	Get(l, c string) (string, error)
	Langs() ([]string, error)
	All(l string) (map[string]string, error)
}
