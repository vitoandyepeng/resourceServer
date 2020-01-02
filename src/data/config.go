package data

import (
	"common/utils"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

var Config *ServerConfig

type ServerConfig struct {
	RunMode     string `json:"run_mode"`
	Host        string `json:"host"`
	Port        string `json:"port"`
	PrivateKey  string `json:"private_key"`
	Size        int    `json:"size"`
	UserRange   string `json:"user_range"`
	ClubRange   string `json:"club_range"`
	AssignRange string `json:"assign_range"`

	Type       map[string]bool   `json:"type"`
	IconRange  map[int]SizeRange `json:"icon_range"`
	StaticPath string            `json:"static_path"`
}

type SizeRange struct {
	ColMin int `json:"col_min"`
	ColMax int `json:"col_max"`
	RowMin int `json:"row_min"`
	RowMax int `json:"row_max"`
}

func Init() bool {
	Config = new(ServerConfig)
	Config.Type = make(map[string]bool)
	Config.StaticPath = "icon"
	Config.Type["jpg"] = true
	Config.Type["jpeg"] = true
	Config.Type["png"] = true
	err := utils.LoadJsonConfig("config.json", Config)
	if err != nil {
		utils.WErr("Load config json file err,please check !", err)
		return false
	}

	Config.IconRange = make(map[int]SizeRange)
	userRange, err := getRange(Config.UserRange)
	if err != nil {
		utils.WErr("Load config get user image range err!", err.Error())
		return false
	}
	Config.IconRange[USER] = userRange

	clubRange, err := getRange(Config.ClubRange)
	if err != nil {
		utils.WErr("Load config get club image range err!", err.Error())
		return false
	}
	Config.IconRange[CLUB] = clubRange

	assignRange, err := getRange(Config.AssignRange)
	if err != nil {
		utils.WErr("Load config get assign image range err!", err.Error())
		return false
	}
	Config.IconRange[ASSIGN] = assignRange

	return true
}

func getRange(str string) (item SizeRange, err error) {
	list := strings.Split(str, ";")
	if len(list) != 2 {
		utils.WErr("load config getRange err.", str)
		err = errors.New("range str err")
		return
	}
	col := strings.Split(list[0], ",")
	if len(col) != 2 {
		utils.WErr("load config getRange col err.", str)
		err = errors.New("range str err")
		return
	}
	row := strings.Split(list[1], ",")
	if len(row) != 2 {
		utils.WErr("load config getRange row err.", str)
		err = errors.New("range str err")
		return
	}
	item.ColMin, err = strconv.Atoi(col[0])
	if err != nil {
		utils.WErr("load config col str 2 int err.", err.Error(), str)
		return
	}
	item.ColMax, err = strconv.Atoi(col[1])
	if err != nil {
		utils.WErr("load config col str 2 int err.", err.Error(), str)
		return
	}

	item.RowMin, err = strconv.Atoi(row[0])
	if err != nil {
		utils.WErr("load config row str 2 int err.", err.Error(), str)
		return
	}
	item.RowMax, err = strconv.Atoi(row[1])
	if err != nil {
		utils.WErr("load config col str 2 int err.", err.Error(), str)
		return
	}
	return
}
