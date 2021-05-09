package recipe

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

type (
	Tags []Tag

	Tag struct {
		Tag     TagInfo  `json:"tag"`
		Recipes []string `json:"recipes"`
	}

	TagInfo struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	TagStorage struct {
		tagsFile string
		tags     Tags
	}
)

func NewTagStorage(tagsFile string) *TagStorage {
	return &TagStorage{tagsFile, []Tag{}}
}

func (t *Tag) ContainsRecipe(recipeID string) (bool, error) {
	for _, id := range t.Recipes {
		if id == recipeID {
			return true, nil
		}
	}
	return false, nil
}

func (t *Tag) AddRecipe(recipeID string) error {
	contains, err := t.ContainsRecipe(recipeID)
	if err != nil {
		return err
	}
	if contains {
		return nil
	}
	t.Recipes = append(t.Recipes, recipeID)
	return nil
}

func (t *Tag) RemoveRecipe(recipeID string) error {
	for i, id := range t.Recipes {
		if id == recipeID {
			t.Recipes = append(t.Recipes[:i], t.Recipes[i+1:]...)
			return nil
		}
	}
	return nil
}

func (t Tags) FindByRecipeID(recipeID string) ([]TagInfo, error) {
	tags := make([]TagInfo, 0)
	for _, tag := range t {
		found, err := tag.ContainsRecipe(recipeID)
		if err != nil {
			return nil, err
		}
		if found {
			tags = append(tags, tag.Tag)
		}
	}
	return tags, nil
}

func (t Tags) FindByTagName(tagName string) (*Tag, error) {
	for _, tag := range t {
		if tag.Tag.Name == tagName {
			return &tag, nil
		}
	}
	return nil, nil
}

func (t *TagStorage) NewTag(tagName string) (*Tag, error) {
	var maxTagID int64 = 0
	for _, tag := range t.tags {
		tagID, err := strconv.ParseInt(tag.Tag.ID, 10, 64)
		if err != nil {
			return nil, err
		}
		if maxTagID < tagID {
			maxTagID = tagID
		}
	}
	newTag := Tag{
		Tag: TagInfo{
			ID:   fmt.Sprint(maxTagID + 1),
			Name: tagName,
		},
		Recipes: []string{},
	}
	t.tags = append(t.tags, newTag)
	return &newTag, nil
}

func (t *TagStorage) Load() error {
	if _, err := os.Stat(t.tagsFile); os.IsNotExist(err) {
		return nil
	}
	f, err := os.Open(t.tagsFile)
	if err != nil {
		return err
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	var tags Tags
	if err := json.Unmarshal(b, &tags); err != nil {
		return err
	}
	t.tags = tags
	return nil
}

func (t *TagStorage) Save() error {
	b, err := json.Marshal(t.tags)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(t.tagsFile, b, 0777); err != nil {
		return err
	}
	return nil
}

func (t *TagStorage) Tags() ([]TagInfo, error) {
	tagInfos := make([]TagInfo, len(t.tags))
	for _, tag := range t.tags {
		tagInfos = append(tagInfos, tag.Tag)
	}
	return tagInfos, nil
}

func (t *TagStorage) FindByRecipeID(recipeID string) ([]TagInfo, error) {
	return t.tags.FindByRecipeID(recipeID)
}

func (t *TagStorage) SetRecipeTags(recipeID string, tagNames []string) error {
	for _, tag := range t.tags {
		tag.RemoveRecipe(recipeID)
	}
	for _, tagName := range tagNames {
		tag, err := t.tags.FindByTagName(tagName)
		if err != nil {
			return err
		}
		if tag == nil {
			newTag, err := t.NewTag(tagName)
			if err != nil {
				return err
			}
			tag = newTag
		}
		if err := tag.AddRecipe(recipeID); err != nil {
			return err
		}
	}
	return t.Save()
}
