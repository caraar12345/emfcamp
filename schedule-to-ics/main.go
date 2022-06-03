package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	ics "github.com/arran4/golang-ical"
	"github.com/joho/godotenv"
)

var emf_session_token = "session=" + envVar("EMF_SESSION_TOKEN")

func envVar(key string) string {
	outVar := os.Getenv(key)
	if outVar == "" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalf("Error loading .env file: %v\n", err)
		}
		outVar = os.Getenv(key)
	}
	return outVar
}

func handleScheduleRequest(w http.ResponseWriter, r *http.Request) {
	scheduleMap := collectScheduleJson()
	icsString := assembleIcalFile(scheduleMap)

	w.Header().Add("Content-Type", "text/calendar; charset=utf-8")
	w.Write([]byte(icsString))
}

func collectScheduleJson() []ScheduleItem {
	httpClient := &http.Client{}
	var schedule []ScheduleItem

	req, err := http.NewRequest("GET", scheduleUrl, nil)
	if err != nil {
		log.Println("error creating request:", err)
	}

	req.Header.Add("Cookie", emf_session_token)

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Println("error getting schedule:", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("error reading request body:", err)
	}

	err = json.Unmarshal(body, &schedule)
	if err != nil {
		log.Println("error unmarshaling json:", err)
	}
	return schedule
}

func assembleIcalFile(scheduleMap []ScheduleItem) string {
	cal := ics.NewCalendar()
	for _, scheduleItem := range scheduleMap {
		if !scheduleItem.IsFave {
			continue
		}

		event := cal.AddEvent(fmt.Sprintf("%v@favourites.emf.adhd.energy", strconv.Itoa(scheduleItem.ID)))

		var newSummary string
		if scheduleItem.MayRecord {
			newSummary = fmt.Sprintf("📹 %v", scheduleItem.Title)
		} else {
			newSummary = scheduleItem.Title
		}

		event.SetSummary(newSummary)

		//var contentNote string
		//if scheduleItem.ContentNote == "" {
		//	contentNote = "No content note provided"
		//} else {
		//	contentNote = scheduleItem.ContentNote
		//}
		//event.SetDescription(fmt.Sprintf("%v\n%v\n%v\n", contentNote, scheduleItem.Description, beingRecorded))
    event.SetDescription(fmt.Sprintf("%v\nURL: %v", scheduleItem.Venue, scheduleItem.Link))
		event.SetURL(scheduleItem.Link)

		if len(scheduleItem.Latlon) != 0 {
			lat := fmt.Sprintf("%f", scheduleItem.Latlon[0])
			lon := fmt.Sprintf("%f", scheduleItem.Latlon[1])
			event.SetLocation(fmt.Sprintf("%v, %v", lat, lon))
		}

		event.SetStartAt(time.Time(scheduleItem.StartDate).Add(-time.Hour))
		event.SetEndAt(time.Time(scheduleItem.EndDate).Add(-time.Hour))
		if scheduleItem.Pronouns == "" {
			scheduleItem.Pronouns = "no pronouns provided"
		}

		event.SetOrganizer(fmt.Sprintf("%v (%v)", scheduleItem.Speaker, scheduleItem.Pronouns))
	}

	return cal.Serialize()
}

func main() {
	var portFlag = flag.Int("p", 8080, "port on which to serve http")
	flag.Parse()
	var port = fmt.Sprintf(":%v", strconv.Itoa(*portFlag))

	log.Println("server started on port", port)
	http.HandleFunc("/schedule.ics", handleScheduleRequest)
	log.Fatal(http.ListenAndServe(port, nil))
}
