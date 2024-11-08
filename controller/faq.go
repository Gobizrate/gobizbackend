package controller

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gocroot/config"
	"github.com/gocroot/helper/at"

	"github.com/gocroot/helper/atdb"
	"github.com/gocroot/helper/watoken"
	"github.com/gocroot/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetAllFAQ(respw http.ResponseWriter, req *http.Request) {
	data, err := atdb.GetAllDoc[[]model.FAQ](config.Mongoconn, "faq", bson.M{})
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Data FAQ tidak ditemukan"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusForbidden, respn)
		return
	}

	if len(data) == 0 {
		var respn model.Response
		respn.Status = "Error: Data FAQ tidak ditemukan"
		at.WriteJSON(respw, http.StatusNotFound, respn)
		return
	}

	at.WriteJSON(respw, http.StatusOK, data)
}

func PostFAQ(respw http.ResponseWriter, req *http.Request) {
	payload, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(req))
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Token Tidak Valid"
		respn.Info = at.GetSecretFromHeader(req)
		respn.Location = "Decode Token Error"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusForbidden, respn)
		return
	}

	var faq model.FAQ
	if err := json.NewDecoder(req.Body).Decode(&faq); err != nil {
		var respn model.Response
		respn.Status = "Error: Bad Request"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	dataInsert := model.FAQ{
		Nama:      faq.Nama,
		Email:     faq.Email,
		Handphone: faq.Handphone,
		Questions: faq.Questions,
	}

	dataFaq, err := atdb.InsertOneDoc(config.Mongoconn, "faq", dataInsert)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Gagal Insert Database"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusNotModified, respn)
		return
	}
	response := map[string]interface{}{
		"status":  "success",
		"message": "FAQ berhasil ditambahkan",
		"name":    payload.Alias,
		"data":    dataFaq,
	}

	at.WriteJSON(respw, http.StatusOK, response)
}

func DeleteFAQByID(respw http.ResponseWriter, req *http.Request) {
	payload, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(req))
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Token Tidak Valid"
		respn.Info = at.GetSecretFromHeader(req)
		respn.Location = "Decode Token Error"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusForbidden, respn)
		return
	}

	at.WriteJSON(respw, http.StatusOK, payload)
}

func UpdateFAQ(respw http.ResponseWriter, req *http.Request) {
	payload, err := watoken.Decode(config.PublicKeyWhatsAuth, at.GetLoginFromHeader(req))
	if err != nil {
		var respn model.Response
		respn.Status = "Error: Token Tidak Valid"
		respn.Info = at.GetSecretFromHeader(req)
		respn.Location = "Decode Token Error"
		respn.Response = err.Error()
		at.WriteJSON(respw, http.StatusForbidden, respn)
		return
	}

	// faqID := r.URL.Query().Get("id")
	// if faqID == "" {
	// 	var respn model.Response
	// 	respn.Status = "Error: ID FAQ tidak ditemukan"
	// 	at.WriteJSON(respw, http.StatusBadRequest, respn)
	// 	return
	// }

	at.WriteJSON(respw, http.StatusOK, payload)
}

func GetFAQByID(respw http.ResponseWriter, req *http.Request) {
	faqID := req.URL.Query().Get("id")
	if faqID == "" {
		var respn model.Response
		respn.Status = "Error: ID FAQ tidak ditemukan"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(faqID)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: ID FAQ tidak valid"
		at.WriteJSON(respw, http.StatusBadRequest, respn)
		return
	}

	var faq model.FAQ
	filter := bson.M{"_id": objectID}
	err = config.Mongoconn.Collection("faq").FindOne(context.TODO(), filter).Decode(&faq)
	if err != nil {
		var respn model.Response
		respn.Status = "Error: FAQ tidak ditemukan"
		at.WriteJSON(respw, http.StatusNotFound, respn)
		return
	}
	
	response := map[string]interface{}{
		"status":  "success",
		"message": "FAQ ditemukan",
		"data":    faq,
	}

	at.WriteJSON(respw, http.StatusOK, response)
}