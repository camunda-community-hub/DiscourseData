package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

const lurkers = "5"
const top50 = "6"
const top50Posters = "4"
const mostSolutions = "7"

// APIKey is your Discpourse API_KEY here
const APIKey = "0476e2a15a1152d8c00c1eeec60bbcec4977e54b87f2d7556a857e814c5a4fb0"

// APIUser is your Discourse User ID
const APIUser = "davidgs"

// MaxItems is how many records you want to fetch each time
const MaxItems = "1000"

// DiscourseUser the full discourse user data
type DiscourseUser struct {
	ID                     int         `json:"id"`
	Username               string      `json:"username"`
	Name                   string      `json:"name"`
	AvatarTemplate         string      `json:"avatar_template"`
	Active                 bool        `json:"active"`
	Admin                  bool        `json:"admin"`
	Moderator              bool        `json:"moderator"`
	LastSeenAt             time.Time   `json:"last_seen_at"`
	LastEmailedAt          time.Time   `json:"last_emailed_at"`
	CreatedAt              time.Time   `json:"created_at"`
	LastSeenAge            float64     `json:"last_seen_age"`
	LastEmailedAge         float64     `json:"last_emailed_age"`
	CreatedAtAge           float64     `json:"created_at_age"`
	TrustLevel             int         `json:"trust_level"`
	ManualLockedTrustLevel interface{} `json:"manual_locked_trust_level"`
	FlagLevel              int         `json:"flag_level"`
	Title                  interface{} `json:"title"`
	TimeRead               int         `json:"time_read"`
	Staged                 bool        `json:"staged"`
	DaysVisited            int         `json:"days_visited"`
	PostsReadCount         int         `json:"posts_read_count"`
	TopicsEntered          int         `json:"topics_entered"`
	PostCount              int         `json:"post_count"`
	CanSendActivationEmail bool        `json:"can_send_activation_email"`
	CanActivate            bool        `json:"can_activate"`
	CanDeactivate          bool        `json:"can_deactivate"`
	IPAddress              string      `json:"ip_address"`
	RegistrationIPAddress  string      `json:"registration_ip_address"`
	CanGrantAdmin          bool        `json:"can_grant_admin"`
	CanRevokeAdmin         bool        `json:"can_revoke_admin"`
	CanGrantModeration     bool        `json:"can_grant_moderation"`
	CanRevokeModeration    bool        `json:"can_revoke_moderation"`
	CanImpersonate         bool        `json:"can_impersonate"`
	LikeCount              int         `json:"like_count"`
	LikeGivenCount         int         `json:"like_given_count"`
	TopicCount             int         `json:"topic_count"`
	FlagsGivenCount        int         `json:"flags_given_count"`
	FlagsReceivedCount     int         `json:"flags_received_count"`
	PrivateTopicsCount     int         `json:"private_topics_count"`
	CanDeleteAllPosts      bool        `json:"can_delete_all_posts"`
	CanBeDeleted           bool        `json:"can_be_deleted"`
	CanBeAnonymized        bool        `json:"can_be_anonymized"`
	CanBeMerged            bool        `json:"can_be_merged"`
	FullSuspendReason      interface{} `json:"full_suspend_reason"`
	SilenceReason          interface{} `json:"silence_reason"`
	PrimaryGroupID         interface{} `json:"primary_group_id"`
	BadgeCount             int         `json:"badge_count"`
	WarningsReceivedCount  int         `json:"warnings_received_count"`
	BounceScore            float64     `json:"bounce_score"`
	ResetBounceScoreAfter  interface{} `json:"reset_bounce_score_after"`
	CanViewActionLogs      bool        `json:"can_view_action_logs"`
	CanDisableSecondFactor bool        `json:"can_disable_second_factor"`
	CanDeleteSsoRecord     bool        `json:"can_delete_sso_record"`
	APIKeyCount            int         `json:"api_key_count"`
	SingleSignOnRecord     interface{} `json:"single_sign_on_record"`
	ApprovedBy             interface{} `json:"approved_by"`
	SuspendedBy            interface{} `json:"suspended_by"`
	SilencedBy             interface{} `json:"silenced_by"`
	Tl3Requirements        struct {
		TimePeriod             int  `json:"time_period"`
		RequirementsMet        bool `json:"requirements_met"`
		RequirementsLost       bool `json:"requirements_lost"`
		TrustLevelLocked       bool `json:"trust_level_locked"`
		OnGracePeriod          bool `json:"on_grace_period"`
		DaysVisited            int  `json:"days_visited"`
		MinDaysVisited         int  `json:"min_days_visited"`
		NumTopicsRepliedTo     int  `json:"num_topics_replied_to"`
		MinTopicsRepliedTo     int  `json:"min_topics_replied_to"`
		TopicsViewed           int  `json:"topics_viewed"`
		MinTopicsViewed        int  `json:"min_topics_viewed"`
		PostsRead              int  `json:"posts_read"`
		MinPostsRead           int  `json:"min_posts_read"`
		TopicsViewedAllTime    int  `json:"topics_viewed_all_time"`
		MinTopicsViewedAllTime int  `json:"min_topics_viewed_all_time"`
		PostsReadAllTime       int  `json:"posts_read_all_time"`
		MinPostsReadAllTime    int  `json:"min_posts_read_all_time"`
		NumFlaggedPosts        int  `json:"num_flagged_posts"`
		MaxFlaggedPosts        int  `json:"max_flagged_posts"`
		NumFlaggedByUsers      int  `json:"num_flagged_by_users"`
		MaxFlaggedByUsers      int  `json:"max_flagged_by_users"`
		NumLikesGiven          int  `json:"num_likes_given"`
		MinLikesGiven          int  `json:"min_likes_given"`
		NumLikesReceived       int  `json:"num_likes_received"`
		MinLikesReceived       int  `json:"min_likes_received"`
		NumLikesReceivedDays   int  `json:"num_likes_received_days"`
		MinLikesReceivedDays   int  `json:"min_likes_received_days"`
		NumLikesReceivedUsers  int  `json:"num_likes_received_users"`
		MinLikesReceivedUsers  int  `json:"min_likes_received_users"`
		PenaltyCounts          struct {
			Silenced  int `json:"silenced"`
			Suspended int `json:"suspended"`
			Total     int `json:"total"`
		} `json:"penalty_counts"`
	} `json:"tl3_requirements"`
	Groups []struct {
		ID                        int         `json:"id"`
		Automatic                 bool        `json:"automatic"`
		Name                      string      `json:"name"`
		DisplayName               string      `json:"display_name"`
		UserCount                 int         `json:"user_count"`
		MentionableLevel          int         `json:"mentionable_level"`
		MessageableLevel          int         `json:"messageable_level"`
		VisibilityLevel           int         `json:"visibility_level"`
		PrimaryGroup              bool        `json:"primary_group"`
		Title                     interface{} `json:"title"`
		GrantTrustLevel           interface{} `json:"grant_trust_level"`
		IncomingEmail             interface{} `json:"incoming_email"`
		HasMessages               bool        `json:"has_messages"`
		FlairURL                  interface{} `json:"flair_url"`
		FlairBgColor              interface{} `json:"flair_bg_color"`
		FlairColor                interface{} `json:"flair_color"`
		BioRaw                    interface{} `json:"bio_raw"`
		BioCooked                 interface{} `json:"bio_cooked"`
		BioExcerpt                interface{} `json:"bio_excerpt"`
		PublicAdmission           bool        `json:"public_admission"`
		PublicExit                bool        `json:"public_exit"`
		AllowMembershipRequests   bool        `json:"allow_membership_requests"`
		FullName                  interface{} `json:"full_name"`
		DefaultNotificationLevel  int         `json:"default_notification_level"`
		MembershipRequestTemplate interface{} `json:"membership_request_template"`
		MembersVisibilityLevel    int         `json:"members_visibility_level"`
		CanSeeMembers             bool        `json:"can_see_members"`
		CanAdminGroup             bool        `json:"can_admin_group"`
		PublishReadState          bool        `json:"publish_read_state"`
	} `json:"groups"`
}

// QueryResult the basic query data returned
type QueryResult struct {
	Success     bool          `json:"success"`
	Errors      []interface{} `json:"errors"`
	Duration    float64       `json:"duration"`
	ResultCount int           `json:"result_count"`
	Params      struct {
		Page string `json:"page"`
	} `json:"params"`
	Columns      []string        `json:"columns"`
	DefaultLimit int             `json:"default_limit"`
	Rows         [][]interface{} `json:"rows"`
}

// Results the actual data in the table
type Results struct {
	ID           int     `json:"user_id"`
	Username     string  `json:"username"`
	Name         string  `json:"name"`
	CreatedAt    time.Time  `json:"created_at"`
	LastSeen     time.Time  `json:"last_seen_at"`
	IPAddress    string  `json:"ip_address"`
	PostCount    int     `json:"post_count"`
	PostsCreated int     `json:"posts_created"`
	PostsRead    int     `json:"posts_read"`
	SolvedCount  int     `json:"solved_count"`
	AverageScore float64 `json:"average_score"`
}

var outFile *os.File
var err error

func trimLastChar(s string) string {
	r, size := utf8.DecodeLastRuneInString(s)
	if r == utf8.RuneError && (size == 0 || size == 1) {
		size = 0
	}
	return s[:len(s)-size]
}

func getUserName(id int) (DiscourseUser, error) {
	du := DiscourseUser{}
	var userURL = "https://forum.camunda.org/admin/users/" + strconv.Itoa(id) + ".json"
	var client = &http.Client{}
	req, err := http.NewRequest("GET", userURL, nil)
	if err != nil {
		return du, err
	}
	req.Header.Set("Content-Type", "multipart/form-data")
	req.Header.Set("Api-Key", APIKey)
	req.Header.Set("Api-Username", APIUser)
	req.Header.Set("Accept", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return du, err
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return du, err
	}
	res.Body.Close()
	_ = json.Unmarshal(data, &du)
	return du, nil

}

func main() {

	outFilePtr := flag.String("out", "out.csv", "an output file.")
	queryType := flag.String("query", "lurker", "Query Type: [lurker|top50|top50poster|solutions")
	durationPtr := flag.Int("months", 1, "Number of Months to query")
	flag.Parse()
	var queryString = ""

	var qt string
	if queryType != nil {
		qt = *queryType
	}
	switch qt {
	case "lurker":
		queryString = lurkers
		break
	case "top50poster":
		queryString = top50Posters
		break
	case "solutions":
		queryString = mostSolutions
		break
	case "top50":
		queryString = top50
		break
	}
	if *outFilePtr == "" {
		outFile, err = os.Create("DiscourseQuery.csv")
		if err != nil {
			log.Fatal("Could not open file ", err)
		}
	} else {
		outFile, err = os.Create(*outFilePtr)
		if err != nil {
			log.Fatal("Open File Error", err)
		}
	}
	defer outFile.Close()

	formValues := url.Values{}
	formValues.Set("months_ago", strconv.Itoa(*durationPtr))
	//data.Set("bar", "baz")

	var DefaultClient = &http.Client{}
	urlPlace := "https://forum.camunda.org/admin/plugins/explorer/queries/" + queryString + "/run"
	request, err := http.NewRequest("POST", urlPlace, strings.NewReader(formValues.Encode()))
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("Content-Type", "multipart/form-data")
	request.Header.Set("Api-Key", APIKey)
	request.Header.Set("Api-Username", APIUser)
	request.Header.Set("Accept", "application/json")
	var FinalResults []Results
	var oData = QueryResult{}

	res, err := DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != 200 {
		log.Fatal(res.StatusCode)
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal("Got no Data")
	}
	res.Body.Close()
	_ = json.Unmarshal(data, &oData)

	var usernameColumn int
	for z := 0; z < len(oData.Columns); z++ {
		if oData.Columns[z] == "id" || oData.Columns[z] == "user_id" {
			usernameColumn = z
			break
		}
	}
	l := len(oData.Rows)
	//fmt.Printf("\nRequest Number:\t\t%d\nNumber of records: %d\n", x-1, l)
	FinalResults = make([]Results, l)
	for i := 0; i < len(oData.Rows); i++ {
		//fmt.Println("Reading Record: " + strconv.Itoa(i))
		rd := oData.Rows[i]
		info, err := getUserName(int(rd[usernameColumn].(float64)))
		if err != nil {
			log.Fatal(err)
		}
		FinalResults[i].Username = info.Username
		FinalResults[i].Name = info.Name
		FinalResults[i].LastSeen = info.LastSeenAt
		FinalResults[i].CreatedAt = info.CreatedAt
		FinalResults[i].IPAddress = info.IPAddress
		FinalResults[i].PostCount = info.PostCount
		FinalResults[i].PostsCreated = info.PostCount
		FinalResults[i].PostsRead = info.PostsReadCount
		FinalResults[i].ID = info.ID
		for b := 0; b < len(rd); b++ {
			dType := fmt.Sprintf("%T", rd[b])
			switch dType {
			case "string":

			case "float64":
				if strings.Contains(oData.Columns[b], "average") {
					FinalResults[i].AverageScore = rd[b].(float64)
				} else if strings.Contains(oData.Columns[b], "solved") {
					FinalResults[i].SolvedCount = int(rd[b].(float64))
				} else {
				}
			case "int":
			default:
				fmt.Printf("Type %T Not defined\n", reflect.TypeOf(rd[b]))
			}

		}
	}
	file, _ := json.MarshalIndent(FinalResults, "", " ")
	_ = ioutil.WriteFile("test.json", file, 0644)
}
