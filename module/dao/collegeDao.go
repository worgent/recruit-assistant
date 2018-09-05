package entity

import "fmt"
import "goBoss/module/entity"

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	cf "goBoss/config"
)

//教育信息
type CollegeDao struct {
	//查询学校信息
}

//dump
func (c *CollegeDao) FindCollege(name string) (entity.College, error) {
	college := entity.College{}

	db, err := sql.Open("mysql", cf.RConfig.MysqlConnectStr)
	if err != nil {
		log.Println(err)
		return college, err
	}
	err = db.Ping()
	if err != nil {
		log.Println(err)
		return college, err

	}
	//在这里进行一些数据库操作
	rows, err := db.Query(
		"select id,name,is211,is985,createType,specialType,belong,educationType,batch from college where name = ? ", name)
	if err != nil {
		fmt.Println(err)
		return college, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&college.Id, &college.Name, &college.Info211,
			&college.Info985, &college.CreateType, &college.SpecialType,
			&college.Belong, &college.EducationType, &college.Batch)
		if err != nil {
			fmt.Println(err)
			return college, err
		}
	}
	err = rows.Err()
	if err != nil {
		fmt.Println(err)
		return college, err
	}
	//fmt.Println("name:", url, "age:", description)

	defer db.Close()

	return college, nil
}
