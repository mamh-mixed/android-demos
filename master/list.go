package master

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/CardInfoLink/quickpay/model"
	"github.com/CardInfoLink/quickpay/mongo"

	"github.com/CardInfoLink/log"
)

// kvListHandle 列表查询
func kvListHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("Read all body error: %s", err)
		w.WriteHeader(501)
		return
	}

	// 从当前session中取得agentCode，subAgentCode，groupCode， merId
	sid, err := r.Cookie(SessionKey)
	keyword := r.FormValue("keyword")

	ret := listFind(sid.Value, data, keyword)
	rdata, err := json.Marshal(ret)
	if err != nil {
		w.Write([]byte("mashal data error"))
	}

	log.Tracef("response message: %s", rdata)
	w.Write(rdata)
}

func listFind(sessionId string, data []byte, keyword string) (result *model.ResultBody) {
	cond := new(mongo.KVListCondition)
	err := json.Unmarshal(data, cond)
	if err != nil {
		log.Errorf("json(%s) unmarshal error: %s", string(data), err)
		return model.NewResultBody(2, "解析失败")
	}

	if cond.ColName == "" {
		log.Error("missing collection name")
		return model.NewResultBody(3, "COLLECTION_NAME_IS_REQUIRED")
	}

	log.Debugf("sessionId is %s", sessionId)
	// 查找session信息
	session, err := mongo.SessionColl.Find(sessionId)
	if err != nil {
		log.Errorf("find session(%s) err: %s", sessionId, err)
		return model.NewResultBody(1, "查询失败")
	}
	user := session.User

	switch user.UserType {
	case "admin":

	case "genAdmin":

	case "agent":
		cond.FilterMap["agentCode"] = user.AgentCode
	case "subAgent":
		cond.FilterMap["subAgentCode"] = user.SubAgentCode
	case "group":
		cond.FilterMap["groupCode"] = user.GroupCode
	case "merchant":
		cond.FilterMap["merId"] = user.MerId
	default:
		return model.NewResultBody(3, "INVALID_USER_TYPE")
	}

	list, err := cond.Find(keyword)
	if err != nil {
		log.Errorf("find list error: %s", err)
		return model.NewResultBody(2, "QUERY_LIST_ERROR")
	}

	result = &model.ResultBody{
		Status:  0,
		Message: "查询成功",
		Data:    list,
	}

	return result
}
