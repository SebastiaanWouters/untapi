package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tools/cron"
	"github.com/pocketbase/pocketbase/tools/types"
	"github.com/spf13/cobra"
)


const (
	baseURL  = "https://api.untappd.com/v4/user/info/"
)

type Response struct {
	Response struct {
		User struct {
			Firstname  string `json:"first_name"`
			Lastname  string `json:"last_name"`
			Avatar  string `json:"user_avatar"`
			Stats struct {
				Beers int `json:"total_beers"`
				Checkins int `json:"total_checkins"`
				Badges  int `json:"total_badges"`
				
				// Add more fields as per your requirement
			} `json:"stats"`
		} `json:"user"`
	} `json:"response"`
}

var _ models.Model = (*User)(nil)

type User struct {
    models.BaseModel

    Id       string         `db:"id" json:"id"`
    Username        string         `db:"username" json:"username"`
	Firstname        string         `db:"firstname" json:"firstname"`
	Lastname        string         `db:"lastname" json:"lastname"`
	Avatar        string         `db:"avatar" json:"avatar"`
    Beers int         `db:"beers" json:"beers"`
	Badges int         `db:"badges" json:"badges"`
	Checkins int         `db:"checkins" json:"checkins"`
}

func (m *User) TableName() string {
    return "untappd_data" // the name of your collection
}

func fetchUserInfo(username string) (*Response, error) {
	access_token := os.Getenv("ACCESS_TOKEN")
	url := fmt.Sprintf("%s%s?access_token=%s", baseURL, username, access_token)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	
	if err != nil {
		return nil, err
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func updateUserData(app *pocketbase.PocketBase) {
		usersEnv := os.Getenv("USERS")
    
		// Split it into a slice of usernames
		users := strings.Split(usersEnv, ",")
		for _, username := range users {
					userInfo, err := fetchUserInfo(username)
					if err != nil {
						fmt.Println("Error:", err)
						continue
					}
					fmt.Println("fetched")

					record, err := app.Dao().FindRecordById("untappd_data", username)

					if err != nil {
						collection, err := app.Dao().FindCollectionByNameOrId("untappd_data")
						if err != nil {
							continue
						}

						fmt.Println("collection found")

						record := models.NewRecord(collection)
						record.Set("id", username)
						record.Set("username", username)
						record.Set("beers", userInfo.Response.User.Stats.Beers)
						record.Set("checkins", userInfo.Response.User.Stats.Checkins)
						record.Set("badges", userInfo.Response.User.Stats.Badges)
						record.Set("firstname", userInfo.Response.User.Firstname)
						record.Set("lastname", userInfo.Response.User.Lastname)
						record.Set("avatar", userInfo.Response.User.Avatar)

						if err := app.Dao().SaveRecord(record); err != nil {
							continue 
						}

						fmt.Println("record saved")

						continue
					}

					record.Set("id", username)
					record.Set("username", username)
					record.Set("beers", userInfo.Response.User.Stats.Beers)
					record.Set("checkins", userInfo.Response.User.Stats.Checkins)
					record.Set("badges", userInfo.Response.User.Stats.Badges)
					record.Set("firstname", userInfo.Response.User.Firstname)
					record.Set("lastname", userInfo.Response.User.Lastname)
					record.Set("avatar", userInfo.Response.User.Avatar)

					if err := app.Dao().SaveRecord(record); err != nil {
						continue 
					}
					
					}
	}


func main() {
    app := pocketbase.New()

	err := godotenv.Load()
  	if err != nil {
    	log.Fatal("Error loading .env file")
  	}

	 app.RootCmd.AddCommand(&cobra.Command{
        Use: "init",
        Run: func(cmd *cobra.Command, args []string) {
			collection := &models.Collection{
            Name:       "untappd_data",
			Type:       models.CollectionTypeBase,
			ListRule:   types.Pointer(""),
			ViewRule:   types.Pointer(""),
			CreateRule: types.Pointer(""),
			UpdateRule: types.Pointer(""),
			DeleteRule: types.Pointer(""),
			Schema:     schema.NewSchema(
				&schema.SchemaField{
					Name:     "username",
					Type:     schema.FieldTypeText,
					Required: true,
				},
				&schema.SchemaField{
					Name:     "lastname",
					Type:     schema.FieldTypeText,
					Required: true,
				},
				&schema.SchemaField{
					Name:     "firstname",
					Type:     schema.FieldTypeText,
					Required: true,
				},
				&schema.SchemaField{
					Name:     "avatar",
					Type:     schema.FieldTypeText,
					Required: true,
				},
				&schema.SchemaField{
					Name:     "beers",
					Type:     schema.FieldTypeNumber,
					Required: true,
				},
				&schema.SchemaField{
					Name:     "checkins",
					Type:     schema.FieldTypeNumber,
					Required: true,
				},
				&schema.SchemaField{
					Name:     "badges",
					Type:     schema.FieldTypeNumber,
					Required: true,
				},
			),
		}

		// the id is autogenerated, but you can set a specific one if you want to:
		// collection.SetId("...")

		if err := app.Dao().SaveCollection(collection); err != nil {
			return
		}

		updateUserData(app)

        },
    })

	

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
        scheduler := cron.New()

        // updates user data from api
        scheduler.MustAdd("updateUserData", "*/5 * * * *", func() {
				updateUserData(app)
		})

        scheduler.Start()

        return nil
    })

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {

		e.Router.GET("/api/stats", func(c echo.Context) error {
            // records, err := app.Dao().FindRecordsByIds("untappd_data", users)

			users := []User{}

            app.Dao().DB().
    		Select("untappd_data.*").
    		From("untappd_data").
			All(&users)

			return c.JSON(200, map[string][]User{"stats": users})
           
        })

        return nil
    })

    if err := app.Start(); err != nil {
        log.Fatal(err)
    }
}