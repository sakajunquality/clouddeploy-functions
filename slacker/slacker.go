package slacker

type Slacker struct {
	token       string
	channel     string
	stateBucket string
}

func NewSlacker(token, channel string) *Slacker {
	return &Slacker{
		token:   token,
		channel: channel,
	}
}

func (s *Slacker) SetStateBucket(bucket string) {
	s.stateBucket = bucket
}

func (*Slacker) post(msg, ts string) (*string, error) {
	return nil, nil
}
