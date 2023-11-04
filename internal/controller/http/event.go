package http

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/internal/entity"
	"github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/internal/usecase/event"
	"github.com/gin-gonic/gin"
)

const (
	startTimeParam = "start"
	endTimeParam   = "end"
	monthFormat    = "2006-01"
)

type eventRoutes struct {
	events event.EventUseCase
}

func initEventRoutes(handler *gin.RouterGroup, events event.EventUseCase) {
	ar := newEventRoutes(events)
	h := handler.Group("/user")

	h.GET(fmt.Sprintf("/:%v/events", userIDParam), ar.getUserEvents)
}

func newEventRoutes(events event.EventUseCase) eventRoutes {
	return eventRoutes{events: events}
}

func (er *eventRoutes) getUserEvents(ctx *gin.Context) {
	in, err := er.retriveGetUserEventsArguments(ctx)
	if err != nil {
		abortWithErrorAnalize(ctx, fmt.Errorf("%w: %w", entity.ErrInvalidArgument, err))
		return
	}
	events, err := er.events.ReadByUserID(ctx.Request.Context(), in)
	if err != nil {
		abortWithErrorAnalize(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, map[string]string{
		"report": er.createCSVReport(events, in.UID),
	})

}

func (er *eventRoutes) retriveGetUserEventsArguments(ctx *gin.Context) (*event.InReadEventsByUserID, error) {
	uid, err := strconv.ParseUint(ctx.Param(userIDParam), 10, 32)
	if err != nil {
		return nil, fmt.Errorf("failed to parse user id %v: %w", ctx.Param(userIDParam), err)
	}
	start, err := time.Parse(monthFormat, ctx.Query(startTimeParam))
	if err != nil {
		return nil, fmt.Errorf("failed to parse start time %v: %w", ctx.Query(startTimeParam), err)
	}
	end, err := time.Parse(monthFormat, ctx.Query(endTimeParam))
	if err != nil {
		return nil, fmt.Errorf("failed to parse end time %v: %w", ctx.Query(endTimeParam), err)
	}
	return &event.InReadEventsByUserID{
		UID:   entity.UserID(uid),
		Start: &start,
		End:   &end,
	}, nil
}

func (er *eventRoutes) createCSVReport(events []*entity.Event, defaultUID entity.UserID) string {
	report := "пользователь;сегмент;операция;дата"
	for _, event := range events {
		uid := fmt.Sprint(defaultUID)
		if event.User != nil {
			uid = fmt.Sprint(*event.User)
		}
		report = fmt.Sprintf("%v\n%v;%v;%v;%v", report, uid, *event.Segment, event.Op, event.Ts)
	}
	return report
}
