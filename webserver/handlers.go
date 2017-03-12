package webserver

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nicelegs/burning.cash/models"
)

func userHandler(c *gin.Context) {
	var (
		userIdStr          = c.Param("id")
		userId, parseIdErr = strconv.ParseInt(userIdStr, 10, 8)

		user    models.User
		userErr error

		positions models.Positions
		events    models.Events

		deltas models.DailyDeltas
	)
	if parseIdErr != nil {
		panic(parseIdErr)
	}
	user, userErr = models.GetUser(models.GetUserOpts{Id: userId})
	if userErr != nil {
		switch userErr {
		case models.UserNotFoundError:
			c.String(404, "NOT FOUND")
			return
		}
	}

	events = models.GetEvents(models.GetEventsOpts{UserId: int(userId)})
	positions = models.GetPositions(models.GetPositionsOpts{UserId: int(userId)})
	deltas = models.GetDailyDeltasForUser(userId)

	var tickerIds = make(map[int32]bool)

	for _, t := range events {
		tid := t.Asset.Ticker.Id
		if tid != 0 {
			tickerIds[tid] = true
		}
	}
	for _, p := range positions {
		tid := p.Asset.Ticker.Id
		if tid != 0 {
			tickerIds[tid] = true
		}
	}

	var tickerIdsSet []int32
	for k, _ := range tickerIds {
		tickerIdsSet = append(tickerIdsSet, k)
	}

	tickers, terr := models.GetTickers(tickerIdsSet)
	if terr != nil {
		panic(terr)
	}

	events = events.FillInTickers(tickers)
	positions = positions.FillInTickers(tickers)

	rc := NewReactComponent("ProfileView", map[string]interface{}{
		"user":        user,
		"events":      events,
		"positions":   positions,
		"dailyDeltas": deltas,
		"tickers":     tickers,
	})

	c.HTML(http.StatusOK, "user:show", struct {
		View      reactComponent
		UserId    string
		Events    models.Events
		Positions models.Positions
		Deltas    models.DailyDeltas
		Tickers   map[int32]models.Ticker
	}{
		View:      rc,
		UserId:    c.Param("id"),
		Events:    events,
		Positions: positions,
		Deltas:    deltas,
		Tickers:   tickers,
	})
}
