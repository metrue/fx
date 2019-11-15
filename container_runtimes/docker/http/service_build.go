package api

// import (
// 	"fmt"
// 	"io/ioutil"
// 	"os"
// 	"path/filepath"
//
// 	"github.com/google/uuid"
// 	"github.com/metrue/fx/types"
// 	"github.com/metrue/fx/utils"
// )
//
// func makeTar(project types.Project, tarFilePath string) error {
// 	dir, err := ioutil.TempDir("/tmp", "fx-build-dir")
// 	if err != nil {
// 		return err
// 	}
//
// 	defer os.RemoveAll(dir)
//
// 	for _, file := range project.Files {
// 		tmpfn := filepath.Join(dir, file.Path)
// 		if err := utils.EnsureFile(tmpfn); err != nil {
// 			return err
// 		}
// 		if err := ioutil.WriteFile(tmpfn, []byte(file.Body), 0666); err != nil {
// 			return err
// 		}
// 	}
//
// 	return utils.TarDir(dir, tarFilePath)
// }
//
// // Build build a project
// func (api *API) Build(project types.Project) (types.Service, error) {
// 	tarDir, err := ioutil.TempDir("/tmp", "fx-tar")
// 	if err != nil {
// 		return types.Service{}, err
// 	}
// 	defer os.RemoveAll(tarDir)
//
// 	imageID := uuid.New().String()
// 	tarFilePath := filepath.Join(tarDir, fmt.Sprintf("%s.tar", imageID))
// 	if err := makeTar(project, tarFilePath); err != nil {
// 		return types.Service{}, err
// 	}
// 	labels := map[string]string{
// 		"belong-to": "fx",
// 	}
// 	if err := api.BuildImage(tarFilePath, imageID, labels); err != nil {
// 		return types.Service{}, err
// 	}
//
// 	return types.Service{
// 		Name:  project.Name,
// 		Image: imageID,
// 	}, nil
// }
