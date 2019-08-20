package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jivalabs/picsum-photos/api/handler"
	"github.com/jivalabs/picsum-photos/api/params"
	"github.com/jivalabs/picsum-photos/database"
	"github.com/jivalabs/picsum-photos/image"
)

// /img/width/height/resize/id
func (a *API) imageHandler(w http.ResponseWriter, r *http.Request) *handler.Error {
	// Get the path and query parameters
	uri := r.RequestURI
	uriSlice := strings.Split(uri, "/")
	if len(uriSlice) < 4 {
		return nil
	}
	width, _ := strconv.Atoi(uriSlice[2])
	height, _ := strconv.Atoi(uriSlice[3])
	cmd := uriSlice[4]

	from := strings.Index(uri, cmd) + len(cmd)
	reminder := uri[from+1:]

	log.Printf("width: %d height: %d cmd: %s reminder: %s", width, height, cmd, reminder)

	databaseImage, err := a.Database.Get(reminder)
	if err != nil {
		if err == database.ErrNotFound {
			return &handler.Error{Message: err.Error(), Code: http.StatusNotFound}
		}

		a.logError(r, "error getting image from database", err)
		return handler.InternalServerError()
	}

	// Validate the parameters
	//if err := params.ValidateParams(a.MaxImageSize, databaseImage, p); err != nil {
	//	return handler.BadRequest(err.Error())
	//}

	// Default to the image width/height if 0 is passed

	if width == 0 {
		width = databaseImage.Width
	}

	if height == 0 {
		height = databaseImage.Height
	}

	// Build the image task
	task := image.NewTask(databaseImage.ID, width, height)

	//if cmd == "Blur" {
	//	task.Blur(p.BlurAmount)
	//}
	//
	//if cmd == "Grayscale" {
	//	task.Grayscale()
	//}

	// Process the image
	processedImage, err := a.ImageProcessor.ProcessImage(r.Context(), task)
	if err != nil {
		a.logError(r, "error processing image", err)
		return handler.InternalServerError()
	}

	// Set the headers
	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", buildFilenameNew(reminder, width, height)))
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Cache-Control", "public, max-age=2592000") // Cache for a month

	// Return the image
	w.Write(processedImage)
	return nil
}

func (a *API) imageHandlerOld(w http.ResponseWriter, r *http.Request) *handler.Error {
	// Get the path and query parameters
	p, err := params.GetParams(r)
	if err != nil {
		return handler.BadRequest(err.Error())
	}

	// Get the image from the database
	vars := mux.Vars(r)
	imageID := vars["id"]
	databaseImage, err := a.Database.Get(imageID)
	if err != nil {
		if err == database.ErrNotFound {
			return &handler.Error{Message: err.Error(), Code: http.StatusNotFound}
		}

		a.logError(r, "error getting image from database", err)
		return handler.InternalServerError()
	}

	// Validate the parameters
	if err := params.ValidateParams(a.MaxImageSize, databaseImage, p); err != nil {
		return handler.BadRequest(err.Error())
	}

	// Default to the image width/height if 0 is passed
	width := p.Width
	height := p.Height

	if width == 0 {
		width = databaseImage.Width
	}

	if height == 0 {
		height = databaseImage.Height
	}

	// Build the image task
	task := image.NewTask(databaseImage.ID, width, height)
	if p.Blur {
		task.Blur(p.BlurAmount)
	}

	if p.Grayscale {
		task.Grayscale()
	}

	// Process the image
	processedImage, err := a.ImageProcessor.ProcessImage(r.Context(), task)
	if err != nil {
		a.logError(r, "error processing image", err)
		return handler.InternalServerError()
	}

	// Set the headers
	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", buildFilename(imageID, p, width, height)))
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Cache-Control", "public, max-age=2592000") // Cache for a month

	// Return the image
	w.Write(processedImage)

	return nil
}

func buildFilenameNew(imageID string, width int, height int) string {
	filename := fmt.Sprintf("%s-%dx%d", imageID, width, height)

	//if p.Blur {
	//	filename += fmt.Sprintf("-blur_%d", p.BlurAmount)
	//}
	//
	//if p.Grayscale {
	//	filename += "-grayscale"
	//}

	//if imageID == "" {
	//	filename += ".jpg"
	//} else {
	//	ext := filepath.Ext(imageID)
	//	filename += ext
	//}

	return filename
}

func buildFilename(imageID string, p *params.Params, width int, height int) string {
	filename := fmt.Sprintf("%s-%dx%d", imageID, width, height)

	if p.Blur {
		filename += fmt.Sprintf("-blur_%d", p.BlurAmount)
	}

	if p.Grayscale {
		filename += "-grayscale"
	}

	if p.Extension == "" {
		filename += ".jpg"
	} else {
		filename += p.Extension
	}

	return filename
}

func (a *API) imageRedirectHandler(w http.ResponseWriter, r *http.Request) *handler.Error {
	// Get the path and query parameters
	p, err := params.GetParams(r)
	if err != nil {
		return handler.BadRequest(err.Error())
	}

	// Get the image ID
	vars := mux.Vars(r)
	imageID := vars["id"]

	// Redirect
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	http.Redirect(w, r, fmt.Sprintf("/id/%s/%d/%d%s%s", imageID, p.Width, p.Height, p.Extension, params.BuildQuery(p.Grayscale, p.Blur, p.BlurAmount)), http.StatusFound)
	return nil
}

func (a *API) randomImageRedirectHandler(w http.ResponseWriter, r *http.Request) *handler.Error {
	// Get the path and query parameters
	p, err := params.GetParams(r)
	if err != nil {
		return handler.BadRequest(err.Error())
	}

	// Get a random image
	randomImage, err := a.Database.GetRandom()
	if err != nil {
		a.logError(r, "error getting random image from database", err)
		return handler.InternalServerError()
	}

	// Redirect
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	http.Redirect(w, r, fmt.Sprintf("/id/%s/%d/%d%s%s", randomImage, p.Width, p.Height, p.Extension, params.BuildQuery(p.Grayscale, p.Blur, p.BlurAmount)), http.StatusFound)
	return nil
}
