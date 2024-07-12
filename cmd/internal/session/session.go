package session

type Session struct {
	SessionName string `json:"sessionName"`
	UserName    string `json:"userName"`
	Host        string `json:"host"`
	Port        int    `json:"port,omitempty"`
	Password    string `json:"password,omitempty"`
}

func New(sessionName, userName, host, password string, port int) Session {
	return Session{
		SessionName: sessionName,
		UserName:    userName,
		Host:        host,
		Port:        port,
		Password:    password,
	}
}
