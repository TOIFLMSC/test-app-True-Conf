package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

//User struct
type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func main() {
	e := echo.New()

	e.POST("/users", addUser)
	e.GET("/users", getUsers)
	e.GET("/users/:id", getUser)
	e.PUT("/users/:id", updateUser)
	e.DELETE("/users/:id", deleteUser)

	e.Logger.Fatal(e.Start(":8080"))
}

//addUser func. Adds new user if it doesn't exist already.
func addUser(c echo.Context) error {
	u := new(User)
	if err := c.Bind(u); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "Unable to fill new User with request body data")
	}

	newUsers := []User{}

	storageFile, err := os.Open("storage.json")
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "Unable to open storage file")
	}

	defer storageFile.Close()

	jsonParser := json.NewDecoder(storageFile)
	jsonParser.Decode(&newUsers)

	for i := range newUsers {
		if newUsers[i].ID == u.ID {
			return c.JSON(http.StatusBadRequest, "Already have user with same ID")
		}
	}

	newUsers = append(newUsers, *u)

	file, err := json.MarshalIndent(newUsers, "", " ")
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "Unable to marshal data in JSON")
	}

	err = ioutil.WriteFile("storage.json", file, 0644)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "Unable to write JSON in file")
	}

	return c.JSON(http.StatusCreated, u)
}

//getUsers func. Returns list of Users from storage.json
func getUsers(c echo.Context) error {
	newUsers := []User{}

	storageFile, err := os.Open("storage.json")
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "Unable to open storage file")
	}

	defer storageFile.Close()

	jsonParser := json.NewDecoder(storageFile)
	jsonParser.Decode(&newUsers)

	return c.JSON(http.StatusOK, newUsers)
}

//getUser func. Returns User entity with selected ID
func getUser(c echo.Context) error {
	id := c.Param("id")

	newUsers := []User{}

	storageFile, err := os.Open("storage.json")
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "Unable to open storage file")
	}

	defer storageFile.Close()

	jsonParser := json.NewDecoder(storageFile)
	jsonParser.Decode(&newUsers)

	for i := range newUsers {
		if newUsers[i].ID == id {
			return c.JSON(http.StatusOK, newUsers[i])
		}
	}

	return c.JSON(http.StatusBadRequest, "Unable to find User with selected ID")
}

//updateUser func. Updates User name with name from request body if User with selected id exists
func updateUser(c echo.Context) error {
	id := c.Param("id")

	u := new(User)
	if err := c.Bind(u); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "Unable to fill new User with request body data")
	}

	newUsers := []User{}

	storageFile, err := os.Open("storage.json")
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "Unable to open storage file")
	}

	defer storageFile.Close()

	jsonParser := json.NewDecoder(storageFile)
	jsonParser.Decode(&newUsers)

	for i := range newUsers {
		if newUsers[i].ID == id {
			newUsers[i].Name = u.Name

			file, err := json.MarshalIndent(newUsers, "", " ")
			if err != nil {
				return c.JSON(http.StatusUnprocessableEntity, "Unable to marshal data in JSON")
			}

			err = ioutil.WriteFile("storage.json", file, 0644)
			if err != nil {
				return c.JSON(http.StatusUnprocessableEntity, "Unable to write JSON in file")
			}

			return c.JSON(http.StatusOK, newUsers[i])
		}
	}
	return c.JSON(http.StatusBadRequest, "Unable to find User with selected ID")
}

//deleteUser func. Deletes User entity with selected id from storage.json if User with selected id exists
func deleteUser(c echo.Context) error {
	id := c.Param("id")

	newUsers := []User{}

	storageFile, err := os.Open("storage.json")
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "Unable to open storage file")
	}

	defer storageFile.Close()

	jsonParser := json.NewDecoder(storageFile)
	jsonParser.Decode(&newUsers)

	for i := range newUsers {
		if newUsers[i].ID == id {
			newUsers[i] = newUsers[len(newUsers)-1]
			newUsers[len(newUsers)-1] = User{}
			newUsers = newUsers[:len(newUsers)-1]

			file, err := json.MarshalIndent(newUsers, "", " ")
			if err != nil {
				return c.JSON(http.StatusUnprocessableEntity, "Unable to marshal data in JSON")
			}

			err = ioutil.WriteFile("storage.json", file, 0644)
			if err != nil {
				return c.JSON(http.StatusUnprocessableEntity, "Unable to write JSON in file")
			}

			return c.JSON(http.StatusOK, "Successfully deleted User with selected ID")
		}
	}
	return c.JSON(http.StatusBadRequest, "Unable to find User with selected ID")
}
