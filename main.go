package main

import (
	"demo-gin/handler"
	"demo-gin/models"
	"encoding/gob"
	"encoding/json"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"net/http"
	"time"
)

// 模拟一些私人数据
var secrets = gin.H{
	"foo":    gin.H{"email": "foo@bar.com", "phone": "123433"},
	"austin": gin.H{"email": "austin@example.com", "phone": "666"},
	"lena":   gin.H{"email": "lena@guapa.com", "phone": "523443"},
}

func main() {

	r := gin.Default()
	// 拦截器 中间件 注册中间件一定要定义在所有请求的前面
	//	r.Use(handler.TestHandler())
	// 禁用控制台颜色,将日志写入文件时不需要控制台颜色
	//gin.DisableConsoleColor()
	//// 日志写入
	//f, _ := os.Create("test.log")
	//gin.DefaultWriter = io.MultiWriter(f)

	// 加载静态网页
	r.LoadHTMLGlob("templates/*")
	// 加载静态资源
	r.Static("/static", "./static")
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "hello gin")
	})

	r.GET("/someJSON", func(c *gin.Context) {
		data := map[string]interface{}{
			"lang": "GO语言",
			"tag":  "<br>",
		}
		c.AsciiJSON(http.StatusOK, data)
	})
	// 跳转静态页面
	r.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"msg": "hello gin",
		})
	})

	r.GET("/user", func(c *gin.Context) {
		user := models.User{Name: "东方不败", Age: 12, Sex: "1"}
		c.JSON(http.StatusOK, user)
	})

	r.GET("/userArr", func(c *gin.Context) {
		users := make([]models.User, 2)
		users[0] = models.User{Name: "东方不败", Age: 12, Sex: "1"}
		users[1] = models.User{Name: "东方不败", Age: 12, Sex: "1"}
		c.JSON(http.StatusOK, users)
	})

	// 在请求路径上获取参数 https://localhost:8088/user/info?userId=1&name=zhangsan
	r.GET("/user/info", func(c *gin.Context) {
		userId := c.Query("userId")
		name := c.Query("name")
		c.JSON(http.StatusOK, gin.H{
			"userId": userId,
			"name":   name,
		})
	})
	// https://localhost:8088/user/info/1/zhangsan
	r.GET("/user/info/:userId/:name", func(c *gin.Context) {
		userId := c.Param("userId")
		name := c.Param("name")
		c.JSON(http.StatusOK, gin.H{
			"userId": userId,
			"name":   name,
		})
	})
	// from 表单形式传参
	r.POST("/user/form", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		c.JSON(http.StatusOK, gin.H{
			"username": username,
			"password": password,
		})
	})

	r.POST("/user/json", func(c *gin.Context) {

		userData, _ := c.GetRawData()
		if err := c.ShouldBindJSON(&userData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		user := map[string]any{}
		json.Unmarshal(userData, &user)
		c.JSON(http.StatusOK, user)
	})
	// 重定向
	r.GET("/redirect", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "https://www.baidu.com")
	})
	// 错误页面
	r.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404.html", gin.H{
			"error": "404 Not Found",
		})
	})
	// 路由组  如果只想在一类业务或某一个请求上加入中间件，gin框架同样可以做到。 在路由组中加中间件
	orderGroup := r.Group("/order", handler.TestHandler())
	orderGroup.GET("/get", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"order": "111",
		})
	})
	orderGroup.GET("/update", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"order": "update",
		})
	})
	// 在某一个请求上加入中间件，中间件作为参数定义在请求接口后，真正的请求之前
	//r.GET("/user", handler.TestHandler(), func(ctx *gin.Context) {
	//	ctx.JSON(http.StatusOK, gin.H{
	//		"msg": "GET",
	//	})
	//})

	// 创建 cookie 存储
	store := cookie.NewStore([]byte("mycookie"))

	// 添加session中间件
	r.Use(sessions.Sessions("mysession", store))
	// 注册用户模型
	gob.Register(models.User{})
	// 登录添加session
	r.POST("/login", func(ctx *gin.Context) {
		// 创建一个session
		session := sessions.Default(ctx)
		// 将用户放入session中
		session.Set("user", models.User{
			Name: "zhangsan",
			Age:  18,
			Sex:  "男",
		})
		// 保存
		session.Save()
	})
	r.POST("/session", func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		// 获取session中信息
		user := session.Get("user")
		ctx.JSON(http.StatusOK, user)
	})

	// 读取网络文件, 返回文件
	r.GET("someDataFromReader", func(c *gin.Context) {

		response, err := http.Get("https://aier-scrm.obs.cn-southwest-2.myhuaweicloud.com/2024/11/04/c4ef7860-b470-4f32-8df8-c09f78eb4f12.jpg")

		if err != nil || response.StatusCode != http.StatusOK {
			c.Status(http.StatusServiceUnavailable)
			return
		}

		reader := response.Body
		contentLength := response.ContentLength
		contentType := response.Header.Get("Content-Type")
		extraHeaders := map[string]string{
			"Content-Disposition": `attachment; filename="gopher.png"`,
		}
		c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)

	})

	// 模型绑定和验证
	r.POST("/loginJSON", func(c *gin.Context) {
		var json models.Login
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if json.Username != "manu" || json.Password != "123" {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})

	})

	// 路由保护
	// 路由组使用 gin.BasicAuth() 中间件
	// gin.Accounts 是 map[string]string 的一种快捷方式
	//authorized := r.Group("/admin", gin.BasicAuth(gin.Accounts{
	//	"foo":    "bar",
	//	"austin": "1234",
	//	"lena":   "hello2",
	//	"manu":   "4321",
	//}))

	// /admin/secrets 端点
	// 触发 "localhost:8080/admin/secrets
	//authorized.GET("/secrets", func(c *gin.Context) {
	//	// 获取用户，它是由 BasicAuth 中间件设置的
	//	user := c.MustGet(gin.AuthUserKey).(string)
	//	if secret, ok := secrets[user]; ok {
	//		c.JSON(http.StatusOK, gin.H{"user": user, "secret": secret})
	//	} else {
	//		c.JSON(http.StatusOK, gin.H{"user": user, "secret": "NO SECRET :("})
	//	}
	//})
	//
	//// Ping handler
	//r.GET("/ping", func(c *gin.Context) {
	//	c.String(200, "pong")
	//})
	//
	//log.Fatal(autotls.Run(r, "example1.com", "example2.com"))

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("bookabledate", bookableDate)
	}

	r.GET("/bookable", func(c *gin.Context) {
		var b Booking
		if err := c.ShouldBindWith(&b, binding.Query); err == nil {
			c.JSON(http.StatusOK, gin.H{"message": "Booking dates are valid!"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	})

	r.Run(":8080")
}

type Booking struct {
	CheckIn  time.Time `form:"check_in" binding:"required,bookabledate" time_format:"2006-01-02"`
	CheckOut time.Time `form:"check_out" binding:"required,gtfield=CheckIn,bookabledate" time_format:"2006-01-02"`
}

var bookableDate validator.Func = func(fl validator.FieldLevel) bool {
	date, ok := fl.Field().Interface().(time.Time)
	if ok {
		today := time.Now()
		if today.After(date) {
			return false
		}
	}
	return true
}
