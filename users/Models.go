package main

import (
	"bytes"
	"encoding/json"
	"log"
	"strings"
	"time"
	"net/http"
	uuid "github.com/satori/go.uuid"
)

/*
	-------------- Constants
*/
var ValidMailProvider = [...]string{
	"outlook.com",
	"outlook.fr",
	"hotmail.fr",
	"hotmail.com",
	"gmail.com",
}

/*
	--------------- MODELS
*/

type AccountModel struct {
	id       string    `json:"id"`
	name     string    `json:"name"`
	username string    `json:"username"`
	description string `json:"description"`
	creation time.Time `json:"creation"`
	emails   []string  `json:"emails"`
	validated bool 	   `json:"validated"`
}

type AccountCreateRequest struct {
	id       string   `json:"id"`
	name     string   `json:"name"`
	description string`json:"description"`
	username string   `json:"username"`
	emails   []string `json:"emails"`

}

type AccountValidationRequest struct {
	id       string   `json:"id"`
	token	 string	  `json:"token"`
}

type AuthorizationToken struct {
	token	 string	  `json:"token"`
}


func Find(id string, channel chan AccountModel) {
	jchannel := make(chan string)

	defer close(jchannel)
	go redisConnector.Get(id, jchannel)

	account := AccountModel{}

	result := <-jchannel
	log.Printf("Found data, %s", result)

	json.Unmarshal([]byte(result), &account)

	channel <- account

}

func FindId(id string, channel chan AccountModel) {
	jchannel := make(chan string)
	defer close(jchannel)
	go redisConnector.Get(id, jchannel)
	//login := LoginModel{}

	//json.Unmarshal([]byte(<-jchannel), &login)

	Find(id, channel)
}

func UpdateModel(account AccountModel, channel chan<- bool) {
	go redisConnector.Set(account.id, account, channel)
}

func Private_ValidateEmails(emails []string) []string {

	var buffer bytes.Buffer

	valid := make([]string, 0)

	for _, email := range emails {
		contained := false
		for _, provider := range ValidMailProvider {

			buffer.WriteString("@")
			buffer.WriteString(provider)

			contained = strings.Contains(email, buffer.String())

			buffer.Reset()
			if contained {
				break
			}
		}

		if contained {
			valid = append(valid, email)
		}
	}

	return valid
}

func UpdateEmails(id string, emails []string, channel chan<- bool) {
	aChannel := make(chan AccountModel)

	defer close(aChannel)
	go Find(id, aChannel)

	valid := Private_ValidateEmails(emails)

	account := <-aChannel
	account.emails = valid

	go redisConnector.Set(account.id, account, channel)
}

func Remove(id string, channel chan<- bool) {
	go redisConnector.Remove(id, channel)
}

func Create(name string, username string, emails []string, channel chan AccountModel) {
	bChannel := make(chan bool)
	defer close(bChannel)
	uid := uuid.NewV4().String()

	account := AccountModel{
		name:     name,
		username: username,
		description: "",
		id:       uid,
		validated:   false,
		emails:   Private_ValidateEmails(emails),
		creation: time.Now(),
	}

	go redisConnector.Set(uid, account, bChannel)
	if <-bChannel == false {
		channel <- AccountModel{id: "ERR_500INTEX"}
		log.Println("CREATE: Could not register user")
		return
	}

	log.Println("CREATE: ", <-bChannel)
	channel <- account
}

func Validate(post AccountValidationRequest, channel chan bool) {
	bChannel := make(chan bool)
	defer close(bChannel)

	cChannel := make(chan AccountModel)
	defer close(cChannel)

	go Find(post.id, cChannel);


	account := <-cChannel;
	account.validated = true;

	go redisConnector.Set(post.id, account, bChannel)
	if <-bChannel == false {
		channel <- false;
		log.Println("VALIDATE : Could not validate user")
		return
	}

	go CreateSpace(post.id);
	log.Println("VALIDATE: ", <-bChannel)
	channel <- <-bChannel
}

func Reset() {
	redisConnector.Flush()
	redisConnector.Save()
}

func CreateSpace(id string){
	url := "http://localhost:10850/create"
    log.Println("URL:>", url)

	
	anonymousStruct := struct {
		id 	    string
		token   string
	}{
		id,
		"",
	}

	b, err := json.Marshal(anonymousStruct)
    if err != nil {
        log.Println(err)
        return
	}
	
    req, err := http.NewRequest("POST", url,  bytes.NewBuffer(b))
    req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    log.Println("response Status:", resp.Status)
    log.Println("response Headers:", resp.Header)
}