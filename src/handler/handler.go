package handler 

// Handler is an adapter to wrap receieve methods 
type Handler struct func Handle(m message.Message) (n int, e error)

func (h Handler) Handle(m message.Message) (n int, e error) { 
	n, e :=  h(m)
	if e != nil { 
		return 0, e
	}
	return n, nil 
}