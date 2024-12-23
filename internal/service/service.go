package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/ninedraft/core-api/internal/api"
	"github.com/ninedraft/core-api/internal/storage"
)

var _ api.StrictServerInterface = (*Service)(nil)

var (
	errBadRequest = errors.New("bad request")
)

type Service struct {
	storage *storage.Storage
}

func New(storage *storage.Storage) *Service {
	return &Service{storage: storage}
}

// Get all events
// (GET /events)
func (service *Service) GetEvents(ctx context.Context, request api.GetEventsRequestObject) (api.GetEventsResponseObject, error) {
	return &api.GetEvents500JSONResponse{
		Message: errors.ErrUnsupported.Error(),
	}, nil
}

// Create an event
// (POST /events)
func (service *Service) PostEvents(ctx context.Context, request api.PostEventsRequestObject) (api.PostEventsResponseObject, error) {
	if request.Body == nil {
		return &api.PostEvents400JSONResponse{
			Message: errBadRequest.Error(),
		}, nil
	}

	created, err := service.storage.CreateEvent(ctx, *request.Body)
	if err != nil {
		return &api.PostEvents500JSONResponse{
			Message: err.Error(),
		}, nil
	}

	return (*api.PostEvents200JSONResponse)(created), nil
}

// Delete an event
// (DELETE /events/{id})
func (service *Service) DeleteEventsId(ctx context.Context, request api.DeleteEventsIdRequestObject) (api.DeleteEventsIdResponseObject, error) {
	// not implemented
	return &api.DeleteEventsId500JSONResponse{
		Message: errors.ErrUnsupported.Error(),
	}, nil
}

// Get an event by ID
// (GET /events/{id})
func (service *Service) GetEventsId(ctx context.Context, request api.GetEventsIdRequestObject) (api.GetEventsIdResponseObject, error) {
	event, err := service.storage.GetEvent(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	return (*api.GetEventsId200JSONResponse)(event), nil
}

// Update an event
// (PATCH /events/{id})
func (service *Service) PatchEventsId(ctx context.Context, request api.PatchEventsIdRequestObject) (api.PatchEventsIdResponseObject, error) {
	updated, err := service.storage.UpdateEvent(ctx, request.Id, *request.Body)
	if err != nil {
		return &api.PatchEventsId500JSONResponse{
			Message: err.Error(),
		}, nil
	}

	return (*api.PatchEventsId200JSONResponse)(updated), nil
}

// Create a user report
// (POST /user-reports)
func (service *Service) PostUserReports(ctx context.Context, request api.PostUserReportsRequestObject) (api.PostUserReportsResponseObject, error) {
	created, err := service.storage.CreateUserReport(ctx, *request.Body)
	if err != nil {
		return &api.PostUserReports500JSONResponse{
			Message: err.Error(),
		}, nil
	}

	slog.InfoContext(ctx, "created user report", "report", created)

	return &api.PostUserReports201JSONResponse{}, nil
}

// Update a user report
// (PATCH /user-reports/{id})
func (service *Service) UpdateUserReport(ctx context.Context, request api.UpdateUserReportRequestObject) (api.UpdateUserReportResponseObject, error) {
	if request.Body == nil {
		return &api.UpdateUserReport400JSONResponse{
			Message: errBadRequest.Error(),
		}, nil
	}

	updated, err := service.storage.UpdateUserReport(ctx, request.Id, *request.Body)
	if err != nil {
		return &api.UpdateUserReport500JSONResponse{
			Message: err.Error(),
		}, nil
	}

	return (*api.UpdateUserReport200JSONResponse)(updated), nil
}

// Get users by filters
// (GET /users)
func (service *Service) GetUsers(ctx context.Context, request api.GetUsersRequestObject) (api.GetUsersResponseObject, error) {
	// not implemented
	return &api.GetUsers500JSONResponse{
		Message: errors.ErrUnsupported.Error(),
	}, nil
}

// Create a new user
// (POST /users)
func (service *Service) PostUsers(ctx context.Context, request api.PostUsersRequestObject) (api.PostUsersResponseObject, error) {
	// not implemented
	return &api.PostUsers500JSONResponse{
		Message: errors.ErrUnsupported.Error(),
	}, nil
}

// Delete a user
// (DELETE /users/{id})
func (service *Service) DeleteUsersId(ctx context.Context, request api.DeleteUsersIdRequestObject) (api.DeleteUsersIdResponseObject, error) {
	// not implemented
	return &api.DeleteUsersId500JSONResponse{
		Message: errors.ErrUnsupported.Error(),
	}, nil
}

// Update user details
// (PATCH /users/{id})
func (service *Service) PatchUsersId(ctx context.Context, request api.PatchUsersIdRequestObject) (api.PatchUsersIdResponseObject, error) {
	// not implemented
	return &api.PatchUsersId500JSONResponse{
		Message: errors.ErrUnsupported.Error(),
	}, nil
}

func (service *Service) PostUserReportsKpi(ctx context.Context, request api.PostUserReportsKpiRequestObject) (api.PostUserReportsKpiResponseObject, error) {
	err := service.storage.AddKPI(ctx, request.Body.Date.Time, float64(request.Body.Roads), float64(request.Body.Engagement))
	if err != nil {
		return &api.PostUserReportsKpi500JSONResponse{
			Message: err.Error(),
		}, nil
	}

	return api.PostUserReportsKpi201JSONResponse{}, nil
}
