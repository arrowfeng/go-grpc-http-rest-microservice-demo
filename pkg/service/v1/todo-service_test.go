package v1

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"

	v1 "github.com/arrowfeng/go-grpc-http-rest-microservice-demo/pkg/api/v1"
	"github.com/golang/protobuf/ptypes"

	"github.com/DATA-DOG/go-sqlmock"
)

func Test_toDoService_Create(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	s := NewToDoServiceServer(db)
	tm := time.Now().In(time.UTC)
	reminder, _ := ptypes.TimestampProto(tm)

	type args struct {
		ctx context.Context
		req *v1.CreateRequest
	}

	tests := []struct {
		name    string
		args    args
		mock    func()
		want    *v1.CreateResponse
		wantErr bool
	}{
		{
			name: "OK",
			args: args{
				ctx: ctx,
				req: &v1.CreateRequest{
					Api: "v1",
					ToDo: &v1.ToDo{
						Title:       "title",
						Description: "description",
						Reminder:    reminder,
					},
				},
			},
			mock: func() {
				mock.ExpectExec("INSERT INTO ToDo").WithArgs("title", "description", tm).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			want: &v1.CreateResponse{
				Api: "v1",
				Id:  1,
			},
		},
		{
			name: "Unsupported API",
			args: args{
				ctx: ctx,
				req: &v1.CreateRequest{
					Api: "v11",
					ToDo: &v1.ToDo{
						Title:       "title",
						Description: "description",
						Reminder: &timestamp.Timestamp{
							Seconds: 1,
							Nanos:   0,
						},
					},
				},
			},
			mock:    func() {},
			wantErr: true,
		},
		{
			name: "Invaild Reminder field format",
			args: args{
				ctx: ctx,
				req: &v1.CreateRequest{
					Api: "v1",
					ToDo: &v1.ToDo{
						Title:       "title",
						Description: "description",
						Reminder: &timestamp.Timestamp{
							Seconds: 1,
							Nanos:   0,
						},
					},
				},
			},
			mock:    func() {},
			wantErr: true,
		},
		{
			name: "Insert failed",
			args: args{
				ctx: ctx,
				req: &v1.CreateRequest{
					Api: "v1",
					ToDo: &v1.ToDo{
						Title:       "title",
						Description: "description",
						Reminder:    reminder,
					},
				},
			},
			mock: func() {
				mock.ExpectExec("INSERT INTO ToDo").WithArgs("title", "description", tm).WillReturnError(errors.New("INSERT failed"))
			},
			wantErr: true,
		},
		{
			name: "LastInsert failed",
			args: args{
				ctx: ctx,
				req: &v1.CreateRequest{
					Api: "v1",
					ToDo: &v1.ToDo{
						Title:       "title",
						Description: "description",
						Reminder:    reminder,
					},
				},
			},
			mock: func() {
				mock.ExpectExec("INSERT INTO ToDo").WithArgs("title", "description", tm).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("LastInsertId failed")))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := s.Create(tt.args.ctx, tt.args.req)
			if err != nil {
				t.Errorf("toDoServiceServer.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toDoServiceServer.Create() = %v. want %v", got, tt.want)
			}
		})
	}
}
