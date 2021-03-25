package common

import (
	"time"
)

type CacheDuration int

const (
	CacheNone CacheDuration = iota
	CacheShort = 1 * time.Minute
	CacheMedium = 10 * CacheShort
	CacheLong = 60 * CacheShort
	CacheDay = CacheLong * 24
	CacheWeek = CacheDay * 7
)

type Action int
const (
	ActionRead Action = iota
	ActionCreate
	ActionUpdate
	ActionDelete
	ActionRestore
)

type Key string
const (
	ActionKey Key = "VALIDATION_ACTION"
	UserKey = "USER_ID"
	RecipeKey = "RECIPE_ID"
	UserSessionKey = "USER_SESSION_ID"
)

const (
	RecipeReference string = "recipes"
)

type Queue string
const (
	QueueDefault Queue = "QUEUE_DEFAULT"
	QueueStats Queue = "QUEUE_STATS"
)

type Job string
const (
	ProcessLike Job = "ProcessLikeJob"
	ProcessDislike Job = "ProcessDislikeJob"
	ProcessView Job = "ProcessViewJob"
)
