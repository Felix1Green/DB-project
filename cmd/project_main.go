package main

import (
	"errors"
	"flag"
	"github.com/Felix1Green/DB-project/internal/app/cleaner"
	"github.com/Felix1Green/DB-project/internal/app/forum"
	"github.com/Felix1Green/DB-project/internal/app/post"
	"github.com/Felix1Green/DB-project/internal/app/thread"
	"github.com/Felix1Green/DB-project/internal/app/users"
	"github.com/Felix1Green/DB-project/internal/pkg/utils"
	"github.com/jackc/pgx"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"strconv"
	"time"
)

type ServerConfig struct{
	StatusService *cleaner.Service
	ForumService *forum.Service
	PostService *post.Service
	ThreadService *thread.Service
	UserService *users.Service
}


func CreateDBConnection(config *utils.ServiceConfig)(*pgx.ConnPool,error){
	pool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			User:     config.DatabaseUser,
			Password: config.DatabasePassword,
			Port:     uint16(config.DatabasePort),
			Database: config.DatabaseName,
			Host: config.DatabaseDomain,
		},
		MaxConnections: 50,
	})
	if err != nil {
		return nil, errors.New("no postgresql connection")
	}

	return pool, nil
}

func ParseConfigPath() string {
	configPath := ""
	flag.StringVar(&configPath, "c", "configuration.json", "set configuration")
	flag.Parse()
	return configPath
}

func InitService(sqlConn *pgx.ConnPool) *ServerConfig{
	statusService := cleaner.Start(sqlConn)
	usersService := users.Start(sqlConn)
	forumService := forum.Start(sqlConn, usersService.Repository)
	threadService := thread.Start(sqlConn)
	postService := post.Start(sqlConn, usersService.Repository, forumService.Repository, threadService.Repository)

	return &ServerConfig{
		StatusService: statusService,
		ForumService:  forumService,
		PostService:   postService,
		ThreadService: threadService,
		UserService:   usersService,
	}
}

func configureMainRouter(application *ServerConfig) http.Handler{
	handler := http.NewServeMux()

	handler.Handle("/api/forum/", application.ForumService.Router)
	handler.Handle("/api/post/", application.PostService.Router)
	handler.Handle("/api/service/", application.StatusService.Router)
	handler.Handle("/api/thread/", application.ThreadService.Router)
	handler.Handle("/api/user/", application.UserService.Router)

	return handler
}


func main(){
	configPath := ParseConfigPath()
	config, configErr := utils.Run(configPath)
	if configErr != nil{
		log.Fatalln(configErr)
	}

	conn, err := CreateDBConnection(config)
	if err != nil{
		log.Fatalln(err)
	}

	conf := InitService(conn)
	URLHandler := configureMainRouter(conf)
	httpServer := &http.Server{
		Addr: config.Domain + ":" + strconv.Itoa(config.Port),
		Handler: URLHandler,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("STARTING SERVER AT PORT: ", config.Port)
	serverErr := httpServer.ListenAndServe()
	if serverErr != nil{
		log.Fatalln(serverErr)
	}

	defer func() {
		if conn != nil{
			conn.Close()
		}
	}()
}
