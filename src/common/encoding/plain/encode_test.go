package plain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type ProgrammerInfo struct {
	Name string `plain:"title=[basic]"`
	ID   uint64
	Age  int `plain:"newline"`

	Skill []*skill `plain:"title=[ability]"`
}

type skill struct {
	Name         string
	PracticeYear int
}

func TestMarshal(t *testing.T) {
	p := ProgrammerInfo{}
	p.Name = "John"
	p.ID = uint64(100)
	p.Age = 30
	p.Skill = make([]*skill, 2)
	p.Skill[0] = &skill{"C", 3}
	p.Skill[1] = &skill{"Golang", 1}

	b, _ := Marshal(p)
	expected :=
		`ProgrammerInfo:
  [basic]
  Name: John
  ID: 100
  Age: 30

  [ability]
  Skill:
    0:
      Name: C
      PracticeYear: 3
    1:
      Name: Golang
      PracticeYear: 1
`
	assert.Equal(t, expected, string(b))
}
