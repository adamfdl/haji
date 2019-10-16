package domain

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

var (
	ContextKeyUser = contextKey("merchant")
)
