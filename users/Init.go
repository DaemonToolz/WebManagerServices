package main

import (
	"log"
	"net/http"
	// Git repos here
)

//var redisConnector *RedisClientWrapper

func main() {

	initRabbitMq()

	Wrapper = ArangoWrapper{}
	log.Println("Preparing a new Arango Driver")

	Wrapper.initDriver("http://localhost:8529", "user-service", "password")

	log.Println("Driver created")
	/*
		var info UserInfo
		log.Println("<<<<<<<< TEST ADD USER")
		for value := 1; value < 50; value++ {

			info = UserInfo{
				Username:  fmt.Sprintf("%s%d", "user", value),
				Email:     fmt.Sprintf("%s%d%s", "user", value, "@email.coooooooooooooooooooooom"),
				RealName:  "REAL NAME",
				CreatedAt: time.Now(),
			}

			log.Println("<<<<<<<< TEST ADD ", info)
			go Wrapper.Create(ObjectToMap(info))
		}

		log.Println("<<<<<<<< TEST ADD RELATIONS")
		var relation RelationModel
		for value := 1; value < 50; value++ {
			if value < 45 {
				for subvalue := value + 1; subvalue < value+5; subvalue++ {
					relation = RelationModel{
						From:     fmt.Sprintf("%s%d", "user", value),
						To:       fmt.Sprintf("%s%d", "user", subvalue),
						Relation: string(RELATION_FRIEND),
					}
					go Wrapper.SetRelation(relation)
					to := relation.To
					relation.To = relation.From
					relation.From = to
					go Wrapper.SetRelation(relation)
				}
			}

		}
	*/
	//log.Printf("Opening the database %s %d", address, dbIndex)
	log.Println("Database ready")
	defer Wrapper.Close()
	router := NewRouter()
	log.Fatal(http.ListenAndServe(":10840", router))

}
