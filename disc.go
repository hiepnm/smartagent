package main
import (
  "time"
  // "fmt"
)
type Disc struct {
  Id    int       `json:"id"`
  KeyId string    `json:"keyid"`
  Key   string    `json:"key"`
  Name  string    `json:"name"`
  Due   time.Time `json:"due"`
  Info  string    `json:"info"`
  CDId  string    `json:"cdid"`
}

// func (d Disc) String() string {
//   return fmt.Sprintf("Id: %d\nName: %s\nDescription: %s\nKey: %s\nDue: %v", d.Id, d.Name, d.Info, d.Key, d.Due);
// }
type Discs []Disc
