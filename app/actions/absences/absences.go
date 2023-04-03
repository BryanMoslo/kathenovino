package absences

import (
	"kathenovino/app/models"
	"kathenovino/app/render"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v6"
)

var (
	// r is a buffalo/render Engine that will be used by actions
	// on this package to render render HTML or any other formats.
	r = render.Engine
)

func List(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.New("transaction failure")
	}

	absences := models.Absences{}
	if err := tx.All(&absences); err != nil {
		return errors.WithStack(errors.Wrap(err, "Error getting absences"))
	}

	c.Set("absences", absences)

	return c.Render(http.StatusOK, r.HTML("absences/index.plush.html"))
}

func New(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.New("transaction failure")
	}

	total, err := models.CalculateSalary(tx)
	if err != nil {
		return err
	}

	c.Set("total", total)
	c.Set("absence", models.Absence{
		Date: time.Now(),
	})

	return c.Render(http.StatusOK, r.HTML("absences/new.plush.html"))
}

func Create(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.New("transaction failure")
	}

	absence := models.Absence{}
	if err := c.Bind(&absence); err != nil {
		return errors.WithStack(errors.Wrap(err, "Error parsing absence"))
	}

	if verrs := absence.Validate(tx); verrs.HasAny() {
		c.Set("errors", verrs)
		c.Set("absence", absence)

		total, err := models.CalculateSalary(tx)
		if err != nil {
			return errors.WithStack(errors.Wrap(err, "Error calculating salary"))
		}

		c.Set("total", total)

		return c.Render(http.StatusUnprocessableEntity, r.HTML("absences/new.plush.html"))
	}

	if err := tx.Create(&absence); err != nil {
		return errors.WithStack(errors.Wrap(err, "Error creating absence"))
	}

	c.Flash().Add("info", "Ausencia agregada, todo bien.")
	return c.Redirect(http.StatusSeeOther, "/")
}

func Delete(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.New("transaction failure")
	}

	absence := models.Absence{}
	if err := tx.Find(&absence, c.Param("id")); err != nil {
		return errors.WithStack(errors.Wrap(err, "Error finding absence"))
	}

	if err := tx.Destroy(&absence); err != nil {
		return errors.WithStack(errors.Wrap(err, "Error removing absence"))
	}

	c.Flash().Add("info", "Eliminado, todo bien.")
	return c.Redirect(http.StatusSeeOther, "/")
}
