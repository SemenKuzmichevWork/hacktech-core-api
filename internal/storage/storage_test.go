package storage_test

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/ninedraft/core-api/internal/api"
	"github.com/ninedraft/core-api/internal/storage"
	"github.com/ninedraft/core-api/internal/testx"
	"github.com/oapi-codegen/runtime/types"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	testx.AssertFast(t)
	_, db := testx.Postgres(t)

	slog.SetLogLoggerLevel(slog.LevelDebug)

	st := storage.New(db)

	ctx := context.Background()

	const slackID = "test slack id"

	created, err := st.UserCreate(ctx, api.UserBody{
		CompanyPosition: "test company position",
		FamilyStatus:    "married",
		Departments: []string{
			"test department 01",
			"test department 02",
		},
		IsActive: true,
		JoinedAt: types.Date{Time: time.Now().AddDate(-1, 0, 0)},
		Role:     "admin",
		Sex:      api.Female,
		SlackId:  slackID,
	})

	require.NoError(t, err)
	require.NotEmpty(t, created.Id, "user id must be not empty")

	report, err := st.CreateUserReport(ctx, api.UserReportBody{
		Kind:       api.UserReportBodyKindBusiness,
		Rating:     5,
		ReportedBy: slackID,
	})

	require.NoError(t, err)
	require.NotEmpty(t, report.Id, "report id must be not empty")

	createdEvent, err := st.CreateEvent(ctx, api.EventBody{
		CreatedBySlackId: slackID,
		Description:      "test event description",
		Identifier:       "test slack message ID",
		Title:            "corporate event",
	})

	require.NoError(t, err)
	require.NotEmpty(t, createdEvent.Id, "event id must be not empty")

	const newDescription = "updated event description"
	updatedEvent, err := st.UpdateEvent(ctx, createdEvent.Id, api.EventBody{
		CreatedBySlackId: slackID,
		Description:      newDescription,
		Identifier:       "updated slack message ID",
		Title:            "updated corporate event",
	})

	require.NoError(t, err)
	require.NotEmpty(t, updatedEvent.Id, "event id must be not empty")
	require.Equal(t, createdEvent.Id, updatedEvent.Id, "event id must be the same")
	require.Equal(t, newDescription, updatedEvent.Body.Description, "event description must be updated")
}
