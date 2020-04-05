package handler

import (
	"database/sql"
	"errors"
	"github.com/cobbinma/example-go-api/pkg/models"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func (h *handler) GetPet(c echo.Context) error {
	ctx := c.Request().Context()

	id, pErr := getID(c)
	if pErr != nil {
		pErr.Wrap("could not get id")
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse(pErr))
	}

	pets, pErr := h.repository.GetPet(ctx, id)
	if pErr != nil {
		if errors.Is(pErr, sql.ErrNoRows) {
			pErr.Wrap("could not find pet in repository")
			logrus.Info(pErr)
			return c.JSON(http.StatusBadRequest, models.NewErrorResponse(pErr))
		}
		pErr.Wrap("error getting pet")
		logrus.Error(pErr)
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse(pErr))
	}
	return c.JSON(http.StatusOK, pets)
}

func getID(c echo.Context) (int, models.PetError) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		pErr := models.NewPetError(err, "could not parse pet id", 0)
		pErr.Wrap("could not parse pet id")
		return 0, pErr
	}

	return id, nil
}
