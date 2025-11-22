package game

import (
	"fmt"
	"regexp"
	"time"

	mdb "github.com/unifuu/hitotose/gin/db/mongo"
	"github.com/unifuu/hitotose/gin/model/game"

	mgo "github.com/unifuu/monggo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service interface {
	Badge(status game.Status) game.Badge
	ByGenre(genre game.Genre) []game.Game
	ByID(id any) game.Game
	ByPlaying() []game.Game
	ByStatus(status game.Status) []game.Game
	Create(g game.Game) error
	Delete(id string) error
	Query(kw string, pf game.Platform, st game.Status, page, limit int) ([]game.Game, int)
	TitleByID(id any) string
	Update(g game.Game) error
}

func NewService() Service {
	return &service{}
}

type service struct{}

func (s *service) Badge(status game.Status) game.Badge {
	playedCnt := countStatus(game.PLAYED)
	playingCnt := countStatus(game.PLAYING)
	toPlayCnt := countStatus(game.TO_PLAY)

	allCnt := countStatus(status)
	pcCnt := countPlatform(game.PC, status)
	nsCnt := countPlatform(game.NINTENDO_SWITCH, status)
	psCnt := countPlatform(game.PLAYSTATION, status)
	xboxCnt := countPlatform(game.XBOX, status)
	mobileCnt := countPlatform(game.MOBILE, status)

	return game.Badge{
		Played:         playedCnt,
		Playing:        playingCnt,
		ToPlay:         toPlayCnt,
		AllPlatform:    allCnt,
		PC:             pcCnt,
		PlayStation:    psCnt,
		NintendoSwitch: nsCnt,
		XBox:           xboxCnt,
		Mobile:         mobileCnt,
	}
}

func (s *service) ByID(id any) game.Game {
	var game game.Game
	mgo.FindID(mdb.Games, id).Decode(&game)
	return game
}

func (s *service) ByGenre(genre game.Genre) []game.Game {
	var games []game.Game
	filter := bson.D{primitive.E{Key: "genre", Value: genre}}
	sort := bson.D{primitive.E{Key: "title", Value: 1}}
	mgo.FindMany(mdb.Games, &games, filter, sort)
	return games
}

func (s *service) ByPlaying() []game.Game {
	var games []game.Game
	filter := bson.D{primitive.E{Key: "status", Value: game.PLAYING}}
	sort := bson.D{primitive.E{Key: "title", Value: 1}}
	mgo.FindMany(mdb.Games, &games, filter, sort)
	return games
}

func byRankingNo(rNo int) game.Game {
	var g game.Game
	filter := bson.D{primitive.E{Key: "ranking", Value: rNo}}
	result := mgo.FindOne(mdb.Games, filter)
	result.Decode(&g)
	return g
}

func (s *service) ByStatus(status game.Status) []game.Game {
	var filter bson.D
	var games []game.Game

	// Default status is "playing"
	if len(status) != 0 {
		filter = bson.D{primitive.E{Key: "status", Value: status}}
	} else {
		filter = bson.D{primitive.E{Key: "status", Value: game.PLAYING}}
	}

	sort := bson.D{primitive.E{Key: "title", Value: 1}}
	mgo.FindMany(mdb.Games, &games, filter, sort)
	return games
}

func count() int {
	filter := bson.D{}
	cnt, _ := mgo.Count(mdb.Games, filter)
	return int(cnt)
}

func countPlatform(platform game.Platform, status game.Status) int {
	filter := bson.D{
		primitive.E{Key: "platform", Value: platform},
		primitive.E{Key: "status", Value: status},
	}
	cnt, _ := mgo.Count(mdb.Games, filter)
	return int(cnt)
}

func countStatus(status game.Status) int {
	filter := bson.D{primitive.E{Key: "status", Value: status}}
	cnt, _ := mgo.Count(mdb.Games, filter)
	return int(cnt)
}

func countInRanking() int {
	filter := bson.D{primitive.E{Key: "ranking", Value: bson.D{primitive.E{Key: "$gt", Value: 0}}}}
	cnt, _ := mgo.Count(mdb.Games, filter)
	return int(cnt)
}

func (s *service) Create(g game.Game) error {
	g.ID = primitive.NewObjectIDFromTimestamp(time.Now())
	g.CreatedAt = time.Now()
	g.UpdatedAt = time.Now()
	return mgo.Insert(mdb.Games, g)
}

func (s *service) Delete(id string) error {
	return mgo.DeleteID(mdb.Games, id)
}

func (s *service) Query(kw string, platform game.Platform, status game.Status, page, limit int) ([]game.Game, int) {
	filter := bson.D{}

	// Check status
	if len(status) != 0 {
		filter = append(filter, bson.E{Key: "status", Value: status})
	}

	if platform == "all" {
		platform = "All"
	}

	// Check platform
	if len(platform) != 0 && platform != "All" {
		filter = append(filter, primitive.E{Key: "platform", Value: platform})
	}

	regex := fmt.Sprintf(".*%s.*", regexp.QuoteMeta(kw))
	filter = append(filter, bson.E{Key: "title", Value: bson.M{"$regex": regex, "$options": "i"}})

	var games []game.Game
	sort := bson.D{primitive.E{Key: "title", Value: 1}}
	totalPages, err := mgo.FindPage(mdb.Games, &games, filter, page, limit, sort)
	if err != nil {
		return nil, 1
	}
	return games, totalPages
}

func (s *service) TitleByID(id any) string {
	game := s.ByID(id)
	return game.Title
}

func (s *service) Update(g game.Game) error {
	update := bson.D{primitive.E{
		Key: "$set", Value: bson.D{
			primitive.E{Key: "title", Value: g.Title},
			primitive.E{Key: "genre", Value: g.Genre},
			primitive.E{Key: "platform", Value: g.Platform},
			primitive.E{Key: "developer", Value: g.Developer},
			primitive.E{Key: "publisher", Value: g.Publisher},
			primitive.E{Key: "status", Value: g.Status},
			primitive.E{Key: "played_time", Value: g.PlayedTime},
			primitive.E{Key: "rating", Value: g.Rating},
			primitive.E{Key: "updated_at", Value: time.Now()},
		}},
	}
	return mgo.Update(mdb.Games, g.ID, update)
}
