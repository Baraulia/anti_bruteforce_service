package handlers

import (
	"errors"
	"github.com/Baraulia/anti_bruteforce_service/internal/models"
)

func (h *Handler) validateData(data models.Data, fullData bool) error {
	var err string
	if data.Login == "" {
		err = "empty login in request"
	}

	if data.Password == "" && fullData {
		err += "empty password in request"
	}

	valid := h.patternIP.MatchString(data.Ip)

	if !valid {
		err += "invalid ip in request"
	}

	if len(err) != 0 {
		return errors.New(err)
	}

	return nil
}

func (h *Handler) validateIP(ip string) error {
	valid := h.patternIP.MatchString(ip)
	if !valid {
		return errors.New("invalid ip in request")
	}

	return nil
}
