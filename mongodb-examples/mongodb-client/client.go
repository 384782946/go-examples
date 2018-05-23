package main

/*
 Docs: https://godoc.org/gopkg.in/mgo.v2

 简单mongodb 操作示例
*/

import (
	"fmt"
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	//mongodb服务器地址
	ADDRESS string = "localhost:32768"
)

//
type Person struct {
	Name  string
	Phone string
}

func main() {
	//创建连接
	session, err := mgo.Dial(ADDRESS)
	if err != nil {
		panic(err)
	}

	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	//获取想要操作的数据集
	c := session.DB("test").C("people")

	//清理之前的所有数据
	c.RemoveAll(nil)

	err = c.Insert(&Person{Name: "zxj", Phone: "123123"},
		&Person{Name: "ltf", Phone: "456456"})

	if err != nil {
		log.Fatal(err)
	}

	var result []Person

	//显示所有数据
	err = c.Find(nil).All(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("all datas ###############################")
	for _, per := range result {
		fmt.Printf("Name:%s,Phone:%s\n", per.Name, per.Phone)
	}

	//单条件精确查询
	err = c.Find(bson.M{"name": "zxj"}).All(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("query one ###############################")
	for _, per := range result {
		fmt.Printf("Name:%s,Phone:%s\n", per.Name, per.Phone)
	}

	//单条件正则查询
	err = c.Find(bson.M{"name": bson.M{"$regex": "zx"}}).All(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("regx query one ###############################")
	for _, per := range result {
		fmt.Printf("Name:%s,Phone:%s\n", per.Name, per.Phone)
	}

	//多条件正则查询
	condition := bson.M{"$or": []bson.M{bson.M{"name": "zxj"}, bson.M{"phone": "456456"}}}
	err = c.Find(condition).All(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("or query ###############################")
	for _, per := range result {
		fmt.Printf("Name:%s,Phone:%s\n", per.Name, per.Phone)
	}
}
