package libs

import (
	utils "jnoronhautils"
	"strings"
)

func getExtensionVsixPath() string {
	return utils.ResolvePath(utils.GetCurrentDir() + "/download")
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
		if !utils.HasInternet() {
			utils.ErrorLog("Not detect internet.", false)
			utils.WaitForAnyKeyPressed("Please, connect to internet(PRESS ANY KEY TO CONTINUE)")
		}
		response := utils.Download(getUrl(extensionId), filePath)
		if !response.Data {
			utils.ErrorLog(response.Error.Error(), false)
		}
		return response.Data
	}
	return false
}
