package main

import (
	"net/http"
	"log"
	// "runtime"
	"flag"
	"github.com/vaughan0/go-ini"
	"os"
)
type ConfigGlobal struct {
	bindAddress			string
	bindAddressSSL	string
	peerAddress			string
	enrollId				string
	enrollSecret		string
	chaincodePath		string
	DEVELOPMENT     bool
}

var (
	config *ConfigGlobal
)

func initialize() {
	config = &ConfigGlobal{
		bindAddress:			"0.0.0.0:44400",
		bindAddressSSL:		"0.0.0.0:44433",
		peerAddress:			"192.168.0.157:7054",
		DEVELOPMENT:			true,
	}
	var conf = flag.String("f", "", "Config file")
	var developmentConf = "0"
	flag.Parse()
	if len(*conf) > 0 {
		confFile, e := ini.LoadFile(*conf)
		if e != nil {
			log.Panicln(e)
		}
		if bindAddr, ok := confFile.Get("", "bindAddress"); ok {
			config.bindAddress = bindAddr
		}
		if bindAddrSSL, ok := confFile.Get("", "bindAddressSSL"); ok {
			config.bindAddressSSL = bindAddrSSL
		}
		if peerAddress, ok := confFile.Get("", "peerAddress"); ok {
			config.peerAddress = peerAddress
		}
		if enrollId, ok := confFile.Get("", "enrollId"); ok {
			config.enrollId = enrollId
		}
		if enrollSecret, ok := confFile.Get("", "enrollSecret"); ok {
			config.enrollSecret = enrollSecret
		}
		if chaincodePath, ok := confFile.Get("", "chaincodePath"); ok {
			config.chaincodePath = chaincodePath
		}
		if dev, ok := confFile.Get("", "DEVELOPMENT"); ok {
			developmentConf = dev
		}
	}
	config.DEVELOPMENT = (os.Getenv("DEVELOPMENT") == "1" || developmentConf == "1")
}

func registerHandler() {
	http.HandleFunc("/api/upkey", handleUpKey)
	http.HandleFunc("/api/getkey", handleGetKey)
	http.HandleFunc("/api/checkalive", handleCheckAlive)
	http.HandleFunc("/api/test", handleTest)
}

func main() {
	// runtime.GOMAXPROCS(runtime.NumCPU())	//multithread configuration.
	initialize()
	log.Println("Registrar to peer");
	// TODO: Registrar login to Peer if fail => exit
	if err:=registrar(); err != nil {
		log.Println(err);
		return;
	}
	// TODO: set chaincodeID global variable by deploy function.
	log.Println("Waiting to get chaincodeID from Peer. Keep patient ....")
	if err:=deploy(); err != nil {
		log.Println(err);
		return;
	}
	log.Printf("chaincodeID = %s", chaincodeID)
	log.Printf("SmartAgent is listening at: %v", config.bindAddress)
	registerHandler()
	if err:=http.ListenAndServe(config.bindAddress, nil); err!=nil {
		log.Panicln(err)
	}
}
