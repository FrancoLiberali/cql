package test

import (
	"gorm.io/gorm"

	"github.com/FrancoLiberali/cql"
	"github.com/FrancoLiberali/cql/test/conditions"
	"github.com/FrancoLiberali/cql/test/models"
)

type SelectIntTestSuite struct {
	testSuite
}

func NewSelectIntTestSuite(
	db *gorm.DB,
) *SelectIntTestSuite {
	return &SelectIntTestSuite{
		testSuite: testSuite{
			db: db,
		},
	}
}

// TODO hacer lo mismo para los selects del groupby

func (ts *SelectIntTestSuite) TestSelect() {
	ts.createProduct("1", 1, 0, false, nil)
	ts.createProduct("2", 1, 1, false, nil)
	ts.createProduct("5", 0, 2, false, nil)

	results, err := cql.Select(
		cql.Query[models.Product](
			ts.db,
		),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *ResultInt) {
			result.Int = int(value)
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []ResultInt{
		{Int: 42},
		{Int: 42},
		{Int: 42},
	}, results)
}

// TODO
// function
// multiple
// joined

// func (ts *QueryIntTestSuite) TestASD() {
// 	// tipo group by
// 	// products, err := cql.Query[models.Product](
// 	// 	ts.db,
// 	// 	conditions.Product.Int.Is().Eq(cql.Int(1)),
// 	// ).Ascending(conditions.Product.Float).Offset(1).Limit(1).Select(
// 	// 	"asd",
// 	// 	"asd2",
// 	// ).Into(
// 	// 	&algo
// 	// )
// 	// // poblema es que depende de los nombres de las cosas

// 	// // poniendo hacia donde quiero el select
// 	// products, err = cql.Query[models.Product](
// 	// 	ts.db,
// 	// 	conditions.Product.Int.Is().Eq(cql.Int(1)),
// 	// ).Ascending(conditions.Product.Float).Offset(1).Limit(1).
// 	// 	Select(
// 	// 		"asd", &algo1
// 	// 	).Select(
// 	// 		"asd2", &algo2,
// 	// 	)
// 	// // problema los metodos no pueden ser de generics, por lo que no se puede chequear el tipo de algo

// 	// cql.Select(
// 	// 	"asd", &algo1,
// 	// 	cql.Query[models.Product](
// 	// 		ts.db,
// 	// 		conditions.Product.Int.Is().Eq(cql.Int(1)),
// 	// 	).Ascending(conditions.Product.Float).Offset(1).Limit(1),
// 	// )
// 	// // problema como pongo multiples

// 	// // con alguna estructura
// 	// cql.Select(
// 	// 	cql.Query[models.Product](
// 	// 		ts.db,
// 	// 		conditions.Product.Int.Is().Eq(cql.Int(1)),
// 	// 	).Ascending(conditions.Product.Float).Offset(1).Limit(1),
// 	// 	cql.Algo[T1]{sql: "asd", value: &algo1},
// 	// 	cql.Algo[T2]{sql: "asd", value: &algo1},
// 	// )
// 	// // no es necesario que sea funcion pero es muy feo

// 	// cql.Query[models.Product](
// 	// 	ts.db,
// 	// 	conditions.Product.Int.Is().Eq(cql.Int(1)),
// 	// ).Ascending(conditions.Product.Float).Offset(1).Limit(1).Select(
// 	// 	cql.Algo("asd").Into(&algo),
// 	// 	cql.Algo("asd2").Into(&algo2),
// 	// )
// 	// creo que es mas lindo

// 	// y las listas que onda?
// 	type Container struct {
// 		algo1 string
// 		algo2 float32
// 	}

// 	a := Container{}
// 	println(a)

// 	b := &a.algo1
// 	println(b)

// 	// a partir de esto devolver una lista de containers?
// 	// enviando una funcion?
// 	func()

// 	// que select devuelva tambien el query
// 	// cql.Select(
// 	// 	"asd", &algo1,
// 	// 	cql.Select(
// 	// 		"asd", &algo1,
// 	// 		cql.Query[models.Product](
// 	// 			ts.db,
// 	// 			conditions.Product.Int.Is().Eq(cql.Int(1)),
// 	// 		).Ascending(conditions.Product.Float).Offset(1).Limit(1),
// 	// 	),
// 	// )
// 	// // queda horrible

// 	// // insert select
// 	// cql.Insert[models.Product](
// 	// 	conditions.Product.Bool, "select ...."
// 	// )
// 	// // pero como seteo multiples

// 	// cql.Insert[models.Product](
// 	// 	cql.Query[models.Product](
// 	// 		ts.db,
// 	// 		conditions.Product.Int.Is().Eq(cql.Int(1)),
// 	// 	).Select(
// 	// 		cql.Algo[T1]{sql: "asd", value: &conditions.Product.Bool},
// 	// 		cql.Algo[T2]{sql: "asd", value: &algo1},
// 	// 	),
// 	// )
// 	// // ponele que si pero no se va a poder porque conditions.Product.Bool es de otro tipo que nada que ver
// 	// // tendria que ser otra estructura o campo, ya depende del usuario que lo use bien, pero puedo hacerlo de alguna forma

// 	// cql.Insert[models.Product](
// 	// 	cql.Query[models.Product](
// 	// 		ts.db,
// 	// 		conditions.Product.Int.Is().Eq(cql.Int(1)),
// 	// 	),
// 	// 	cql.Algo2[T1]{sql: "asd", value: &conditions.Product.Bool},
// 	// 	cql.Algo2[T2]{sql: "asd", value: &algo1},
// 	// )

// 	// cql.Insert[models.Product](
// 	// 	cql.Query[models.Product](
// 	// 		ts.db,
// 	// 		conditions.Product.Int.Is().Eq(cql.Int(1)),
// 	// 	),
// 	// 	conditions.Product.Bool.Set("asd"),
// 	// 	conditions.Product.Float.Set(conditions.Brand.UpdatedAt),
// 	// )
// 	// // esta esta buena, despues tengo que tener un validador en cqllint que se fije:
// 	// // 1. que estan todos los campos en los inserts (que el null explicito sea obligatorio)
// 	// // 2. que las cosas de la derecha esten en la query

// 	// // metodo select que devulve una estructura con metodos?
// 	// .Select().Bool("select ...")
// 	// // vuelvo a tener el mismo problema del multiple

// 	// ts.Require().NoError(err)

// 	// EqualList(&ts.Suite, []*models.Product{product2}, products)
// }
