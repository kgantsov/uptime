package handler

import (
	"io/fs"

	rice "github.com/GeertJohan/go.rice"
	"github.com/labstack/echo/v4"
)

type HTTPBox struct {
	*rice.Box
}

func (hb *HTTPBox) Open(name string) (fs.File, error) {
	return hb.Box.Open(name)
}

func (h *Handler) InitStaticServer(e *echo.Echo) {
	appStaticBox, err := rice.FindBox("../../frontend/build/static/")
	if err != nil {
		e.Logger.Fatal(err)
	}

	appIndexBox, err := rice.FindBox("../../frontend/build/")
	if err != nil {
		e.Logger.Fatal(err)
	}

	e.StaticFS("/static/", &HTTPBox{appStaticBox})
	e.GET("/*", echo.StaticFileHandler("index.html", &HTTPBox{appIndexBox}))
}
