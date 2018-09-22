package main

import (
	"sort"
)

type AllCourses struct {
	Courses map[string]*Course
}

type Course struct {
	ID               string   `json:"id"`
	Name             string   `json:"name"`
	AvailableSems    []int    `json:"available_sems"`
	PreReqs          []string `json:"prereqs"`
	Excluded         []string `json:"excluded"`
	GeneralEducation bool     `json:"GENED"`
	UOC              int      `json:"UOC"`
	Capacity         []int    `json:"capacity"`
}

type Filter struct {
	Semester []int
}

type User struct {
	Completed []string
}

func (c *Course) Fulfilled(user *User) bool {
	if len(c.PreReqs) > 0 {
		return c.NumFulfilled(user) > 0
	}

	return true
}

func (c *Course) NumFulfilled(user *User) int {
	count := 0
	for _, course := range c.PreReqs {
		if user.HasCompleted(course) {
			count++
		}
	}

	return count
}

func (c *Course) HasExcluded(user *User) bool {
	if user.HasCompleted(c.ID) {
		return true
	}

	for _, course := range c.Excluded {
		if user.HasCompleted(course) {
			return true
		}
	}

	return false
}

func (c *Course) IsEligible(user *User) bool {
	return c.Fulfilled(user) && !c.HasExcluded(user)
}

func (c *Course) Score(user *User) int {
	if len(c.Capacity) == 0 {
		return 20
	}

	return c.Capacity[0] * (c.NumFulfilled(user) + 1)
}

func (u *User) HasCompleted(course string) bool {
	for _, c := range u.Completed {
		if c == course {
			return true
		}
	}

	return false
}

func (a *AllCourses) EligibleCourses(user *User) []*Course {
	var results []*Course
	for _, course := range a.Courses {
		if course.IsEligible(user) {
			results = append(results, course)
		}
	}

	return results
}

func (a *AllCourses) Suggestions(user *User, filter func(c *Course) bool) []*Course {
	eligible := a.EligibleCourses(user)

	results := eligible
	if filter != nil {
		results = []*Course{}
		for _, course := range eligible {
			if filter(course) {
				results = append(results, course)
			}
		}
	}

	sort.Slice(results, func(i int, j int) bool {
		return results[i].Score(user) > results[j].Score(user)
	})

	return results
}
