package query

type SQLKeyWord string

const (
	kwSelect SQLKeyWord = "SELECT"
	kwUpdate SQLKeyWord = "UPDATE"
	kwInsert SQLKeyWord = "INSERT"
	kwInto   SQLKeyWord = "INTO"
	kwDelete SQLKeyWord = "DELETE"
	kwFrom   SQLKeyWord = "FROM"
	kwWhere  SQLKeyWord = "WHERE"
	kwCount  SQLKeyWord = "COUNT"
	kwAs     SQLKeyWord = "AS"
	kwOn     SQLKeyWord = "ON"
	kwValues SQLKeyWord = "VALUES"
	kwSet    SQLKeyWord = "SET"
)

// String returns string representation of SQLKeyWord.
func (kw SQLKeyWord) String() string {
	return string(kw)
}
