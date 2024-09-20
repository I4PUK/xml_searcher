package main

import (
	"encoding/xml"
	"slices"
	"sort"
	"strings"
)

type Root struct {
	XMLName  xml.Name      `xml:"root"`
	Profiles []UserProfile `xml:"row"`
}

type UserProfile struct {
	ID            int    `xml:"id"`
	GUID          string `xml:"guid"`
	IsActive      bool   `xml:"isActive"`
	Balance       string `xml:"balance"`
	Picture       string `xml:"picture"`
	Age           int    `xml:"age"`
	EyeColor      string `xml:"eye_color"`
	FirstName     string `xml:"first_name"`
	LastName      string `xml:"last_name"`
	Gender        string `xml:"gender"`
	Company       string `xml:"company"`
	Email         string `xml:"email"`
	Phone         string `xml:"phone"`
	Address       string `xml:"address"`
	About         string `xml:"about"`
	Registered    string `xml:"registered"`
	FavoriteFruit string `xml:"favorite_fruit"`
}

func findSubstringInFirstNameLastNameAndAbout(users []UserProfile, query string) []UserProfile {
	var res []UserProfile
	for _, user := range users {
		var name = user.FirstName + " " + user.LastName

		if strings.Contains(name, query) || strings.Contains(user.About, query) {
			if slices.Contains(res, user) {
				continue
			}
			res = append(res, user)
		}
	}

	if len(res) == 0 {
		return users
	}

	return res
}

func sortUsers(users []UserProfile, orderField string, orderBy int) {
	switch orderField {
	case "Id":
		sortById(users, orderBy)
	case "Age":
		sortByAge(users, orderBy)
	case "Name":
		sortByName(users, orderBy)
	case "":
		sortByName(users, orderBy)
	default:
		panic("Unknown field " + orderField)
	}
}

func sortByName(users []UserProfile, orderBy int) {
	switch orderBy {
	case OrderByAsc:
		sort.Slice(users, func(i, j int) bool {
			return users[i].FirstName < users[j].FirstName
		})
	case OrderByAsIs:
		break
	case OrderByDesc:
		sort.Slice(users, func(i, j int) bool {
			return users[i].FirstName > users[j].FirstName
		})
	default:
		panic("Bad orderBy value")
	}
}

func sortByAge(users []UserProfile, orderBy int) {
	switch orderBy {
	case OrderByAsc:
		sort.Slice(users, func(i, j int) bool {
			return users[i].Age < users[j].Age
		})
	case OrderByAsIs:
		break
	case OrderByDesc:
		sort.Slice(users, func(i, j int) bool {
			return users[i].Age > users[j].Age
		})
	default:
		panic("Bad orderBy value")
	}
}

func sortById(users []UserProfile, orderBy int) {
	switch orderBy {
	case OrderByAsc:
		sort.Slice(users, func(i, j int) bool {
			return users[i].ID < users[j].ID
		})
	case OrderByAsIs:
		break
	case OrderByDesc:
		sort.Slice(users, func(i, j int) bool {
			return users[i].ID > users[j].ID
		})
	default:
		panic("Bad orderBy value")
	}
}

func limitUsers(users []UserProfile, limit int) []UserProfile {
	if limit == 0 {
		return users
	}

	return users[:limit]
}

func offsetUsers(users []UserProfile, offset int) []UserProfile {
	offset = offset - 1

	if offset > len(users) || offset < 0 {
		offset = 0
	}

	result := users[offset:]
	return result
}
