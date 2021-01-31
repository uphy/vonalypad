package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/uphy/vonalypad/api"
	"github.com/uphy/vonalypad/pkg/recipe"

	"github.com/urfave/cli"
)

type App struct {
	storage *recipe.RecipeStorage
}

func main() {
	a := &App{nil}
	app := cli.NewApp()
	app.Commands = []cli.Command{
		a.parseRecipes(),
		a.search(),
		a.startServer(),
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "r,recipe-dir",
			Value: "./recipe",
		},
	}
	app.Before = func(c *cli.Context) error {
		recipeDir := c.String("recipe-dir")
		a.storage = recipe.NewStorage(recipeDir)
		return nil
	}
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}

func (a *App) startServer() cli.Command {
	return cli.Command{
		Name: "start",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "s,static-dir",
				Value: "./static",
			},
		},
		Action: func(c *cli.Context) error {
			staticDir := c.String("static-dir")
			api := api.New(a.storage, staticDir)
			return api.Start(8080)
		},
	}
}

func (a *App) parseRecipes() cli.Command {
	return cli.Command{
		Name: "parse",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "i,input-dir",
				Value: "html",
			},
		},
		Action: func(c *cli.Context) error {
			inputDir := c.String("input-dir")
			files, err := ioutil.ReadDir(inputDir)
			if err != nil {
				return err
			}
			for _, file := range files {
				if file.Name() == ".DS_Store" {
					continue
				}
				log.Printf("Parse recipe: file=%s", file.Name())
				inputRecipeDir := filepath.Join(inputDir, file.Name())
				indexFile := filepath.Join(inputRecipeDir, "index.html")
				err := a.storage.ImportRecipeHTML(indexFile)
				if err != nil {
					return err
				}
			}
			return nil
		},
	}
}

func (a *App) search() cli.Command {
	return cli.Command{
		Name: "search",
		Action: func(c *cli.Context) error {
			args := c.Args()
			if len(args) == 0 {
				return nil
			}

			result, err := a.storage.Find(&recipe.RecipeQuery{
				Keyword: args,
			})
			if err != nil {
				return err
			}
			sort.Slice(result, func(i, j int) bool {
				return result[i].Tsukurepo > result[j].Tsukurepo
			})
			for _, r := range result {
				fmt.Printf("%s (つくれぽ: %d, コメント: %d, id: %s)\n", r.Title, r.Tsukurepo, r.Comment, r.ID)
			}
			fmt.Println("Total: ", len(result))
			return nil
		},
	}
}
