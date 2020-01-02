package data

type Req struct {
	Type  int    `json:"type" form:"type"`   // 1 玩家 2 俱乐部 3 指定图片名
	Id    int    `json:"id" form:"id"`       // 玩家或俱乐部id.当type=3时id=3
	Name  string `json:"name" form:"name"`   // 指定图片名字
	Index int    `json:"index" form:"index"` // 图片索引 1 玩家固定为0 2 俱乐部0-6
	Md5   string `json:"md5" form:"md5"`     // PrivateKey+id
	File  string `json:"file" form:"file"`   // 图片内容
}

//{"file":"1","id":1,"type:1","md5":"1", "index":0}
