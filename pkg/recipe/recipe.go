package recipe

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type (
	Recipe struct {
		ID          string       `json:"id"`
		Title       string       `json:"title"`
		Image       string       `json:"image"`
		Description string       `json:"description"`
		ServingsFor string       `json:"servingsFor"`
		Ingredients []Ingredient `json:"ingredients"`
		Steps       []Step       `json:"steps"`
		Advice      string       `json:"advice"`
		History     string       `json:"history"`
		Categories  []Category   `json:"categories"`
		Tsukurepo   int          `json:"tsukurepo"`
		Comment     int          `json:"comment"`
		Video       *string      `json:"video"`
		Published   *time.Time   `json:"published"`
		Updated     *time.Time   `json:"updated"`
	}
	Ingredient struct {
		Name     string `json:"name"`
		Quantity string `json:"quantity"`
	}
	VideoSource struct {
		URL  string
		Type string
	}
	Category []string
	Step     struct {
		Step  string `json:"step"`
		Image string `json:"image"`
	}
	RecipeStorage struct {
		Dir     string
		Recipes []Recipe
	}
	RecipeQuery struct {
		Category string
		Keyword  []string
	}
)

func NewStorage(dir string) *RecipeStorage {
	return &RecipeStorage{dir, nil}
}

func (s *RecipeStorage) Find(query *RecipeQuery) ([]Recipe, error) {
	recipes, err := s.load()
	if err != nil {
		return nil, err
	}
	result := make([]Recipe, 0)
	for _, recipe := range recipes {
		if query.Matches(&recipe) {
			result = append(result, recipe)
		}
	}
	return result, nil
}

func (s *RecipeStorage) Random(count int) ([]Recipe, error) {
	recipes, err := s.load()
	if err != nil {
		return nil, err
	}
	c := make([]Recipe, len(recipes))
	copy(c, recipes)
	result := make([]Recipe, int(math.Min(float64(count), float64(len(c)))))
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for i := 0; i < count; i++ {
		index := r.Intn(len(c))
		result[i] = c[index]
		// delete appended element.
		c = append(c[:index], c[index+1:]...)
	}
	return result, nil
}

func (q *RecipeQuery) Matches(r *Recipe) bool {
	if len(q.Category) > 0 {
		for _, c := range r.Categories {
			for _, name := range c {
				if q.Category == name {
					return true
				}
			}
		}
		return false
	}
	if len(q.Keyword) > 0 {
		for _, k := range q.Keyword {
			if len(k) == 0 {
				continue
			}
			if strings.Contains(r.Title, k) {
				return true
			}
			if strings.Contains(r.Description, k) {
				return true
			}
			for _, c := range r.Categories {
				for _, name := range c {
					if strings.Contains(name, k) {
						return true
					}
				}
			}
		}
	}
	return false
}

func (s *RecipeStorage) load() ([]Recipe, error) {
	if len(s.Recipes) > 0 {
		return s.Recipes, nil
	}
	files, err := ioutil.ReadDir(s.Dir)
	if err != nil {
		return nil, err
	}
	var recipes []Recipe
	for _, file := range files {
		if file.Name() == ".DS_Store" {
			continue
		}
		r, err := s.loadFile(filepath.Join(s.Dir, file.Name(), "data.json"))
		if err != nil {
			return nil, err
		}
		recipes = append(recipes, *r)
	}
	s.Recipes = recipes
	return recipes, nil
}

func (s *RecipeStorage) loadFile(file string) (*Recipe, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	var r Recipe
	if err := json.Unmarshal(b, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *RecipeStorage) ImportRecipeHTML(file string) error {
	r, err := parseRecipeHTML(file)
	if err != nil {
		return err
	}
	b, err := json.Marshal(r)
	if err != nil {
		return err
	}
	dir := filepath.Join(s.Dir, r.ID)
	os.MkdirAll(dir, 0777)
	dest := filepath.Join(dir, "data.json")
	if err := ioutil.WriteFile(dest, b, 0777); err != nil {
		return err
	}
	return nil
}

func parseRecipeHTML(file string) (*Recipe, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		return nil, err
	}
	id := regexp.MustCompile(`.*/(\d+)/.*`).FindStringSubmatch(file)
	title := doc.Find("#recipe-title>h1").Text()
	image, _ := doc.Find("#main-photo>img").Attr("src")
	description, _ := doc.Find("meta[property='og:description']").Attr("content")
	rep := strings.NewReplacer("（", "", "）", "")
	servingFor := rep.Replace(trim(doc.Find("span.servings_for").Text()))
	tsukurepo := parseInt(doc.Find("#tsukurepo_tab .count").Text())
	comment := parseInt(doc.Find("#comment_tab .count").Text())
	ingredients := make([]Ingredient, 0)
	doc.Find("div.ingredient_row").Each(func(i int, s *goquery.Selection) {
		ingredients = append(ingredients, Ingredient{
			Name:     trim(s.Find(".ingredient_name").Text()),
			Quantity: trim(s.Find(".ingredient_quantity").Text()),
		})
	})
	var steps []Step
	doc.Find("div.step dd").Each(func(i int, s *goquery.Selection) {
		step := trim(s.Text())
		image, _ := s.Find(".instruction img").Attr("src")
		steps = append(steps, Step{
			Step:  step,
			Image: image,
		})
	})
	advice := trim(doc.Find("#advice").Text())
	history := trim(doc.Find("#history").Text())
	categories := make([]Category, 0)
	doc.Find("#category_list li").Each(func(i int, s *goquery.Selection) {
		category := make(Category, 0)
		s.Find("a").Each(func(i int, s *goquery.Selection) {
			category = append(category, s.Text())
		})
		categories = append(categories, category)
	})
	var published *time.Time = nil
	var updated *time.Time = nil
	doc.Find("#recipe_id_and_published_date").Children().Each(func(n int, s *goquery.Selection) {
		text := trim(s.Text())
		if strings.HasPrefix(text, "公開日 : ") {
			published = parseDate(text[11:])
		} else if strings.HasPrefix(text, "更新日 : ") {
			updated = parseDate(text[11:])
		}
	})
	var video *string = nil
	doc.Find("iframe").Each(func(i int, s *goquery.Selection) {
		src, exist := s.Attr("src")
		if exist {
			if strings.Contains(src, "cookpad-video") {
				video = &src
			}
		}
	})
	return &Recipe{
		ID:          id[1],
		Title:       trim(title),
		Image:       image,
		Description: trim(description),
		ServingsFor: servingFor,
		Ingredients: ingredients,
		Steps:       steps,
		Advice:      advice,
		History:     history,
		Categories:  categories,
		Tsukurepo:   tsukurepo,
		Comment:     comment,
		Video:       video,
		Published:   published,
		Updated:     updated,
	}, nil
}

func parseDate(s string) *time.Time {
	t, err := time.Parse("06/01/02", trim(s))
	if err != nil {
		panic(err)
	}
	return &t
}

func parseInt(s string) int {
	replaced := strings.ReplaceAll(trim(s), ",", "")
	replaced = strings.ReplaceAll(replaced, "(", "")
	replaced = strings.ReplaceAll(replaced, ")", "")
	replaced = strings.ReplaceAll(replaced, "件", "")
	n, err := strconv.Atoi(replaced)
	if err != nil {
		return 0
	}
	return n
}

func trim(s string) string {
	return strings.Trim(s, "\n ")
}
