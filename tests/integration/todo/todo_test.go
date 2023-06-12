package todo

import (
	"testing"
	"time"
	"todo-service/src/infrastructure/authentication"
	"todo-service/src/infrastructure/delivery/graphql"
	"todo-service/src/infrastructure/storage"
	"todo-service/src/models"
	"todo-service/src/registry"
	"todo-service/tests/tools"
	"todo-service/utils"

	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var signInUser1Resp userInfo
var signInUser2Resp userInfo

type userInfo struct {
	User  models.User
	Todos []models.Todo
	signIn
}

type signIn struct {
	Auth struct {
		Data struct {
			AccessToken  string `json:"accessToken"`
			RefreshToken string `json:"refreshToken"`
		} `graphql:"signIn(email:$email, password:$password)"`
	} `json:"auth"`
}

type createTodo struct {
	Data struct {
		ID   uuid.UUID `json:"id"`
		Text string    `json:"text"`
		Done bool      `json:"done"`
		User struct {
			ID    uuid.UUID `json:"id"`
			Name  string    `json:"name"`
			Email string    `json:"email"`
		}
	} `graphql:"createTodo(input: {text:$text})"`
}

type markCompleteTodo struct {
	Data struct {
		ID   uuid.UUID `json:"id"`
		Text string    `json:"text"`
		Done bool      `json:"done"`
		User struct {
			ID    uuid.UUID `json:"id"`
			Name  string    `json:"name"`
			Email string    `json:"email"`
		}
	} `graphql:"markCompleteTodo(todoID: $todoID)"`
}

type listTodos struct {
	Todos []struct {
		ID   uuid.UUID `json:"id"`
		Text string    `json:"text"`
		Done bool      `json:"done"`
		User struct {
			ID    uuid.UUID `json:"id"`
			Name  string    `json:"name"`
			Email string    `json:"email"`
		}
	} `graphql:"todos"`
}

type deleteTodo struct {
	DeleteTodoRes bool `graphql:"deleteTodo(todoID: $todoID)"`
}

func initCreateTodoStruct(todoText, userName, userEmail string, todoDone bool, userID uuid.UUID) createTodo {
	return createTodo{
		Data: struct {
			ID   uuid.UUID `json:"id"`
			Text string    `json:"text"`
			Done bool      `json:"done"`
			User struct {
				ID    uuid.UUID `json:"id"`
				Name  string    `json:"name"`
				Email string    `json:"email"`
			}
		}{
			Text: todoText,
			Done: todoDone,
			User: struct {
				ID    uuid.UUID `json:"id"`
				Name  string    `json:"name"`
				Email string    `json:"email"`
			}{
				ID:    userID,
				Name:  userName,
				Email: userEmail,
			},
		},
	}
}

func initListTodosStruct(ui userInfo) listTodos {
	var list listTodos

	for _, todo := range ui.Todos {
		Data := struct {
			ID   uuid.UUID `json:"id"`
			Text string    `json:"text"`
			Done bool      `json:"done"`
			User struct {
				ID    uuid.UUID `json:"id"`
				Name  string    `json:"name"`
				Email string    `json:"email"`
			}
		}{
			ID:   todo.ID,
			Text: todo.Text,
			Done: todo.Done,
			User: struct {
				ID    uuid.UUID `json:"id"`
				Name  string    `json:"name"`
				Email string    `json:"email"`
			}{
				ID:    ui.User.ID,
				Name:  ui.User.Name,
				Email: ui.User.Email,
			},
		}
		list.Todos = append(list.Todos, Data)
	}
	return list
}

func initMarkCompleteTodoStruct(todoText, userName, userEmail string, todoDone bool, userID, todoID uuid.UUID) markCompleteTodo {
	return markCompleteTodo{
		Data: struct {
			ID   uuid.UUID `json:"id"`
			Text string    `json:"text"`
			Done bool      `json:"done"`
			User struct {
				ID    uuid.UUID `json:"id"`
				Name  string    `json:"name"`
				Email string    `json:"email"`
			}
		}{
			ID:   todoID,
			Text: todoText,
			Done: todoDone,
			User: struct {
				ID    uuid.UUID `json:"id"`
				Name  string    `json:"name"`
				Email string    `json:"email"`
			}{
				ID:    userID,
				Name:  userName,
				Email: userEmail,
			},
		},
	}
}

func TestTodo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Todo Suite")
}

var router *echo.Echo
var db *gorm.DB

var _ = BeforeSuite(func() {
	viper.AddConfigPath("../../../conf")
	viper.SetConfigName("test_config")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zapcore.PanicLevel)
	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	db = storage.InitPostgres(logger)

	jc := authentication.NewJwtConfigurator(logger, "../../../rsa_keys/private_key.pem",
		"../../../rsa_keys/public_key.pem")

	// Register and create controller
	useCase := registry.NewRegistry(db, jc).NewUseCase()

	router = echo.New()

	// Initialize Echo instance
	graphql.NewGraphqlRouter(router, useCase)
})

func checkCreateTodoResult(err error, wantErr error, q createTodo, wantQ createTodo) {
	if wantErr == nil {
		Expect(err).To(BeNil())
	} else {
		Expect(err).ToNot(BeNil())
	}

	Expect(q.Data.Text).To(Equal(wantQ.Data.Text))
	Expect(q.Data.Done).To(Equal(wantQ.Data.Done))
	Expect(q.Data.User.ID).To(Equal(wantQ.Data.User.ID))
	Expect(q.Data.User.Name).To(Equal(wantQ.Data.User.Name))
	Expect(q.Data.User.Email).To(Equal(wantQ.Data.User.Email))
}

func SignIn(email, pwd string) (signIn, error) {
	var q signIn
	variables := map[string]interface{}{
		"email":    email,
		"password": pwd,
	}

	err := tools.DoMutate(&q, variables, "", router)
	return q, err
}

func AddUsersToDb() {
	users := []models.User{
		{
			ID:       uuid.MustParse("48f875c5-4d1f-4eb6-abbc-5e85dae826af"),
			Name:     "test_name",
			Email:    "test@gmail.com",
			Password: utils.HashPwd("12345"),
		},
		{
			ID:       uuid.MustParse("2fd3d635-3a46-4085-a854-81f76a25cfd0"),
			Name:     "test_name2",
			Email:    "test2@gmail.com",
			Password: utils.HashPwd("12345"),
		},
	}
	result := db.Create(&users)
	if result.Error != nil {
		panic(result.Error)
	}

	signInUser1Resp.User = users[0]
	signInUser1Resp.User.Password = "12345"

	signInUser2Resp.User = users[1]
	signInUser2Resp.User.Password = "12345"

}

func AddTodosToDb() {
	signInUser1Resp.Todos = make([]models.Todo, 0)
	signInUser2Resp.Todos = make([]models.Todo, 0)
	todos := []models.Todo{
		//User ID = 48f875c5-4d1f-4eb6-abbc-5e85dae826af
		{
			ID:      uuid.MustParse("3f6f02ff-b44a-467c-a0eb-649f7f0683e1"),
			Text:    "text_todo_1",
			Done:    false,
			UserID:  uuid.MustParse("48f875c5-4d1f-4eb6-abbc-5e85dae826af"),
			Created: time.Now(),
		},
		{
			ID:      uuid.MustParse("e51baf4c-cb01-4371-a422-7d94b70951d7"),
			Text:    "text_todo_2",
			Done:    false,
			UserID:  uuid.MustParse("48f875c5-4d1f-4eb6-abbc-5e85dae826af"),
			Created: time.Now(),
		},
		{
			ID:      uuid.MustParse("9e701781-8e26-492b-9cf2-f67b4faf24f2"),
			Text:    "text_todo_3",
			Done:    false,
			UserID:  uuid.MustParse("48f875c5-4d1f-4eb6-abbc-5e85dae826af"),
			Created: time.Now(),
		},
		//User ID = 2fd3d635-3a46-4085-a854-81f76a25cfd0
		{
			ID:      uuid.MustParse("bbf14a0c-cfb4-4b91-9fbc-6f0a93fc0b32"),
			Text:    "text_todo_4",
			Done:    false,
			UserID:  uuid.MustParse("2fd3d635-3a46-4085-a854-81f76a25cfd0"),
			Created: time.Now(),
		},
		{
			ID:      uuid.MustParse("3203170d-6ed5-4a0f-afcf-9337f079590b"),
			Text:    "text_todo_5",
			Done:    false,
			UserID:  uuid.MustParse("2fd3d635-3a46-4085-a854-81f76a25cfd0"),
			Created: time.Now(),
		},
		{
			ID:      uuid.MustParse("8258b38d-8c90-4226-8e16-129088704914"),
			Text:    "text_todo_6",
			Done:    false,
			UserID:  uuid.MustParse("2fd3d635-3a46-4085-a854-81f76a25cfd0"),
			Created: time.Now(),
		},
	}

	result := db.Create(&todos)
	if result.Error != nil {
		panic(result.Error)
	}

	signInUser1Resp.Todos = append(signInUser1Resp.Todos, todos[0], todos[1], todos[2])
	signInUser2Resp.Todos = append(signInUser2Resp.Todos, todos[3], todos[4], todos[5])
}

var _ = Describe("Todo", func() {

	Describe("Create todo", func() {
		BeforeEach(func() {
			err := db.Migrator().DropTable(&models.Todo{})
			if err != nil {
				panic(err)
			}

			err = db.Migrator().DropTable(&models.User{})
			if err != nil {
				panic(err)
			}

			err = db.AutoMigrate(&models.User{})
			if err != nil {
				panic(err)
			}
			err = db.AutoMigrate(&models.Todo{})
			if err != nil {
				panic(err)
			}

			AddUsersToDb()

			signInUser1Resp.signIn, err = SignIn(signInUser1Resp.User.Email, signInUser1Resp.User.Password)
			if err != nil {
				panic(err)
			}

		})

		Context("With valid access token and valid param", func() {
			It("success create todo", func() {

				var q createTodo

				variables := map[string]interface{}{
					"text": "text_test",
				}

				wantQ := initCreateTodoStruct("text_test", signInUser1Resp.User.Name, signInUser1Resp.User.Email, false, signInUser1Resp.User.ID)

				err := tools.DoMutate(&q, variables, signInUser1Resp.Auth.Data.AccessToken, router)
				checkCreateTodoResult(err, nil, q, wantQ)

			})
		})

		Context("With valid access token and invalid name of param (text1)", func() {
			It("error: error type", func() {

				var q createTodo
				var wantQ createTodo
				variables := map[string]interface{}{
					"text1": 1,
				}

				err := tools.DoMutate(&q, variables, signInUser1Resp.Auth.Data.AccessToken, router)
				checkCreateTodoResult(err, err, q, wantQ)

			})
		})

		Context("With valid access token and invalid type of param (int)", func() {
			It("error: error type", func() {

				var q createTodo
				var wantQ createTodo
				variables := map[string]interface{}{
					"text": 1,
				}

				err := tools.DoMutate(&q, variables, signInUser1Resp.Auth.Data.AccessToken, router)
				checkCreateTodoResult(err, err, q, wantQ)

			})
		})

		Context("With valid access token and invalid type of param (bool)", func() {
			It("error: error type", func() {

				var q createTodo
				var wantQ createTodo
				variables := map[string]interface{}{
					"text": true,
				}

				err := tools.DoMutate(&q, variables, signInUser1Resp.Auth.Data.AccessToken, router)
				checkCreateTodoResult(err, err, q, wantQ)

			})
		})

		Context("With valid access token and invalid type of param (float)", func() {
			It("error: error type", func() {

				var q createTodo
				var wantQ createTodo
				variables := map[string]interface{}{
					"text": 1.5,
				}

				err := tools.DoMutate(&q, variables, signInUser1Resp.Auth.Data.AccessToken, router)
				checkCreateTodoResult(err, err, q, wantQ)

			})
		})

		Context("With invalid access token and valid param", func() {
			It("error: 401 Unauthorized", func() {

				var q createTodo
				var wantQ createTodo
				variables := map[string]interface{}{
					"text": "text_test2",
				}

				err := tools.DoMutate(&q, variables, signInUser1Resp.Auth.Data.AccessToken+"invalid", router)
				checkCreateTodoResult(err, err, q, wantQ)

			})
		})

		Context("Without access token and valid param", func() {
			It("error: Access Denied", func() {

				var q createTodo
				var wantQ createTodo
				variables := map[string]interface{}{
					"text": "text_test2",
				}

				err := tools.DoMutate(&q, variables, "", router)
				checkCreateTodoResult(err, err, q, wantQ)

			})
		})
	})

	Describe("Mark to complete todo", func() {

		BeforeEach(func() {
			err := db.Migrator().DropTable(&models.Todo{})
			if err != nil {
				panic(err)
			}

			err = db.Migrator().DropTable(&models.User{})
			if err != nil {
				panic(err)
			}

			err = db.AutoMigrate(&models.User{})
			if err != nil {
				panic(err)
			}
			err = db.AutoMigrate(&models.Todo{})
			if err != nil {
				panic(err)
			}

			AddUsersToDb()
			AddTodosToDb()

			signInUser1Resp.signIn, err = SignIn(signInUser1Resp.User.Email, signInUser1Resp.User.Password)
			if err != nil {
				panic(err)
			}

			signInUser2Resp.signIn, err = SignIn(signInUser2Resp.User.Email, signInUser2Resp.User.Password)
			if err != nil {
				panic(err)
			}

		})

		Context("With valid access token and valid todo id", func() {
			It("success mark complete todo", func() {

				var q markCompleteTodo
				wantQ := initMarkCompleteTodoStruct(signInUser1Resp.Todos[0].Text, signInUser1Resp.User.Name,
					signInUser1Resp.User.Email, true, signInUser1Resp.User.ID, signInUser1Resp.Todos[0].ID)

				variables := map[string]interface{}{
					"todoID": signInUser1Resp.Todos[0].ID.String(),
				}
				err := tools.DoMutate(&q, variables, signInUser1Resp.Auth.Data.AccessToken, router)
				Expect(err).To(BeNil())
				Expect(q).To(Equal(wantQ))

			})
		})

		Context("With valid access token and valid todo id", func() {
			It("success mark complete todo", func() {

				var q markCompleteTodo
				wantQ := initMarkCompleteTodoStruct(signInUser1Resp.Todos[0].Text, signInUser1Resp.User.Name,
					signInUser1Resp.User.Email, true, signInUser1Resp.User.ID, signInUser1Resp.Todos[0].ID)

				variables := map[string]interface{}{
					"todoID": signInUser1Resp.Todos[0].ID.String(),
				}

				err := tools.DoMutate(&q, variables, signInUser1Resp.Auth.Data.AccessToken, router)
				Expect(err).To(BeNil())
				Expect(q).To(Equal(wantQ))

				err = tools.DoMutate(&q, variables, signInUser1Resp.Auth.Data.AccessToken, router)
				Expect(err).To(BeNil())
				Expect(q).To(Equal(wantQ))

			})
		})

		Context("With valid access token and invalid todo id", func() {
			It("error: invalid input syntax for type uuid", func() {

				var q markCompleteTodo
				var wantQ markCompleteTodo
				//wantQ := initMarkCompleteTodoStruct(signInUser1Resp.Todos[0].Text, signInUser1Resp.User.Name,
				//	signInUser1Resp.User.Email, true, signInUser1Resp.User.ID, signInUser1Resp.Todos[0].ID)

				variables := map[string]interface{}{
					"todoID": signInUser1Resp.Todos[0].ID.String() + "rr",
				}
				err := tools.DoMutate(&q, variables, signInUser1Resp.Auth.Data.AccessToken, router)
				Expect(err.Error()).To(Equal("Message: ERROR: invalid input syntax for type uuid: \"" +
					signInUser1Resp.Todos[0].ID.String() + "rr" + "\" (SQLSTATE 22P02), Locations: [], Extensions: map[]"))
				Expect(q).To(Equal(wantQ))

			})
		})

		Context("With valid access token and todo id doesn't exist in db", func() {
			It("error: record not found", func() {

				var q markCompleteTodo
				var wantQ markCompleteTodo
				//wantQ := initMarkCompleteTodoStruct(signInUser1Resp.Todos[0].Text, signInUser1Resp.User.Name,
				//	signInUser1Resp.User.Email, true, signInUser1Resp.User.ID, signInUser1Resp.Todos[0].ID)

				variables := map[string]interface{}{
					"todoID": "50e89ad6-5365-4574-be31-b16eed742c09",
				}
				err := tools.DoMutate(&q, variables, signInUser1Resp.Auth.Data.AccessToken, router)
				Expect(err.Error()).To(Equal("Message: record not found, Locations: [], Extensions: map[]"))
				Expect(q).To(Equal(wantQ))

			})
		})
		Context("With valid access token and todo id another user id", func() {
			It("error: record not found", func() {

				var q markCompleteTodo
				var wantQ markCompleteTodo
				variables := map[string]interface{}{
					"todoID": signInUser2Resp.Todos[0].ID.String(),
				}
				err := tools.DoMutate(&q, variables, signInUser1Resp.Auth.Data.AccessToken, router)
				Expect(err.Error()).To(Equal("Message: record not found, Locations: [], Extensions: map[]"))
				Expect(q).To(Equal(wantQ))

			})
		})

		Context("With valid access token and invalid type of param (int)", func() {
			It("error: error type", func() {

				var q markCompleteTodo
				var wantQ markCompleteTodo
				variables := map[string]interface{}{
					"todoID": 1,
				}

				err := tools.DoMutate(&q, variables, signInUser1Resp.Auth.Data.AccessToken, router)
				Expect(err).ToNot(BeNil())
				Expect(q).To(Equal(wantQ))

			})
		})

		Context("With valid access token and invalid type of param (bool)", func() {
			It("error: error type", func() {

				var q markCompleteTodo
				var wantQ markCompleteTodo
				variables := map[string]interface{}{
					"todoID": true,
				}

				err := tools.DoMutate(&q, variables, signInUser1Resp.Auth.Data.AccessToken, router)
				Expect(err).ToNot(BeNil())
				Expect(q).To(Equal(wantQ))

			})
		})

		Context("With valid access token and invalid type of param (float)", func() {
			It("error: error type", func() {

				var q markCompleteTodo
				var wantQ markCompleteTodo
				variables := map[string]interface{}{
					"todoID": 1.5,
				}

				err := tools.DoMutate(&q, variables, signInUser1Resp.Auth.Data.AccessToken, router)
				Expect(err).ToNot(BeNil())
				Expect(q).To(Equal(wantQ))

			})
		})

		Context("With invalid access token and valid param", func() {
			It("error: 401 Unauthorized", func() {

				var q markCompleteTodo
				var wantQ markCompleteTodo
				variables := map[string]interface{}{
					"todoID": signInUser1Resp.Todos[0].ID.String(),
				}

				err := tools.DoMutate(&q, variables, signInUser1Resp.Auth.Data.AccessToken+"r", router)
				Expect(err).ToNot(BeNil())
				Expect(q).To(Equal(wantQ))

			})
		})

		Context("Without access token and valid param", func() {
			It("error: Access Denied", func() {

				var q markCompleteTodo
				var wantQ markCompleteTodo
				variables := map[string]interface{}{
					"todoID": signInUser1Resp.Todos[0].ID.String(),
				}

				err := tools.DoMutate(&q, variables, "", router)
				Expect(err).ToNot(BeNil())
				Expect(q).To(Equal(wantQ))

			})
		})
	})

	Describe("Delete todo", func() {

		BeforeEach(func() {
			err := db.Migrator().DropTable(&models.Todo{})
			if err != nil {
				panic(err)
			}

			err = db.Migrator().DropTable(&models.User{})
			if err != nil {
				panic(err)
			}

			err = db.AutoMigrate(&models.User{})
			if err != nil {
				panic(err)
			}
			err = db.AutoMigrate(&models.Todo{})
			if err != nil {
				panic(err)
			}

			AddUsersToDb()
			AddTodosToDb()

			signInUser1Resp.signIn, err = SignIn(signInUser1Resp.User.Email, signInUser1Resp.User.Password)
			if err != nil {
				panic(err)
			}

			signInUser2Resp.signIn, err = SignIn(signInUser2Resp.User.Email, signInUser2Resp.User.Password)
			if err != nil {
				panic(err)
			}

		})

		Context("With valid access token and valid todo id", func() {
			It("success delete todo", func() {

				var q deleteTodo
				wantQ := deleteTodo{
					DeleteTodoRes: true,
				}

				variables := map[string]interface{}{
					"todoID": signInUser1Resp.Todos[0].ID.String(),
				}
				err := tools.DoMutate(&q, variables, signInUser1Resp.Auth.Data.AccessToken, router)

				Expect(err).To(BeNil())
				Expect(q).To(Equal(wantQ))

			})
		})

		Context("With valid access token and valid todo id and delete the same task again", func() {
			It("error: record not found", func() {

				var q deleteTodo
				wantQ := deleteTodo{
					DeleteTodoRes: false,
				}

				variables := map[string]interface{}{
					"todoID": signInUser1Resp.Todos[0].ID.String(),
				}

				err := tools.DoMutate(&q, variables, signInUser1Resp.Auth.Data.AccessToken, router)
				Expect(err).To(BeNil())

				err = tools.DoMutate(&q, variables, signInUser1Resp.Auth.Data.AccessToken, router)
				Expect(err).To(BeNil())
				Expect(q).To(Equal(wantQ))

			})
		})

		Context("With valid access token and invalid todo id", func() {
			It("error: invalid input syntax for type uuid", func() {

				var q deleteTodo
				wantQ := deleteTodo{
					DeleteTodoRes: false,
				}

				variables := map[string]interface{}{
					"todoID": signInUser1Resp.Todos[0].ID.String() + "rr",
				}
				err := tools.DoMutate(&q, variables, signInUser1Resp.Auth.Data.AccessToken, router)
				Expect(err.Error()).To(Equal("Message: ERROR: invalid input syntax for type uuid: \"" +
					signInUser1Resp.Todos[0].ID.String() + "rr" + "\" (SQLSTATE 22P02), Locations: [], Extensions: map[]"))
				Expect(q).To(Equal(wantQ))

			})
		})

		Context("With valid access token and todo id doesn't exist in db", func() {
			It("error: record not found", func() {

				var q deleteTodo
				wantQ := deleteTodo{
					DeleteTodoRes: false,
				}

				variables := map[string]interface{}{
					"todoID": "50e89ad6-5365-4574-be31-b16eed742c09",
				}
				err := tools.DoMutate(&q, variables, signInUser1Resp.Auth.Data.AccessToken, router)
				Expect(err).To(BeNil())
				Expect(q).To(Equal(wantQ))

			})
		})
		Context("With valid access token and todo id another user id", func() {
			It("return false", func() {

				var q deleteTodo
				wantQ := deleteTodo{
					DeleteTodoRes: false,
				}

				variables := map[string]interface{}{
					"todoID": signInUser2Resp.Todos[0].ID.String(),
				}
				err := tools.DoMutate(&q, variables, signInUser1Resp.Auth.Data.AccessToken, router)
				Expect(err).To(BeNil())
				Expect(q).To(Equal(wantQ))

			})
		})

		Context("With valid access token and invalid type of id (int)", func() {
			It("error: error type", func() {

				var q deleteTodo
				wantQ := deleteTodo{
					DeleteTodoRes: false,
				}

				variables := map[string]interface{}{
					"todoID": 1,
				}

				err := tools.DoMutate(&q, variables, signInUser1Resp.Auth.Data.AccessToken, router)
				Expect(err).ToNot(BeNil())
				Expect(q).To(Equal(wantQ))

			})
		})

		Context("With valid access token and invalid type of id (bool)", func() {
			It("error: error type", func() {

				var q deleteTodo
				wantQ := deleteTodo{
					DeleteTodoRes: false,
				}

				variables := map[string]interface{}{
					"todoID": true,
				}

				err := tools.DoMutate(&q, variables, signInUser1Resp.Auth.Data.AccessToken, router)
				Expect(err).ToNot(BeNil())
				Expect(q).To(Equal(wantQ))

			})
		})

		Context("With valid access token and invalid type of id (float)", func() {
			It("error: error type", func() {

				var q deleteTodo
				wantQ := deleteTodo{
					DeleteTodoRes: false,
				}

				variables := map[string]interface{}{
					"todoID": 1.5,
				}

				err := tools.DoMutate(&q, variables, signInUser1Resp.Auth.Data.AccessToken, router)
				Expect(err).ToNot(BeNil())
				Expect(q).To(Equal(wantQ))

			})
		})

		Context("With invalid access token and valid id", func() {
			It("error: 401 Unauthorized", func() {

				var q deleteTodo
				wantQ := deleteTodo{
					DeleteTodoRes: false,
				}

				variables := map[string]interface{}{
					"todoID": signInUser1Resp.Todos[0].ID.String(),
				}

				err := tools.DoMutate(&q, variables, signInUser1Resp.Auth.Data.AccessToken+"r", router)
				Expect(err).ToNot(BeNil())
				Expect(q).To(Equal(wantQ))

			})
		})

		Context("Without access token and valid param", func() {
			It("error: Access Denied", func() {

				var q deleteTodo
				wantQ := deleteTodo{
					DeleteTodoRes: false,
				}

				variables := map[string]interface{}{
					"todoID": signInUser1Resp.Todos[0].ID.String(),
				}

				err := tools.DoMutate(&q, variables, "", router)
				Expect(err).ToNot(BeNil())
				Expect(q).To(Equal(wantQ))

			})
		})
	})

	Describe("List todos", func() {

		BeforeEach(func() {
			err := db.Migrator().DropTable(&models.Todo{})
			if err != nil {
				panic(err)
			}

			err = db.Migrator().DropTable(&models.User{})
			if err != nil {
				panic(err)
			}

			err = db.AutoMigrate(&models.User{})
			if err != nil {
				panic(err)
			}
			err = db.AutoMigrate(&models.Todo{})
			if err != nil {
				panic(err)
			}

			AddUsersToDb()
			AddTodosToDb()

			signInUser1Resp.signIn, err = SignIn(signInUser1Resp.User.Email, signInUser1Resp.User.Password)
			if err != nil {
				panic(err)
			}
		})

		Context("With valid access token", func() {
			It("success list todos", func() {

				var q listTodos
				wantQ := initListTodosStruct(signInUser1Resp)

				err := tools.DoQuery(&q, nil, signInUser1Resp.Auth.Data.AccessToken, router)
				Expect(err).To(BeNil())
				Expect(q).To(Equal(wantQ))

			})
		})

		Context("With invalid access token", func() {
			It("error: 401 Unauthorized", func() {

				var q listTodos
				var wantQ listTodos

				err := tools.DoQuery(&q, nil, signInUser1Resp.Auth.Data.AccessToken+"r", router)
				Expect(err).ToNot(BeNil())
				Expect(q).To(Equal(wantQ))

			})
		})

		Context("Without access token and valid param", func() {
			It("error: Access Denied", func() {

				var q listTodos
				var wantQ listTodos

				err := tools.DoQuery(&q, nil, "", router)
				Expect(err).ToNot(BeNil())
				Expect(q).To(Equal(wantQ))

			})
		})
	})
})
