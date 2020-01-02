package handler

import (
	"bytes"
	"common/utils"
	"data"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"net/url"
	"os"
	"strings"
)


func UploadIcon(c *gin.Context) {
	var req data.Req
	req = c.MustGet("data").(data.Req)

	if req.Type != data.USER && req.Type != data.CLUB && req.Type != data.ASSIGN{
		utils.WErr("UploadIcon type err.", req.Type)
		Echo(c, http.StatusBadRequest, "")
		return
	}

	file, err := url.QueryUnescape(req.File)
	if err != nil {
		utils.WErr("UploadIcon unescape err.", err.Error())
		Echo(c, http.StatusBadRequest, "")
		return
	}

	list := strings.Split(file, ";")
	if len(list) != 2 {
		utils.WErr("UploadIcon file err.")
		Echo(c, http.StatusBadRequest, "")
		return
	}
	l := strings.Split(list[0], "/")
	if len(list) != 2 {
		utils.WErr("UploadIcon file type data err.", list[0])
		Echo(c, http.StatusBadRequest, "")
		return
	}
	//iconType := l[1]
	//if _, ok := data.Config.Type[iconType]; !ok {
	//	utils.WErr("UploadIcon file type not right.", iconType, data.Config.Type)
	//	Echo(c, http.StatusBadRequest, "")
	//	return
	//}

	l = strings.Split(list[1], ",")
	if len(list) != 2 {
		utils.WErr("UploadIcon file data err.", list[0])
		Echo(c, http.StatusBadRequest, "")
		return
	}

	iconString := l[1]
	iconByte, err := base64.StdEncoding.DecodeString(iconString)
	fmt.Println(hex.Dump(iconByte[1:22]))
	if err != nil {
		utils.WErr("UploadIcon file decodes err.", err.Error())
		Echo(c, http.StatusBadRequest, "")
		return
	}

	if len(iconByte) > data.Config.Size {
		utils.WErr("UploadIcon file too much err.", len(iconByte), data.Config.Size)
		Echo(c, http.StatusBadRequest, "")
		return
	}

	var name string
	if req.Type == data.USER || req.Type == data.CLUB {
		name = fmt.Sprintf("U_%d_%d.png", req.Id+731, req.Index)
	} else if req.Type == data.ASSIGN {
		name = req.Name
	}

	buff := bytes.NewBuffer(iconByte)
	var img image.Image

	if strings.Index(hex.Dump(iconByte[1:22]), "PNG") >= 0 {
		img, err = png.Decode(buff)
	} else if strings.Index(hex.Dump(iconByte[1:22]), "JFIF") >= 0 ||
		strings.Index(hex.Dump(iconByte[1:22]), "JPEG") >= 0||
		strings.Index(hex.Dump(iconByte[1:22]), "JPE") >= 0 ||
		strings.Index(hex.Dump(iconByte[1:22]), "JIF") >= 0 {
		img, err = jpeg.Decode(buff)

	} else {
		utils.WErr("UploadIcon image type err.")
		Echo(c, http.StatusBadRequest, "")
		return
	}
	if err != nil {
		utils.WErr("UploadIcon decode image err.", err.Error())
		Echo(c, http.StatusBadRequest, "")
		return
	}
	col := img.Bounds().Dx()
	row := img.Bounds().Dy()
	if !checkSize(req.Type, col, row) {
		utils.WErr("UploadIcon image size err.", col, row, data.Config.IconRange[req.Type])
		Echo(c, http.StatusBadRequest, "")
		return
	}

	filepath := fmt.Sprintf("%s/static/icon/%s", data.Config.Host, name)
	f, err := os.Create(filepath)
	if err != nil {
		utils.WErr("UploadIcon open file err.", err.Error())
		Echo(c, http.StatusBadRequest, "")
		return
	}
	defer f.Close()

	err = png.Encode(f, img)
	if err != nil {
		utils.WErr("UploadIcon png err.", err.Error())
		Echo(c, http.StatusBadRequest, "")
		return
	}

	filePath := fmt.Sprintf("%s/%s", data.Config.StaticPath, name)
	Echo(c, http.StatusOK, filePath)
	return
}

func Echo(c *gin.Context, code int, str string)  {
	c.Writer.Header().Set("Content-type", "application/text")
	c.Writer.WriteHeader(code)
	n, err := c.Writer.Write([]byte(str))
	if err != nil {
		utils.WErr("Echo err", n, err.Error())
	}
}

func checkSize(kind int, col, row int) bool {
	item, ok := data.Config.IconRange[kind]
	if !ok {
		return false
	}
	if !(item.ColMin <= col && col <= item.ColMax) {
		return false
	}
	if !(item.RowMin <= row && row <= item.RowMax) {
		return false
	}
	return true
}

