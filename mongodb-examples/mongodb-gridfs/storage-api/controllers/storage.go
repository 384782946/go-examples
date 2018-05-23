package controllers

import (
	"io"

	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type StorageController struct {
	beego.Controller
}

type fileinfo struct {
	//文件大小
	LENGTH int32
	//md5
	MD5 string
	//文件名
	FILENAME string
}

// @Title Get
// @Description delete the object
// @Param	filename		query 	string	true		"The objectId you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 objectId is empty
// @router / [get]
func (c *StorageController) Get() { //下载
	filename := c.GetString("filename")

	session, err := mgo.Dial("localhost:32769")
	defer session.Close()
	datas := map[string]string{}

	if err != nil {
		datas["errorString"] = "can not open mongodb " + err.Error()
		c.Data["json"] = datas
		c.ServeJSON()
		return
	}

	gfs := session.DB("test").GridFS("fs")
	gfsFile, err := gfs.Open(filename)
	if err != nil {
		datas["errorString"] = "can not open file " + filename + ",  " + err.Error()
		datas["filename"] = filename
		c.Data["json"] = datas
		c.ServeJSON()
		return
	}

	c.Ctx.Output.Header("Content-type", "application/jpeg")
	c.Ctx.Output.Header("Content-Disposition", "attachment; filename=test.jpeg")

	buff := make([]byte, 1024)

	for {
		n, err := gfsFile.Read(buff)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if 0 == n {
			break
		}
		c.Ctx.ResponseWriter.Write(buff)
	}
}

// @Title GetAll
// @Description get all files in database
// @Success 200 {string} query success!
// @Failure 403 connot find files
// @router /GetAll [get]
func (c *StorageController) GetAll() { //列出所有文件
	session, err := mgo.Dial("localhost:32769")
	defer session.Close()
	datas := map[string]string{}

	if err != nil {
		datas["errorString"] = "can not open mongodb " + err.Error()
		c.Data["json"] = datas
		c.ServeJSON()
		return
	}

	gfs := session.DB("test").GridFS("fs")
	iter := gfs.Find(nil).Iter()

	result := new(fileinfo)

	filenames := ""
	for iter.Next(&result) {
		filenames += result.FILENAME + "  [" + result.MD5 + "]  "
	}

	datas["files"] = filenames
	c.Data["json"] = datas
	c.ServeJSON()
}

// @Title Delete
// @Description delete the object
// @Param	uploadFile			formData 	file	true		"The file you want upload"
// @Success 200 {string} delete success!
// @Failure 403 objectId is empty
// @router /upload [post]
func (c *StorageController) Upload() { //上传
	fromFile, _, _ := c.GetFile("uploadFile")

	session, err := mgo.Dial("localhost:32769")
	defer session.Close()
	datas := map[string]string{}

	if err != nil {
		datas["errorString"] = "can not open mongodb " + err.Error()
		c.Data["json"] = datas
		c.ServeJSON()
		return
	}

	gfs := session.DB("test").GridFS("fs")
	gfsFile, err := gfs.Create("test")
	if err != nil {
		datas["errorString"] = "can not create file " + ",  " + err.Error()
		c.Data["json"] = datas
		c.ServeJSON()
		return
	}

	//c.Ctx.Output.Header("Content-type", "application/jpeg")
	//c.Ctx.Output.Header("Content-Disposition", "attachment; filename=test.jpeg")

	io.Copy(gfsFile, fromFile)

	err = gfsFile.Close()

	fromFile.Close()

	if err != nil {
		datas["errorString"] = "close gridfs error " + ",  " + err.Error()
		c.Data["json"] = datas
		c.ServeJSON()
		return
	}

	// buff := make([]byte, 1024)

	// for {
	// 	fromFile.Read()
	// 	n, err := gfsFile.Read(buff)
	// 	if err != nil && err != io.EOF {
	// 		panic(err)
	// 	}
	// 	if 0 == n {
	// 		break
	// 	}
	// 	c.Ctx.ResponseWriter.Write(buff)
	// }
}

// func (c *StorageController) Put() [ //更新

// ]

// @Title Delete
// @Description delete the object
// @Param	filename		query 	string	true		"The file you want remove"
// @Success 200 {string} delete success!
// @Failure 403 objectId is empty
// @router /delete [delete]
func (c *StorageController) Delete() { //删除
	filename := c.GetString("filename")
	session, err := mgo.Dial("localhost:32769")
	defer session.Close()
	datas := map[string]string{}

	if err != nil {
		datas["errorString"] = "can not open mongodb " + err.Error()
		c.Data["json"] = datas
		c.ServeJSON()
		return
	}

	gfs := session.DB("test").GridFS("fs")
	//filename:"/Users/zhangxiaojian/Downloads/WechatIMG31.jpeg"
	n, err1 := gfs.Find(bson.M{"filename": filename}).Count()
	if n == 0 || err1 != nil {
		datas["errorString"] = "file not exist, " + filename //+ ",  " + err1.Error()
		datas["filename"] = filename
		c.Data["json"] = datas
		c.ServeJSON()
		return
	}

	err = gfs.Remove(filename)
	if err != nil {
		datas["errorString"] = "can not remove file " + filename + ",  " + err.Error()
		datas["filename"] = filename
		c.Data["json"] = datas
		c.ServeJSON()
		return
	}

	datas["ret"] = "success remove file" + filename
	c.Data["json"] = datas
	c.ServeJSON()
}
