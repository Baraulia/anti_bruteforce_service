package scripts

import (
	"bytes"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/cucumber/godog"
)

var ip = "192.1.1.0/25"
var login = "TestLogin"
var password = "TestPassword"
var urlClearAllBuckets = "http://ab_service:8085/clearAllBuckets"

var dataBaseConnectionString = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
	"postgres", "password", os.Getenv("POSTGRES_HOST"), 5432, "postgres")

type serviceTest struct {
	httpClient http.Client
	db         *sql.DB

	responseCode int
	responseBody []byte
}

func (test *serviceTest) setupTest() error {
	db, err := sql.Open("postgres", dataBaseConnectionString)
	if err != nil {
		return err
	}

	test.db = db
	test.httpClient = http.Client{}

	return nil
}

func (test *serviceTest) tearDownTest() error {
	log.Println("Clearing database from service test...")
	query := `DELETE FROM black_list`
	_, err := test.db.Exec(query)
	if err != nil {
		return err
	}

	query = `DELETE FROM white_list`
	_, err = test.db.Exec(query)
	if err != nil {
		return err
	}

	err = test.db.Close()
	if err != nil {
		return err
	}

	request, err := http.NewRequest("PUT", urlClearAllBuckets, nil)
	if err != nil {
		return err
	}

	_, err = test.httpClient.Do(request)
	if err != nil {
		return err
	}

	return nil
}

func (test *serviceTest) iSendRequestTo(httpMethod, addr string, count int) error {
	var response *http.Response
	body := fmt.Sprintf(`{"ip":"%s","login":"%s","password":"%s"}`, ip, login, password)
	switch httpMethod {
	case http.MethodPost:
		for i := 0; i < count; i++ {
			request, err := http.NewRequest("POST", addr, bytes.NewBuffer([]byte(body)))
			if err != nil {
				return err
			}

			response, err = test.httpClient.Do(request)
			if err != nil {
				return err
			}

			if i == count-1 && response.StatusCode == http.StatusOK {
				responseBody, err := io.ReadAll(response.Body)
				if err != nil {
					return err
				}

				test.responseBody = responseBody
			}

			if err = response.Body.Close(); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("unsupported method: %s", httpMethod)
	}

	test.responseCode = response.StatusCode

	return nil
}

func (test *serviceTest) iSendRequestWith(addr string, count int, field string) error {
	var requests []*http.Request
	switch field {
	case "ip":
		for i := 0; i < count; i++ {
			request, err := http.NewRequest("POST", addr, bytes.NewBuffer([]byte(
				fmt.Sprintf(`{"ip":"%s","login":"%s%d","password":"%s%d"}`, ip, login, i, password, i))))
			if err != nil {
				return err
			}
			requests = append(requests, request)
		}
	case "password":
		for i := 0; i < count; i++ {
			request, err := http.NewRequest("POST", addr, bytes.NewBuffer([]byte(
				fmt.Sprintf(`{"ip":"%s%d","login":"%s%d","password":"%s"}`, ip[:len(ip)-1], i, login, i, password))))
			if err != nil {
				return err
			}
			requests = append(requests, request)
		}
	case "login":
		for i := 0; i < count; i++ {
			request, err := http.NewRequest("POST", addr, bytes.NewBuffer([]byte(
				fmt.Sprintf(`{"ip":"%s%d","login":"%s","password":"%s%d"}`, ip[:len(ip)-1], i, login, password, i))))
			if err != nil {
				return err
			}
			requests = append(requests, request)
		}
	default:
		return fmt.Errorf("unsupported field: %s", field)
	}

	err := test.doRequest(requests)
	if err != nil {
		return err
	}

	return nil
}

func (test *serviceTest) doRequest(requests []*http.Request) error {
	var response *http.Response
	var err error
	for i, request := range requests {
		response, err = test.httpClient.Do(request)
		if err != nil {
			return err
		}

		if i == len(requests)-1 && response.StatusCode == http.StatusOK {
			responseBody, err := io.ReadAll(response.Body)
			if err != nil {
				return err
			}

			test.responseBody = responseBody
			test.responseCode = response.StatusCode
		}

		if err = response.Body.Close(); err != nil {
			return err
		}
	}

	return nil
}

func (test *serviceTest) theResponseCodeShouldBe(code int) error {
	if test.responseCode != code {
		return fmt.Errorf("unexpected status code: %d != %d", test.responseCode, code)
	}
	return nil
}

func (test *serviceTest) iReceiveResponse(response string) error {
	if string(test.responseBody) != response {
		return fmt.Errorf("unexpected response: %s != %s", test.responseBody, response)
	}

	return nil
}

func InitializeServiceScenario(ctx *godog.ScenarioContext) {
	test := &serviceTest{httpClient: http.Client{}}

	ctx.Step(`Setup test for service`, test.setupTest)
	ctx.Step(`^I send "([^"]*)" request to "([^"]*)" (\d+) times$`, test.iSendRequestTo)
	ctx.Step(`^The response code should be (\d+)$`, test.theResponseCodeShouldBe)
	ctx.Step(`^I send "POST" request to "([^"]*)" (\d+) times with the same ([^"]*)$`, test.iSendRequestWith)
	ctx.Step(`^I receive response - "([^"]*)"$`, test.iReceiveResponse)
	ctx.Step(`Teardown test for service`, test.tearDownTest)
}
