package models

type DeviceInfo struct {
	Manufacturer          string `json:"manufacturer" bson:"manufacturer"`
	Model                 string `json:"model" bson:"model"`
	Board                 string `json:"board" bson:"board"`
	FCMToken              string `json:"fcm_token" bson:"fcm_token"`
	CurrentAndroidVersion string `json:"current_android_version" bson:"current_android_version"`
	SDKLevel              int    `json:"sdk_level" bson:"sdk_level"`
	Device                string `json:"device" bson:"device"`
}
