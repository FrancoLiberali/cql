package conditions

import (
	"strings"

	"github.com/fatih/structtag"
)

type GormTag string

const (
	embeddedTagName       GormTag = "embedded"
	embeddedPrefixTagName GormTag = "embeddedPrefix"
	columnTagName         GormTag = "column"
	foreignKeyTagName     GormTag = "foreignKey"
	referencesTagName     GormTag = "references"
)

type GormTags map[GormTag]string

func (tags GormTags) getEmbeddedPrefix() string {
	embeddedPrefix, isPresent := tags[embeddedPrefixTagName]
	if !isPresent {
		return ""
	}

	return embeddedPrefix
}

func (tags GormTags) hasEmbedded() bool {
	_, isPresent := tags[embeddedTagName]
	return isPresent
}

func getGormTags(tag string) GormTags {
	tagMap := GormTags{}

	allTags, err := structtag.Parse(tag)
	if err != nil {
		return tagMap
	}

	gormTag, err := allTags.Get("gorm")
	if err != nil {
		return tagMap
	}

	gormTags := strings.Split(gormTag.Name, ";")
	for _, tag := range gormTags {
		splitted := strings.Split(tag, ":")
		tagName := GormTag(splitted[0])

		if len(splitted) == 1 {
			tagMap[tagName] = ""
		} else {
			tagMap[tagName] = splitted[1]
		}
	}

	return tagMap
}
