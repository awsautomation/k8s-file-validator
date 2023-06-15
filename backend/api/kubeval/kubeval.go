package kubeval

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
	"github.com/komodorio/validkube/backend/api/utils"
	"github.com/komodorio/validkube/backend/internal/routing"
	"sigs.k8s.io/yaml"
)

const Path = "/kubeval"
const Method = routing.POST

func kubevalWrapper(inputYaml []byte) ([]byte, error) {
	if err := utils.CreateDirectory("/tmp/yaml"); err != nil {
		fmt.Printf("got error 1: %s \n", err.Error())
		return nil, err
	}

	if err := utils.WriteFile("/tmp/yaml/target_yaml.yaml", inputYaml); err != nil {
		fmt.Printf("got error 2: %s \n", err.Error())
		return nil, err
	}
	outputFromKubevalAsJson, err := exec.Command("kubeval", "-o", "json", "/tmp/yaml/target_yaml.yaml").Output()
	if err != nil {
		fmt.Printf("got error 3: %s \n", err.Error())
		return nil, err
	}

	outputFromKubevalAsYaml, err := yaml.JSONToYAML(outputFromKubevalAsJson)
	if err != nil {
		fmt.Printf("got error 4: %s \n", err.Error())
		return nil, err
	}
	return outputFromKubevalAsYaml, nil
}

func ProcessRequest(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Printf("got error 5: %s \n", err.Error())
		fmt.Printf("Erorr has with reading request body: %v", err)
		c.JSON(http.StatusOK, gin.H{"data": "", "err": err.Error()})
		return
	}
	bodyAsMap, err := utils.JsonToMap(body)
	if err != nil {
		fmt.Printf("got error 6: %s \n", err.Error())
		c.JSON(http.StatusOK, gin.H{"data": "", "err": err.Error()})
		return
	}
	yamlAsInterface := bodyAsMap["yaml"]
	kubevalOutput, err := kubevalWrapper(utils.InterfaceToBytes(yamlAsInterface))
	if err != nil {
		fmt.Printf("got error 7: %s \n", err.Error())
		c.JSON(http.StatusOK, gin.H{"data": "", "err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": string(kubevalOutput), "err": nil})
}
