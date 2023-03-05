package sstr
import(
	"io/ioutil"
	"strings"
	"log"
	"path/filepath"
	"strconv"
	"io/fs"
	"path"
)
func HTMLDirList(pathto string, a string) string {
	filesListRaw, err := ioutil.ReadDir("./" + pathto + a)
	if err != nil {
		log.Printf("%v %v\n" + err.Error())
	}
	if len(filesListRaw) == 0 {
		return "<!DOCTYPE html><body style=\"background-color:black\"><p style=\"color: white;font-size:1cm;\"><b>No files are found in " + pathto + " </b></p></body>"
	}
	filesList := "<!DOCTYPE html><body style=\"background-color:black\"><p style=\"color: white;font-size:1cm;\"><b>Index of " + a + "</b></p>"
	for index, file := range filesListRaw {
		link := a
		if !strings.HasSuffix(link, "/") {
			link += "/"
		}
		link += file.Name()

		filesList += "<a href=\"" + link + "\"><u style=\"text-decoration-color: black;line-height: 0.01;\"><p style=\"font-size: 0.7cm;color:white;line-height: 0.01;\">" + strconv.Itoa(index+1) + ". " + func(currfile fs.FileInfo) string {
			currfile.IsDir()
			if currfile.IsDir() {
				return "\U0001F4C1"
			} else {
				for _, ex := range Imgex {
					if filepath.Ext(currfile.Name()) == ex {
						return "\U0001f5bc"
					}
				}
				for _, ex := range Mediasoundext {
					if filepath.Ext(currfile.Name()) == ex {
						return "\U0001f3b5"
					}
				}
				for _, ex := range Vidext {
					if filepath.Ext(currfile.Name()) == ex {
						return "▶️"
					}
				}
				return "\U0001F4C4"
			}
		}(file) + "<b>" + file.Name() + " </b>" + "</p></u></a><br>"
	}
	filesList += "</body>"
	return filesList
}
func BadRequest400() string {
	return "<!DOCTYPE html><head><title>400 Bad Request</title></head><body style=\"background-color: black;\"></body><h1 style=\"text-align: center;color:white;font-size: 90px\">400</h1><br><p style=\"text-align: center;color:white;font-size: 30px\"> Bad Request " + "</p></body>"
}
func NotFound404(pathto string) string {
	return "<!DOCTYPE html><head><title>404 Not Found</title></head><body style=\"background-color: black;\"></body><h1 style=\"text-align: center;color:white;font-size: 90px\">404</h1><br><p style=\"text-align: center;color:white;font-size: 30px\"> The specified content is not found: " + path.Base(pathto) + "</p></body>"
}
func ServerErr500() string {
	return "<!DOCTYPE html><head><title>500 Server Error!></title></head><body style=\"background-color: black;\"></body><h1 style=\"text-align: center;color:white;font-size: 90px\">404</h1><br><p style=\"text-align: center;color:white;font-size: 30px\">Internal Server Error" + "</p></body>"
}


