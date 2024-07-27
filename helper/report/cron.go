package report

import (
	"encoding/base64"
	"errors"
	"net/http"

	"github.com/gocroot/config"
	"github.com/gocroot/helper/at"
	"github.com/gocroot/helper/atapi"
	"github.com/gocroot/helper/atdb"
	"github.com/gocroot/helper/whatsauth"
	"github.com/gocroot/model"
	"github.com/whatsauth/itmodel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func RekapMeetingKemarin(db *mongo.Database) (err error) {
	filter := bson.M{"_id": YesterdayFilter()}
	wagroupidlist, err := atdb.GetAllDistinctDoc(db, filter, "project.wagroupid", "uxlaporan")
	if err != nil {
		return
	}
	if len(wagroupidlist) == 0 {
		return
	}
	for _, gid := range wagroupidlist { //iterasi di setiap wa group
		// Type assertion to convert any to string
		groupID, ok := gid.(string)
		if !ok {
			err = errors.New("wagroupid is not a string")
			return
		}
		filter := bson.M{"wagroupid": groupID}
		var projectDocuments []model.Project
		projectDocuments, err = atdb.GetAllDoc[[]model.Project](db, "project", filter)
		if err != nil {
			return
		}
		for _, project := range projectDocuments {
			var base64pdf, md string
			base64pdf, md, err = GetPDFandMDMeeting(db, project.Name)
			if err != nil {
				return
			}
			dt := &itmodel.DocumentMessage{
				To:        groupID,
				IsGroup:   true,
				Base64Doc: base64pdf,
				Filename:  project.Name + ".pdf",
				Caption:   "Berikut ini rekap rapat kemaren ya kak untuk project " + project.Name,
			}
			_, _, err = atapi.PostStructWithToken[model.Response]("Token", config.WAAPIToken, dt, config.WAAPIDocMessage)
			if err != nil {
				return
			}
			//upload file markdown ke log repo untuk tipe rapat
			if project.RepoLogName != "" {
				// Encode string ke base64
				encodedString := base64.StdEncoding.EncodeToString([]byte(md))

				// Format markdown dengan base64 string
				//markdownContent := fmt.Sprintf("```base64\n%s\n```", encodedString)
				dt := model.LogInfo{
					PhoneNumber: project.Owner.PhoneNumber,
					Alias:       project.Owner.Name,
					FileName:    "README.md",
					RepoOrg:     project.RepoOrg,
					RepoName:    project.RepoLogName,
					Base64Str:   encodedString,
				}
				var conf model.Config
				conf, err = atdb.GetOneDoc[model.Config](db, "config", bson.M{"phonenumber": "62895601060000"})
				if err != nil {
					return
				}
				var statuscode int
				var loginf model.LogInfo
				statuscode, loginf, err = atapi.PostStructWithToken[model.LogInfo]("secret", conf.LeaflySecret, dt, conf.LeaflyURL)
				if err != nil {
					return
				}
				if statuscode != http.StatusOK {
					err = errors.New(loginf.Error)
					return

				}
			}
		}
	}

	return

}

func RekapPagiHari(respw http.ResponseWriter, req *http.Request) {
	var resp model.Response
	filter := bson.M{"_id": YesterdayFilter()}
	wagroupidlist, err := atdb.GetAllDistinctDoc(config.Mongoconn, filter, "project.wagroupid", "pushrepo")
	if err != nil {
		resp.Info = "Gagal Query Distincs project.wagroupid"
		resp.Response = err.Error()
		at.WriteJSON(respw, http.StatusUnauthorized, resp)
		return
	}
	for _, gid := range wagroupidlist { //iterasi di setiap wa group
		// Type assertion to convert any to string
		groupID, ok := gid.(string)
		if !ok {
			resp.Info = "wagroupid is not a string"
			resp.Response = "wagroupid is not a string"
			at.WriteJSON(respw, http.StatusUnauthorized, resp)
			return
		}
		msg, err := GenerateRekapMessageKemarinPerWAGroupID(config.Mongoconn, groupID)
		if err != nil {
			resp.Info = "Gagal Membuat Rekapitulasi perhitungan per wa group id"
			resp.Response = err.Error()
			at.WriteJSON(respw, http.StatusExpectationFailed, resp)
			return
		}
		dt := &whatsauth.TextMessage{
			To:       groupID,
			IsGroup:  true,
			Messages: msg,
		}
		_, resp, err = atapi.PostStructWithToken[model.Response]("Token", config.WAAPIToken, dt, config.WAAPIMessage)
		if err != nil {
			resp.Info = "Tidak berhak"
			resp.Response = err.Error()
			at.WriteJSON(respw, http.StatusUnauthorized, resp)
			return
		}
	}

}
