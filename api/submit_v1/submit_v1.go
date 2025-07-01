package submit_v1;
// MC Server to Shard



type Submit struct {
	Timestamp  int64  `json:"Timestamp"`
	ServerUID  string `json:"ServerUID"`
	Platform   string `json:"Platform"`
	NumPlayers int16  `json:"NumPlayers"`
}

type SubmitReply struct {
	Status uint8 `json:"Status"`
}
