package brainwave_yagpdb

import (
	"encoding/json"
	"testing"

	"github.com/jbowens/codenames"
	"github.com/jinzhu/gorm"
	"github.com/jonas747/yagpdb/bot"
	"github.com/jonas747/yagpdb/commands"
	"github.com/jonas747/yagpdb/common"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var _ bot.BotInitHandler = (*Plugin)(nil)
var _ commands.CommandProvider = (*Plugin)(nil)
var _ common.Plugin = (*Plugin)(nil)

func TestGameSave(t *testing.T) {
	db, mock, _ := sqlmock.New()
	orm, _ := gorm.Open("postgres", db)
	write := brainwaveGame{
		GuildID:   1057,
		ChannelID: "378365823572",
		Game: codenames.Game{
			StartingTeam: codenames.Blue,
			Layout:       []codenames.Team{codenames.Blue, codenames.Black, codenames.Red, codenames.Neutral},
		},
	}
	gameSave, _ := json.Marshal(write.Game)

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "brainwave_games"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, 1057, "378365823572", gameSave).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(1))
	mock.ExpectCommit()
	if err := orm.Create(&write).Error; err != nil {
		t.Errorf("failed to write game: %v", err)
	}

	mock.ExpectQuery(`SELECT \* FROM "brainwave_games"`).
		WithArgs(1057, "378365823572").
		WillReturnRows(sqlmock.NewRows([]string{"game_save"}).
			AddRow(gameSave))
	var read brainwaveGame
	if err := orm.Where(&brainwaveGame{GuildID: write.GuildID, ChannelID: write.ChannelID}).
		First(&read).Error; err != nil {
		t.Errorf("failed to read game: %v", err)
	}

	if read.Game.StartingTeam != codenames.Blue {
		t.Errorf("want %d; got %d", codenames.Blue, read.Game.StartingTeam)
	}
}
