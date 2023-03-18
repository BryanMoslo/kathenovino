package render

import (
	"kathenovino/app/models"
	"kathenovino/app/templates"
	"kathenovino/public"
	"time"

	"github.com/gobuffalo/buffalo/render"
	"github.com/jinzhu/now"
	"github.com/leekchan/accounting"
	"github.com/wawandco/ox/pkg/buffalotools"
)

// Engine for rendering across the app, it provides
// the base for rendering HTML, JSON, XML and other formats
// while also defining thing like the base layout.
var Engine = render.New(render.Options{
	HTMLLayout:  "application.plush.html",
	TemplatesFS: templates.FS(),
	AssetsFS:    public.FS(),
	Helpers:     Helpers,
})

// Helpers available for the plush templates, there are
// some helpers that are injected by Buffalo but this is
// the list of custom Helpers.
var Helpers = map[string]interface{}{
	// partialFeeder is the helper used by the render engine
	// to find the partials that will be used, this is important
	"partialFeeder": buffalotools.NewPartialFeeder(templates.FS()),
	"formatCurrency": func(value int) string {
		ac := accounting.Accounting{Symbol: "$", Precision: 0}
		return ac.FormatMoney(value)
	},

	"isPaymentDay": func() bool {
		today := time.Now()
		tomorrow := today.Add(time.Hour * 24).Day()
		tomorrowIsSunday := today.Day() == 14 && models.IsSunday(tomorrow)

		if today.Day() == 15 || tomorrowIsSunday {
			return true
		}

		tomorrowIsSunday = now.EndOfMonth().Day() -1 == today.Day() && models.IsSunday(tomorrow)
		if today.Day() == now.EndOfMonth().Day() || tomorrowIsSunday {
			return true
		}

		return false
	},
}
