package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/ninedraft/core-api/internal/api"
	"github.com/ninedraft/core-api/internal/storage/core-api/public/model"
	"github.com/ninedraft/core-api/internal/storage/core-api/public/table"
)

type Storage struct {
	db *sql.DB
}

func New(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func DialPGX(ctx context.Context, dsn string) (*sql.DB, error) {
	p, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.New: %w", err)
	}

	return stdlib.OpenDBFromPool(p), nil
}

type ctxKey struct{}

func (storage *Storage) txOrConn(ctx context.Context) qrm.DB {
	tx, ok := ctx.Value(ctxKey{}).(*sql.Tx)
	if ok {
		return tx
	}

	return storage.db
}

func (storage *Storage) InTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, ok := ctx.Value(ctxKey{}).(*sql.Tx)
	if ok {
		return fn(ctx)
	}

	tx, err := storage.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("db.BeginTx: %w", err)
	}

	err = fn(context.WithValue(ctx, ctxKey{}, tx))
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (storage *Storage) UserCreate(ctx context.Context, input api.UserBody) (*api.User, error) {
	upsertRole := table.CompanyRoles.INSERT(
		table.CompanyRoles.RoleName,
	).VALUES(input.Role).
		ON_CONFLICT(table.CompanyRoles.RoleName).
		DO_NOTHING()

	upsertPosition := table.CompanyPositions.INSERT(
		table.CompanyPositions.PositionName,
	).VALUES(input.CompanyPosition).
		ON_CONFLICT(table.CompanyPositions.PositionName).
		DO_UPDATE(postgres.SET(table.CompanyPositions.PositionName.SET(postgres.String(input.CompanyPosition))))

	departments := make([]model.Departments, 0, len(input.Departments))
	for _, dep := range input.Departments {
		departments = append(departments, model.Departments{
			DepartmentName: dep,
		})
	}

	upsertDepartment := table.Departments.INSERT(
		table.Departments.DepartmentName,
	).MODELS(departments).
		ON_CONFLICT(table.Departments.DepartmentName).
		DO_NOTHING()

	connectUserDepartment := table.DepartmentUsers.INSERT(
		table.DepartmentUsers.UserID,
		table.DepartmentUsers.DepartmentID,
	).ON_CONFLICT(
		table.DepartmentUsers.UserID,
		table.DepartmentUsers.DepartmentID,
	).DO_NOTHING()

	insertUser := table.Users.INSERT(
		table.Users.AllColumns.Except(table.Users.ID),
	).MODEL(model.Users{
		SlackID:      input.SlackId,
		IsAdmin:      input.IsAdmin,
		Role:         input.Role,
		Position:     input.CompanyPosition,
		FamilyStatus: model.FamilyStatus(input.FamilyStatus),
		Sex:          model.UserSex(input.Sex),
		JoinedAt:     input.JoinedAt.Time,
		IsActive:     input.IsActive,
	}).RETURNING(table.Users.ID)

	slog.DebugContext(ctx, "storage.UserCreate",
		"insertUser", insertUser.DebugSql(),
		"upsertRole", upsertRole.DebugSql(),
		"upsertPosition", upsertPosition.DebugSql(),
		"upsertDepartment", upsertDepartment.DebugSql(),
		"connectUserDepartment", connectUserDepartment.DebugSql(),
	)

	var id uuid.UUID

	err := storage.InTx(ctx, func(ctx context.Context) error {
		q := storage.txOrConn(ctx)

		_, errDepartment := upsertDepartment.MODELS(departments).ExecContext(ctx, q)
		if errDepartment != nil {
			return fmt.Errorf("upsertDepartment.ExecContext: %w", errDepartment)
		}

		_, errRole := upsertRole.ExecContext(ctx, q)
		if errRole != nil {
			return fmt.Errorf("upsertRole.ExecContext: %w", errRole)
		}

		_, errPosition := upsertPosition.ExecContext(ctx, q)
		if errPosition != nil {
			return fmt.Errorf("upsertPosition.ExecContext: %w", errPosition)
		}

		var generated model.Users
		errUser := insertUser.QueryContext(ctx, q, &generated)
		if errUser != nil {
			return fmt.Errorf("insertUser.ExecContext: %w", errUser)
		}

		id = generated.ID

		departmentUser := make([]model.DepartmentUsers, 0, len(input.Departments))
		for _, dep := range input.Departments {
			departmentUser = append(departmentUser, model.DepartmentUsers{
				UserID:       generated.ID,
				DepartmentID: dep,
			})
		}

		_, errConnect := connectUserDepartment.
			MODELS(departmentUser).
			ExecContext(ctx, q)
		if errConnect != nil {
			return fmt.Errorf("connectUserDepartment.ExecContext: %w", errConnect)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("storage.InTx: %w", err)
	}

	return &api.User{
		Id:   id,
		Body: input,
	}, nil
}

func (storage *Storage) CreateEvent(ctx context.Context, input api.EventBody) (*api.Event, error) {
	getUserID := table.Users.SELECT(table.Users.ID).
		WHERE(table.Users.SlackID.EQ(postgres.String(input.CreatedBySlackId)))

	slog.DebugContext(ctx, "")

	insertEvent := table.Events.INSERT(
		table.Events.AllColumns.Except(table.Events.ID),
	).RETURNING(table.Events.ID)

	slog.DebugContext(ctx, "storage.CreateEvent",
		"insertEvent", insertEvent.DebugSql(),
		"getUserID", getUserID.DebugSql())

	var id uuid.UUID

	err := storage.InTx(ctx, func(ctx context.Context) error {
		q := storage.txOrConn(ctx)

		var existingUser model.Users
		if err := getUserID.QueryContext(ctx, q, &existingUser); err != nil {
			return fmt.Errorf("getUserID.QueryContext: %w", err)
		}

		var generated model.Events
		errEvent := insertEvent.MODEL(model.Events{
			IsSent:      input.IsSent,
			CreatedBy:   existingUser.ID,
			Description: input.Description,
			Title:       input.Title,
			Identifier:  input.Identifier,
		}).QueryContext(ctx, q, &generated)
		if errEvent != nil {
			return fmt.Errorf("insertEvent.ExecContext: %w", errEvent)
		}

		id = generated.ID

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("storage.InTx: %w", err)
	}

	return &api.Event{
		Id:   id,
		Body: input,
	}, nil
}

func (storage *Storage) GetEvent(ctx context.Context, slackID string) (*api.Event, error) {
	getEvent := table.Events.SELECT(
		table.Events.AllColumns,
	).WHERE(
		table.Events.Identifier.EQ(postgres.String(slackID)),
	)

	slog.DebugContext(ctx, "storage.GetEvent", "getEvent", getEvent.DebugSql())

	var event model.Events
	err := getEvent.QueryContext(ctx, storage.txOrConn(ctx), &event)
	if err != nil {
		return nil, fmt.Errorf("getEvent.QueryContext: %w", err)
	}

	return &api.Event{
		Id: event.ID,
		Body: api.EventBody{
			Title:       event.Title,
			Description: event.Description,
			IsSent:      event.IsSent,
			Identifier:  event.Identifier,
		},
	}, nil
}

func (storage *Storage) UpdateEvent(ctx context.Context, id api.EventID, input api.EventBody) (*api.Event, error) {
	getUserID := table.Users.SELECT(table.Users.ID).
		WHERE(table.Users.SlackID.EQ(postgres.String(input.CreatedBySlackId)))

	updateEvent := table.Events.UPDATE(
		table.Events.IsSent,
		table.Events.Title,
		table.Events.Description,
	).SET(
		input.IsSent,
		input.Title,
		input.Description,
	).WHERE(
		table.Events.ID.EQ(postgres.UUID(id)),
	).RETURNING(table.Events.ID)

	slog.DebugContext(ctx, "storage.UpdateEvent",
		"updateEvent", updateEvent.DebugSql(),
		"getUserID", getUserID.DebugSql())

	err := storage.InTx(ctx, func(ctx context.Context) error {
		q := storage.txOrConn(ctx)

		var existingUser model.Users
		if err := getUserID.QueryContext(ctx, q, &existingUser); err != nil {
			return fmt.Errorf("getUserID.QueryContext: %w", err)
		}

		var generated model.Events
		errEvent := updateEvent.QueryContext(ctx, q, &generated)
		if errEvent != nil {
			return fmt.Errorf("updateEvent.ExecContext: %w", errEvent)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("storage.InTx: %w", err)
	}

	return &api.Event{
		Id:   id,
		Body: input,
	}, nil
}

func (storage *Storage) CreateCompanyDepartment(ctx context.Context, departmentName string) error {
	insertDepartment := table.Departments.INSERT(
		table.Departments.DepartmentName,
	).VALUES(departmentName).
		ON_CONFLICT(table.Departments.DepartmentName).
		DO_NOTHING()

	slog.DebugContext(ctx, "storage.CreateCompanyDepartment", "insertDepartment", insertDepartment.DebugSql())

	err := storage.InTx(ctx, func(ctx context.Context) error {
		q := storage.txOrConn(ctx)

		_, errDepartment := insertDepartment.ExecContext(ctx, q)
		if errDepartment != nil {
			return fmt.Errorf("insertDepartment.ExecContext: %w", errDepartment)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("storage.InTx: %w", err)
	}

	return nil
}

func (storage *Storage) CreateUserReport(ctx context.Context, input api.UserReportBody) (*api.UserReport, error) {
	getUserID := table.Users.SELECT(table.Users.ID).
		WHERE(table.Users.SlackID.EQ(postgres.String(input.ReportedBy)))

	insertReport := table.UserReports.INSERT(
		table.UserReports.AllColumns.Except(table.UserReports.ID),
	).RETURNING(table.UserReports.ID)

	slog.DebugContext(ctx, "storage.CreateUserReport",
		"insertReport", insertReport.DebugSql(),
		"getUserID", getUserID.DebugSql())

	var id uuid.UUID

	err := storage.InTx(ctx, func(ctx context.Context) error {
		q := storage.txOrConn(ctx)

		var existingUser model.Users
		if err := getUserID.QueryContext(ctx, q, &existingUser); err != nil {
			return fmt.Errorf("getUserID.QueryContext: %w", err)
		}

		var generated model.UserReports
		errReport := insertReport.MODEL(model.UserReports{
			ReportedBy: existingUser.ID,
			Kind:       model.ReportKind(input.Kind),
			Rating:     int32(input.Rating),
		}).QueryContext(ctx, q, &generated)
		if errReport != nil {
			return fmt.Errorf("insertReport.ExecContext: %w", errReport)
		}

		id = generated.ID
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("storage.InTx: %w", err)
	}

	return &api.UserReport{
		Id:   id,
		Body: input,
	}, nil
}

func (storage *Storage) UpdateUserReport(ctx context.Context, id uuid.UUID, input api.UserReportBody) (*api.UserReport, error) {
	getUserID := table.Users.SELECT(table.Users.ID).
		WHERE(table.Users.SlackID.EQ(postgres.String(input.ReportedBy)))

	updateReport := table.UserReports.UPDATE(
		table.UserReports.Kind,
		table.UserReports.Rating,
	).SET(
		model.ReportKind(input.Kind),
		int32(input.Rating),
	).WHERE(
		table.UserReports.ID.EQ(postgres.UUID(id)),
	).RETURNING(table.UserReports.ID)

	slog.DebugContext(ctx, "storage.UpdateUserReport",
		"updateReport", updateReport.DebugSql(),
		"getUserID", getUserID.DebugSql())

	err := storage.InTx(ctx, func(ctx context.Context) error {
		q := storage.txOrConn(ctx)

		var existingUser model.Users
		if err := getUserID.QueryContext(ctx, q, &existingUser); err != nil {
			return fmt.Errorf("getUserID.QueryContext: %w", err)
		}

		var generated model.UserReports
		errReport := updateReport.QueryContext(ctx, q, &generated)
		if errReport != nil {
			return fmt.Errorf("updateReport.ExecContext: %w", errReport)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("storage.InTx: %w", err)
	}

	return &api.UserReport{
		Id:   id,
		Body: input,
	}, nil
}

func (storage *Storage) AddKPI(ctx context.Context, date time.Time, roads, engagement float64) error {
	insertKPI := table.Kpi.
		INSERT(table.Kpi.AllColumns.Except(table.Kpi.ID)).
		MODEL(model.Kpi{
			Date:             date,
			Roads:            roads,
			EmploeEngagement: engagement,
		})

	slog.DebugContext(ctx, "storage.AddKPI", "insertKPI", insertKPI.DebugSql())

	q := storage.txOrConn(ctx)

	_, errKPI := insertKPI.ExecContext(ctx, q)
	if errKPI != nil {
		return fmt.Errorf("insertKPI.ExecContext: %w", errKPI)
	}

	return nil
}
