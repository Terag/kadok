package graphql

type Query struct {
	Query     string            `json:"query"`
	Variables map[string]string `json:"variables"`
}
