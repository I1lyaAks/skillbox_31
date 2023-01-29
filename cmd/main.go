package main

import (
	"bufio"
	"encoding/json"
	"finish/proxy/proxy"
	"fmt"
	"strconv"

	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
)

func main() {
	OpenFileJson()
	defer CloseFileJson()
	go func() {
		ReplicaOne()
	}()
	go func() {
		ReplicaTwo()
	}()
	go func() {
		proxy.ProxyTwoReplicasRun()
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

}
func ReplicaOne() {
	go func() {
		log.Fatalln(http.ListenAndServe(":8080", handler()))
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

}
func ReplicaTwo() {
	go func() {
		log.Fatalln(http.ListenAndServe(":8081", handler()))
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
}

func handler() http.Handler {
	r := chi.NewRouter()
	r.Post("/test", TestProxy)
	r.Get("/users", Users)
	r.Post("/create", Create)
	r.Delete("/user", DeleteUser)
	r.Post("/make_friends", MakeFriends)
	r.Get("/friends/{id}", FriendsList)
	r.Put("/{id}", UpdateAge)
	return r
}
func OpenFileJson() {
	jsonFile, err := os.Open("user.json")
	if err != nil {
		log.Fatal("Unable to read input file ", err)
	}
	defer jsonFile.Close()
	fileScanner := bufio.NewScanner(jsonFile)
	Id += 1
	for fileScanner.Scan() {
		json.Unmarshal([]byte(fileScanner.Text()), &StorageMap)
	}
	if err := fileScanner.Err(); err != nil {
		log.Fatalf("Error while reading file: %s", err)
	}
}

func CloseFileJson() {
	jsonFile, err := os.Open("user.json")
	if err != nil {
		log.Fatal("Unable to read input file ", err)
	}
	defer jsonFile.Close()
	k, err := json.Marshal(StorageMap)
	if err != nil {
		log.Print(err)
	}
	err = ioutil.WriteFile("user.json", k, 0666)
	if err != nil {
		log.Print(err)
	}
}

var Id int = 0

var StorageMap = make(map[int]*User, 500)

type User struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Frend []int  `json:"frends"`
}
type Frends struct {
	SourceId int `json:"sourceid"`
	TargetId int `json:"targetid"`
}
type Delete struct {
	Id int `json:"id"`
}
type NewAge struct {
	NovelAge int `json:"new"`
}

func TestProxy(w http.ResponseWriter, r *http.Request) {
	content, _ := ioutil.ReadAll(r.Body)
	str := string(content) + "   " + r.Host
	w.Write([]byte(str))
}

func Create(w http.ResponseWriter, r *http.Request) {
	content, _ := ioutil.ReadAll(r.Body)
	var s User
	json.Unmarshal(content, &s)
	fmt.Println(s)
	Id++
	StorageMap[Id] = &s
	str := fmt.Sprintf("Name: %s, Age: %d, Frends:%v", StorageMap[Id].Name, s.Age, s.Frend)
	w.Write([]byte(str))
}
func Users(w http.ResponseWriter, r *http.Request) {
	for key, val := range StorageMap {
		s := val.Str(key)
		w.Write([]byte(s))
	}
}
func MakeFriends(w http.ResponseWriter, r *http.Request) {
	content, _ := ioutil.ReadAll(r.Body)
	var f Frends
	json.Unmarshal(content, &f)

	f.Friendship()
	str := fmt.Sprintf("user %d and user %d are friends now", f.SourceId, f.TargetId)
	w.Write([]byte(str))
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}
	data := Delete{}
	json.Unmarshal(content, &data)
	id := data.Id
	if _, ok := StorageMap[id]; !ok {
		w.Write([]byte("Пользователь отсутствует"))
	}
	deletefrends := StorageMap[id].Frend
	if len(deletefrends) > 0 {
		for _, i := range deletefrends {
			data.Delite(i)
		}
	}
	str := StorageMap[id].Name

	w.Write([]byte(str))
	delete(StorageMap, id)
}

func FriendsList(w http.ResponseWriter, r *http.Request) {
	idString := chi.URLParam(r, "id")
	idInt, err := strconv.Atoi(idString)
	if err != nil {
		fmt.Println(err)
	}
	for _, i := range StorageMap[idInt].Frend {
		str := StorageMap[i].Name + " "
		w.Write([]byte(str))
	}
}
func UpdateAge(w http.ResponseWriter, r *http.Request) {

	idString := chi.URLParam(r, "id")
	idInt, er := strconv.Atoi(idString)
	if er != nil {
		fmt.Println(er)
	}
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	var a NewAge
	json.Unmarshal(content, &a)
	a.Novel(idInt)
}

func (a NewAge) Novel(id int) {
	StorageMap[id].Age = a.NovelAge
}

func (d Delete) Delite(id int) {
	deleteuser := StorageMap[id].Frend
	if len(deleteuser) > 1 {
		for i := range deleteuser {
			if deleteuser[i] == id {
				deleteuser = append(deleteuser[:i], deleteuser[i+1:]...)
			}
		}
	} else {
		var x []int
		StorageMap[id].Frend = x
	}

}
func (u *User) Str(id int) string {
	s := fmt.Sprintf("id:%d, name:%s, age:%d, frends:%v \n", id, u.Name, u.Age, u.Frend)
	return s
}

func (f Frends) Friendship() {
	StorageMap[f.SourceId].Frend = append(StorageMap[f.SourceId].Frend, f.TargetId)
	StorageMap[f.TargetId].Frend = append(StorageMap[f.TargetId].Frend, f.SourceId)
}
