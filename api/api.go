package api

import (
	"fmt"
	"sort"
	"strings"

	"github.com/uphy/vonalypad/pkg/recipe"

	"github.com/labstack/echo"
)

type API struct {
	storage *recipe.RecipeStorage
	app     *echo.Echo
}

func New(storage *recipe.RecipeStorage, staticDir string) *API {
	app := echo.New()
	a := &API{storage, app}
	app.GET("/api/search", a.Search)
	app.Static("/", staticDir)
	return &API{storage, app}
}

func (a *API) Search(ctx echo.Context) error {
	q := ctx.QueryParams().Get("q")
	query := &recipe.RecipeQuery{
		Keyword: strings.Split(q, " "),
	}
	result, err := a.storage.Find(query)
	if err != nil {
		return echo.NewHTTPError(500, fmt.Errorf("failed to search: %w", err))
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Tsukurepo > result[j].Tsukurepo
	})
	return ctx.JSON(200, result)
}

func (a *API) Start(port int) error {
	return a.app.Start(fmt.Sprintf(":%d", port))
}
