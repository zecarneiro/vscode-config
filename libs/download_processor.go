package libs

import (
	utils "jnoronha_golangutils"
	"os"
	"strings"
)

func getExtensionVsixPath() string {
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}
	return utils.ResolvePath(wd + "/download")
}

func getExtensionVsixFile(extensionId string) string {
	path := getExtensionVsixPath()
	if len(path) == 0 {
		return path
	}
	extensionDest := path + "/" + extensionId + ".vsix"
	return utils.ResolvePath(extensionDest)
}

func getUrl(extensionId string) string {
	delimiter := "."
	url := "https://{publisher}.gallery.vsassets.io/_apis/public/gallery/publisher/{publisher}/extension/{package}/latest/assetbyname/Microsoft.VisualStudio.Services.VSIXPackage"

	// Split the string into substrings using the delimiter
	publisherPackage := strings.Split(extensionId, delimiter)
	publisherExtension := publisherPackage[0]
	packageExtension := publisherPackage[1]
	newUrl := strings.Replace(url, "{publisher}", publisherExtension, -1)
	newUrl = strings.Replace(newUrl, "{package}", packageExtension, -1)
	return newUrl
}

func download(extensionId string) bool {
	utils.InfoLog("Download extension: "+extensionId, false)
	filePath := getExtensionVsixFile(extensionId)
	if len(filePath) > 0 {
		response := utils.Download(getUrl(extensionId), filePath)
		if !response.Data {
			utils.ErrorLog(response.Error.Error(), false)
		}
		return response.Data
	}
	return false

}
