package domain

import (
	"context"
	"dailyalu-server/pkg/response"
)

type ILogService interface {
	CreateInformationLog(ctx context.Context, appRequest *response.AppError) (error)
}