module github.com/FrancoLiberali/cql/cqllint/pkg/analyzer/testdata

go 1.22

require (
	appearance v0.0.1
	github.com/FrancoLiberali/cql v0.0.1
	not_concerned v0.0.1
	repeated v0.0.1
)

require (
	github.com/elliotchance/pie/v2 v2.8.0 // indirect
	github.com/google/uuid v1.5.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	golang.org/x/exp v0.0.0-20231219180239-dc181d75b848 // indirect
	gorm.io/gorm v1.25.5 // indirect
)

replace not_concerned => ./src/not_concerned

replace repeated => ./src/repeated

replace appearance => ./src/appearance

replace github.com/FrancoLiberali/cql => ../../../..
