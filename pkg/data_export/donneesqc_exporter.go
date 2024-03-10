package internal_data

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ClubCedille/hackqc2024/pkg/event"
)

func PostJsonEventsToDQ(apiKey string, datasetIdentifier string, filePath string) error {
    // Check if the resource already exists
    checkUrl := "https://pab.donneesquebec.ca/recherche/api/3/action/package_show?id=" + datasetIdentifier
    checkReq, err := http.NewRequest("GET", checkUrl, nil)
    if err != nil {
        return err
    }

    checkReq.Header.Set("Authorization", apiKey)
    client := &http.Client{}
    checkResp, err := client.Do(checkReq)
    if err != nil {
        return err
    }
    defer checkResp.Body.Close()

    var result map[string]interface{}
    if err := json.NewDecoder(checkResp.Body).Decode(&result); err != nil {
        return err
    }

    resources := result["result"].(map[string]interface{})["resources"].([]interface{})
    var resourceId string
    for _, resource := range resources {
        res := resource.(map[string]interface{})
        if res["name"] == filepath.Base(filePath) {
            resourceId = res["id"].(string)
            break
        }
    }

    if resourceId != "" {
        // Resource exists, update it
        return patchJsonResource(apiKey, resourceId, filePath)
    } else {
        // Resource does not exist, create it
        return postJsonResource(apiKey, datasetIdentifier, filePath)
    }
}

func patchJsonResource(apiKey string, resourceId string, filePath string) error {
    url := "https://pab.donneesquebec.ca/recherche/api/3/action/resource_patch"
    reqBody := &bytes.Buffer{}
    writer := multipart.NewWriter(reqBody)

    writer.WriteField("id", resourceId)
    writer.WriteField("description", "JSON containing event list")
    writer.WriteField("taille_entier", fileSizeInMB(filePath))

    file, err := os.Open(filePath)
    if err != nil {
        return err
    }
    defer file.Close()

    part, err := writer.CreateFormFile("upload", filepath.Base(filePath))
    if err != nil {
        return err
    }
    _, err = io.Copy(part, file)
    if err != nil {
        return err
    }

    err = writer.Close()
    if err != nil {
        return err
    }

    req, err := http.NewRequest("POST", url, reqBody)
    if err != nil {
        return err
    }

    req.Header.Set("Authorization", apiKey)
    req.Header.Set("Content-Type", writer.FormDataContentType())

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
    }

    return nil
}

func postJsonResource(apiKey string, datasetIdentifier string, filePath string) error {
    url := "https://pab.donneesquebec.ca/recherche/api/3/action/resource_create"
    reqBody := &bytes.Buffer{}
    writer := multipart.NewWriter(reqBody)

    writer.WriteField("package_id", datasetIdentifier)
    writer.WriteField("name", filepath.Base(filePath))
    writer.WriteField("url", "upload")
    writer.WriteField("description", "JSON containing event list")
    writer.WriteField("taille_entier", fileSizeInMB(filePath))
    writer.WriteField("format", "JSON")
    writer.WriteField("relidi_config_separateur_virgule", "n/a")
    writer.WriteField("resource_type", "donnees")
    writer.WriteField("relidi_condon_valinc", "oui")

    file, err := os.Open(filePath)
    if err != nil {
        return err
    }
    defer file.Close()

    part, err := writer.CreateFormFile("upload", filepath.Base(filePath))
    if err != nil {
        return err
    }
    _, err = io.Copy(part, file)
    if err != nil {
        return err
    }
    
    err = writer.Close()
    if err != nil {
        return err
    }
      

    req, err := http.NewRequest("POST", url, reqBody)
    if err != nil {
        return err
    }

    req.Header.Set("Authorization", apiKey)
    req.Header.Set("Content-Type", writer.FormDataContentType())

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
    }

	return nil
}

func PostGeoJsonEventsToDQ(apiKey string, datasetIdentifier string, events []*event.Event) error {
    geoJsonData, err := ConvertEventsToGeoJSON(events)
    if err != nil {
        return err
    }

    checkUrl := "https://pab.donneesquebec.ca/recherche/api/3/action/package_show?id=" + datasetIdentifier
    checkReq, err := http.NewRequest("GET", checkUrl, nil)
    if err != nil {
        return err
    }

    checkReq.Header.Set("Authorization", apiKey)
    client := &http.Client{}
    checkResp, err := client.Do(checkReq)
    if err != nil {
        return err
    }
    defer checkResp.Body.Close()

    var result map[string]interface{}
    if err := json.NewDecoder(checkResp.Body).Decode(&result); err != nil {
        return err
    }

    // resources := result["result"].(map[string]interface{})["resources"].([]interface{})
    // var resourceId string
    // for _, resource := range resources {
    //     res := resource.(map[string]interface{})
    //     if res["name"] == "events.geojson" {
    //         resourceId = res["id"].(string)
    //         break
    //     }
    // }

    // If the resourceId is not empty, patch the existing resource
    // if resourceId != "" {
    //     return patchGeoJsonResource(apiKey, resourceId, geoJsonData)
    // } else {
    //     return postGeoJsonResource(apiKey, datasetIdentifier, geoJsonData)
    // }
    return postGeoJsonResource(apiKey, datasetIdentifier, geoJsonData)

}

// func patchGeoJsonResource(apiKey string, resourceId string, geoJsonData []byte) error {
//     // todo
//     return nil
// }

func postGeoJsonResource(apiKey string, datasetIdentifier string, geoJsonData []byte) error {
    url := "https://pab.donneesquebec.ca/recherche/api/3/action/resource_create"
    reqBody := &bytes.Buffer{}
    writer := multipart.NewWriter(reqBody)

    writer.WriteField("package_id", datasetIdentifier)
    writer.WriteField("name", "events.geojson")
    writer.WriteField("url", "upload")
    writer.WriteField("description", "GeoJSON data containing events spatial information")
    writer.WriteField("format", "GeoJSON")
    writer.WriteField("resource_type", "donnees")

    part, err := writer.CreateFormField("upload")
    if err != nil {
        return err
    }

    _, err = part.Write(geoJsonData)
    if err != nil {
        return err
    }

    err = writer.Close()
    if err != nil {
        return err
    }

    req, err := http.NewRequest("POST", url, reqBody)
    if err != nil {
        return err
    }

    req.Header.Set("Authorization", apiKey)
    req.Header.Set("Content-Type", writer.FormDataContentType())

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
    }

    return nil
}
