package dbscylla

import (
    "errors"
    "github.com/gocql/gocql"
    _ "github.com/gocql/gocql"
    . "github.com/harishb2k/easy-go/db"
    . "github.com/harishb2k/easy-go/errors"
)

type Context struct {
    HostList []string
    Keyspace string
    Cluster  *gocql.ClusterConfig
    *gocql.Session
}

func (context *Context) InitDatabase() (err error) {

    context.Cluster = gocql.NewCluster(context.HostList...)
    context.Cluster.Keyspace = context.Keyspace

    context.Session, err = context.Cluster.CreateSession()
    if err != nil {
        return &ErrorObj{
            Name: "failed_to_create_db_session",
            Err:  err,
        }
    }
    return
}

func (context *Context) ensureSession() (err error) {
    if context.Session == nil || context.Session.Closed() {
        return &ErrorObj{
            Name: "db_session_is_null",
            Err:  errors.New("session is not created"),
        }
    }
    return
}

func (context *Context) FindAll(queryString string, mapper RowMapper, val ...interface{}) (result []interface{}, e error) {

    // Ensure we have a session before we make any call
    if e = context.ensureSession(); e != nil {
        return
    }

    // Make a iterator
    iterator := context.Query(queryString, val...).Consistency(gocql.LocalQuorum).Iter()

    // Make a array with default size and fill it
    results := make([]interface{}, 0, 4)
    for ; ; {
        if item, err := mapper.Map(iterator); err != nil {
            return nil, err
        } else if item == nil {
            return results, nil
        } else {
            results = append(results, item)
        }
    }
}

func (context *Context) FindOne(queryString string, mapper RowMapper, val ...interface{}) (result interface{}, e error) {

    // Ensure we have a session before we make any call
    if e = context.ensureSession(); e != nil {
        return
    }

    query := context.Query(queryString, val...).Consistency(gocql.LocalQuorum);
    return mapper.Map(query)
}

func (context *Context) Persist(sql string, val ...interface{}) (e error) {

    // Ensure we have a session before we make any call
    if e = context.ensureSession(); e != nil {
        return
    }

    if err := context.Query(sql, val...).Exec(); err != nil {
        return &ErrorObj{
            Name: "failed_to_persist_to_db",
            Err:  err,
        }
    }
    return
}
