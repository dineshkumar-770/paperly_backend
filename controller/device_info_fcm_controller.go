package controller

import (
	"encoding/json"
	"mongo_api/database"
	"mongo_api/models"
	"mongo_api/response"
	"net/http"
)

type DeviceInfoAndFCM struct {
}

var myDatabseInst = database.DataBase{}

func (dv *DeviceInfoAndFCM) SaveUSerDeviceInfoWithFCM(w http.ResponseWriter, r *http.Request) {
	myDatabseInst.InitDataBase()
	resp := response.SuccessResponse{
		Status: "Failed",
	}
	var deviceInfo models.DeviceInfo
	json.NewDecoder(r.Body).Decode(&deviceInfo)
	_, err := myDatabseInst.SaveDeviceInfo(deviceInfo)
	if err != nil {
		w.WriteHeader(500)
		resp.Message = "Unable to Save User Info"
		resp.Data = err
		_ = json.NewEncoder(w).Encode(resp)
		return
	}

	w.WriteHeader(200)
	resp.Status = "Success"
	resp.Message = "Device Info Saved SuccessFully"
	resp.Data = nil
	json.NewEncoder(w).Encode(resp)

}
