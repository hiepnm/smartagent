//a crude mock database
package main

import "fmt"

var currentId int

var discs Discs

// Give us some seed data
func init() {

    RepoCreateDisc(Disc{
      KeyId: "12345678901234567890123456789012",
      Key: "09876543211234567890098765432112",
      Name: "Fast and Furious",
      Info: "This is a compact disc infomation",
      CDId: "1",
    })
    RepoCreateDisc(Disc{
      KeyId: "ID_Romeo_and_Juliet",
      Key: "ThisIsASecretKeyToEncryptTheCDFile",
      Name: "Romeo and Juliet",
      Info: "This is a compact disc infomation",
      CDId: "ID_Romeo_and_Juliet",
    })
}

func RepoFindDisc(cdid string) Disc {
    for _, d := range discs {
        if d.CDId == cdid {
            return d
        }
    }
    // return empty Disc if not found
    return Disc{}
}

func RepoCreateDisc(d Disc) Disc {
    currentId += 1
    d.Id = currentId
    discs = append(discs, d)
    return d
}

func RepoDestroyDisc(id int) error {
    for i, d := range discs {
        if d.Id == id {
            discs = append(discs[:i], discs[i+1:]...)
            return nil
        }
    }
    return fmt.Errorf("Could not find Disc with id of %d to delete", id)
}
