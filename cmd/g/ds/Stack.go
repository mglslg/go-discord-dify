package ds

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

type Stack struct {
	elements []*discordgo.Message
}

func NewStack() *Stack {
	return &Stack{elements: make([]*discordgo.Message, 0)}
}

func (s *Stack) Push(element *discordgo.Message) {
	s.elements = append(s.elements, element)
}

func (s *Stack) Pop() (*discordgo.Message, error) {
	if len(s.elements) == 0 {
		return nil, fmt.Errorf("stack is empty")
	}

	topIndex := len(s.elements) - 1
	topElement := s.elements[topIndex]
	s.elements = s.elements[:topIndex]

	return topElement, nil
}

// GetBottomElement returns the bottom element of the stack
func (s *Stack) GetBottomElement() (*discordgo.Message, error) {
	if s.Size() != 0 {
		return s.elements[0], nil
	} else {
		return nil, nil
	}
}

func (s *Stack) Size() int {
	return len(s.elements)
}

func (s *Stack) IsEmpty() bool {
	if s.Size() == 0 {
		return true
	}
	return false
}
