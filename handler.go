package main
import (
	"encoding/json"
	"net/http"
	"io"
	"io/ioutil"
	"log"
)
type ResponseCurator struct {
	Status string `json:"status"`
	Data interface{}		`json:"data"`
}

func handleUpKey(w http.ResponseWriter, r *http.Request) {
	var disc Disc
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
  	panic(err)
  }
	if err := r.Body.Close(); err != nil {
  	panic(err)
	}

	if err := json.Unmarshal(body, &disc); err != nil {
	  w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	  w.WriteHeader(422) // unprocessable entity
	  if err := json.NewEncoder(w).Encode(err); err != nil {
      panic(err)
	  }
	}
	d := RepoCreateDisc(disc)
	// TODO: invoke request to addAsset to ledger.
	// ... addAsset
	//
	var ret ResponseCurator
	if err := addAsset(d.CDId); err != nil {
		log.Println(err);
		log.Println("Jump here");
		ret = ResponseCurator{
			"500",
			"Can not Add Asset to Blockchain network",
		}
	} else {
		ret = ResponseCurator{
			"200",
			d,
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(ret); err != nil {
		panic(err)
	}
}

func handleCheckAlive(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	ret := ResponseCurator{
		"200",
		discs,
	}
	if err := json.NewEncoder(w).Encode(ret); err != nil {
		panic(err)
	}
}

func handleGetKey(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
  	panic(err)
  }
	if err := r.Body.Close(); err != nil {
  	panic(err)
	}
	type RequestGetKey struct {
		CDId string    `json:"cdid"`
	}
	var data RequestGetKey
	if err := json.Unmarshal(body, &data); err != nil {
	  w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	  w.WriteHeader(422) // unprocessable entity
	  if err := json.NewEncoder(w).Encode(err); err != nil {
      panic(err)
	  }
	}

	disc := RepoFindDisc(data.CDId)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	ret := ResponseCurator{
		"200",
		disc,
	}
	if err := json.NewEncoder(w).Encode(ret); err != nil {
		panic(err)
	}
}

func handleTest(w http.ResponseWriter, r *http.Request) {
	ret := ResponseCurator{
		"200",
		discs,
	}
	makeJsonResponse(&w, &ret)
}

func handleEnroll(w http.ResponseWriter, r *http.Request) {
	var ret ResponseCurator
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
  	panic(err)
  }

	if err := r.Body.Close(); err != nil {
  	panic(err)
	}

	if err := json.Unmarshal(body, &curator); err != nil {
	  w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	  w.WriteHeader(422) // unprocessable entity
	  if err := json.NewEncoder(w).Encode(err); err != nil {
      panic(err)
	  }
	}
	// TODO: registrar
	if err:=registrar(); err != nil {
		log.Println(err);
		ret = ResponseCurator{"500", err}
		makeJsonResponse(&w, &ret)
		return
	}
	// TODO: get chaincodeID
	if err:=getChaincodeID(); err != nil {
		log.Println(err);
		ret = ResponseCurator{"500",err}
		makeJsonResponse(&w, &ret)
		return
	}

	ret = ResponseCurator{
		"200",
		curator,
	}
	w.WriteHeader(http.StatusCreated)
	makeJsonResponse(&w, &ret)
}

func makeJsonResponse(w *http.ResponseWriter, ret *ResponseCurator) {
	(*w).Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(*w).Encode(*ret); err != nil {
		panic(err)
	}
}
