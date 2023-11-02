package utils

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/valyala/fasthttp"
)

func (c *Config) SolveRecaptcha() (string, error) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURIBytes([]byte(`https://api.capsolver.com/createTask`))
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.Set("Content-Type", "application/json")
	req.SetBodyRaw([]byte(`{"clientKey": "` + c.CapSolverKey + `","task": {"type": "ReCaptchaV3M1TaskProxyLess","websiteURL": "https://picsart.com","websiteKey": "6LdM2s8cAAAAAN7jqVXAqWdDlQ3Qca88ke3xdtpR","pageAction": "signup"}}`))
	resp := fasthttp.AcquireResponse()
	for {
		if err := fasthttp.Do(req, resp); err != nil {
			continue
		}
		break
	}
	defer fasthttp.ReleaseResponse(resp)
	defer fasthttp.ReleaseRequest(req)
	switch resp.StatusCode() {
	case 200:
		var taskId struct {
			TaskId string `json:"taskId"`
		}
		if err := json.Unmarshal(resp.Body(), &taskId); err != nil {
			return "", err
		}
		solution, err := c.getSolutionSolv(taskId.TaskId)
		if err != nil {
			return "", err
		}
		return solution, nil
	default:
		return "", errors.New("error creating task")
	}
}

func (c *Config) getSolutionSolv(taskID string) (string, error) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURIBytes([]byte(`https://api.capsolver.com/getTaskResult`))
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.Set("Content-Type", "application/json")
	req.SetBodyRaw([]byte(`{"clientKey": "` + c.CapSolverKey + `","taskId": "` + taskID + `"}`))
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	defer fasthttp.ReleaseRequest(req)
	for {
		time.Sleep(200 * time.Millisecond)
		if err := fasthttp.Do(req, resp); err != nil {
			continue
		}
		switch resp.StatusCode() {
		case 200:
			var solution struct {
				Solution struct {
					CaptchaKey string `json:"gRecaptchaResponse"`
				} `json:"solution"`
			}
			if err := json.Unmarshal(resp.Body(), &solution); err != nil {
				return "", err
			}
			if solution.Solution.CaptchaKey == "" {
				continue
			}
			return solution.Solution.CaptchaKey, nil
		default:
			return "", errors.New("errors in taskID")
		}
	}
}
