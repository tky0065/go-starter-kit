package main

import (
	_ "embed"
	"strings"
)

var (
	//go:embed assets/graphql_generated.go.tpl
	graphQLGeneratedTemplate string

	//go:embed assets/graphql_models_gen.go.tpl
	graphQLModelsGenTemplate string
)

func encodeGraphQLModulePath(module string) string {
	replacer := strings.NewReplacer(
		"-", "ᚑ",
		"/", "ᚋ",
		".", "ᚐ",
	)
	return replacer.Replace(module)
}

// GraphQLGeneratedTemplate returns a pre-generated gqlgen executable schema.
// It is scaffolded up front so the generated project can run immediately.
func (t *ProjectTemplates) GraphQLGeneratedTemplate() string {
	replacer := strings.NewReplacer(
		"__MODULE__", t.projectName,
		"__ENCODED_MODULE__", encodeGraphQLModulePath(t.projectName),
	)
	return replacer.Replace(graphQLGeneratedTemplate)
}

// GraphQLModelsGenTemplate returns the gqlgen-generated root GraphQL types.
func (t *ProjectTemplates) GraphQLModelsGenTemplate() string {
	return graphQLModelsGenTemplate
}
