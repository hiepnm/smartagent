package main;

import (
  "net/http"
	"bytes"
	"fmt"
  "io/ioutil"
  "encoding/json"
)
var client = &http.Client{}
var chaincodeID = "mycc"

type ErrBlockChain string;
func (e ErrBlockChain) Error() string {
  return fmt.Sprintf("Error on BlockChain Network: %s", string(e));
}
type ChainCodeSuccess struct {
  Status string `json:"status"`
  Message string  `json:"message"`
}
type ChainCodeError struct {
  Code int `json:"code"`
  Message string  `json:"message"`
  Data string `json:"data"`
}
type ChainCodeResponse struct {
  Jsonrpc string  `json:"jsonrpc"`
  Result ChainCodeSuccess `json:"result"`
  Error ChainCodeError `json:"error"`
  Id int `json:"id"`
}
func postData(url, s string) ([]byte, error) {
  jsonStr := []byte(s);
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

  resp, err := client.Do(req)
  if err != nil {
    return nil, err;
  }
	defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return nil, err;
  }
  if (resp.StatusCode != 200) {
    return nil, ErrBlockChain("\n"+resp.Status+"\n"+string(body));
  }
  return body, nil;
}


func registrar() error {
  url := "http://" + config.peerAddress + "/registrar";
  s := "{\"enrollId\":\"" + curator.EnrollID + "\",\"enrollSecret\":\"" + curator.EnrollSecret + "\"}";
  _, err := postData(url, s);
  return err;
}

func getChaincodeID() error {
  url := "http://" + config.peerAddress + "/chaincode";
  template :=
  `
  {
    "jsonrpc": "2.0",
    "method": "deploy",
    "params": {
      "type": 1,
      "chaincodeID": {
        "path": "%s"
      },
      "ctorMsg": {
        "function": "Init"
      },
      "secureContext": "%s"
    },
    "id": 1
  }
  `
  s := fmt.Sprintf(template, config.chaincodePath, curator.EnrollID)
  body, err := postData(url, s);
  if err != nil {
    return err;
  }

  var resp ChainCodeResponse
  if err := json.Unmarshal(body, &resp); err != nil {
    return err
  }
  if resp.Error.Code != 0 {
    return ErrBlockChain(resp.Error.Data);
  }

  curator.ChaincodeID = resp.Result.Message;
  return nil
}

func addAsset(cdId string) error  {
  url := "http://" + config.peerAddress + "/chaincode";
  template :=
  `
  {
    "jsonrpc": "2.0",
    "method": "invoke",
    "params": {
      "type": 1,
      "chaincodeID": {
        "name": "%s"
      },
      "ctorMsg": {
        "function": "assign",
        "args": [
          "%s", "%s"
        ]
      },
      "secureContext": "%s"
    },
    "id": 2
  }
  `
  s := fmt.Sprintf(template, curator.ChaincodeID, cdId, curator.EnrollID, curator.EnrollID)
  var err error;
  _, err = postData(url, s);
  if (err != nil) {
    return err;
  }
  // time.Sleep(5000 * time.Millisecond)
  err = queryAsset(cdId);
  return err;
}

func queryAsset(cdId string) error {
  url := "http://" + config.peerAddress + "/chaincode";
  template :=
  `{
    "jsonrpc": "2.0",
    "method": "query",
    "params": {
      "type": 1,
      "chaincodeID": {
        "name": "%s"
      },
      "ctorMsg": {
        "function": "query",
        "args": [
          "%s"
        ]
      },
      "secureContext": "%s"
    },
    "id": 2
  }`
  s := fmt.Sprintf(template, curator.ChaincodeID, cdId, curator.EnrollID)
  body, err := postData(url, s);
  if err != nil {
    return err;
  }
  var resp ChainCodeResponse
  if err := json.Unmarshal(body, &resp); err != nil {
    return err
  }
  if resp.Error.Code != 0 {
    return ErrBlockChain(resp.Error.Data);
  }
  return nil;
}
