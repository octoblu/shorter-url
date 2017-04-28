package main

import (
	"fmt"
	"log"
	"os"

	"github.com/coreos/go-semver/semver"
	"github.com/urfave/cli"
	De "github.com/visionmedia/go-debug"
)

var debug = De.Debug("shorter-url:main")

func main() {
	app := cli.NewApp()
	app.Name = "shorter-url"
	app.Version = version()
	app.Action = run
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "auth, a",
			EnvVar: "AUTH",
			Usage:  "<username>:<password> that is allowed create access",
			Value:  "user:pass",
		},
		cli.StringFlag{
			Name:   "mongodb-url, m",
			EnvVar: "MONGODB_URL",
			Usage:  "Mongo db url to use for data persistence",
			Value:  "mongodb://localhost:27017/shorter-url",
		},
		cli.IntFlag{
			Name:   "port, p",
			EnvVar: "PORT",
			Usage:  "Port to listen on for incoming HTTP requests",
			Value:  80,
		},
		cli.StringFlag{
			Name:   "redis-namespace, n",
			EnvVar: "REDIS_NAMESPACE",
			Usage:  "Redis namespace to use when caching",
			Value:  "redis://localhost:6379",
		},
		cli.StringFlag{
			Name:   "redis-url, r",
			EnvVar: "REDIS_URL",
			Usage:  "Redis db url to use for data caching",
			Value:  "redis://localhost:6379",
		},
	}
	app.Run(os.Args)
}

func run(context *cli.Context) {
	auth, mongoDBURL, port, redisNamespace, redisURL := getOpts(context)

	server := shorterurl.New(auth, mongoDBURL, port, redisNamespace, redisURL)
	fmt.Printf("Listening on 0.0.0.0:%v\n", port)
	err := server.Run()
	log.Fatalln(err.Error())
}

func getOpts(context *cli.Context) (string, string, int, string, string) {
	auth := context.String("auth")
	mongoDBURL := context.String("mongodb-url")
	port := context.Int("port")
	redisNamespace := context.String("redis-namespace")
	redisURL := context.String("redis-url")

	return auth, mongoDBURL, port, redisNamespace, redisURL
}

func version() string {
	version, err := semver.NewVersion(VERSION)
	if err != nil {
		errorMessage := fmt.Sprintf("Error with version number: %v", VERSION)
		log.Panicln(errorMessage, err.Error())
	}
	return version.String()
}
