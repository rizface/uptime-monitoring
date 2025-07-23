package services

import (
	"fmt"
	"net/http"
	"time"
	"uptime-monitoring/models"
)

type UptimeChecker struct {
	urlRepo *models.URLRepository
}

func NewUptimeChecker(urlRepo *models.URLRepository) *UptimeChecker {
	return &UptimeChecker{
		urlRepo: urlRepo,
	}
}

func (uc *UptimeChecker) CheckURL(url string) (string, int, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	start := time.Now()
	resp, err := client.Get(url)
	responseTime := int(time.Since(start).Milliseconds())

	if err != nil {
		return "down", responseTime, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return "up", responseTime, nil
	}

	// Handle specific HTTP status codes
	switch resp.StatusCode {
	case 401:
		return "unauthorized", responseTime, nil
	case 403:
		return "forbidden", responseTime, nil
	default:
		return "down", responseTime, fmt.Errorf("HTTP status: %d", resp.StatusCode)
	}
}

func (uc *UptimeChecker) CheckAllURLs() error {
	urls, err := uc.urlRepo.GetAll()
	if err != nil {
		return err
	}

	for _, url := range urls {
		status, responseTime, _ := uc.CheckURL(url.URL)
		err := uc.urlRepo.UpdateStatus(url.ID, status, responseTime)
		if err != nil {
			return err
		}
	}

	return nil
}

func (uc *UptimeChecker) CheckSingleURL(id uint) error {
	url, err := uc.urlRepo.GetByID(id)
	if err != nil {
		return err
	}

	status, responseTime, _ := uc.CheckURL(url.URL)
	return uc.urlRepo.UpdateStatus(id, status, responseTime)
}
