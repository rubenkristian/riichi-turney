package riichicommand

type SidData struct {
	Data string `json:"data"`
}

type Domain struct {
	DomainName string `json:"domain_name"`
	IsOpen     bool   `json:"is_open"`
}

type Credential struct {
	Password string `json:"passwd"`
	Email    string `json:"email"`
}

type ResponseLogin struct {
	Code         int       `json:"code"`
	AnalysisData string    `json:"analysis_data"`
	CheckData    string    `json:"check_data"`
	Message      string    `json:"message"`
	Data         LoginData `json:"data"`
}

type LoginData struct {
	BanStartAt         int      `json:"banStartAt"`
	BanUntil           int      `json:"banUnti"`
	CancelContactEmail string   `json:"cancelContactEmail"`
	CancelEndAt        int      `json:"cancelEndAt"`
	Country            string   `json:"country"`
	IsCompleteCourse   bool     `json:"isCompleteCourse"`
	IsCompleteGive     bool     `json:"isCompleteGive"`
	IsCompleteNew      bool     `json:"isCompleteNew"`
	IsCompleteNewRole  bool     `json:"isCompleteNewRole"`
	ServerTime         int      `json:"serverTime"`
	TokenTypes         []int    `json:"tokenTypes"`
	User               UserData `json:"user"`
}

type UserData struct {
	Avatar     string `json:"avatar"`
	Email      string `json:"email"`
	Id         int64  `json:"id"`
	Nickname   string `json:"nickname"`
	RegisterAt int64  `json:"registerAt"`
	Status     int    `json:"status"`
}

type ResponseTournamentInfo struct {
	Code         int            `json:"code"`
	AnalysisData string         `json:"analysis_data"`
	CheckData    string         `json:"check_data"`
	Message      string         `json:"message"`
	Data         TournamentInfo `json:"data"`
}

type TournamentInfo struct {
	ClassifyID     string            `json:"classifyID"`
	IsAdmin        bool              `json:"isAdmin"`
	IsCanReady     bool              `json:"isCanReady"`
	IsCollect      bool              `json:"isCollect"`
	LateRoomRank   []any             `json:"lateRoomRank"`
	LateTotalScore int               `json:"lateTotalScore"`
	MatchID        int               `json:"matchID"`
	MatchInfo      TournamentSetting `json:"matchInfo"`
	OnlineSize     int               `json:"onlineSize"`
	OwnerID        int               `json:"ownerID"`
	Rank           int               `json:"rank"`
	RankScore      int               `json:"rankScore"`
	RoomSize       int               `json:"roomSize"`
	StartedSize    int               `json:"startedSize"`
	Status         int               `json:"status"`
	TeamId         int               `json:"teamID"`
	TeamName       string            `json:"teamName"`
}

type TournamentSetting struct {
	BriefIntroduction string        `json:"briefIntroduction"`
	CardType          int           `json:"CardType"`
	ChangBang         int           `json:"ChangBang"`
	ContinuousX       int           `json:"ContinuousX"`
	EndTime           int64         `json:"endTime"`
	EnterPassWord     string        `json:"enterPassWord"`
	EnterRulesType    int           `json:"enterRulesType"`
	FangFu            int           `json:"FangFu"`
	FightType         int           `json:"fightType"`
	FristReqPoints    int           `json:"fristReqPoints"`
	GameMode          int           `json:"GameMode"`
	HighStageRank     int           `json:"highStageRank"`
	InitialPoints     int           `json:"initialPoints"`
	IsAddUpYakuman    bool          `json:"isAddUpYakuman"`
	IsAutoLi          bool          `json:"IsAutoLi"`
	IsChiDuan         bool          `json:"isChiDuan"`
	IsChiTi           bool          `json:"isChiTi"`
	IsConvenientTips  bool          `json:"isConvenientTips"`
	IsCutOff          bool          `json:"isCutOff"`
	IsDeclareNotTing  bool          `json:"IsDeclareNotTing"`
	IsFourGang        bool          `json:"isFourGang"`
	IsFourRiichi      bool          `json:"isFourRiichi"`
	IsFourWinds       bool          `json:"isFourWinds"`
	IsGangDora        bool          `json:"isGangDora"`
	IsGangLiDora      bool          `json:"isGangLiDora"`
	IsGangOpen        bool          `json:"isGangOpen"`
	IsGangPay         bool          `json:"isGangPay"`
	IsGeMu            bool          `json:"IsGeMu"`
	IsHandFondle      bool          `json:"isHandFondle"`
	IsJinChiTi        bool          `json:"isJinChiTi"`
	IsKaiLiZhi        bool          `json:"isKaiLiZhi"`
	IsKnock           bool          `json:"isKnock"`
	IsLastHeEnd       bool          `json:"isLastHeEnd"`
	IsLastTingEnd     bool          `json:"isLastTingEnd"`
	IsLiDora          bool          `json:"isLiDora"`
	IsLuck            bool          `json:"isLuck"`
	IsMinusRiichi     bool          `json:"isMinusRiichi"`
	IsMultipleYakuman bool          `json:"isMultipleYakuman"`
	IsNanXiRu         bool          `json:"isNanXiRu"`
	IsNineCards       bool          `json:"isNineCards"`
	IsNotEffect       bool          `json:"IsNotEffect"`
	IsNotKeepDealer   bool          `json:"IsNotKeepDealer"`
	IsNotLuJvManGuan  bool          `json:"IsNotLuJvManGuan"`
	IsOfficial        bool          `json:"isOfficial"`
	IsOpenFace        bool          `json:"isOpenFace"`
	IsOpenGuFang      bool          `json:"IsOpenGuFang"`
	IsPayYakuman      bool          `json:"isPayYakuman"`
	IsPublic          bool          `json:"isPublic"`
	IsQiangGang       bool          `json:"isQiangGang"`
	IsQieShang        bool          `json:"isQieShang"`
	IsRenHe           bool          `json:"isRenHe"`
	IsSameOrder       bool          `json:"isSameOrder"`
	IsShaoJi          bool          `json:"isShaoJi"`
	IsThreeHe         bool          `json:"isThreeHe"`
	IsTimeLimit       bool          `json:"isTimeLimit"`
	IsTopReward       bool          `json:"isTopReward"`
	IsYiFa            bool          `json:"isYiFa"`
	LimitRoomSize     int           `json:"limitRoomSize"`
	LimitTime         int           `json:"limitTime"`
	LowStageRank      int           `json:"lowStageRank"`
	MatchRulesType    int           `json:"matchRulesType"`
	MinimumPoints     int           `json:"minimumPoints"`
	Name              string        `json:"name"`
	NumRedCard        int           `json:"numRedCard"`
	OperFixedTime     int           `json:"operFixedTime"`
	OperVarTime       int           `json:"operVarTime"`
	OrderPoints       []int         `json:"orderPoints"`
	PenaltyInfo       []PenaltyInfo `json:"penaltyInfo"`
	PlayerCount       int           `json:"playerCount"`
	RankType          int           `json:"RankType"`
	Round             int           `json:"round"`
	StartTime         int64         `json:"startTime"`
	ThreeZiMoType     int           `json:"threeZiMoType"`
	TimeLimitType     int           `json:"timeLimitType"`
	Type              int           `json:"type"`
}

type PenaltyInfo struct {
	Num         int `json:"num"`
	PointsCount int `json:"pointsCount"`
}

type TournamentLogList struct {
	EndTime       int       `json:"endTime"`
	GamePlay      int       `json:"gamePlay"`
	IsClear       bool      `json:"isClear"`
	IsCollect     bool      `json:"isCollect"`
	IsMiddlePause bool      `json:"isMiddlePause"`
	MatchStage    bool      `json:"matchStage"`
	MatchType     int       `json:"matchType"`
	PaiPuId       string    `json:"paiPuId"`
	PaiPuNotId    string    `json:"paiPuNotId"`
	Period        int       `json:"period"`
	PlayerCount   int       `json:"playerCount"`
	Players       []Players `json:"players"`
	Remark        string    `json:"remark"`
	RoomID        string    `json:"roomID"`
	Round         int       `json:"round"`
	SignItemID    int       `json:"signItemID"`
	StageNum      int       `json:"stageNum"`
	StageType     int       `json:"stageType"`
}

type Players struct {
	Identity     int    `json:"identity"`
	IsExistYiMan bool   `json:"isExistYiMan"`
	Nickname     string `json:"nickname"`
	Points       int64  `json:"points"`
	Rank         int    `json:"rank"`
	RoleID       int    `json:"roleID"`
	Serial       int    `json:"serial"`
	SkinID       int    `json:"skinID"`
	UserId       uint64 `json:"userId"`
}

type ResponseTournameLogList struct {
	Code         int                 `json:"code"`
	AnalysisData string              `json:"analysis_data"`
	CheckData    string              `json:"check_data"`
	Message      string              `json:"message"`
	Data         []TournamentLogList `json:"data"`
}

type TournamentPlayer struct {
	IsInvited bool   `json:"isInvited"`
	Nick      string `json:"nick"`
	Status    int    `json:"status"`
	UserID    uint64 `json:"userID"`
}

type ResponseTournamentPlayers struct {
	Code         int                `json:"code"`
	AnalysisData string             `json:"analysisData"`
	CheckData    string             `json:"checkData"`
	Message      string             `json:"message"`
	Data         []TournamentPlayer `json:"data"`
}

type TournamentLiveGame struct {
	IsEnd     bool            `json:"isEnd"`
	IsPause   bool            `json:"isPause"`
	NowTime   int             `json:"nowTime"`
	Players   []PlayersIngame `json:"players"`
	RoomId    string          `json:"roomId"`
	StartTime int             `json:"startTime"`
}

type PlayersIngame struct {
	Headtag        int    `json:"headtag"`
	Nickname       string `json:"nickname"`
	Position       int    `json:"position"`
	ProfileFrameId int    `json:"profileFrameId"`
	RoleID         int    `json:"roleID"`
	SkinID         int    `json:"skinID"`
	StageLevel     int    `json:"stageLevel"`
	UserId         int    `json:"userId"`
}

type ResponseTournamentLiveGame struct {
	Code         int                  `json:"code"`
	AnalysisData string               `json:"analysis_data"`
	CheckData    string               `json:"check_data"`
	Message      string               `json:"message"`
	Data         []TournamentLiveGame `json:"data"`
}

type Log struct {
	fangFu      int          `json:"fangFu"`
	GameMode    int          `json:"gameMode"`
	GamePlay    int          `json:"gamePlay"`
	HandRecords []HandRecord `json:"handRecords"`
	InitPoints  int          `json:"initPoints"`
	KeyValue    string       `json:"keyValue"`
	NowTime     int          `json:"nowTime"`
	Remark      string       `json:"remark"`
	RoomID      string       `json:"roomId"`
}

type HandRecord struct {
	BenChangNum        int               `json:"benChangNum"`
	ChangCi            int               `json:"changCi"`
	HandPos            int               `json:"handPos"`
	QuanFeng           int               `json:"quanFeng"`
	HandCardEncode     string            `json:"handCardEncode"`
	HandCardsSHA256    string            `json:"handCardsSHA256"`
	HandID             string            `json:"handID"`
	HandEventRecords   []HandEventRecord `json:"handEventRecord"`
	PaiShan            []int             `json:"paiShan"`
	Players            []Player          `json:"players"`
	SymbolEventRequest []any             `json:"symbolEventRecord"`
}

type Player struct {
	CardBackID         int    `json:"cardBackID"`
	GameMusicId        int    `json:"gameMusicId"`
	HandId             int    `json:"handId"`
	HeadTag            int    `json:"headTag"`
	Identity           int    `json:"identity"`
	MatchMusicId       int    `json:"matchMusicId"`
	Model              int    `json:"model"`
	MortalModule       string `json:"mortalModule"`
	Nickname           string `json:"nickname"`
	Points             int    `json:"points"`
	Position           int    `json:"position"`
	ProfileFrameId     int    `json:"profileFrameId"`
	RandGameMusicIds   []int  `json:"randGameMusicIDs"`
	RandGameMusicModel int    `json:"randGameMusicModel"`
	RiichiEffectId     int    `json:"riichiEffectId"`
	RiichiMusicId      int    `json:"riichiMusicId"`
	RiichiStickId      int    `json:"riichiStickID"`
	RobotLevel         int    `json:"robotLevel"`
	RoleID             int    `json:"roleID"`
	SkinID             int    `json:"skinID"`
	SpecialEffectID    int    `json:"specialEffectID"`
	TableClothId       int    `json:"tableClothID"`
	TitleId            int    `json:"titleID"`
	UserTag            int    `json:"userTag"`
}

type HandEventRecord struct {
	/**
	 * Events, defined by event_type:
	 * 1. **Start** - Fires at the start of the round, once for each player.
	 * 2. **Draw tile** - Fires whenever a player draws from the wall.
	 * 3. **Call** - Fires whenever a player makes a call (to draw a tile).
	 * 4. **Discard** - Fires whenever a player discards a tile.
	 * 5. **End** - Fires when the round ends.
	 * 6. **Game End** - Fires when the game ends.
	 * 7. **Dora reveal** - Fires when a dora is revealed.
	 * 8. **Auto kan** - Fires when a player enables auto kan.
	 * 9. **Empty** - Unknown what this event does.
	 * 10. **Pause** - Fires when a pause is issued.
	 * 11. **Tenpai** - Fires when a player has tenpai.
	 */
	Data      string `json:"data"`
	EventPos  int    `json:"eventPos"`
	EventType int    `json:"eventType"`
	HandId    string `json:"handId"`
	StartTime int    `json:"startTime"`
	UserId    int    `json:"userId"`
}

type FindPlayer struct {
	FindType   int          `json:"findType"`
	FriendList []FriendList `json:"friendList"`
}

type FriendList struct {
	CreateTime        int64  `json:"createTime"`
	CurrentStatus     int    `json:"currentStatus"`
	FourLevel         int    `json:"fourLevel"`
	IsCanApply        bool   `json:"isCanApply"`
	IsMutual          bool   `json:"isMutual"`
	Model             int    `json:"model"`
	Nickname          string `json:"nickname"`
	NowTime           int64  `json:"nowTime"`
	ObserveGamePlay   int    `json:"observeGamePlay"`
	ObserveGameRoomID int    `json:"observeGameRoomID"`
	PlayerCount       int    `json:"playerCount"`
	ProfileFrameID    int    `json:"profileFrameID"`
	Remark            string `json:"remark"`
	RoleID            int    `json:"roleID"`
	RoomStartTime     int64  `json:"roomStartTime"`
	ShowGameToFriend  bool   `json:"showGameToFriend"`
	SkinID            int    `json:"skinID"`
	StatusTime        int64  `json:"statusTime"`
	ThreeLevel        int    `json:"threeLevel"`
	TitleID           int    `json:"titleID"`
	UserID            int    `json:"userID"`
}

type ResponseFindPlayer struct {
	Code         int        `json:"code"`
	AnalysisData string     `json:"analysis_data"`
	CheckData    string     `json:"check_data"`
	Message      string     `json:"message"`
	Data         FindPlayer `json:"data"`
}
