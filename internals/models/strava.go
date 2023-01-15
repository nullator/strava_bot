package models

type StravaUser struct {
	Token_type    string        `json:"token_type"`
	Expires_at    int           `json:"expires_at"`
	Expires_in    int           `json:"expires_in"`
	Refresh_token string        `json:"refresh_token"`
	Access_token  string        `json:"access_token"`
	Athlete       StravaAthlete `json:"athlete"`
}

type AuthHandler struct {
	ID    string `form:"state"`
	Code  string `form:"code"`
	Scope string `form:"scope"`
}

type StravaAthlete struct {
	Id             int      `json:"id"`             // 2534833,
	Username       string   `json:"username"`       // "nullator_n",
	Resource_state int      `json:"resource_state"` // 2,
	Firstname      string   `json:"firstname"`      // "Vyacheslav",
	Lastname       string   `json:"lastname"`       // "O",
	Bio            string   `json:"bio"`            // "",
	City           string   `json:"city"`           // "Samara",
	State          string   `json:"state"`          // "Samarskaya oblast",
	Country        string   `json:"country"`        // "Russia",
	Sex            string   `json:"sex"`            // "M",
	Premium        bool     `json:"premium"`        // false,
	Summit         bool     `json:"summit"`         // false,
	Created_at     string   `json:"created_at"`     // "2013-07-11T14:30:18Z",
	Updated_at     string   `json:"updated_at"`     // "2021-09-10T16:14:11Z",
	Badge_type_id  int      `json:"badge_type_id"`  // 0,
	Weight         float32  `json:"weight"`         // 62.0,
	Profile_medium string   `json:"profile_medium"` // "https://dgalywyr863hv.cloudfront.net/pictures/athletes/2534833/6923497/2/medium.jpg",
	Profile        string   `json:"profile"`        // "https://dgalywyr863hv.cloudfront.net/pictures/athletes/2534833/6923497/2/large.jpg",
	Friend         []string `json:"friend"`         // null,
	Follower       []string `json:"follower"`       // null
}
