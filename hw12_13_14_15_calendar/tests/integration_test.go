//go:build integration
// +build integration

package tests

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/ozontech/cute"
	cutejson "github.com/ozontech/cute/asserts/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

const uriEvents = "http://localhost:8888/events/"
const timeout = 300 * time.Millisecond

func getValidEvent() map[string]interface{} {
	return map[string]interface{}{
		"id":          uuid.New().String(),
		"startTime":   "2022-01-01T15:00:00.0Z",
		"endTime":     "2022-01-01T16:00:00.0Z",
		"title":       "Event title",
		"ownerId":     1,
		"description": "Event description",
	}
}

func getValidEventWithNotification() map[string]interface{} {
	return map[string]interface{}{
		"id":               uuid.New().String(),
		"startTime":        "2022-09-17T15:00:00.0Z",
		"endTime":          "2022-09-17T16:00:00.0Z",
		"notificationTime": "2022-09-17T10:00:00.0Z",
		"title":            "Event title",
		"ownerId":          1,
		"description":      "Event description",
	}
}

func getValidEventForDate(date string) map[string]interface{} {
	return map[string]interface{}{
		"id":          uuid.New().String(),
		"startTime":   date + "T15:00:00.0Z",
		"endTime":     date + "T16:00:00.0Z",
		"title":       "Event title for " + date,
		"ownerId":     1,
		"description": "Event description",
	}
}

func getUriEventsDay(date string) string {
	return "http://localhost:8888/events/day/" + date
}

func getUriEventsWeek(date string) string {
	return "http://localhost:8888/events/week/" + date
}

func getUriEventsMonth(date string) string {
	return "http://localhost:8888/events/month/" + date
}

func getInvalidEvent() map[string]interface{} {
	return map[string]interface{}{
		"id":          uuid.New().String(),
		"startTime":   "2022-09-17T15:00:00.0Z",
		"endTime":     "2022-09-17T14:00:00.0Z",
		"title":       "Event with endTime before startTime",
		"ownerId":     1,
		"description": "Event description",
	}
}

func TestCreateDuplicateEvent(t *testing.T) {
	evBody, err := json.Marshal(getValidEvent())
	assert.NoError(t, err)

	cute.NewTestBuilder().
		Title("Add new event").
		CreateStep("Create new event").
		RequestBuilder(
			cute.WithURI(uriEvents),
			cute.WithMethod(http.MethodPost),
			cute.WithBody(evBody),
		).
		ExpectExecuteTimeout(timeout).
		ExpectStatus(http.StatusCreated).
		NextTest().
		CreateStep("Create event with same id fails").
		RequestBuilder(
			cute.WithURI(uriEvents),
			cute.WithMethod(http.MethodPost),
			cute.WithBody(evBody),
		).
		ExpectExecuteTimeout(timeout).
		ExpectStatus(http.StatusBadRequest).
		AssertBody(
			cutejson.Equal("$.error", "id duplicate"),
		).
		ExecuteTest(context.Background(), t)
}

func TestCreateEventWithInvalidStartEndTime(t *testing.T) {
	evBody, err := json.Marshal(getInvalidEvent())
	assert.NoError(t, err)

	cute.NewTestBuilder().
		Title("Add new event").
		CreateStep("Create new event").
		RequestBuilder(
			cute.WithURI(uriEvents),
			cute.WithMethod(http.MethodPost),
			cute.WithBody(evBody),
		).
		ExpectExecuteTimeout(timeout).
		ExpectStatus(http.StatusBadRequest).
		AssertBody(
			cutejson.Equal("$.error", "end time should be after start time"),
		).
		ExecuteTest(context.Background(), t)
}

func TestEventNotificationHandledBySender(t *testing.T) {
	evBody, err := json.Marshal(getValidEventWithNotification())
	assert.NoError(t, err)

	ctx := context.Background()

	cute.NewTestBuilder().
		Title("Add new event with notification").
		CreateStep("Create event").
		RequestBuilder(
			cute.WithURI(uriEvents),
			cute.WithMethod(http.MethodPost),
			cute.WithBody(evBody),
		).
		ExpectExecuteTimeout(timeout).
		ExpectStatus(http.StatusCreated)

	time.Sleep(10 * time.Millisecond)

	cute.NewTestBuilder().
		Title("Check if event handled successfully").
		CreateStep("List events").
		RequestBuilder(
			cute.WithURI(getUriEventsDay("2022-09-17")),
			cute.WithMethod(http.MethodGet),
		).
		AssertBody(
			cutejson.Equal("$[0].notificationTime", "2022-09-17T10:00:00Z"),
			cutejson.Present("$[0].notifiedAt"),
			cutejson.Present("$[0].notifiedHandledAt"),
		).
		ExecuteTest(ctx, t)

}

func TestListEventDay(t *testing.T) {
	evBody1, err := json.Marshal(getValidEventForDate("2022-05-01"))
	assert.NoError(t, err)
	evBody2, err := json.Marshal(getValidEventForDate("2022-05-02"))
	assert.NoError(t, err)

	cute.NewTestBuilder().
		Title("Add new event").
		CreateStep("Create first event").
		RequestBuilder(
			cute.WithURI(uriEvents),
			cute.WithMethod(http.MethodPost),
			cute.WithBody(evBody1),
		).
		ExpectExecuteTimeout(timeout).
		ExpectStatus(http.StatusCreated).
		NextTest().
		CreateStep("Create second event").
		RequestBuilder(
			cute.WithURI(uriEvents),
			cute.WithMethod(http.MethodPost),
			cute.WithBody(evBody2),
		).
		ExpectExecuteTimeout(timeout).
		ExpectStatus(http.StatusCreated).
		NextTest().
		CreateStep("List events per date 1").
		RequestBuilder(
			cute.WithURI(getUriEventsDay("2022-05-01")),
			cute.WithMethod(http.MethodGet),
		).
		AssertBody(
			cutejson.Equal("$[0].title", "Event title for 2022-05-01"),
		).
		NextTest().
		CreateStep("List events per date 2").
		RequestBuilder(
			cute.WithURI(getUriEventsDay("2022-05-02")),
			cute.WithMethod(http.MethodGet),
		).
		AssertBody(
			cutejson.Equal("$[0].title", "Event title for 2022-05-02"),
		).
		ExecuteTest(context.Background(), t)
}

func TestListEventWeek(t *testing.T) {
	evBody1, err := json.Marshal(getValidEventForDate("2022-07-01"))
	assert.NoError(t, err)
	evBody2, err := json.Marshal(getValidEventForDate("2022-07-02"))
	assert.NoError(t, err)

	cute.NewTestBuilder().
		Title("Add new event").
		CreateStep("Create first event").
		RequestBuilder(
			cute.WithURI(uriEvents),
			cute.WithMethod(http.MethodPost),
			cute.WithBody(evBody1),
		).
		ExpectExecuteTimeout(timeout).
		ExpectStatus(http.StatusCreated).
		NextTest().
		CreateStep("Create second event").
		RequestBuilder(
			cute.WithURI(uriEvents),
			cute.WithMethod(http.MethodPost),
			cute.WithBody(evBody2),
		).
		ExpectExecuteTimeout(timeout).
		ExpectStatus(http.StatusCreated).
		NextTest().
		CreateStep("List events per week").
		RequestBuilder(
			cute.WithURI(getUriEventsWeek("2022-07-01")),
			cute.WithMethod(http.MethodGet),
		).
		AssertBody(
			cutejson.Equal("$[0].title", "Event title for 2022-07-01"),
			cutejson.Equal("$[1].title", "Event title for 2022-07-02"),
		).
		ExecuteTest(context.Background(), t)
}

func TestListEventMonth(t *testing.T) {
	evBody1, err := json.Marshal(getValidEventForDate("2022-06-01"))
	assert.NoError(t, err)
	evBody2, err := json.Marshal(getValidEventForDate("2022-06-02"))
	assert.NoError(t, err)

	cute.NewTestBuilder().
		Title("Add new event").
		CreateStep("Create first event").
		RequestBuilder(
			cute.WithURI(uriEvents),
			cute.WithMethod(http.MethodPost),
			cute.WithBody(evBody1),
		).
		ExpectExecuteTimeout(timeout).
		ExpectStatus(http.StatusCreated).
		NextTest().
		CreateStep("Create second event").
		RequestBuilder(
			cute.WithURI(uriEvents),
			cute.WithMethod(http.MethodPost),
			cute.WithBody(evBody2),
		).
		ExpectExecuteTimeout(timeout).
		ExpectStatus(http.StatusCreated).
		NextTest().
		CreateStep("List events per month").
		RequestBuilder(
			cute.WithURI(getUriEventsMonth("2022-06-01")),
			cute.WithMethod(http.MethodGet),
		).
		AssertBody(
			cutejson.Equal("$[0].title", "Event title for 2022-06-01"),
			cutejson.Equal("$[1].title", "Event title for 2022-06-02"),
		).
		ExecuteTest(context.Background(), t)
}
