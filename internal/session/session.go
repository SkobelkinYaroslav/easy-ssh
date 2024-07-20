package session

type Session struct {
	SessionName   string `json:"sessionName"`
	UserName      string `json:"userName"`
	Host          string `json:"host"`
	Port          int    `json:"port,omitempty"`
	Password      string `json:"password,omitempty"`
	IsConnectable bool   `json:"isConnectable"`
}

func New(sessionName, userName, host, password string, port int) Session {

	return Session{
		SessionName:   sessionName,
		UserName:      userName,
		Host:          host,
		Port:          port,
		Password:      password,
		IsConnectable: true,
	}
}

func NewDefault() Session {
	return Session{
		SessionName: "Test Session",
		UserName:    "yourName",
		Host:        "1.1.1.1",
		Port:        22,
		Password:    "password",
	}
}

func (s *Session) SetSessionName(sessionName string) {
	s.SessionName = sessionName
}

func (s *Session) SetUserName(userName string) {
	s.UserName = userName
}

func (s *Session) SetHost(host string) {
	s.Host = host
}

func (s *Session) SetPort(port int) {
	s.Port = port
}

func (s *Session) SetPassword(password string) {
	s.Password = password
}

func (s *Session) SetConnectable(state bool) {
	s.IsConnectable = state
}
