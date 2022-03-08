package service

import (
    "fmt"
    "reflect"

    "around/backend"
    "around/constants"
    "around/model"

    "github.com/olivere/elastic/v7"
)

func CheckUser(username, password string) (bool, error) {  //去数据库读一下，看用户是否存在
    query := elastic.NewBoolQuery() // select * from userTable where username = xxx AND password = xxx
    query.Must(elastic.NewTermQuery("username", username))
    query.Must(elastic.NewTermQuery("password", password))
    searchResult, err := backend.ESBackend.ReadFromES(query, constants.USER_INDEX)
    if err != nil {
        return false, err
    }

    var utype model.User
    for _, item := range searchResult.Each(reflect.TypeOf(utype)) {
        u := item.(model.User)
        if u.Password == password {
            fmt.Printf("Login as %s\n", username)
            return true, nil
        }
    }
    return false, nil
}

func AddUser(user *model.User) (bool, error) { //插入新用户
    query := elastic.NewTermQuery("username", user.Username) //先做一个query操作看user存在不存在数据库里
    searchResult, err := backend.ESBackend.ReadFromES(query, constants.USER_INDEX)
    if err != nil {
        return false, err
    }

    if searchResult.TotalHits() > 0 {
        return false, nil
    }

    err = backend.ESBackend.SaveToES(user, constants.USER_INDEX, user.Username)
    if err != nil {
        return false, err
    }
    fmt.Printf("User is added: %s\n", user.Username)
    return true, nil
}



