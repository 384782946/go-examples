/*
利用mongodb-GridFS 实现文件存储，并提供http上传下载接口
*/

package main

import (
	"github.com/astaxie/beego"
)

type HomeController struct {
	beego.Controller
}

func (c *HomeController) Get() {
	c.Ctx.WriteString("hello world")
}

func main() {
	beego.Router("/", &HomeController{})
	beego.Run()
}
