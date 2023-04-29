package models

type PushContentReq struct {
	Token   string `json:"token"`
	Content string `json:"content"`
}

type PushContentRes struct {
	Hash    string `json:"hash,omitempty"`
	ErrNo   int    `json:"errno,omitempty"`
	Message string `json:"message,omitempty"`
}

type TopicActionReq struct {
	Action string `json:"action"`
	Param  string `json:"param"`
	Token  string `json:"token"`
}

type AskQuestionRes struct {
	Content string `json:"content,omitempty"`
	ErrNo   int    `json:"errno,omitempty"`
	Message string `json:"message,omitempty"`
}

type GetTopicRes struct {
	Topics  []TopicWithRelative `json:"topics,omitempty"`
	ErrNo   int                 `json:"errno,omitempty"`
	Message string              `json:"message,omitempty"`
}

type GetTopicRelativeRes struct {
	Topics  []TopicWithRelative `json:"topics,omitempty"`
	ErrNo   int                 `json:"errno,omitempty"`
	Message string              `json:"message,omitempty"`
}

type SummaryRes struct {
	Content string `json:"content,omitempty"`
	ErrNo   int    `json:"errno,omitempty"`
	Message string `json:"message,omitempty"`
}

type TopicWithRelative struct {
	Topic    string
	Relative string
}
