package engine

import (
	"net/http"
	"io/ioutil"
	"github.com/pquerna/ffjson/ffjson"
	"models"
	"resultor"
	"github.com/jacoblai/httprouter"
	"strconv"
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

	if err := d.Engine.Create(&user).Error; err != nil {
		resultor.RetErr(w, err.Error())
		return
	}

	resultor.RetOk(w, &user)
}

func (d *DbEngine) GetUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	defer r.Body.Close()

	var users []models.User
	if err := d.Engine.Find(&users).Error; err != nil {
		resultor.RetErr(w, err.Error())
		return
	}

	resultor.RetOk(w, &users)
}

func (d *DbEngine) PutUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	defer r.Body.Close()
	uid := ps.ByName("uid")
	if uid == "" {
		resultor.RetErr(w, "params err")
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		resultor.RetErr(w, err.Error())
		return
	}
	if len(body) == 0 {
		resultor.RetErr(w, "body null err")
		return
	}
	var user map[string]interface{}
	err = ffjson.Unmarshal(body, &user)
	if err != nil {
		resultor.RetErr(w, err.Error())
		return
	}
	res := d.Engine.Model(&models.User{}).Where("id = ?", uid).Updates(user)
	if res.Error != nil {
		resultor.RetErr(w, res.Error.Error())
		return
	}

	resultor.RetChanges(w, res.RowsAffected)
}

func (d *DbEngine) GetUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	defer r.Body.Close()
	uid := ps.ByName("uid")
	id, err := strconv.ParseInt(uid,10,64)
	if err != nil {
		resultor.RetErr(w, "params err")
		return
	}

	var user models.User
	res := d.Engine.First(&user, id)
	if res.Error != nil {
		resultor.RetErr(w, res.Error.Error())
		return
	}

	resultor.RetOk(w, res)
}

func (d *DbEngine) DelUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	defer r.Body.Close()
	uid := ps.ByName("uid")
	id, err := strconv.ParseInt(uid,10,64)
	if err != nil {
		resultor.RetErr(w, "params err")
		return
	}

	res := d.Engine.Where("id = ?", id).Delete(&models.User{})
	if res.Error != nil {
		resultor.RetErr(w, res.Error.Error())
		return
	}

	resultor.RetChanges(w, res.RowsAffected)
}