package domain

import (
	"encoding/json"
	"fmt"
)

type User struct {
	Name string `json:"name"`
}

type Played struct {
	Score int `json:"score"`
}

func (user User) ToCypher(char string) Cypher {
	q, _ := json.Marshal(user)
	properties := string(q)
	return fmt.Sprintf("(%s:User %s)", char, properties)
}

func (r Played) ToCypherRight(char string) Cypher {
	q, _ := json.Marshal(r)
	properties := string(q)
	return fmt.Sprintf("-[%s:Played %s]->", char, properties)
}

func (r Played) ToCypherLeft(char string) Cypher {
	q, _ := json.Marshal(r)
	properties := string(q)
	return fmt.Sprintf("<-[%s:Played %s]-", char, properties)
}
