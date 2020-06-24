package db

type IDb interface {
    InitDatabase() (err error)
    Persist(queryString string, val ...interface{}) (e error)
    FindOne(queryString string, mapper RowMapper, val ...interface{}) (result interface{}, e error)
    FindAll(queryString string, mapper RowMapper, val ...interface{}) (result []interface{}, e error)
}

type RowMapper interface {
    Map(row interface{}) (result interface{}, err error)
}
