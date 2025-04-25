package riichicommand

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"os"
	"time"

	"github.com/rubenkristian/riichi-turney/database"
)

type RiichiApi struct {
	deviceId      string
	domain        string
	sid           string
	uid           string
	version       string
	defaultHeader map[string]string
	credential    Credential
	isLoggedIn    bool
	dbGame        *database.DatabaseGame
}

type SIDResponse struct {
}

func CreateRiichiApi(dbGame *database.DatabaseGame) *RiichiApi {
	return &RiichiApi{
		dbGame: dbGame,
	}
}

func (ra *RiichiApi) SetupRiichi(mainHost string, email string, password string) error {
	domain, err := ra.getDomain(mainHost)

	if err != nil {
		return err
	}

	ra.domain = domain
	hash := md5.Sum([]byte(password))
	ra.credential = Credential{
		Email:    email,
		Password: hex.EncodeToString(hash[:]),
	}

	if err := ra.login(); err != nil {
		return err
	}

	return nil
}

func (ra *RiichiApi) getDomain(mainHost string) (string, error) {
	res, err := http.Get(mainHost)

	if err != nil {
		return "", err
	}

	var domain Domain
	if err := json.NewDecoder(res.Body).Decode(&domain); err != nil {
		return "", err
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return "", err
	}

	if err := os.WriteFile("./temp/domain.json", body, 0644); err != nil {
		return "", err
	}

	return domain.DomainName, nil
}

func (ra *RiichiApi) login() error {
	sid, err := ra.fetchSID()

	if err != nil {
		return err
	}

	uid, err := ra.fetchLogin()

	if err != nil {
		return err
	}

	ra.sid = sid
	ra.uid = uid
	ra.isLoggedIn = true
	ra.refreshHeader()

	return nil
}

func (ra *RiichiApi) fetchSID() (string, error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s/users/initSession", ra.domain), nil)

	if err != nil {
		return "", err
	}

	req.Header.Add("Cookies", fmt.Sprintf(`{"channel":"default","deviceid":"%s","lang":"en","version":"%s","platform":"pc"}`, ra.deviceId, ra.version))

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	var result SidData
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Data, nil
}

func (ra *RiichiApi) fetchLogin() (string, error) {
	credential, err := json.MarshalIndent(ra.credential, "", " ")

	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s/users/emailLogin", ra.domain), bytes.NewBuffer(credential))

	if err != nil {
		return "", err
	}

	req.Header.Add("Cookies", fmt.Sprintf(`{"channel":"default","deviceid":"%s","lang":"en","sid":"%s","version":"%s","platform":"pc"}`, ra.deviceId, ra.sid, ra.version))
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	var loginData ResponseLogin
	if err := json.NewDecoder(res.Body).Decode(&loginData); err != nil {
		return "", err
	}

	body, err := json.MarshalIndent(loginData, "", " ")

	if err != nil {
		return "", err
	}

	if err := os.WriteFile(fmt.Sprintf("./temp/user-%d.json", time.Now().UnixNano()), body, 0644); err != nil {
		return "", err
	}

	return loginData.Data.User.Id, nil
}

func (ra *RiichiApi) refreshHeader() {
	ra.defaultHeader = map[string]string{
		"User-Agent":      "UnityPlayer/2020.3.42f1c1 (UnityWebRequest/1.0, libcurl/7.84.0-DEV)",
		"Cookies":         fmt.Sprintf(`{"channel":"default","lang":"en","deviceid":"%s","sid":"%s","uid":%s,"region":"cn","platform":"pc","version":"%s"}`, ra.deviceId, ra.sid, ra.uid, ra.version),
		"Content-Type":    "application/json",
		"Accept":          "application/json",
		"X-Unity-Version": "2021.3.38f1",
	}
}

func (ra *RiichiApi) FetchTournamentInfo(turneyId int) (*TournamentInfo, error) {
	if !ra.isLoggedIn {
		return nil, fmt.Errorf("failed to fetch tournament info, you're not login yet")
	}

	body, err := json.Marshal(map[string]int{
		"id": turneyId,
	})

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s/lobbys/enterSelfBuild", ra.domain), bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	for name, value := range ra.defaultHeader {
		req.Header.Add(name, value)
	}

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var responseTournamentInfo ResponseTournamentInfo
	if err := json.NewDecoder(res.Body).Decode(&responseTournamentInfo); err != nil {
		return nil, err
	}

	return &responseTournamentInfo.Data, nil
}

func (ra *RiichiApi) UpdateTournamentInfo(turneyId int, tournamentSetting TournamentSetting) (bool, error) {
	if !ra.isLoggedIn {
		return false, fmt.Errorf("failed to update tournament info, you're not login yet")
	}

	body, err := json.Marshal(map[string]any{
		"Info":    tournamentSetting,
		"matchID": turneyId,
	})

	if err != nil {
		return false, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s/lobbys/changeSelfCfg", ra.domain), bytes.NewBuffer(body))

	if err != nil {
		return false, err
	}

	for name, value := range ra.defaultHeader {
		req.Header.Add(name, value)
	}

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return false, err
	}

	defer res.Body.Close()

	return true, nil
}

func (ra *RiichiApi) FetchTournamentLogList(classifyId string, lastId int) ([]TournamentLogList, error) {
	if !ra.isLoggedIn {
		return nil, fmt.Errorf("failed to fetch tournament info, you're not login yet")
	}

	body, err := json.Marshal(map[string]any{
		"endTime":    0,
		"skip":       lastId,
		"startTime":  0,
		"classifyID": classifyId,
		"isSelf":     true,
		"limit":      20,
		"gamePlay":   1002,
	})

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s/record/readPaiPuList", ra.domain), bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	for name, value := range ra.defaultHeader {
		req.Header.Add(name, value)
	}

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var responseTournameLogList ResponseTournameLogList
	if err := json.NewDecoder(res.Body).Decode(&responseTournameLogList); err != nil {
		return nil, err
	}

	return responseTournameLogList.Data, nil
}

func (ra *RiichiApi) FetchTournamentPlayers(tourneyId int) ([]TournamentPlayer, error) {
	if !ra.isLoggedIn {
		return nil, fmt.Errorf("failed to fetch tournament info, you're not login yet")
	}

	body, err := json.Marshal(map[string]any{
		"matchID": tourneyId,
	})

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s/lobbys/getSelfManageInfo", ra.domain), bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	for name, value := range ra.defaultHeader {
		req.Header.Add(name, value)
	}

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var responseTournamentPlayers ResponseTournamentPlayers
	if err := json.NewDecoder(res.Body).Decode(&responseTournamentPlayers); err != nil {
		return nil, err
	}

	return responseTournamentPlayers.Data, nil
}

func (ra *RiichiApi) StartTournamentGame(tourneyId int, players []int, randomSeat bool) (bool, error) {
	if !ra.isLoggedIn {
		return false, fmt.Errorf("failed to update tournament info, you're not login yet")
	}

	if randomSeat {
		players = ra.shuffleArray(players)
	}

	body, err := json.Marshal(map[string]any{
		"usersID":   players,
		"matchID":   tourneyId,
		"table_idx": 1,
	})

	if err != nil {
		return false, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s/lobbys/allocateSelfUser", ra.domain), bytes.NewBuffer(body))

	if err != nil {
		return false, err
	}

	for name, value := range ra.defaultHeader {
		req.Header.Add(name, value)
	}

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return false, err
	}

	defer res.Body.Close()

	return true, nil
}

func (ra *RiichiApi) AdvStartTournamentGame(tourneyId int, players []int, scores []int, randomSeat bool) (bool, error) {
	if !ra.isLoggedIn {
		return false, fmt.Errorf("failed to update tournament info, you're not login yet")
	}

	if randomSeat {
		players, scores = ra.shuffleArrays(players, scores)
	}

	body, err := json.Marshal(map[string]any{
		"usersID":    players,
		"matchID":    tourneyId,
		"table_idx":  1,
		"initPoints": scores,
	})

	if err != nil {
		return false, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s/lobbys/allocateSelfUser", ra.domain), bytes.NewBuffer(body))

	if err != nil {
		return false, err
	}

	for name, value := range ra.defaultHeader {
		req.Header.Add(name, value)
	}

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return false, err
	}

	defer res.Body.Close()

	return true, nil
}

func (ra *RiichiApi) FetchOngoingGames(classifyId string) ([]TournamentLiveGame, error) {
	if !ra.isLoggedIn {
		return nil, fmt.Errorf("failed to fetch tournament info, you're not login yet")
	}

	body, err := json.Marshal(map[string]any{
		"stageType":  0,
		"round":      2,
		"matchType":  1,
		"classifyID": classifyId,
		"gamePlay":   1002,
	})

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s/record/readOnlineRoom", ra.domain), bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	for name, value := range ra.defaultHeader {
		req.Header.Add(name, value)
	}

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var responseTournamentLiveGame ResponseTournamentLiveGame
	if err := json.NewDecoder(res.Body).Decode(&responseTournamentLiveGame); err != nil {
		return nil, err
	}

	return responseTournamentLiveGame.Data, nil
}

func (ra *RiichiApi) ManageTournamentGame(tourneyId int, roomId string, Type int) (bool, error) {
	if !ra.isLoggedIn {
		return false, fmt.Errorf("failed to update tournament info, you're not login yet")
	}

	body, err := json.Marshal(map[string]any{
		"roomID":  roomId,
		"type":    Type,
		"matchID": tourneyId,
	})

	if err != nil {
		return false, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s/lobbys/controlSelfRoom", ra.domain), bytes.NewBuffer(body))

	if err != nil {
		return false, err
	}

	for name, value := range ra.defaultHeader {
		req.Header.Add(name, value)
	}

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return false, err
	}

	defer res.Body.Close()

	return true, nil
}

func (ra *RiichiApi) FetchLog(paifu string) (*Log, error) {
	if !ra.isLoggedIn {
		return nil, fmt.Errorf("failed to fetch tournament info, you're not login yet")
	}

	body, err := json.Marshal(map[string]any{
		"isObserve": false,
		"keyValue":  paifu,
	})

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s/record/getRoomData", ra.domain), bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	for name, value := range ra.defaultHeader {
		req.Header.Add(name, value)
	}

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var log Log
	if err := json.NewDecoder(res.Body).Decode(&log); err != nil {
		return nil, err
	}

	return &log, nil
}

func (ra *RiichiApi) FindPlayer(userId string) (*FindPlayer, error) {
	if !ra.isLoggedIn {
		return nil, fmt.Errorf("failed to fetch tournament info, you're not login yet")
	}

	body, err := json.Marshal(map[string]any{
		"findType": 2,
		"content":  userId,
	})

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s/mixed_client/findFriend", ra.domain), bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	for name, value := range ra.defaultHeader {
		req.Header.Add(name, value)
	}

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var players FindPlayer
	if err := json.NewDecoder(res.Body).Decode(&players); err != nil {
		return nil, err
	}

	return &players, nil
}

func (ra *RiichiApi) shuffleArray(players []int) []int {
	for i := len(players) - 1; i > 0; i-- {
		j := rand.IntN(i + 1)
		players[i], players[j] = players[j], players[i]
	}
	return players
}

func (ra *RiichiApi) shuffleArrays(players []int, scores []int) ([]int, []int) {
	var index = len(players)
	var rnd int = 0
	var tmp1 int = 0
	var tmp2 int = 0

	for index > 0 {
		rnd = int(rand.IntN(index))
		index -= 1
		tmp1 = players[index]
		tmp2 = scores[index]
		players[index] = players[rnd]
		scores[index] = scores[rnd]
		players[rnd] = tmp1
		scores[rnd] = tmp2
	}

	return players, scores
}
