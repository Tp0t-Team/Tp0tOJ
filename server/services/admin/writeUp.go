package admin

import (
	"archive/zip"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"server/services/database/resolvers"
	"server/services/types"
	"server/utils/configure"
	"strconv"
	"strings"
)

func DownloadAllWP(w http.ResponseWriter, req *http.Request) {
	//TODO: pack writeup folder and download
	var zipFileName = "WP.zip"
	err := os.RemoveAll(zipFileName)
	if err != nil {
		log.Panicln("DownloadAllWP: ", err)
		return
	}
	zipFile, _ := os.Create(zipFileName)
	//defer func(zipFile *os.File) {
	//	err := zipFile.Close()
	//	if err != nil {
	//		log.Panicln("DownloadAllWP: ", err)
	//		return
	//	}
	//}(zipFile)

	archive := zip.NewWriter(zipFile)
	//defer func(archive *zip.Writer) {
	//	err := archive.Close()
	//	if err != nil {
	//		log.Panicln("DownloadAllWP: ", err)
	//		return
	//	}
	//}(archive)
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
		header.Name = file.Name()
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
	err = zipFile.Close()
	if err != nil {
		log.Panicln("DownloadAllWP: ", err)
	}
	zipFile, _ = os.Open(zipFileName)
	defer func(zipFile *os.File) {
		err := zipFile.Close()
		if err != nil {
			log.Panicln("DownloadAllWP: ", err)
			return
		}
	}(zipFile)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	_, err = io.Copy(w, zipFile)
	if err != nil {
		log.Panicln("DownloadAllWP: ", err)
	}
}

func DownloadWPByUserId(w http.ResponseWriter, req *http.Request, userId string) {
	var zipFileName = "WP_" + userId + ".zip"
	err := os.RemoveAll(zipFileName)
	if err != nil {
		log.Panicln("DownloadWP: ", err)
		return
	}
	zipFile, _ := os.Create(zipFileName)
	//defer func(zipFile *os.File) {
	//	err := zipFile.Close()
	//	if err != nil {
	//		log.Panicln("DownloadAllWP: ", err)
	//		return
	//	}
	//}(zipFile)

	archive := zip.NewWriter(zipFile)
	//defer func(archive *zip.Writer) {
	//	err := archive.Close()
	//	if err != nil {
	//		log.Panicln("DownloadAllWP: ", err)
	//		return
	//	}
	//}(archive)
	filepath.Walk(configure.WriteUpPath, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		_, fileName := filepath.Split(path)
		id := strings.Split(fileName, ".")[0]
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
		header.Name = file.Name()
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
	err = zipFile.Close()
	if err != nil {
		log.Panicln("DownloadWP: ", err)
	}
	zipFile, _ = os.Open(zipFileName)
	defer func(zipFile *os.File) {
		err := zipFile.Close()
		if err != nil {
			log.Panicln("DownloadWP: ", err)
			return
		}
		err = os.RemoveAll(zipFileName)
		if err != nil {
			log.Println("remove ", zipFileName, " error: ", err)
			return
		}
	}(zipFile)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	_, err = io.Copy(w, zipFile)
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
			UserId: string(userId),
			Name:   user.Name,
			Mail:   user.Mail,
			Solved: int32(solved),
		})
		return nil
	})
	return writeUpInfos
}
