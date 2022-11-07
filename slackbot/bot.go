package slackbot

type Slackbot struct {
	token       string
	channel     string
	stateBucket string
}

func NewSlackbot(token, channel string) *Slackbot {
	return &Slackbot{
		token:   token,
		channel: channel,
	}
}

func (s *Slackbot) SetStateBucket(bucket string) {
	s.stateBucket = bucket
}
