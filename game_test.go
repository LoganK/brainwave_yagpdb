package brainwave_yagpdb

import (
	"encoding/json"
	"testing"

	"github.com/jbowens/codenames"
	"github.com/jinzhu/gorm"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestGameSave(t *testing.T) {
	db, mock, _ := sqlmock.New()
	orm, _ := gorm.Open("postgres", db)
	write := Game{
		GuildID:   1057,
		ChannelID: 378365823572,
		Game: codenames.Game{
			StartingTeam: codenames.Blue,
			Layout:       []codenames.Team{codenames.Blue, codenames.Black, codenames.Red, codenames.Neutral},
		},
	}
	gameSave, _ := json.Marshal(write.Game)

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "brainwave_games"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, 1057, 378365823572, sqlmock.AnyArg(), sqlmock.AnyArg(), gameSave).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(1))
	mock.ExpectCommit()
	if err := orm.Create(&write).Error; err != nil {
		t.Errorf("failed to write game: %v", err)
	}

	mock.ExpectQuery(`SELECT \* FROM "brainwave_games"`).
		WithArgs(1057, 378365823572).
		WillReturnRows(sqlmock.NewRows([]string{"game_save"}).
			AddRow(gameSave))
	var read Game
	if err := orm.Where(&Game{GuildID: write.GuildID, ChannelID: write.ChannelID}).
		First(&read).Error; err != nil {
		t.Errorf("failed to read game: %v", err)
	}

	if read.Game.StartingTeam != codenames.Blue {
		t.Errorf("want %d; got %d", codenames.Blue, read.Game.StartingTeam)
	}
}

func TestStartChecksCaptains(t *testing.T) {
	g := Game{
		Captains: Captains{
			Red: 378365823572,
		},
	}

	err := g.Start(defaultWordsEnUs)
	if err != ErrNoCaptains {
		t.Errorf("expected start to fail without captains: %v", err)
	}
}

func TestStart(t *testing.T) {
	g := Game{
		Captains: Captains{
			Red:  378365823572,
			Blue: 7398572852523,
		},
	}

	err := g.Start(defaultWordsEnUs)
	if err != nil {
		t.Errorf("expected start to succeed: %v", err)
	}
}
