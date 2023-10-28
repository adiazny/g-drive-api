package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/driveactivity/v2"
	"google.golang.org/api/option"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// Saves a file to a file path.
func saveFile(path string, file []byte) {
	fmt.Printf("Saving file to: %s\n", path)
	err := os.WriteFile("strong.csv", file, 0666)
	if err != nil {
		log.Fatalf("Unable to write file: %v", err)
	}
}

// Returns a string representation of the first elements in a list.
func truncated(array []string) string {
	return truncatedTo(array, 2)
}

// Returns a string representation of the first elements in a list.
func truncatedTo(array []string, limit int) string {
	var contents string
	var more string
	if len(array) <= limit {
		contents = strings.Join(array, ", ")
		more = ""
	} else {
		contents = strings.Join(array[0:limit], ", ")
		more = ", ..."
	}
	return fmt.Sprintf("[%s%s]", contents, more)
}

// Returns the name of a set property in an object, or else "unknown".
func getOneOf(m interface{}) string {
	v := reflect.ValueOf(m)
	for i := 0; i < v.NumField(); i++ {
		if !v.Field(i).IsNil() {
			return v.Type().Field(i).Name
		}
	}
	return "unknown"
}

// Returns a time associated with an activity.
func getTimeInfo(activity *driveactivity.DriveActivity) string {
	if activity.Timestamp != "" {
		return activity.Timestamp
	}
	if activity.TimeRange != nil {
		return activity.TimeRange.EndTime
	}
	return "unknown"
}

// Returns the type of action.
func getActionInfo(action *driveactivity.ActionDetail) string {
	return getOneOf(*action)
}

// Returns user information, or the type of user if not a known user.
func getUserInfo(user *driveactivity.User) string {
	if user.KnownUser != nil {
		if user.KnownUser.IsCurrentUser {
			return "people/me"
		}
		return user.KnownUser.PersonName
	}
	return getOneOf(*user)
}

// Returns actor information, or the type of actor if not a user.
func getActorInfo(actor *driveactivity.Actor) string {
	if actor.User != nil {
		return getUserInfo(actor.User)
	}
	return getOneOf(*actor)
}

// Returns information for a list of actors.
func getActorsInfo(actors []*driveactivity.Actor) []string {
	actorsInfo := make([]string, len(actors))
	for i := range actors {
		actorsInfo[i] = getActorInfo(actors[i])
	}
	return actorsInfo
}

// Returns the type of a target and an associated title.
func getTargetInfo(target *driveactivity.Target) string {
	if target.DriveItem != nil {
		return fmt.Sprintf("driveItem:\"%s\"", target.DriveItem.Title)
	}
	if target.Drive != nil {
		return fmt.Sprintf("drive:\"%s\"", target.Drive.Title)
	}
	if target.FileComment != nil {
		parent := target.FileComment.Parent
		if parent != nil {
			return fmt.Sprintf("fileComment:\"%s\"", parent.Title)
		}
		return "fileComment:unknown"
	}
	return getOneOf(*target)
}

// Returns information for a list of targets.
func getTargetsInfo(targets []*driveactivity.Target) []string {
	targetsInfo := make([]string, len(targets))
	for i := range targets {
		targetsInfo[i] = getTargetInfo(targets[i])
	}
	return targetsInfo
}

func main() {
	ctx := context.Background()

	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	//google.CredentialsFromJSON()
	//google.ConfigFromJSON()
	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b,
		drive.DriveReadonlyScope,
		driveactivity.DriveActivityReadonlyScope,
		drive.DriveFileScope,
		drive.DriveMetadataReadonlyScope,
		drive.DriveMetadataScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	// config := &oauth2.Config{
	// 	ClientID:     "XYZ",
	// 	ClientSecret: "XYZ",
	// 	Scopes: []string{
	// 		drive.DriveReadonlyScope,
	// 		driveactivity.DriveActivityReadonlyScope,
	// 	},
	// 	Endpoint: oauth2.Endpoint{
	// 		AuthURL:  "https://accounts.google.com/o/oauth2/auth",
	// 		TokenURL: "https://oauth2.googleapis.com/token",
	// 	},
	// 	RedirectURL: "http://localhost",
	// }

	//===================================
	// Drive API

	client := getClient(config)

	driveService, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}

	// Search for files and folders guide: https://developers.google.com/drive/api/guides/search-files
	filesByTime, err := driveService.Files.List().PageSize(5).OrderBy("createdTime desc").Q("name = 'strong.csv'").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve file list: %v", err)
	}

	fmt.Println("**FILES COUNT**", len(filesByTime.Files))

	fmt.Println("**FILES BY TIME**", filesByTime.Files[0].Id)

	// token, err := driveService.Changes.GetStartPageToken().Do()
	// if err != nil {
	// 	log.Fatalf("Unable to retrieve Drive client: %v", err)
	// }

	//fmt.Println("***StartPageToken***", token)

	// changes, _ := driveService.Changes.List("547782").Do()

	// fmt.Println("CHANGES", changes.Changes)

	// for _, c := range changes.Changes {
	// 	fmt.Println("***Changes***")
	// 	fmt.Println("File name: " + c.File.Name)
	// 	fmt.Println("Drive Id: " + c.DriveId)
	// }

	// fileList, err := driveService.Files.List().
	// 	Q("mimeType= 'application/vnd.google-apps.folder'").
	// 	SupportsAllDrives(true).PageSize(10).Do()
	// // Fields("mimeType= 'application/vnd.google-apps.folder'").Do()
	// if err != nil {
	// 	log.Fatalf("Unable to retrieve files: %v", err)
	// }

	// fmt.Println("Files:")
	// if len(fileList.Files) == 0 {
	// 	fmt.Println("No files found.")
	// } else {
	// 	for _, i := range fileList.Files {
	// 		fmt.Printf("Name: %s, ID: %s\n", i.Name, i.Id)
	// 	}
	// }

	downloadedFile, err := driveService.Files.Get("1iDQ7iawuFJax2-fvuwhoWMQSuM0DfBuh").Download()
	if err != nil {
		log.Fatalf("Unable to download files: %v", err)
	}

	defer downloadedFile.Body.Close()

	fileBytes, err := io.ReadAll(downloadedFile.Body)
	if err != nil {
		log.Fatalf("Unable to download files: %v", err)
	}

	saveFile("strong.csv", fileBytes)

	//===================================
	// Drive Activity API

	driveActivityClient := getClient(config)

	driveActivityService, err := driveactivity.NewService(ctx, option.WithHTTPClient(driveActivityClient))
	if err != nil {
		log.Fatalf("Unable to retrieve driveactivity Client %v", err)
	}

	//q := driveactivity.QueryDriveActivityRequest{PageSize: 5, Filter: "detail.action_detail_case:CREATE"}

	query := driveactivity.QueryDriveActivityRequest{Filter: "time >= '2023-01-01T00:00:00Z' AND detail.create.file.name = 'strong.csv'"}

	resp, err := driveActivityService.Activity.Query(&query).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve list of activities. %v", err)
	}

	fmt.Println("")
	// for _, activity := range resp.Activities {
	// 	fmt.Printf("Timestamp: %+v ", activity.Timestamp)

	// 	for _, target := range activity.Targets {
	// 		fmt.Printf("%+v\n", target.DriveItem.Title)
	// 	}

	// 	if activity.PrimaryActionDetail.Create != nil {
	// 		if activity.PrimaryActionDetail.Create.Upload != nil {
	// 			//fmt.Printf("%+v\n", activity.Targets)
	// 		}
	// 	}

	//}

	// fmt.Println("Recent Activity:")
	// if len(resp.Activities) > 0 {
	// 	for _, a := range resp.Activities {
	// 		time := getTimeInfo(a)
	// 		action := getActionInfo(a.PrimaryActionDetail)
	// 		actors := getActorsInfo(a.Actors)
	// 		targets := getTargetsInfo(a.Targets)
	// 		fmt.Printf("%s: %s, Action: %s, %s\n", time, truncated(actors), action, truncated(targets))
	// 	}
	// } else {
	// 	fmt.Print("No activity.")
	// }

	fmt.Println("Recent Activity:")
	if len(resp.Activities) > 0 {
		json, err := resp.Activities[0].MarshalJSON()
		if err != nil {
			log.Fatalf("Unable to marshall json. %v", err)
		}

		fmt.Print(string(json))

		// for _, a := range resp.Activities {

		// 	s, err := a.MarshalJSON()
		// 	if err != nil {
		// 		log.Fatalf("Unable to marshall json. %v", err)
		// 	}

		// 	//fmt.Printf("activities: %s\n", string(s))
		// 	fmt.Print(string(s))

		// }
	} else {
		fmt.Print("No activity.")
	}
}
