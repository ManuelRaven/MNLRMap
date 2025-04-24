package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"sync"

	"time"

	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/protomaps/go-pmtiles/pmtiles"
)

type Builds []Build

type Build struct {
	Key      string    `json:"key"`
	Size     int64     `json:"size"`
	Uploaded time.Time `json:"uploaded"`
	Version  string    `json:"version"`
	Md5Sum   *string   `json:"md5sum,omitempty"`
	B3Sum    *string   `json:"b3sum,omitempty"`
}

// Define a structure to track extraction status
type ExtractionStatus struct {
	Status    string // "pending", "processing", "completed", "failed"
	StartTime time.Time
	EndTime   time.Time
	Error     string
	RecordID  string
}

var useExperimental = false                                                // Set to true to use experimental features
var buildsMetadataUrl = "https://build-metadata.protomaps.dev/builds.json" // URL to fetch the latest builds metadata

// fs to maps folder

// mountFs configures the embedded file system for the application's
// front-end assets when building for production.
func (app *application) usePMTiles() {

	// Create a map to store extraction statuses with mutex for concurrent access
	extractionStatuses := make(map[string]*ExtractionStatus)
	var extractionMutex = &sync.RWMutex{}

	pbdbdir := app.pb.App.DataDir()

	mapsFS := os.DirFS(pbdbdir + "/maps")

	// Check if maps folder exists, if not create it
	if _, err := os.Stat(pbdbdir + "/maps"); os.IsNotExist(err) {
		err := os.Mkdir(pbdbdir+"/maps", 0755)
		if err != nil {
			log.Fatalf("Failed to create maps directory: %v", err)
		}
	}

	// register "GET /hello/{name}" route (allowed for everyone)
	app.pb.OnServe().BindFunc(func(se *core.ServeEvent) error {
		// register "GET /hello/{name}" route (allowed for everyone)
		se.Router.GET("/maps", func(e *core.RequestEvent) error {

			// Serve a json list of all files in the maps directory
			files, err := os.ReadDir(pbdbdir + "/maps")
			if err != nil {
				return e.String(http.StatusInternalServerError, "Failed to read maps directory")
			}
			var fileList []string
			for _, file := range files {
				if file.IsDir() {
					continue
				}
				fileList = append(fileList, file.Name())
			}
			return e.JSON(http.StatusOK, fileList)

		})

		// map file size info
		se.Router.GET("/maps/size/{name}", func(e *core.RequestEvent) error {
			name := e.Request.PathValue("name")
			if name == "" {
				return e.BadRequestError("Missing name parameter", nil)
			}
			// Check if the file exists
			filePath := pbdbdir + "/maps/" + name + ".pmtiles"
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				return e.NotFoundError("File not found", err)
			}
			// Get the file info
			fileInfo, err := os.Stat(filePath)
			if err != nil {
				return e.InternalServerError("Failed to get file info", err)
			}
			// Return the file size as json
			e.Response.Header().Set("Content-Type", "application/json")
			return e.JSON(http.StatusOK, map[string]any{
				"name":      name,
				"sizeBytes": fileInfo.Size(),
			})
		})

		// map info
		se.Router.GET("/maps/info/{name}", func(e *core.RequestEvent) error {
			name := e.Request.PathValue("name")
			if name == "" {
				return e.BadRequestError("Missing name parameter", nil)
			}

			// Check if the file exists
			filePath := pbdbdir + "/maps/" + name + ".pmtiles"
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				return e.NotFoundError("File not found", err)
			}
			logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

			// Create a io.writer to catch the output in a variable
			buf := new(bytes.Buffer)

			// Get the file info
			pmtiles.Show(logger, buf, "", filePath, true, false, false, "", false, 0, 0, 0)

			e.Response.Header().Set("Content-Type", "application/json")
			// return the file info as json
			return e.String(http.StatusOK, buf.String())
		})

		// register "POST /maps/create" route for creating map extractions (allowed only for authenticated users)
		se.Router.POST("/maps/create", func(e *core.RequestEvent) error {
			data := struct {
				Name string `json:"name" form:"name"`
				Bbox string `json:"bbox" form:"bbox"`
			}{}
			if err := e.BindBody(&data); err != nil {
				return e.BadRequestError("Failed to read request data", err)
			}

			if data.Name == "" {
				return e.BadRequestError("Missing name parameter", nil)
			}
			if data.Bbox == "" {
				return e.BadRequestError("Missing bbox parameter", nil)
			}

			// Create the file in the maps directory with name.pmtiles
			filePath := pbdbdir + "/maps/" + data.Name + ".pmtiles"
			if _, err := os.Stat(filePath); !os.IsNotExist(err) {
				return e.BadRequestError("File already exists", nil)
			}

			// Add Record to "maps" collection
			collection, err := app.pb.FindCollectionByNameOrId("maps")
			if err != nil {
				return err
			}

			record := core.NewRecord(collection)
			record.Set("name", data.Name)
			record.Set("bbox", data.Bbox)
			record.Set("status", "pending")

			if err := app.pb.Save(record); err != nil {
				return err
			}

			// Create extraction status entry
			extractionMutex.Lock()
			extractionStatuses[record.Id] = &ExtractionStatus{
				Status:    "pending",
				StartTime: time.Now(),
				RecordID:  record.Id,
			}
			extractionMutex.Unlock()

			// Start extraction process in background
			startMapExtraction(extractionMutex, extractionStatuses, app, record, data, filePath)

			// Return immediately with record ID for status checking
			return e.JSON(http.StatusOK, map[string]string{
				"id":      record.Id,
				"status":  "pending",
				"message": "Map extraction started in background",
			})
		})

		se.Router.GET("/maps/serve/{path...}", apis.Static(mapsFS, true))

		se.Router.GET("/maps/status/{id}", func(e *core.RequestEvent) error {
			id := e.Request.PathValue("id")
			if id == "" {
				return e.BadRequestError("Missing id parameter", nil)
			}

			// Check if extraction exists in our status map
			extractionMutex.RLock()
			status, exists := extractionStatuses[id]
			extractionMutex.RUnlock()

			if !exists {
				// If not in our map, try to get from database
				record, err := app.pb.FindRecordById("maps", id)
				if err != nil {
					return e.NotFoundError("Extraction not found", err)
				}

				// Return status from database
				return e.JSON(http.StatusOK, map[string]any{
					"id":     id,
					"status": record.GetString("status"),
					"error":  record.GetString("error"),
					"name":   record.GetString("name"),
				})
			}

			// Return status from our map
			response := map[string]any{
				"id":     id,
				"status": status.Status,
			}

			// Add additional info based on status
			if status.Status == "failed" {
				response["error"] = status.Error
			}
			if !status.EndTime.IsZero() {
				response["duration"] = status.EndTime.Sub(status.StartTime).String()
			}

			return e.JSON(http.StatusOK, response)
		})

		// Recreate a extraction with the information from the record
		se.Router.POST("/maps/recreate/{id}", func(e *core.RequestEvent) error {
			id := e.Request.PathValue("id")
			if id == "" {
				return e.BadRequestError("Missing id parameter", nil)
			}
			// Get the record from the database
			record, err := app.pb.FindRecordById("maps", id)
			if err != nil {
				return e.NotFoundError("Extraction not found", err)
			}

			// If status is not failed or completed, return an error
			if record.GetString("status") != "failed" && record.GetString("status") != "completed" {
				return e.BadRequestError("Extraction is not in a failed or completed state", nil)
			}

			// Delete the old file if it exists
			filePath := pbdbdir + "/maps/" + record.GetString("name") + ".pmtiles"
			if _, err := os.Stat(filePath); !os.IsNotExist(err) {
				os.Remove(filePath)
			}

			// Start the extraction process again with the same parameters in the background
			data := struct {
				Name string `json:"name" form:"name"`
				Bbox string `json:"bbox" form:"bbox"`
			}{
				Name: record.GetString("name"),
				Bbox: record.GetString("bbox"),
			}
			extractionMutex.Lock()
			extractionStatuses[record.Id] = &ExtractionStatus{
				Status:    "pending",
				StartTime: time.Now(),
				RecordID:  record.Id,
			}
			extractionMutex.Unlock()

			// When in failed state make normal extraction, if in completed state make sync extraction
			if record.GetString("status") == "failed" {
				// Start extraction process in background
				startMapExtraction(extractionMutex, extractionStatuses, app, record, data, filePath)
			} else {
				// Start sync process in background

				if useExperimental {
					startSyncMap(extractionMutex, extractionStatuses, app, record, data.Name, filePath)
				} else {
					startMapExtraction(extractionMutex, extractionStatuses, app, record, data, filePath)
				}

			}

			// Update the record status in the database
			record.Set("status", "pending")
			record.Set("error", nil)
			_ = app.pb.Save(record)

			// Return immediately with record ID for status checking
			return e.JSON(http.StatusAccepted, map[string]string{
				"id":      record.Id,
				"status":  "pending",
				"message": "Map extraction started in background",
			})
		})
		return se.Next()
	})

	// Delete a file if they exist after a record is deleted
	app.pb.OnRecordAfterDeleteSuccess("maps").BindFunc(func(e *core.RecordEvent) error {

		// Get the file path from the record
		filePath := pbdbdir + "/maps/" + e.Record.GetString("name") + ".pmtiles"
		// Check if the file exists
		if _, err := os.Stat(filePath); !os.IsNotExist(err) {
			// Delete the file
			os.Remove(filePath)
		}

		return e.Next()
	})
}

func startSyncMap(extractionMutex *sync.RWMutex, extractionStatuses map[string]*ExtractionStatus, app *application, record *core.Record, data string, filePath string) {
	go func(recordID string, name string, filePath string) {
		// Update status to processing
		extractionMutex.Lock()
		extractionStatuses[recordID].Status = "processing"
		extractionMutex.Unlock()

		// Update the record status in the database
		record, err := app.pb.FindRecordById("maps", recordID)
		if err == nil {
			record.Set("status", "processing")
			_ = app.pb.Save(record)
		}

		buildurl, err := FetchLatestBuildURL()
		if err != nil {
			fmt.Print(err.Error())
			// Update status to failed
			extractionMutex.Lock()
			extractionStatuses[recordID].Status = "failed"
			extractionStatuses[recordID].Error = err.Error()
			extractionMutex.Unlock()
			if record, err := app.pb.FindRecordById("maps", recordID); err == nil {
				record.Set("status", "failed")
				record.Set("error", extractionStatuses[recordID].Error)
				_ = app.pb.Save(record)
			}
			return
		}

		logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
		// Run the extraction process
		err = pmtiles.Sync(logger, filePath, buildurl, false)

		extractionMutex.Lock()
		extractionStatuses[recordID].EndTime = time.Now()
		if err != nil {

			fmt.Print(err.Error())

			// Update status to failed
			extractionStatuses[recordID].Status = "failed"
			extractionStatuses[recordID].Error = err.Error()

			if record, err := app.pb.FindRecordById("maps", recordID); err == nil {
				record.Set("status", "failed")
				record.Set("error", extractionStatuses[recordID].Error)
				_ = app.pb.Save(record)
			}
		} else {
			// Update status to completed
			extractionStatuses[recordID].Status = "completed"

			if record, err := app.pb.FindRecordById("maps", recordID); err == nil {
				record.Set("status", "completed")
				_ = app.pb.Save(record)
			}
		}
		extractionMutex.Unlock()
	}(record.Id, data, filePath)
}

func startMapExtraction(extractionMutex *sync.RWMutex, extractionStatuses map[string]*ExtractionStatus, app *application, record *core.Record, data struct {
	Name string "json:\"name\" form:\"name\""
	Bbox string "json:\"bbox\" form:\"bbox\""
}, filePath string) {
	go func(recordID string, name string, bbox string, filePath string) {
		// Update status to processing
		extractionMutex.Lock()
		extractionStatuses[recordID].Status = "processing"
		extractionMutex.Unlock()

		// Update the record status in the database
		record, err := app.pb.FindRecordById("maps", recordID)
		if err == nil {
			record.Set("status", "processing")
			_ = app.pb.Save(record)
		}

		logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

		buildurl, err := FetchLatestBuildURL()
		if err != nil {
			fmt.Print(err.Error())
			// Update status to failed
			extractionMutex.Lock()
			extractionStatuses[recordID].Status = "failed"
			extractionStatuses[recordID].Error = err.Error()
			extractionMutex.Unlock()
			if record, err := app.pb.FindRecordById("maps", recordID); err == nil {
				record.Set("status", "failed")
				record.Set("error", extractionStatuses[recordID].Error)
				_ = app.pb.Save(record)
			}
			return
		}

		// Run the extraction process
		err = pmtiles.Extract(logger, "", buildurl, -1, -1, "", bbox, filePath, 4, 0.05, false)

		extractionMutex.Lock()
		extractionStatuses[recordID].EndTime = time.Now()
		if err != nil {

			fmt.Print(err.Error())

			// Update status to failed
			extractionStatuses[recordID].Status = "failed"
			extractionStatuses[recordID].Error = err.Error()

			if record, err := app.pb.FindRecordById("maps", recordID); err == nil {
				record.Set("status", "failed")
				record.Set("error", extractionStatuses[recordID].Error)
				_ = app.pb.Save(record)
			}
		} else {
			// Update status to completed
			extractionStatuses[recordID].Status = "completed"

			if record, err := app.pb.FindRecordById("maps", recordID); err == nil {
				record.Set("status", "completed")
				_ = app.pb.Save(record)
			}
		}
		extractionMutex.Unlock()
	}(record.Id, data.Name, data.Bbox, filePath)
}

// FetchLatestBuildURL fetches the latest build URL from the metadata server (sort on uploaded date)
func FetchLatestBuildURL() (string, error) {
	resp, err := http.Get(buildsMetadataUrl)
	if err != nil {
		return "", fmt.Errorf("failed to fetch builds metadata: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch builds metadata: %s", resp.Status)
	}

	var builds Builds
	if err := json.NewDecoder(resp.Body).Decode(&builds); err != nil {
		return "", fmt.Errorf("failed to decode builds metadata: %w", err)
	}

	if len(builds) == 0 {
		return "", fmt.Errorf("no builds found")
	}

	// Sort builds by uploaded date (newest first)
	sort.Slice(builds, func(i, j int) bool {
		return builds[i].Uploaded.After(builds[j].Uploaded)
	})

	var latestBuildURL string = fmt.Sprintf("https://build.protomaps.com/%s", builds[0].Key)
	fmt.Printf("Latest build URL: %s\n", latestBuildURL)

	return latestBuildURL, nil
}
