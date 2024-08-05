package srv

import (
	"time"

	"github.com/unifuu/hitotose/gin/db"
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
	PageByPlatformStatus(status game.Status, platform game.Platform, page, limit int) ([]game.Game, int)
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
	psCnt := countPlatform(game.PLAY_STATION, status)
	xboxCnt := countPlatform(game.XBOX, status)
	mobileCnt := countPlatform(game.MOBILE, status)

	return game.Badge{
		Played:       playedCnt,
		Playing:      playingCnt,
		ToPlay:       toPlayCnt,
		AllPlatform:  allCnt,
		PC:           pcCnt,
		PlayStaion:   psCnt,
		NintenSwitch: nsCnt,
		XBox:         xboxCnt,
		Mobile:       mobileCnt,
	}
}

func (s *service) ByID(id any) game.Game {
	var game game.Game
	mgo.FindID(db.Games, id).Decode(&game)
	return game
}

func (s *service) ByGenre(genre game.Genre) []game.Game {
	var games []game.Game
	filter := bson.D{primitive.E{Key: "genre", Value: genre}}
	sort := bson.D{primitive.E{Key: "title", Value: 1}}
	mgo.FindMany(db.Games, &games, filter, sort)
	return games
}

func (s *service) ByPlaying() []game.Game {
	var games []game.Game
	filter := bson.D{primitive.E{Key: "status", Value: game.PLAYING}}
	sort := bson.D{primitive.E{Key: "title", Value: 1}}
	mgo.FindMany(db.Games, &games, filter, sort)
	return games
}

func byRankingNo(rNo int) game.Game {
	var g game.Game
	filter := bson.D{primitive.E{Key: "ranking", Value: rNo}}
	result := mgo.FindOne(db.Games, filter)
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
	mgo.FindMany(db.Games, &games, filter, sort)
	return games
}

func count() int {
	filter := bson.D{}
	cnt, _ := mgo.Count(db.Games, filter)
	return int(cnt)
}

func countPlatform(platform game.Platform, status game.Status) int {
	filter := bson.D{
		primitive.E{Key: "platform", Value: platform},
		primitive.E{Key: "status", Value: status},
	}
	cnt, _ := mgo.Count(db.Games, filter)
	return int(cnt)
}

func countStatus(status game.Status) int {
	filter := bson.D{primitive.E{Key: "status", Value: status}}
	cnt, _ := mgo.Count(db.Games, filter)
	return int(cnt)
}

func countInRanking() int {
	filter := bson.D{primitive.E{Key: "ranking", Value: bson.D{primitive.E{Key: "$gt", Value: 0}}}}
	cnt, _ := mgo.Count(db.Games, filter)
	return int(cnt)
}

func (s *service) Create(g game.Game) error {
	g.ID = primitive.NewObjectIDFromTimestamp(time.Now())
	g.CreatedAt = time.Now()
	g.UpdatedAt = time.Now()
	return mgo.Insert(db.Games, g)
}

func (s *service) Delete(id string) error {
	err := mgo.DeleteID(db.Games, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) PageByPlatformStatus(status game.Status, platform game.Platform, page, limit int) ([]game.Game, int) {
	var filter bson.D

	// Check status
	if len(status) != 0 {
		filter = bson.D{primitive.E{Key: "status", Value: status}}
	} else {
		// If empty then set status to "Playing"
		filter = bson.D{primitive.E{Key: "status", Value: game.PLAYING}}
	}

	// Check platform
	if len(platform) != 0 && platform != "All" {
		filter = append(filter, primitive.E{Key: "platform", Value: platform})
	}

	var games []game.Game
	sort := bson.D{primitive.E{Key: "title", Value: 1}}
	totalPages, err := mgo.FindPage(db.Games, &games, filter, page, limit, sort)
	if err != nil {
		return nil, 1
	}
	return games, totalPages
}

func target(g game.Game) game.Game {
	// count of games in ranking
	inRankCnt := countInRanking()

	if inRankCnt == 0 {
		return g
	}

	// Get current ranking number
	curRankNo := g.Ranking

	targetRankNo := 0

	if curRankNo == 0 {
		targetRankNo = (inRankCnt / 2) + (inRankCnt % 2)
	} else {
		targetRankNo = (curRankNo / 2) + (curRankNo % 2)
	}

	for {
		t := byRankingNo(targetRankNo)

		if t.Ranking != 0 {
			return t
		}

		targetRankNo = (targetRankNo / 2) + (targetRankNo % 2)

		if targetRankNo == 0 {
			return g
		}
	}
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
			primitive.E{Key: "time_to_beat", Value: g.TimeToBeat},
			primitive.E{Key: "rating", Value: g.Rating},
			primitive.E{Key: "ranking", Value: g.Ranking},
			primitive.E{Key: "updated_at", Value: time.Now()},
		}},
	}
	return mgo.Update(db.Games, g.ID, update)
}
