package examples

import (
    "github.com/gocql/gocql"
    "github.com/google/uuid"
    "github.com/harishb2k/easy-go/db"
    "github.com/harishb2k/easy-go/dbscylla"
)

import (
    "fmt"
    "github.com/harishb2k/gox-errors"
)

// CREATE KEYSPACE test_me WITH replication = {'class':'SimpleStrategy', 'replication_factor':1};
// use test_me ;
// CREATE TABLE users ( id text, name text, age int,  primary key (id) );

type User struct {
    Id   string
    Name string
    Age  int
}

func ScyllaMain() {
    scyllaExample()
}

func scyllaExample() (err error) {

    var context db.IDb
    context = &dbscylla.Context{
        Keyspace: "test_me",
        HostList: []string{"127.0.0.1"},
    }
    if err := context.InitDatabase(); err != nil {
        return errors.Wrap(err, "Failed to init database")
    }

    uid := uuid.New().String()
    if err := context.Persist(
        "INSERT INTO users (id, name, age) VALUES(?, ?, ?)",
        uid,
        "user_name",
        30,
    ); err != nil {
        return errors.Wrap(err, "Failed to persist")
    }

    if result, err := context.FindOne("SELECT id, name, age FROM users WHERE id=?", &internalRowMapper{}, uid); err != nil {
        return errors.Wrap(err, "Failed to select")
    } else {
        fmt.Println(result)
    }

    return
}

type internalRowMapper struct {
}

func (internalRowMapper) Map(row interface{}) (result interface{}, err error) {
    user := User{}
    if query, ok := row.(*gocql.Query); ok {
        if err = query.Scan(&user.Id, &user.Name, &user.Age); err != nil {
            return
        }
    } else if itr, ok := row.(*gocql.Iter); ok {
        if itr.Scan(&user.Id, &user.Name, &user.Age) == false {
            return
        }
    }
    result = &user
    return
}
