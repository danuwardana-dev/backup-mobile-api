package dto

type HandlerLog struct {
	RequestId string                 `json:"request_id"`
	Timestamp string                 `json:"timestamp"`
	UserUUID  string                 `json:"user_uuid"`
	Email     string                 `json:"email"`
	Remarks   string                 `json:"remarks"`
	Url       string                 `json:"path"`
	Method    string                 `json:"method"`
	Device    Device                 `json:"device"`
	Error     string                 `json:"error"`
	Data      map[string]interface{} `json:"data"`
	Success   bool                   `json:"success"`
}
type Device struct {
	DeviceID  string `json:"device_id"`
	Longitude string `json:"longitude"`
	Latitude  string `json:"latitude"`
	Ip        string `json:"ip"`
}
type CustomLoggerRequest struct {
	UserUUID string
	Email    string
	Remarks  string
	Data     map[string]interface{}
	Error    string
	Success  bool
}
