package Token_Manage

import (
	"sync"
	"time"
)

type Token struct {
	address        string    //客户端地址
	login_time     time.Time //登录时间
	user_id        int       //用户id
	generate_token string    //生成token
	primary        int       //权限等级
}

type Token_Manage struct {
	sync.Mutex                  //锁
	mapToken   map[string]Token //token映射
}

var ins_Token_Manage *Token_Manage

// 初始化
func Instance() *Token_Manage {
	if ins_Token_Manage == nil {
		ins_Token_Manage = new(Token_Manage)
		ins_Token_Manage.Map_Token_Init()
	}
	return ins_Token_Manage
}

func (obj *Token_Manage) Map_Token_Init() {
	obj.mapToken = make(map[string]Token)
}

func (obj *Token_Manage) Add(address, token string, primary, userid int) {
	obj.Lock()
	defer obj.Unlock()
	obj.mapToken[token] = Token{
		address:        address,
		login_time:     time.Now(),
		generate_token: token,
		primary:        primary,
		user_id:        userid,
	}
}

func (obj *Token_Manage) Remove(token string) {
	obj.Lock()
	defer obj.Unlock()
	delete(obj.mapToken, token)
}
func (obj *Token_Manage) Verify_Token(token string) (bool, Token) {
	obj.Lock()
	defer obj.Unlock()
	v, ok := obj.mapToken[token]
	return ok, v
}

func (obj *Token) Get_User_Id() int {
	return obj.user_id
}

func (obj *Token) Get_Primary() int {
	return obj.primary
}

func (obj *Token) Get_Address() string {
	return obj.address
}
