package main

import (
	_ "bytes"
	"fmt"
	_ "image"
	_ "image/jpeg"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	SenderName string
	Message    string
	Time       string
	Date       string
}

func formatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d/%02d/%02d", year, month, day)
}

func formatAsTime(t time.Time) string {
	hour, min, _ := t.Clock()
	if hour > 12 {

		return fmt.Sprintf("%d:%02d:pm", hour-12, min)
	} else {
		return fmt.Sprintf("%d:%02d:am", hour, min)
	}

}

func main() {
	var ips []string
	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)
	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {

			ips = append(ips, fmt.Sprintf("%s", ipv4))

			fmt.Println("IPv4: ", ipv4)
		}
	}

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&Message{})

	// Create
	// db.Create(&Message{SenderName: "khan", Message: "Hello world!", Time: "2:25pm", Date: "11oct2022"})

	// Read
	//   var product Product
	//   db.First(&product, 1) // find product with integer primary key
	//   db.First(&product, "code = ?", "D42") // find product with code D42

	// Update - update product's price to 200
	//   db.Model(&product).Update("Price", 200)
	// Update - update multiple fields
	//   db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // non-zero fields
	//   db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

	// Delete - delete product
	//   db.Delete(&product, 1)

	// github.com/mattn/go-sqlite3
	// db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./assets")

	r.GET("/connect", func(c *gin.Context) {
		//c.Header("Content-Type", "application/json")
		err := os.RemoveAll("/assets/qr") // delete an entire directory
		if err != nil {
			fmt.Println(err)
		}
		c.Header("Access-Control-Allow-Origin", "*")
		//Create a 256x256 PNG image:
		// var png []byte
		// png, _ = qrcode.Encode("https://example.org", qrcode.Medium, 256)
		//Create a 256x256 PNG image and write to a file:
		for i := range ips {
			_ = qrcode.WriteFile(ips[i], qrcode.Medium, 256, fmt.Sprintf("assets/qr/qr%d.png", i))
		}

		// img, s, err := image.Decode(bytes.NewReader(png))
		fmt.Println(err)
		// fmt.Println(s)
		// buf := new(bytes.Buffer)
		// buf2 := new(png.Buffer)
		// original_image, _, err := image.Decode(bytes.NewReader(image_data))
		// err := jpeg.Encode(buf, new_image, nil)
		// send_s3 := buf.Bytes()
		//Create a 256x256 PNG image with custom colors and write to file:
		//err := qrcode.WriteColorFile("https://example.org", qrcode.Medium, 256, color.Black, color.White, "qr.png")
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			// "data": img,
			"ips": ips,
		})
	})

	r.GET("/shutdown", func(c *gin.Context) {
		// db.Find(&user, "id = ?", 10)
		// http.
		// r.Shutdown(context.Background())
		// c.HTML(http.StatusOK, "index2.tmpl", gin.H{
		// 	"data": data,
		// 	"ips":  ips,
		// })
	})

	r.GET("/", func(c *gin.Context) {
		// db.Find(&user, "id = ?", 10)
		var data []Message
		db.Find(&data)
		c.HTML(http.StatusOK, "index2.tmpl", gin.H{
			"data": data,
			"ips":  ips,
		})
	})

	// start api mobile
	r.GET("/getallmsg", func(c *gin.Context) {
		// db.Find(&user, "id = ?", 10)
		var data []Message
		db.Find(&data)
		c.Header("Content-Type", "application/json")
		c.Header("Access-Control-Allow-Origin", "*")
		c.JSON(http.StatusOK, data)
		// c.HTML(http.StatusOK, "index2.tmpl", gin.H{
		// 	"data": data,
		// })
	})

	// end api for mobile
	//*******************************************************************
	// start api mobile
	r.POST("/post", func(c *gin.Context) {

		//router.LoadHTMLFiles("templates/template1.html", "templates/template2.html")
		message := c.PostForm("message")
		sender_name := c.PostForm("sender")
		// Create
		db.Create(&Message{
			SenderName: sender_name,
			Message:    message,
			Time:       formatAsDate(time.Now()),
			Date:       formatAsTime(time.Now())})
		// c.JSON(http.StatusOK, gin.H{
		// 	"message": "pong",
		// })
		c.Redirect(301, "/")
	})
	// end api for mobile

	// start api mobile
	r.POST("/delete", func(c *gin.Context) {
		var requestBody Message

		if err := c.BindJSON(&requestBody); err != nil {
			// DO SOMETHING WITH THE ERROR
		}

		fmt.Println(requestBody.ID)
		fmt.Println(requestBody.SenderName)
		// Delete - delete product
		// var delMsg Message
		//   db.Delete(&delMsg, 1)
		db.Delete(&Message{}, requestBody.ID)
		c.Header("Content-Type", "application/json; charset=utf-8")
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("charset", "utf-8")
		c.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
		//router.LoadHTMLFiles("templates/template1.html", "templates/template2.html")
		// message := c.PostForm("id")
		// sender_name := c.PostForm("sender")
		// Create
		// db.Create(&Message{
		// 	SenderName: sender_name,
		// 	Message:    message,
		// 	Time:       formatAsDate(time.Now()),
		// 	Date:       formatAsTime(time.Now())})
		// c.JSON(http.StatusOK, gin.H{
		// 	"message": "pong",
		// })
		// c.Redirect(301, "/")
	})
	// end api for mobile

	// start api mobile
	r.POST("/postmsg", func(c *gin.Context) {

		//router.LoadHTMLFiles("templates/template1.html", "templates/template2.html")
		var requestBody Message

		if err := c.BindJSON(&requestBody); err != nil {
			// DO SOMETHING WITH THE ERROR
		}

		fmt.Println(requestBody.Message)
		fmt.Println(requestBody.SenderName)

		// c.JSON()
		// message := c.PostForm("message")
		// sender_name := c.PostForm("sender")
		// Create
		db.Create(&Message{
			SenderName: requestBody.SenderName,
			Message:    requestBody.Message,
			Time:       formatAsDate(time.Now()),
			Date:       formatAsTime(time.Now())})

		c.Header("Content-Type", "application/json; charset=utf-8")
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("charset", "utf-8")
		c.JSON(http.StatusOK, gin.H{
			"message": "success",
		})

	})
	// end api for mobile
	r.GET("/index", func(c *gin.Context) {
		// Read
		var message Message
		db.First(&message, 1) // find product with integer primary key
		//db.First(&message, "SenderName = ?", "khan") // find product with code D42
		c.JSON(http.StatusOK, gin.H{
			"message": message,
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
