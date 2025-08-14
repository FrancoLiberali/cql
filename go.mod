module github.com/ditrit/badaas

go 1.18

require (
	github.com/Masterminds/semver/v3 v3.1.1
	github.com/ditrit/verdeter v0.4.0
	github.com/elliotchance/pie/v2 v2.7.0
	github.com/google/go-cmp v0.5.9
	github.com/google/uuid v1.3.0
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/magiconair/properties v1.8.7
	github.com/noirbizarre/gonja v0.0.0-20200629003239-4d051fd0be61
	github.com/spf13/cobra v1.7.0
	github.com/spf13/viper v1.16.0
	github.com/stretchr/testify v1.8.4
	go.uber.org/fx v1.19.3
	go.uber.org/zap v1.24.0
	golang.org/x/crypto v0.9.0
	gorm.io/driver/mysql v1.5.1
	gorm.io/driver/postgres v1.5.2
	gorm.io/driver/sqlite v1.5.2
	gorm.io/driver/sqlserver v1.5.1
	gorm.io/gorm v1.25.2-0.20230610234218-206613868439
	gotest.tools v2.2.0+incompatible
)

require github.com/felixge/httpsnoop v1.0.1 // indirect

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/felixge/httpsnoop v1.0.1 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/go-sql-driver/mysql v1.7.0 // indirect
	github.com/golang-sql/civil v0.0.0-20220223132316-b832511892a9 // indirect
	github.com/golang-sql/sqlexp v0.1.0 // indirect
	github.com/goph/emperror v0.17.2 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.3.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-sqlite3 v1.14.17 // indirect
	github.com/microsoft/go-mssqldb v1.1.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/pelletier/go-toml/v2 v2.0.8 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/sirupsen/logrus v1.9.2 // indirect
	github.com/spf13/afero v1.9.5 // indirect
	github.com/spf13/cast v1.5.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/objx v0.5.0 // indirect
	github.com/subosito/gotenv v1.4.2 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/dig v1.17.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/exp v0.0.0-20220321173239-a90fa8a75705 // indirect
	golang.org/x/sys v0.8.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// TODO agregar comentario de que los de los dialectors se pueden sacar cuando se mergeen los prs
// gorm por el contrario no porque el de los joins no se va a hacer pr porque es muy dificil hacer que ande en todos los casos
replace gorm.io/driver/postgres v1.5.2 => github.com/ditrit/postgres v0.0.0-20230906140800-b3d5f9d4b6ad

replace gorm.io/driver/sqlite v1.5.2 => github.com/ditrit/sqlite v0.0.0-20230906140046-2f37a3f972de

replace gorm.io/driver/sqlserver v1.5.1 => github.com/ditrit/sqlserver v0.0.0-20230908120642-af1820b994f4

replace gorm.io/gorm => github.com/ditrit/gorm v0.0.0-20230912092052-cfff75e01a3a
