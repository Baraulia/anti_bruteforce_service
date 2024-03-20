package handlers

//nolint:depguard
import (
	"net/http"
	"regexp"

	"github.com/Baraulia/anti_bruteforce_service/internal/api"
	"github.com/Baraulia/anti_bruteforce_service/internal/app"
)

type Handler struct {
	logger    app.Logger
	app       api.ApplicationInterface
	patternIP *regexp.Regexp
}

func NewHandler(logger app.Logger, app api.ApplicationInterface) *Handler {
	validatePatternIP := regexp.MustCompile(`^(?:[0-9]{1,3}\.){3}[0-9]{1,3}/([0-9]|[1-2][0-9]|3[0-2])$`)
	return &Handler{logger: logger, app: app, patternIP: validatePatternIP}
}

func (h *Handler) InitRoutes() *http.ServeMux {
	r := http.NewServeMux()

	r.HandleFunc("/check", h.check)
	r.HandleFunc("/blacklist/add", h.addToBlacklist)         // method POST
	r.HandleFunc("/blacklist/remove", h.deleteFromBlacklist) // method DELETE
	r.HandleFunc("/whitelist/add", h.addToWhitelist)         // method POST
	r.HandleFunc("/whitelist/remove", h.deleteFromWhitelist) // method DELETE
	r.HandleFunc("/clear", h.clearBuckets)

	return r
}
