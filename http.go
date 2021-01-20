package main

import (
	"encoding/json"
	"fmt"
	"github.com/vysiondev/qstats-go/bot_constants"
	"github.com/vysiondev/qstats-go/err"
	"github.com/vysiondev/qstats-go/quaverapi_structs"
	"github.com/vysiondev/qstats-go/utils"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"strconv"
)

func MakeHTTPRequest(url string, thisType string) ([]byte, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "QStats Discord Bot")
	resp, reqErr := http.DefaultClient.Do(req)

	if reqErr != nil {
		return nil, &err.SafeError{Message: "Quaver API cannot be reached."}
	}
	if resp.StatusCode == 404 {
		return nil, &err.SafeError{Message: fmt.Sprintf("%s not found.", thisType)}
	}

	defer func() {
		closeErr := resp.Body.Close()
		if closeErr != nil {
			return
		}
	}()
	responseData, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return nil, &err.ReadError{Message: "The response body from Quaver API could not be read."}
	}
	return responseData, nil
}

func (b *BaseHandler) GetUserSpecificOnlineStatus(userID int) (string, error) {
	bytes, reqErr := MakeHTTPRequest(bot_constants.QuaverEndpoint+"/server/users/online/"+strconv.Itoa(userID), "endpoint")
	if reqErr != nil {
		return "", reqErr
	}
	var on *quaverapi_structs.UserOnline
	jsonErr := json.Unmarshal(bytes, &on)
	if jsonErr != nil {
		return "", &err.ReadError{Message: "Could not parse JSON into struct for online status."}
	}

	if !on.IsOnline {
		return b.Config.Emoji.Offline + " Offline", nil
	} else {
		var returnStr string

		switch on.CurrentStatus.Status {
		case 0:
			returnStr = "In menus"
			break
		case 1:
			returnStr = "In track selection"
			break
		case 2:
			returnStr = "Playing " + on.CurrentStatus.Content
			break
		case 3:
			returnStr = "Paused"
			break
		case 4:
			returnStr = fmt.Sprintf("Spectating [%s](%s/user/%s)", on.CurrentStatus.Content, bot_constants.QuaverMainSite, on.CurrentStatus.Content)
			break
		case 5:
			returnStr = "Editing " + utils.RemoveFormattingCharacters(on.CurrentStatus.Content)
			break
		case 6:
			returnStr = "In a multiplayer lobby"
			break
		case 7:
			returnStr = "Playing in a multiplayer match"
			break
		case 8:
			returnStr = "Listening to " + on.CurrentStatus.Content
			break
		default:
			returnStr = "Unknown status"
		}

		return b.Config.Emoji.Online + " " + returnStr, nil
	}
}

func SearchUser(username string) (quaverapi_structs.UserSearch, error) {
	bytes, reqErr := MakeHTTPRequest(bot_constants.QuaverEndpoint+"/users/search/"+url.PathEscape(username), "user")
	if reqErr != nil {
		return quaverapi_structs.UserSearch{}, reqErr
	}
	var resp *quaverapi_structs.UserSearch
	jsonErr := json.Unmarshal(bytes, &resp)
	if jsonErr != nil {
		return quaverapi_structs.UserSearch{}, &err.ReadError{Message: "Could not parse JSON into struct for user search."}
	}
	if len(resp.Users) == 0 {
		return quaverapi_structs.UserSearch{}, &err.SafeError{Message: "No users matched this query."}
	}
	return *resp, nil
}

func GetUserData(userID int) (quaverapi_structs.UserData, error) {
	bytes, reqErr := MakeHTTPRequest(bot_constants.QuaverEndpoint+"/users/full/"+strconv.Itoa(userID), "user")
	if reqErr != nil {
		return quaverapi_structs.UserData{}, reqErr
	}
	var resp *quaverapi_structs.UserData
	jsonErr := json.Unmarshal(bytes, &resp)
	if jsonErr != nil {
		return quaverapi_structs.UserData{}, &err.ReadError{Message: "Could not parse JSON into struct for full user data."}
	}
	return *resp, nil
}

func GetScores(userID int, is7K bool, topPlays bool, page int, oneScoreOnly bool) (quaverapi_structs.UserScoresResponse, error) {
	user, reqErr := GetUserData(userID)
	if reqErr != nil {
		return quaverapi_structs.UserScoresResponse{}, reqErr
	}
	urlPlayType := "recent"
	if topPlays {
		urlPlayType = "best"
	}

	urlScoreLimit := 5
	if oneScoreOnly {
		urlScoreLimit = 1
	}
	bytes, scoreErr := MakeHTTPRequest(bot_constants.QuaverEndpoint+"/users/scores/"+urlPlayType+"?id="+strconv.Itoa(userID)+"&mode="+utils.GetKeymodeIntAsStr(is7K)+"&limit="+strconv.Itoa(urlScoreLimit)+"&page="+strconv.Itoa(page-1), "scores for this user")
	if scoreErr != nil {
		return quaverapi_structs.UserScoresResponse{}, scoreErr
	}
	var resp *quaverapi_structs.UserScores
	jsonErr := json.Unmarshal(bytes, &resp)
	if jsonErr != nil {
		return quaverapi_structs.UserScoresResponse{}, &err.ReadError{Message: "Could not parse JSON into struct for user scores."}
	}

	if len(resp.Scores) == 0 {
		errT := "recent"
		if topPlays {
			errT = "best"
		}
		playCt := fmt.Sprintf("more than %d %s", page, errT)
		if !oneScoreOnly {
			playCt = "any"
		}
		return quaverapi_structs.UserScoresResponse{}, &err.SafeError{Message: fmt.Sprintf("%s doesn't have %s %s scores.",
			user.User.Info.Username,
			playCt,
			utils.GetKeymodeString(is7K),
		)}
	}
	return quaverapi_structs.UserScoresResponse{
		User:   user,
		Scores: resp.Scores,
	}, nil
}

type FirstScoreResponse struct {
	BestScoreIndex int
	Score          quaverapi_structs.Score
	UserData       quaverapi_structs.UserData
}

func FindFirstScoreInMapIDArray(userID int, mapIDs []string, is7K bool) (FirstScoreResponse, error) {
	user, reqErr := GetUserData(userID)
	if reqErr != nil {
		return FirstScoreResponse{}, reqErr
	}

	bytes, scoreErr := MakeHTTPRequest(bot_constants.QuaverEndpoint+"/users/scores/best?id="+strconv.Itoa(user.User.Info.ID)+"&mode="+utils.GetKeymodeIntAsStr(is7K)+"&limit=50", "scores for this user")
	if scoreErr != nil {
		return FirstScoreResponse{}, scoreErr
	}
	var scores *quaverapi_structs.UserScores
	jsonErr := json.Unmarshal(bytes, &scores)
	if jsonErr != nil {
		return FirstScoreResponse{}, &err.ReadError{Message: "Failed to parse JSON for user's first top score in array of scores."}
	}

	var score quaverapi_structs.Score
	bestScoreIndex := -1

	for _, id := range mapIDs {
		for scoreIndex, thisScore := range scores.Scores {
			if strconv.Itoa(int(thisScore.Map.ID)) == id {
				score = thisScore
				bestScoreIndex = scoreIndex
				break
			}
		}
		if bestScoreIndex != -1 {
			break
		}
	}

	if bestScoreIndex == -1 {
		return FirstScoreResponse{}, &err.SafeError{Message: fmt.Sprintf("%s doesn't have any scores in their top 50 for any map in the conversation.", user.User.Info.Username)}
	}
	return FirstScoreResponse{
		BestScoreIndex: bestScoreIndex,
		Score:          score,
		UserData:       user,
	}, nil
}

func GetMap(mapID int) (quaverapi_structs.QuaverMapResponse, error) {
	bytes, reqErr := MakeHTTPRequest(bot_constants.QuaverEndpoint+"/maps/"+strconv.Itoa(mapID), "map")
	if reqErr != nil {
		return quaverapi_structs.QuaverMapResponse{}, reqErr
	}
	var resp *quaverapi_structs.QuaverMapResponse
	jsonErr := json.Unmarshal(bytes, &resp)
	if jsonErr != nil {
		return quaverapi_structs.QuaverMapResponse{}, &err.SafeError{Message: "Failed to parse JSON for map data."}
	}
	return *resp, nil
}

func GetMapset(mapsetID int) (quaverapi_structs.QuaverMapsetResponse, error) {
	bytes, reqErr := MakeHTTPRequest(bot_constants.QuaverEndpoint+"/mapsets/"+strconv.Itoa(mapsetID), "mapset")
	if reqErr != nil {
		return quaverapi_structs.QuaverMapsetResponse{}, reqErr
	}
	var resp *quaverapi_structs.QuaverMapsetResponse
	jsonErr := json.Unmarshal(bytes, &resp)
	if jsonErr != nil {
		return quaverapi_structs.QuaverMapsetResponse{}, &err.SafeError{Message: "Failed to parse JSON for mapset data."}
	}
	return *resp, nil
}

func GetLeaderboard(page int, country string, is7K bool) ([]quaverapi_structs.LeaderboardUser, error) {
	pageToSearch := int(math.Floor(float64(page-1) / 2))

	bytes, reqErr := MakeHTTPRequest(bot_constants.QuaverEndpoint+"/leaderboard/?mode="+utils.GetKeymodeIntAsStr(is7K)+"&page="+strconv.Itoa(pageToSearch)+"&country="+country, "leaderboards for this country")
	if reqErr != nil {
		return nil, reqErr
	}
	var lbData *quaverapi_structs.QuaverLeaderboardResponse
	jsonErr := json.Unmarshal(bytes, &lbData)
	if jsonErr != nil {
		return nil, &err.ReadError{Message: "Failed to parse JSON for mapset data."}
	}

	if len(lbData.Users) == 0 {
		return nil, &err.SafeError{Message: fmt.Sprintf("Leaderboards for this country were not found/there are not any more scores for this country.")}
	}
	if page%2 == 0 {
		if len(lbData.Users) <= 25 {
			return nil, &err.SafeError{Message: "No more players beyond this page."}
		}
		return lbData.Users[25:], nil
	} else {
		if len(lbData.Users) >= 25 {
			return lbData.Users[:25], nil
		} else {
			return lbData.Users, nil
		}
	}
}
