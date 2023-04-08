package types

type Points struct {
	Invite    float32 `envconfig:"INVITE" default:"50"`
	AddReview float32 `envconfig:"ADD_REVIEW" default:"10"`
}

type Config struct {
	Port   int    `envconfig:"PORT" default:"8000"`
	Points Points `envconfig:"POINTS"`
	DB     struct {
		Host     string `envconfig:"HOST" default:"127.0.0.1"`
		Port     int    `envconfig:"PORT" default:"3306"`
		DB       string `envconfig:"DB" default:"rewards"`
		Username string `envconfig:"USERNAME" default:"user"`
		Password string `envconfig:"PASSWORD" default:"pass"`
	} `envconfig:"DB"`
}
