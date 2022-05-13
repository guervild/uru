package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/guervild/uru/pkg/builder"
	"github.com/guervild/uru/pkg/logger"

	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (a *App) Run(addr string) {
	logger.Logger.Info().Str("addr", addr).Msg("Starting http server")
	logger.Logger.Fatal().Msgf("%s", http.ListenAndServe(addr, a.Router))
}

func (a *App) generatePayload(w http.ResponseWriter, r *http.Request) {

	ip := r.RemoteAddr
	xforward := r.Header.Get("X-Forwarded-For")
	logger.Logger.Info().Str("IP", ip).Str("x-forwarded-for", xforward).Msg("New connection to the /generate endpoint")

	payloadFile, payloadFileHeader, err := r.FormFile("payload")
	if err != nil {
		logger.Logger.Error().Err(err).Msg("Error Retrieving the payload")
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error retrieving payload: %s", err.Error()))
		return
	}

	logger.Logger.Info().Str("IP", ip).Str("filename", payloadFileHeader.Filename).Msg("Received a payload")

	configFile, configFileHeader, err := r.FormFile("config")
	if err != nil {
		logger.Logger.Error().Err(err).Msg("Error Retrieving the config")
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error retrieving config: %s", err.Error()))
		return
	}

	logger.Logger.Info().Str("IP", ip).Str("filename", configFileHeader.Filename).Msg("Received a config file")

	fDonut := r.URL.Query().Get("donut")
	godonut := false

	if strings.ToLower(fDonut) == "true" {
		logger.Logger.Info().Msg("Payload sent to api is a binary. donut flag has been passed")
		godonut = true
	}

	fSrdi := r.URL.Query().Get("srdi")
	srdi := false

	if strings.ToLower(fSrdi) == "true" {
		logger.Logger.Info().Msg("Payload sent to api is a dll. srdi flag has been passed")
		srdi = true
	}

	fParameters := r.URL.Query().Get("parameters")

	if (godonut || srdi) && fParameters != "" {
		logger.Logger.Info().Msgf("The following parameters will be passed to the payload: %s", fParameters)
	}

	fFunctionName := r.URL.Query().Get("functionName")

	if (godonut || srdi) && fFunctionName != "" {
		logger.Logger.Info().Msgf("The following functionName will be used: %s", fFunctionName)
	}

	fClass := r.URL.Query().Get("class")

	if (godonut || srdi) && fClass != "" {
		logger.Logger.Info().Msgf("The following Class will be used: %s", fClass)
	}

	fClearHeader := r.URL.Query().Get("clearHeader")
	clearHeader := false

	if srdi && strings.ToLower(fClearHeader) == "true" {
		logger.Logger.Info().Msg("clearHeader is passed, PE header will be removed")
		clearHeader = true
	}

	payloadData := bytes.NewBuffer(nil)
	if _, err := io.Copy(payloadData, payloadFile); err != nil {
		logger.Logger.Error().Err(err).Msg("Error copying payload file")
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	configData := bytes.NewBuffer(nil)
	if _, err := io.Copy(configData, configFile); err != nil {
		logger.Logger.Error().Err(err).Msg("Error copying config file")
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var payloadConfig builder.PayloadConfig

	payloadConfig, err = builder.NewPayloadConfigFromFile(configData.Bytes())

	if err != nil {
		logger.Logger.Error().Err(err).Msg("Error while PayloadConfig struc")
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	//Process payload file
	payloadPath, _, err := payloadConfig.GeneratePayload(payloadFileHeader.Filename, payloadData.Bytes(), godonut, srdi, true, fParameters, fFunctionName, fClass, clearHeader)
	if err != nil {
		logger.Logger.Error().Err(err).Msg("Error while generating the payload")
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	name := filepath.Base(payloadPath)
	w.Header().Set("Content-Disposition", "attachment; filename="+name)
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, payloadPath)
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/generate", a.generatePayload).Methods("POST")
}

func (a *App) Initialize() {

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}
