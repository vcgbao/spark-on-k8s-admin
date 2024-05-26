package sparkoperator

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"spark-on-k8s-admin/config"
	"spark-on-k8s-admin/utils"
	"strconv"
	"strings"
)

type SparkOperatorService struct {
	Cfg *config.Config
}

func Init(cfg *config.Config) *SparkOperatorService {
	return &SparkOperatorService{
		Cfg: cfg,
	}
}

func (s *SparkOperatorService) GetClusterInfo() (string, int, string) {
	if s.Cfg.K8sConfig.InCluster {
		token, _ := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
		port, _ := strconv.Atoi(os.Getenv("KUBERNETES_SERVICE_PORT"))

		return os.Getenv("KUBERNETES_SERVICE_HOST"), port, string(token)
	} else {
		return s.Cfg.K8sConfig.ServiceHost, s.Cfg.K8sConfig.ServicePort, s.Cfg.K8sConfig.Token
	}
}

func (s *SparkOperatorService) GetResourceByName(name string, resourceType string) (*map[string]any, error) {
	host, port, token := s.GetClusterInfo()
	url := fmt.Sprintf("https://%s:%d/apis/sparkoperator.k8s.io/v1beta2/namespaces/%s/%s/%s", host, port, s.Cfg.SparkConfig.Namespace, resourceType, name)
	authorization := fmt.Sprintf("Bearer %s", token)
	cli := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	req, reqErr := http.NewRequest("GET", url, nil)

	if reqErr != nil {
		utils.Sugar.Errorf("Error when init request ", reqErr)
		return nil, reqErr
	}

	req.Header.Set("Authorization", authorization)

	resp, err := cli.Do(req)

	if err != nil {
		utils.Sugar.Errorf("Error when call K8s api", err)
		return nil, err
	}

	if resp.StatusCode == 200 {
		defer resp.Body.Close()
		sparkApplicationResponse := map[string]any{}
		err := json.NewDecoder(resp.Body).Decode(&sparkApplicationResponse)
		if err != nil {
			utils.Sugar.Errorf("Can not unmarshal result from K8s", err)
			return nil, err
		}
		return &sparkApplicationResponse, nil
	} else {
		utils.Sugar.Errorf("Error when call K8s api, StatusCode = %d", resp.StatusCode)
		return nil, fmt.Errorf("error when call K8s api, StatusCode = %d", resp.StatusCode)
	}
}

func (s *SparkOperatorService) GetResource(resourceType string) (*map[string]any, error) {
	host, port, token := s.GetClusterInfo()
	utils.Sugar.Infof("Calling to host: %s, port: %s, resourceType: %s", host, port, resourceType)
	url := fmt.Sprintf("https://%s:%d/apis/sparkoperator.k8s.io/v1beta2/namespaces/%s/%s", host, port, s.Cfg.SparkConfig.Namespace, resourceType)
	authorization := fmt.Sprintf("Bearer %s", token)
	cli := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	req, reqErr := http.NewRequest("GET", url, nil)

	if reqErr != nil {
		utils.Sugar.Errorf("Error when init request ", reqErr)
		return nil, reqErr
	}

	req.Header.Set("Authorization", authorization)

	resp, err := cli.Do(req)

	if err != nil {
		utils.Sugar.Errorf("Error when call K8s api", err)
		return nil, err
	}

	if resp.StatusCode == 200 {
		defer resp.Body.Close()
		sparkApplicationResponse := map[string]any{}
		err := json.NewDecoder(resp.Body).Decode(&sparkApplicationResponse)
		if err != nil {
			utils.Sugar.Errorf("Can not unmarshal result from K8s", err)
			return nil, err
		}
		return &sparkApplicationResponse, nil
	} else {
		utils.Sugar.Errorf("Error when call K8s api, StatusCode = %d", resp.StatusCode)
		return nil, fmt.Errorf("error when call K8s api, StatusCode = %d", resp.StatusCode)
	}
}

func (s *SparkOperatorService) DeleteResource(appName string, resourceType string) error {
	if len(strings.TrimSpace(appName)) == 0 {
		return fmt.Errorf("application name must not empty")
	}
	host, port, token := s.GetClusterInfo()
	url := fmt.Sprintf("https://%s:%d/apis/sparkoperator.k8s.io/v1beta2/namespaces/%s/%s/%s", host, port, s.Cfg.SparkConfig.Namespace, resourceType, appName)
	authorization := fmt.Sprintf("Bearer %s", token)
	cli := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	req, reqErr := http.NewRequest("DELETE", url, nil)

	if reqErr != nil {
		utils.Sugar.Errorf("Error when init request ", reqErr)
		return reqErr
	}

	req.Header.Set("Authorization", authorization)

	resp, err := cli.Do(req)

	if err != nil {
		utils.Sugar.Errorf("Error when call K8s api", err)
		return err
	}

	if resp.StatusCode == 200 {
		return nil
	} else {
		utils.Sugar.Errorf("Error when call K8s api, StatusCode = %d", resp.StatusCode)
		return fmt.Errorf("error when call K8s api, StatusCode = %d", resp.StatusCode)
	}
}

func (s *SparkOperatorService) CreateResource(content map[string]any, resourceType string) error {
	utils.Sugar.Infof("content = %v", content)
	err := AddAnnotations(content, "")
	if err != nil {
		return err
	}
	utils.Sugar.Infof("content = %v", content)

	body, err := json.Marshal(content)

	if err != nil {
		utils.Sugar.Errorf("Error when serialize spark application content", err)
		return err
	}

	host, port, token := s.GetClusterInfo()

	url := fmt.Sprintf(
		"https://%s:%d/apis/sparkoperator.k8s.io/v1beta2/namespaces/%s/%s",
		host,
		port,
		s.Cfg.SparkConfig.Namespace,
		resourceType,
	)
	authorization := fmt.Sprintf("Bearer %s", token)
	cli := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	req, reqErr := http.NewRequest("POST", url, bytes.NewBuffer(body))

	if reqErr != nil {
		utils.Sugar.Errorf("Error when init request ", reqErr)
		return reqErr
	}

	req.Header.Set("Authorization", authorization)
	req.Header.Set("Content-Type", "application/json")

	resp, err := cli.Do(req)

	if err != nil {
		utils.Sugar.Errorf("Error when call K8s api", err)
		return err
	}

	if resp.StatusCode == 200 || resp.StatusCode == 201 {
		return nil
	} else {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		utils.Sugar.Errorf("Error when call K8s api, StatusCode = %d, Body = %s", resp.StatusCode, string(body))
		return fmt.Errorf("error when call K8s api, StatusCode = %d, Body = %s", resp.StatusCode, string(body))
	}
}

func (s *SparkOperatorService) UpdateResource(content map[string]any, resourceType string) error {
	if _, ok := content["metadata"].(map[string]any); !ok {
		return fmt.Errorf("sparkapplication must contains metadata key with map type")
	}
	name, ok := content["metadata"].(map[string]any)["name"]

	if !ok {
		return fmt.Errorf("metadata.name must defined")
	}

	app, err := s.GetResourceByName(name.(string), resourceType)

	if err != nil {
		return err
	}

	resourceVersion, err := GetResourceVersion(*app)
	if err != nil {
		return err
	}

	utils.Sugar.Infof("content = %v", content)
	err = AddAnnotations(content, resourceVersion)
	if err != nil {
		return err
	}
	utils.Sugar.Infof("content = %v", content)

	body, err := json.Marshal(content)

	if err != nil {
		utils.Sugar.Errorf("Error when serialize spark application content", err)
		return err
	}

	host, port, token := s.GetClusterInfo()
	url := fmt.Sprintf("https://%s:%d/apis/sparkoperator.k8s.io/v1beta2/namespaces/%s/%s/%s", host, port, s.Cfg.SparkConfig.Namespace, resourceType, name)
	authorization := fmt.Sprintf("Bearer %s", token)
	cli := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	req, reqErr := http.NewRequest("PUT", url, bytes.NewBuffer(body))

	if reqErr != nil {
		utils.Sugar.Errorf("Error when init request ", reqErr)
		return reqErr
	}

	req.Header.Set("Authorization", authorization)
	req.Header.Set("Content-Type", "application/json")

	resp, err := cli.Do(req)

	if err != nil {
		utils.Sugar.Errorf("Error when call K8s api", err)
		return err
	}

	if resp.StatusCode == 200 || resp.StatusCode == 201 {
		return nil
	} else {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		utils.Sugar.Errorf("Error when call K8s api, StatusCode = %d, Body = %s", resp.StatusCode, string(body))
		return fmt.Errorf("error when call K8s api, StatusCode = %d, Body = %s", resp.StatusCode, string(body))
	}
}

func GetResourceVersion(content map[string]any) (string, error) {
	if _, ok := content["metadata"].(map[string]any); !ok {
		return "", fmt.Errorf("sparkapplication must contains metadata key with map type")
	}

	if resourceVersion, ok := content["metadata"].(map[string]any)["resourceVersion"]; ok {
		return resourceVersion.(string), nil
	}
	return "", fmt.Errorf("can not get resourceVersion")
}

func AddAnnotations(content map[string]any, resourceVersion string) error {
	if _, ok := content["metadata"].(map[string]any); !ok {
		return fmt.Errorf("sparkapplication must contains metadata key with map type")
	}

	if _, ok := content["metadata"].(map[string]any)["annotations"]; ok {
		if _, ok := content["metadata"].(map[string]any)["annotations"].(map[string]any); !ok {
			return fmt.Errorf("anotation must be a map")
		}
	} else {
		content["metadata"].(map[string]any)["annotations"] = map[string]any{}
	}
	metadata := content["metadata"].(map[string]any)
	annotation := metadata["annotations"].(map[string]any)

	delete(annotation, "kubectl.kubernetes.io/last-applied-configuration")

	lastAppliedConfig, err := json.Marshal(content)
	if err != nil {
		return fmt.Errorf("can not marshal content")
	}
	annotation["kubectl.kubernetes.io/last-applied-configuration"] = string(lastAppliedConfig)

	if resourceVersion != "" {
		metadata["resourceVersion"] = resourceVersion
	}

	return nil
}
