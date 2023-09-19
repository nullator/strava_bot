package models

type StravaUser struct {
	Token_type    string  `json:"token_type"`
	Expires_at    int     `json:"expires_at"`
	Expires_in    int     `json:"expires_in"`
	Refresh_token string  `json:"refresh_token"`
	Access_token  string  `json:"access_token"`
	Athlete       Athlete `json:"athlete"`
}

type AuthHandler struct {
	ID    string `form:"state"`
	Code  string `form:"code"`
	Scope string `form:"scope"`
}

type Athlete struct {
	Id             int64    `json:"id"`             // id в страве
	Username       string   `json:"username"`       // ник
	Resource_state int      `json:"resource_state"` //
	Firstname      string   `json:"firstname"`      //
	Lastname       string   `json:"lastname"`       //
	Bio            string   `json:"bio"`            //
	City           string   `json:"city"`           //
	State          string   `json:"state"`          //
	Country        string   `json:"country"`        //
	Sex            string   `json:"sex"`            //
	Premium        bool     `json:"premium"`        //
	Summit         bool     `json:"summit"`         //
	Created_at     string   `json:"created_at"`     //
	Updated_at     string   `json:"updated_at"`     //
	Badge_type_id  int      `json:"badge_type_id"`  //
	Weight         float32  `json:"weight"`         //
	Profile_medium string   `json:"profile_medium"` // ссылка на фото профиля
	Profile        string   `json:"profile"`        // ссылка на фото профиля
	Friend         []string `json:"friend"`         //
	Follower       []string `json:"follower"`       //
}

type RespondRefreshToken struct {
	Token_type    string `json:"token_type"`
	Access_token  string `json:"access_token"`
	Expires_at    int64  `json:"expires_at"`
	Expires_in    int64  `json:"expires_in"`
	Refresh_token string `json:"refresh_token"`
}

type UploadActivity struct {
	File        string `json:"file"`
	Data_type   string `json:"data_type"`
	Description string `json:"description"`
	Name        string `json:"name"`
}

type RespondUploadActivity struct {
	Id          int64       `json:"id"`
	Id_str      string      `json:"id_str"`
	External_id interface{} `json:"external_id"`
	Error       interface{} `json:"error"`
	Status      string      `json:"status"`
	Activity_id int64       `json:"activity_id"`
}
