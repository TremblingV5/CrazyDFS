package service

type Response struct {
	Code int64       `json:"code"`
	Msg  string      `json:"message"`
	Data interface{} `json:"data"`
}

type ClientApi interface {
	Put(local string, remote string) Response
	Get(local string, remote string) Response
	Delete(remote string) Response
	Stat(remote string) Response
	Rename(src string, target string) Response
	Mkdir(remote string) Response
	List(remote string) Response
}
