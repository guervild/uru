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

	payloadFile, _, err := r.FormFile("payload")
	if err != nil {
		logger.Logger.Error().Err(err).Msg("Error Retrieving the payload")
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error retrieving payload: %s", err.Error()))
		return
	}

	configFile, _, err := r.FormFile("config")
	if err != nil {
		logger.Logger.Error().Err(err).Msg("Error Retrieving the config")
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error retrieving config: %s", err.Error()))
		return
	}

	fExe := r.URL.Query().Get("exe")
	exe := false

	if strings.ToLower(fExe) == "true" {
		logger.Logger.Info().Msg("Payload sent to api is a binary. exe flag has been passed")
		exe = true
	}

	fParameters := r.URL.Query().Get("parameters")

	if exe && fParameters != "" {
		logger.Logger.Info().Msgf("The following parameters will be passed to the payload: %s", fParameters)
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
	payloadPath, _, err := payloadConfig.GeneratePayload(payloadData.Bytes(), exe, true, fParameters)
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
