package utils

import (
	"github.com/mattermost/mattermost-server/model"
	"testing"
)

func TestMattermost(t *testing.T) {
	mattermostUrl := "http://192.168.131.139:8065"
	client := model.NewAPIv4Client(mattermostUrl)
	user, resp := client.Login("yasin", "yasin")
	//client.SetToken()
	if resp.StatusCode != 200 {
		t.Error(resp.StatusCode, resp.Error)
		return
	}
	t.Log(user)
}
