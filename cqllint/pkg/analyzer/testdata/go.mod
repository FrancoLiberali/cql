module github.com/FrancoLiberali/cql/cqllint/pkg/analyzer/testdata

go 1.18

require (
	not_concerned v0.0.1
	repeated v0.0.1
	appearance v0.0.1
)

replace not_concerned => ./src/not_concerned

replace repeated => ./src/repeated

replace appearance => ./src/appearance
