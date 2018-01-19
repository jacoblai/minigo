package engine

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"github.com/pquerna/ffjson/ffjson"
	"models"
	"resultor"
)

func (d *DbEngine) AddUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		resultor.RetErr(w, err.Error())
		return
	}
	if len(body) == 0 {
		resultor.RetErr(w, "body null err")
		return
	}
	var user models.User
	err = ffjson.Unmarshal(body, &user)
	if err != nil{
		resultor.RetErr(w, err.Error())
		return
	}

	//插入数据库
	if err := d.Engine.Create(&user).Error; err != nil {
		resultor.RetErr(w, err.Error())
		return
	}

	resultor.RetOk(w, &user)
}