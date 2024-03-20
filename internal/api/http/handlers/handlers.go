package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Baraulia/anti_bruteforce_service/internal/models"
)

type DataIP struct {
	IP string `json:"ip"`
}

func (h *Handler) check(response http.ResponseWriter, request *http.Request) {
	if !h.isMethodAllowed(response, request, http.MethodPost) {
		return
	}

	var input models.Data

	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(&input); err != nil {
		h.logger.Error("Error while decoding request", map[string]interface{}{"error": err})
		http.Error(response, err.Error(), 400)
		return
	}

	err := h.validateData(input, true)
	if err != nil {
		h.logger.Error("Invalid input request received from client", map[string]interface{}{"error": err})
		http.Error(response, err.Error(), 400)
		return
	}

	allowed, err := h.app.Check(request.Context(), input)
	if err != nil {
		h.logger.Error("server error", map[string]interface{}{"error": err})
		http.Error(response, err.Error(), 500)
		return
	}

	response.WriteHeader(http.StatusOK)
	_, err = response.Write([]byte(fmt.Sprintf("ok=%t", allowed)))
	if err != nil {
		h.logger.Error("error while writing response", map[string]interface{}{"error": err})
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
}

//nolint:dupl
func (h *Handler) addToBlacklist(response http.ResponseWriter, request *http.Request) {
	if !h.isMethodAllowed(response, request, http.MethodPost) {
		return
	}

	var input DataIP

	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(&input); err != nil {
		h.logger.Error("Error while decoding request", map[string]interface{}{"error": err})
		http.Error(response, err.Error(), 400)
		return
	}

	err := h.validateIP(input.IP)
	if err != nil {
		h.logger.Error("Invalid ip received from client", map[string]interface{}{"error": err})
		http.Error(response, err.Error(), 400)
		return
	}

	err = h.app.AddToBlackList(request.Context(), input.IP)
	if err != nil {
		h.logger.Error("server error", map[string]interface{}{"error": err})
		http.Error(response, err.Error(), 500)
		return
	}

	response.WriteHeader(http.StatusCreated)
}

func (h *Handler) deleteFromBlacklist(response http.ResponseWriter, request *http.Request) {
	if !h.isMethodAllowed(response, request, http.MethodDelete) {
		return
	}

	params := request.URL.Query()
	ip := params.Get("ip")
	err := h.validateIP(ip)
	if err != nil {
		h.logger.Error("Invalid ip received from client", map[string]interface{}{"error": err})
		http.Error(response, err.Error(), 400)
		return
	}

	err = h.app.RemoveFromBlackList(request.Context(), ip)
	if err != nil {
		h.logger.Error("server error", map[string]interface{}{"error": err})
		http.Error(response, err.Error(), 500)
		return
	}

	response.WriteHeader(http.StatusNoContent)
}

//nolint:dupl
func (h *Handler) addToWhitelist(response http.ResponseWriter, request *http.Request) {
	if !h.isMethodAllowed(response, request, http.MethodPost) {
		return
	}

	var input DataIP

	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(&input); err != nil {
		h.logger.Error("Error while decoding request", map[string]interface{}{"error": err})
		http.Error(response, err.Error(), 400)
		return
	}

	err := h.validateIP(input.IP)
	if err != nil {
		h.logger.Error("Invalid ip received from client", map[string]interface{}{"error": err})
		http.Error(response, err.Error(), 400)
		return
	}

	err = h.app.AddToWhiteList(request.Context(), input.IP)
	if err != nil {
		h.logger.Error("server error", map[string]interface{}{"error": err})
		http.Error(response, err.Error(), 500)
		return
	}

	response.WriteHeader(http.StatusCreated)
}

func (h *Handler) deleteFromWhitelist(response http.ResponseWriter, request *http.Request) {
	if !h.isMethodAllowed(response, request, http.MethodDelete) {
		return
	}

	params := request.URL.Query()
	ip := params.Get("ip")
	err := h.validateIP(ip)
	if err != nil {
		h.logger.Error("Invalid ip received from client", map[string]interface{}{"error": err})
		http.Error(response, err.Error(), 400)
		return
	}

	err = h.app.RemoveFromWhiteList(request.Context(), ip)
	if err != nil {
		h.logger.Error("server error", map[string]interface{}{"error": err})
		http.Error(response, err.Error(), 500)
		return
	}

	response.WriteHeader(http.StatusNoContent)
}

func (h *Handler) clearBuckets(response http.ResponseWriter, request *http.Request) {
	if !h.isMethodAllowed(response, request, http.MethodPost) {
		return
	}

	var input models.Data

	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(&input); err != nil {
		h.logger.Error("Error while decoding request", map[string]interface{}{"error": err})
		http.Error(response, err.Error(), 400)
		return
	}

	err := h.validateData(input, false)
	if err != nil {
		h.logger.Error("Invalid input request received from client", map[string]interface{}{"error": err})
		http.Error(response, err.Error(), 400)
		return
	}

	err = h.app.ClearBuckets(request.Context(), input)
	if err != nil {
		h.logger.Error("server error", map[string]interface{}{"error": err})
		http.Error(response, err.Error(), 500)
		return
	}

	response.WriteHeader(http.StatusNoContent)
}

func (h *Handler) isMethodAllowed(response http.ResponseWriter, request *http.Request, allowedMethod string) bool {
	if request.Method != allowedMethod {
		h.logger.Error(fmt.Sprintf(
			"method %s not allowed", request.Method), map[string]interface{}{"want method": allowedMethod})
		http.Error(response, fmt.Sprintf(
			"Method Not Allowed(want method %s)", allowedMethod), http.StatusMethodNotAllowed)
		return false
	}

	return true
}
