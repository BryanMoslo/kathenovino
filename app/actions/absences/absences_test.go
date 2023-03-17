package absences_test

import (
	"kathenovino/app/models"
	"net/http"
	"net/url"
	"time"

	"github.com/gobuffalo/suite"
	"github.com/icrowley/fake"
)

type AbsenceSuite struct {
	*suite.Action
}

func (abs *AbsenceSuite) Test_Create_Absence() {
	
	form := url.Values{
		"Date":   []string{time.Now().Format("2006-01-02")},
		"Reason":   []string{fake.WordsN(5)},
	}

	res := abs.HTML("/absences/create").Post(form)
	abs.Equal(http.StatusSeeOther, res.Code)
}

func (abs *AbsenceSuite) Test_Validate_Absence_Date() {
	abs.NoError(abs.DB.Create(&models.Absence{
		Date: time.Now(),
		Reason: fake.WordsN(5),
	}))
	
	form := url.Values{
		"Date":   	[]string{time.Now().Format("2006-01-02")},
		"Reason":   []string{fake.WordsN(5)},
	}

	res := abs.HTML("/absences/create").Post(form)
	abs.Equal(http.StatusUnprocessableEntity, res.Code)
	abs.Contains(res.Body.String(), "Esta fecha ya la agregaste")
}

func (abs *AbsenceSuite) Test_Validate_Absence_Reason() {
	form := url.Values{
		"Date":   	[]string{time.Now().Format("2006-01-02")},
		"Reason":   []string{""},
	}

	res := abs.HTML("/absences/create").Post(form)
	abs.Equal(http.StatusUnprocessableEntity, res.Code)
	abs.Contains(res.Body.String(), "Hey, y la razón papi?")
}

