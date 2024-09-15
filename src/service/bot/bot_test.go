package bot

import (
	"testing"
	"workflow/config"

	"github.com/huylqvn/httpserver/log"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

func TestBot_SendGroup(t *testing.T) {

	cfg := config.Load()

	b := Get()

	type fields struct {
		b   *tele.Bot
		log *logrus.Logger
	}
	type args struct {
		id  int64
		msg string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "happy case",
			fields: fields{
				b:   b.GetBot(),
				log: log.Get(),
			},
			args: args{
				id:  cfg.GroupID,
				msg: "test",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Bot{
				b:   tt.fields.b,
				log: tt.fields.log,
			}
			if err := b.SendGroup(tt.args.id, tt.args.msg); (err != nil) != tt.wantErr {
				t.Errorf("Bot.SendGroup() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
