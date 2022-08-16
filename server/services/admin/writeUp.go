package admin

import (
	"archive/zip"
	"bytes"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"server/services/database/resolvers"
	"server/services/types"
	"server/utils/configure"
	"strconv"
	"strings"
)

var usernameRE = regexp.MustCompile("([^\\p{L}\\p{M}\\p{N}\\p{P}\\p{S}]|[/\\\\'\"`?*:@<>|])")

func usernameFilter(name string) string {
	return usernameRE.ReplaceAllString(name, "")
}

func DownloadAllWP(w http.ResponseWriter, req *http.Request) {
	//pack writeup folder and download
	var zipFileName = "WP.zip"
	err := os.RemoveAll(zipFileName)
	if err != nil {
		log.Panicln("DownloadAllWP: ", err)
		return
	}
	zipFile := bytes.Buffer{}

	archive := zip.NewWriter(&zipFile)
	err = filepath.Walk(configure.WriteUpPath, func(path string, info os.FileInfo, _ error) error {
		if info.IsDir() {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Method = zip.Deflate
		parts := strings.Split(info.Name(), ".")
		id, err := strconv.ParseUint(parts[0], 10, 64)
		if err != nil {
			log.Println("non-writeup file in writeup folder: ", info.Name())
			return nil
		}
		user, err := resolvers.FindUser(id)
		if err != nil {
			return err
		}
		header.Name = "writeup/" + parts[0] + "_" + usernameFilter(user.Name) + "." + strings.Join(parts[1:], ".")
		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				log.Panicln("DownloadAllWP: ", err)
				return
			}
		}(file)

		_, err = io.Copy(writer, file)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Panicln("DownloadAllWP: ", err)
		return
	}
	err = archive.Close()
	if err != nil {
		log.Panicln("DownloadAllWP: ", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	_, err = io.Copy(w, &zipFile)
	if err != nil {
		log.Panicln("DownloadAllWP: ", err)
	}
}

func DownloadWPByUserId(w http.ResponseWriter, req *http.Request, userId string) {
	//TODO: maybe the archive should not include folder
	parsedUserId, err := strconv.ParseUint(userId, 10, 64)
	if err != nil {
		log.Println(err)
		return
	}
	user, err := resolvers.FindUser(parsedUserId)
	if err != nil {
		log.Println(err)
		return
	}
	var username = usernameFilter(user.Name)
	var zipFileName = "WP_" + userId + "_" + username + ".zip"
	err = os.RemoveAll(zipFileName)
	if err != nil {
		log.Panicln("DownloadWP: ", err)
		return
	}
	zipFile := bytes.Buffer{}

	archive := zip.NewWriter(&zipFile)

	filepath.Walk(configure.WriteUpPath, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		_, fileName := filepath.Split(path)
		parts := strings.Split(fileName, ".")
		id := parts[0]
		if id != userId {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Method = zip.Deflate
		header.Name = id + "_" + username + "." + strings.Join(parts[1:], ".")
		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				log.Panicln("DownloadWP: ", err)
				return
			}
		}(file)
		_, err = io.Copy(writer, file)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Panicln("DownloadWP: ", err)
		return
	}
	err = archive.Close()
	if err != nil {
		log.Panicln("DownloadWP: ", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	_, err = io.Copy(w, &zipFile)
	if err != nil {
		log.Panicln("DownloadWP: ", err)
	}
}

func GetWriteUpInfos() []types.WriteUpInfo {
	var writeUpInfos []types.WriteUpInfo
	filepath.Walk(configure.WriteUpPath, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		_, file := filepath.Split(path)
		userId, err := strconv.ParseUint(strings.Split(file, ".")[0], 10, 64)
		if err != nil {
			log.Panicln("writeup parse file error", err)
			return nil
		}
		user, err := resolvers.FindUser(userId)
		if err != nil {
			log.Panicln("writeup can't find user", err)
			return nil
		}
		var solved int
		submits := resolvers.FindSubmitCorrectByUserId(userId)
		if submits == nil {
			solved = 0
		}
		solved = len(submits)
		writeUpInfos = append(writeUpInfos, types.WriteUpInfo{
			UserId: strconv.FormatUint(userId, 10),
			Name:   user.Name,
			Mail:   user.Mail,
			Solved: int32(solved),
		})
		return nil
	})
	return writeUpInfos
}
