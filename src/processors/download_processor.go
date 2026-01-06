package processors

import (
	"golangutils"
	"main/libs"
	"strings"
)

type DownloadProcessor struct {
	force bool
}

func (dp *DownloadProcessor) getExtensionVsixPath() string {
	dir := golangutils.ResolvePath(libs.SystemUtils.TempDir() + "/vscode-config-download")
	golangutils.CreateDirectory(dir, true)
	return dir
}

func (dp *DownloadProcessor) getExtensionVsixFile(extensionId string) string {
	path := dp.getExtensionVsixPath()
	if len(path) == 0 {
		return path
	}
	extensionDest := path + "/" + extensionId + ".vsix"
	return golangutils.ResolvePath(extensionDest)
}

func (dp *DownloadProcessor) getUrl(extensionId string) string {
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

func (dp *DownloadProcessor) download(extensionId string) bool {
	processor := &DownloadProcessor{}
	status := false
	filePath := processor.getExtensionVsixFile(extensionId)
	if dp.force {
		golangutils.DeleteFile(filePath)
	}
	if len(filePath) > 0 && !golangutils.IsFile(filePath) {
		if !golangutils.HasInternet() {
			libs.LoggerUtils.Error("Not detect internet.")
			libs.ConsoleUtils.WaitForAnyKeyPressed("Please, connect to internet(PRESS ANY KEY TO CONTINUE)")
		}
		responseErr := golangutils.Download(processor.getUrl(extensionId), filePath)
		if responseErr != nil {
			libs.LoggerUtils.Error(responseErr.Error())
		} else {
			status = true
		}
	}
	return status
}

func (dp *DownloadProcessor) downloadList(extensionIdList []string) {
	for _, id := range extensionIdList {
		dp.download(id)
	}
}
