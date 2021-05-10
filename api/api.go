package api

import (
	"fmt"
	"sort"
	"strings"

	"github.com/uphy/vonalypad/pkg/recipe"

	"github.com/labstack/echo"
)

type API struct {
	storage    *recipe.RecipeStorage
	tagStorage *recipe.TagStorage
	app        *echo.Echo
}

type TagRequest struct {
	Tag string `json:"tag"`
}

func New(storage *recipe.RecipeStorage, tagStorage *recipe.TagStorage, staticDir string) *API {
	app := echo.New()
	a := &API{storage, tagStorage, app}
	app.GET("/api/search", a.Search)
	app.GET("/api/random", a.SearchRandom)
	app.GET("/api/recipes/:id", a.Get)
	app.GET("/api/recipes/:id/tags", a.GetTagsByRecipeID)
	app.PUT("/api/recipes/:id/tags", a.UpdateTagsByRecipeID)
	app.POST("/api/recipes/:id/add-tag", a.AddTag)
	app.POST("/api/recipes/:id/remove-tag", a.RemoveTag)
	app.GET("/api/tags/", a.GetTags)
	app.GET("/api/tags/:id", a.GetRecipesByTagID)
	app.Static("/", staticDir)
	return &API{storage, tagStorage, app}
}

func (a *API) Get(ctx echo.Context) error {
	id := ctx.Param("id")
	r, err := a.storage.FindByID(id)
	if err != nil {
		if err == recipe.ErrNotFound {
			return echo.ErrNotFound
		}
		return err
	}
	return ctx.JSON(200, r)
}

func (a *API) GetTags(ctx echo.Context) error {
	tags, err := a.tagStorage.Tags()
	if err != nil {
		return err
	}
	return ctx.JSON(200, tags)
}

func (a *API) GetTagsByRecipeID(ctx echo.Context) error {
	id := ctx.Param("id")
	tags, err := a.tagStorage.FindByRecipeID(id)
	if err != nil {
		return err
	}
	return ctx.JSON(200, tags)
}

func (a *API) GetRecipesByTagID(ctx echo.Context) error {
	id := ctx.Param("id")
	tag, err := a.tagStorage.FindTagByTagID(id)
	if err != nil {
		return err
	}
	recipes := make([]*recipe.Recipe, len(tag.Recipes))
	for i, recipeID := range tag.Recipes {
		recipe, err := a.storage.FindByID(recipeID)
		if err != nil {
			return err
		}
		recipes[i] = recipe
	}
	return ctx.JSON(200, recipes)
}

func (a *API) UpdateTagsByRecipeID(ctx echo.Context) error {
	var tagNames []string
	if err := ctx.Bind(&tagNames); err != nil {
		return err
	}
	id := ctx.Param("id")
	err := a.tagStorage.SetRecipeTags(id, tagNames)
	if err != nil {
		return err
	}
	return ctx.NoContent(200)
}

func (a *API) AddTag(ctx echo.Context) error {
	recipeID := ctx.Param("id")
	var tag TagRequest
	if err := ctx.Bind(&tag); err != nil {
		fmt.Println(err)
		return err
	}

	err := a.tagStorage.AddRecipeTag(recipeID, tag.Tag)
	if err != nil {
		return err
	}
	return ctx.NoContent(200)
}

func (a *API) RemoveTag(ctx echo.Context) error {
	recipeID := ctx.Param("id")
	var tag TagRequest
	if err := ctx.Bind(&tag); err != nil {
		return err
	}

	err := a.tagStorage.RemoveRecipeTag(recipeID, tag.Tag)
	if err != nil {
		return err
	}
	return ctx.NoContent(200)
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

func (a *API) SearchRandom(ctx echo.Context) error {
	result, err := a.storage.Random(100)
	if err != nil {
		return echo.NewHTTPError(500, fmt.Errorf("failed to search random: %w", err))
	}
	return ctx.JSON(200, result)
}

func (a *API) Start(port int) error {
	return a.app.Start(fmt.Sprintf(":%d", port))
}
