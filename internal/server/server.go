package server

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/negarciacamilo/tcp_chat/internal/format"
	"github.com/negarciacamilo/tcp_chat/internal/logger"
	"go.uber.org/zap"
	"io"
	"net"
	"regexp"
	"strings"
)

var joinChan chan *User
var messageChan chan *Message

var users []*User

type User struct {
	id         string
	username   string
	connection net.Conn
}

type Message struct {
	userID     string
	username   string
	message    string
	isPrivate  bool
	toUsername string
}

func StartServer() {
	conn, err := net.Listen("tcp", ":8080")
	if err != nil {
		logger.Panic("error listening", zap.Error(err))
	}
	defer conn.Close()

	fmt.Println("Server listening port 8080")
	joinChan = make(chan *User, 1000)
	messageChan = make(chan *Message, 1000)
	defer close(joinChan)
	defer close(messageChan)

	for {
		c, err := conn.Accept()
		if err != nil {
			logger.Error("tcp server accept error", zap.Error(err))
		}

		logger.Info(fmt.Sprintf("Serving address: %s", c.RemoteAddr().String()))

		u := &User{}
		u.connection = c
		go u.handleConnection()
	}
}

func parseUsername(c net.Conn) (*User, error) {
	var username string
	username, err := bufio.NewReader(c).ReadString('\n')
	if err != nil {
		c.Write(format.RedMessage("There was an error with your username. Error: %s", err.Error()))
		return nil, err
	}

	err = usernameAlreadyTaken(username)
	if err != nil {
		attempts := 0
		for err != nil && attempts != 3 {
			c.Write(format.RedMessage(err.Error() + "\n Please, try with another username:"))
			username, err = bufio.NewReader(c).ReadString('\n')
			err = usernameAlreadyTaken(username)
			attempts++
		}
		if attempts == 3 {
			c.Write(format.RedMessage("Closing connection"))
			c.Close()
			return nil, err
		}
	}

	u := &User{username: strings.ReplaceAll(strings.ReplaceAll(username, "\n", ""), " ", ""), connection: c, id: uuid.New().String()}

	go func() { joinChan <- u }()
	users = append(users, u)
	return u, nil
}

func usernameAlreadyTaken(usr string) error {
	for _, u := range users {
		usr = strings.ToLower(strings.ReplaceAll(usr, "\n", ""))
		if strings.ToLower(u.username) == usr {
			return errors.New(fmt.Sprintf("%s is already taken", usr))
		}
	}
	return nil
}

func (u *User) handleConnection() {
	defer u.connection.Close()
	u.connection.Write(format.ToByte("Welcome to this simple TCP Chat \nYour address is: %s | Online users: %d\nPlease type your desired username: ", u.connection.RemoteAddr().String(), len(users)))
	u, err := parseUsername(u.connection)
	if err != nil {
		return
	}

	go joinMessage()
	go messenger()
	u.connection.Write(format.YellowMessage("[INFO] - Welcome %s\n", u.username))
	fmt.Println(fmt.Sprintf("User %s joined (from %s)", u.username, u.connection.RemoteAddr().String()))
	u.connection.Write(format.PurpleMessage("[INFO] You can type your message anytime you want\n"))
	for {
		msg, err, killCon := parseMessage(u.connection)
		if err != nil {
			u.connection.Write(format.ToByte("%s\n", err.Error()))
			if killCon {
				return
			}
			continue
		}

		msg.userID = u.id
		msg.username = u.username
		go func() { messageChan <- msg }()
	}
}

// This will return the message to send, error if any and a boolean to kill the server if needed
func parseMessage(c net.Conn) (*Message, error, bool) {
	msg, err := bufio.NewReader(c).ReadString('\n')

	switch err {
	case nil:
		m := Message{message: msg}
		err := m.isPrivateMessage()
		if err != nil {
			return nil, err, false
		}
		return &m, nil, false
	case io.EOF:
		logger.Info(fmt.Sprintf("%s closes the connection", c.RemoteAddr().String()))
		return nil, io.EOF, true
	default:
		logger.Error("something unexpected happened parsing the message", zap.String("message", msg))
		return nil, err, true
	}
}

func (m *Message) isPrivateMessage() error {
	if strings.HasPrefix(m.message, "@") {
		pattern := "@([a-zA-Z])\\S*"
		r := regexp.MustCompile(pattern)
		m.toUsername = r.FindString(m.message)
		m.toUsername = strings.ReplaceAll(m.toUsername, "@", "")
		for _, u := range users {
			if strings.ToLower(u.username) == strings.ToLower(m.toUsername) {
				m.isPrivate = true
				m.message = strings.TrimPrefix(m.message, fmt.Sprintf("@%s", m.toUsername))
				return nil
			}
		}
		return errors.New(fmt.Sprintf("Username %s not online", m.toUsername))
	}
	m.isPrivate = false
	return nil
}

func joinMessage() {
	for {
		usr := <-joinChan
		for _, user := range users {
			if user.id != usr.id {
				user.connection.Write(format.CyanMessage("\n[INFO] - %s joined the chat\n", usr.username))
			}
		}
	}
}

func messenger() {
	for {
		msg := <-messageChan
		for _, m := range users {
			if m.id != msg.userID {
				if msg.isPrivate && m.username == msg.toUsername {
					m.connection.Write(format.PurpleMessage("\n[PRIVATE - %s]: %s", msg.username, msg.message))
					continue
				} else if !msg.isPrivate {
					m.connection.Write(format.GrayMessage("\n%s said: %s", msg.username, msg.message))
				}
			}
		}
	}
}
