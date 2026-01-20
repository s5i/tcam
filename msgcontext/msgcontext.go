package msgcontext

type Message struct {
	Name    string
	Message string
}

type Context struct {
	Messages chan Message
}

func NewContext(size int) *Context {
	return &Context{
		Messages: make(chan Message, size),
	}
}

func (c *Context) Put(name, msg string) {
	select {
	case c.Messages <- Message{Name: name, Message: msg}:
	default:
		<-c.Messages
		c.Messages <- Message{Name: name, Message: msg}
	}
}

func (c *Context) Pop() []Message {
	var msgs []Message
	for {
		select {
		case msg := <-c.Messages:
			msgs = append(msgs, msg)
		default:
			return msgs
		}
	}
}

func (c *Context) Clear() {
	for {
		select {
		case <-c.Messages:
		default:
			return
		}
	}
}
