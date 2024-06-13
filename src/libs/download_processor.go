package libs

import (
	"jnoronhautils"
	"strings"
)

func getExtensionVsixPath() string {
	return jnoronhautils.ResolvePath(jnoronhautils.GetCurrentDir() + "/vscode-config-download")
}

func getExtensionVsixFile(extensionId string) string {
	path := getExtensionVsixPath()
	if len(path) == 0 {
		return path
	}
	extensionDest := path + "/" + extensionId + ".vsix"
	return jnoronhautils.ResolvePath(extensionDest)
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
	jnoronhautils.InfoLog("Download extension: "+extensionId, false)
	filePath := getExtensionVsixFile(extensionId)
	if len(filePath) > 0 {
		if !jnoronhautils.HasInternet() {
			jnoronhautils.ErrorLog("Not detect internet.", false)
			jnoronhautils.WaitForAnyKeyPressed("Please, connect to internet(PRESS ANY KEY TO CONTINUE)")
		}
		response := jnoronhautils.Download(getUrl(extensionId), filePath)
		if !response.Data {
			jnoronhautils.ErrorLog(response.Error.Error(), false)
		}
		return response.Data
	}
	return false
}
