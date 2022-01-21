package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sonemaro/slow/internal/config"
	"github.com/sonemaro/slow/internal/slowmgr"
)

type Pagination struct {
	Page  int `query:"page"`
	Items int `query:"items"`
}

func setupRoutes(c *config.Config, mgr slowmgr.Manager) {
	app := fiber.New()

	app.Get("/api/v1/filter/:operation", func(c *fiber.Ctx) error {
		q := getPagination(c)
		op := c.Params("operation")
		data := mgr.Filter(op, q.Page, q.Items)
		return c.JSON(data)
	})

	app.Get("/api/v1/sort", func(c *fiber.Ctx) error {
		q := getPagination(c)
		data := mgr.Sort(q.Page, q.Items)
		return c.JSON(data)
	})

	addr := fmt.Sprintf("%s:%d", c.AppAddress, c.AppPort)
	app.Listen(addr)
}

func getPagination(c *fiber.Ctx) Pagination {
	var q Pagination
	err := c.QueryParser(&q)
	if err != nil {
		q = Pagination{
			Page:  1,
			Items: 10,
		}
	}
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.Items <= 0 {
		q.Items = 10
	}
	return q
}

func main() {
	cl := config.ViperLoader{}
	c, err := cl.Load()
	if err != nil {
		panic(err)
	}

	mgr := slowmgr.NewManagerPg(c.LogFile, slowmgr.ManagerPgOptions{Parser: slowmgr.NewParserPg()})
	err = mgr.Start()
	if err != nil {
		panic(err)
	}
	setupRoutes(c, mgr)
}
