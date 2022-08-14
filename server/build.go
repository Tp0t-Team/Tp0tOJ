package main

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func main() {
	postgres := flag.Bool("postgres", false, "")
	sqlite := flag.Bool("sqlite", false, "")
	flag.Parse()

	tags := []string{"WithFrontEnd"}

	if *postgres == *sqlite {
		log.Panicln("you must choose one and only one database type.")
	}

	if *postgres {
		tags = append(tags, "DatabasePostgres")
	}
	if *sqlite {
		tags = append(tags, "DatabaseSqlite")
	}

	err := os.Chdir("../app")
	if err != nil {
		log.Panicln(err)
	}
	frontendCmd := exec.Command("npm", "run", "build")
	log.Println("build frontend...")
	err = frontendCmd.Run()
	if err != nil {
		log.Panicln(err)
	}

	err = os.Chdir("../server/services")
	if err != nil {
		log.Panicln(err)
	}

	err = os.MkdirAll("static", 0755)
	if err != nil {
		log.Panicln(err)
	}

	err = os.Chdir("../../app/dist")
	if err != nil {
		log.Panicln(err)
	}
	log.Println("copy static...")
	err = filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			if path == "." {
				return nil
			}
			err := os.MkdirAll("../../server/services/static/"+path, 0755)
			if err != nil {
				return err
			}
			return nil
		}
		_, filename := filepath.Split(path)
		parts := strings.Split(filename, ".")
		if parts[len(parts)-1] == "map" {
			return nil
		}
		dst, err := os.Create("../../server/services/static/" + path)
		if err != nil {
			return err
		}
		src, err := os.Open(path)
		if err != nil {
			return err
		}
		_, err = io.Copy(dst, src)
		if err != nil {
			return err
		}
		err = src.Close()
		if err != nil {
			return err
		}
		err = dst.Close()
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Panicln(err)
	}

	err = os.Chdir("../../server")
	if err != nil {
		log.Panicln(err)
	}
	goModCmd := exec.Command("go", "mod", "download")
	log.Println("prepare go mod...")
	err = goModCmd.Run()
	if err != nil {
		log.Panicln(err)
	}
	goModTidyCmd := exec.Command("go", "mod", "tidy")
	err = goModTidyCmd.Run()
	if err != nil {
		log.Panicln(err)
	}
	os.Setenv("CGO_ENABLED", "0")
	serverCmd := exec.Command("go", "build", "-tags", strings.Join(tags, ","), "-o", fmt.Sprintf("OJ_%s_%s", runtime.GOOS, runtime.GOARCH), "main.go")
	log.Println("build server...")
	err = serverCmd.Run()
	if err != nil {
		log.Panicln(err)
	}

	err = os.RemoveAll("services/static")
	if err != nil {
		log.Panicln(err)
	}
	log.Println("success")
}
