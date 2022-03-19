package handlers

//WsJsonResponse define the response that are sent form websocket
type WsJsonResponse struct {
	Action        string   `json:"action"`
	MessageType   string   `json:"message_type"`
	Message       string   `json:"message"`
	ConnectedUser []string `json:"connected_user"`
}

//WsPayload define the payload that are sent to websocket
type WsPayload struct {
	Action   string              `json:"action"`
	UserName string              `json:"username"`
	Message  string              `json:"message"`
	Conn     WebScoketConnection `json:"_"`
}
