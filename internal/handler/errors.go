package handler

import "github.com/Heiji57/ETB-BE/internal/domain"

func toAppErr(err error) *domain.AppError {
	if ae, ok := err.(*domain.AppError); ok {
		return ae
	}
	return domain.ErrInternalError
}
